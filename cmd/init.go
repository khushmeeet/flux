package cmd

import (
	"github.com/khushmeeet/flux/fluxgen"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init <name>",
	Short: "Creates a new flux project",
	Long: `Flux init command creates a folder structure to be used for making static websites. 
Just cd into the location, where you want to create the project. 
init command will create root folder, sub-folders and config.json file with prefilled data.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fluxgen.FluxInit(args[0])
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
