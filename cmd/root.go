package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "hugo",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func ExecuteRootCmd() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
