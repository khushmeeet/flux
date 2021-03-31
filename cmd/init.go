package cmd

import (
	"github.com/khushmeeet/flux/fluxgen"
	"github.com/spf13/cobra"
	"log"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a new flux project",
	Long: `Flux init command creates a folder structure to be used for making static websites. 
Just cd into the location, where you want to create the project. 
init command will generate the root folder and sub-folders inside it.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("Only one arg must be supplied to init command")
		}
		fluxgen.FluxInit(args[0])
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
