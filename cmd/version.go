package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints version number of Flux",
	Long:  `Current version of Flux installed on the system`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("flux version 0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
