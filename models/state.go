package models

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Snapshot [32]byte

type State struct {
	Balances  map[Account]uint
	TxMempool []Tx
	dbFile    *os.File
	Snapshot  Snapshot
}

func NewStateFromDisk() (*State, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	gen, err := loadGenesis(filepath.Join(cwd, "database", "genesis.json"))
	if err != nil {
		return nil, err
	}

	balances := make(map[Account]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	f, err := os.OpenFile(filepath.Join(cwd, "database", "tx.db"), os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)

	state := &State{
		Balances:  balances,
		TxMempool: make([]Tx, 0),
		dbFile:    f,
	}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		var tx Tx
		json.Unmarshal(scanner.Bytes(), &tx)

		if err := state.apply(tx); err != nil {
			return nil, err
		}
	}

	return state, nil
}

func (s *State) Add(tx Tx) error {
	if err := s.apply(tx); err != nil {
		return err
	}

	s.TxMempool = append(s.TxMempool, tx)

	return nil
}

func (s *State) Persist() (Snapshot, error) {
	mempool := make([]Tx, len(s.TxMempool))
	copy(mempool, s.TxMempool)

	for i := 0; i < len(mempool); i++ {
		txJson, err := json.Marshal(s.TxMempool[i])
		if err != nil {
			return Snapshot{}, err
		}
		fmt.Printf("Persisting new TX to disk:\n")
		fmt.Printf("\t%s\n", txJson)

		if _, err = s.dbFile.Write(append(txJson, '\n')); err != nil {
			return Snapshot{}, err
		}

		if err = s.doSnapshot(); err != nil {
			return Snapshot{}, err
		}
		fmt.Printf("New DB Snapshot: %x\n", s.Snapshot)
		s.TxMempool = append(s.TxMempool[:i], s.TxMempool[i+1:]...)
	}

	return s.Snapshot, nil
}

func (s *State) Close() {
	s.dbFile.Close()
}

func (s *State) apply(tx Tx) error {
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}

	if s.Balances[tx.From] < tx.Value {
		return fmt.Errorf("Insufficient balance")
	}

	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value

	return nil
}

func (s *State) LatestSnapshot() Snapshot {
	return s.Snapshot
}

func (s *State) doSnapshot() error {
	if _, err := s.dbFile.Seek(0, 0); err != nil {
		return err
	}

	txsData, err := ioutil.ReadAll(s.dbFile)
	if err != nil {
		return err
	}
	s.Snapshot = sha256.Sum256(txsData)

	return nil
}
