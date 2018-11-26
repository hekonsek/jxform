package cmd

import (
	"github.com/hekonsek/jxform/forms"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(provisionCmd)
}

var provisionCmd = &cobra.Command{
	Use: "provision",
	Run: func(cmd *cobra.Command, args []string) {
		err := forms.Provision()
		if err != nil {
			panic(err)
		}
	},
}
