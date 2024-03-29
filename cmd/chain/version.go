package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Major = "0"
const Minor = "1"
const Fix = "0"
const Verbal = "TX Add And Balances List"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Describes version",
	Run:   Version,
}

func Version(cmd *cobra.Command, args []string) {
	fmt.Printf("Version %s.%s.%s-beta %s\n", Major, Minor, Fix, Verbal)
}
