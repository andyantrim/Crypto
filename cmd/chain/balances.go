package main

import (
	"fmt"
	"os"

	"github.com/andyantrim/Crypto/models"
	"github.com/spf13/cobra"
)

func balancesCmd() *cobra.Command {
	var balancesCmd = &cobra.Command{
		Use:   "balances",
		Short: "Interact with balances",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	balancesCmd.AddCommand(balancesListCmd)

	return balancesCmd
}

var balancesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all balances",
	Run:   listAll,
}

func listAll(cmd *cobra.Command, args []string) {
	state, err := models.NewStateFromDisk()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	defer state.Close()

	fmt.Println("Account balances:")
	fmt.Println("_________________")
	fmt.Println("")

	for account, balance := range state.Balances {
		fmt.Printf("%s: %d\n", account, balance)
	}
}
