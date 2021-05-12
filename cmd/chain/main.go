package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var chainCmd = &cobra.Command{
		Use:   "chain",
		Short: "Chain CLI",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	chainCmd.AddCommand(versionCmd)
	chainCmd.AddCommand(balancesCmd())
	chainCmd.AddCommand(txCmd())

	err := chainCmd.Execute()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func incorrectUsageErr() error {
	return fmt.Errorf("incorrect usage")
}
