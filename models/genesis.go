package models

import (
	"encoding/json"
	"io/ioutil"
)

type Genesis struct {
	Balances map[Account]uint `json:"balances"`
}

func loadGenesis(path string) (Genesis, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return Genesis{}, err
	}

	var loadedGenesis Genesis
	if err = json.Unmarshal(content, &loadedGenesis); err != nil {
		return Genesis{}, err
	}

	return loadedGenesis, nil
}
