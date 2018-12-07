package cmd

import (
	"github.com/hekonsek/jxform/forms"
	"github.com/spf13/cobra"
)

var provisionCmdFlagVerbose bool

func init() {
	provisionCmd.Flags().BoolVarP(&provisionCmdFlagVerbose, "verbose", "", false, "")
	rootCmd.AddCommand(provisionCmd)
}

var provisionCmd = &cobra.Command{
	Use: "provision",
	Run: func(cmd *cobra.Command, args []string) {
		err := forms.Provision(provisionCmdFlagVerbose)
		if err != nil {
			panic(err)
		}
	},
}
