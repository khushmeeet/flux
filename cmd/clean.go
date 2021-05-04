package cmd

import (
	"github.com/khushmeeet/flux/fluxgen"

	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean command will clean the output directory (_site)",
	Long:  `Clean command will clean the output directory (_site)`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fluxgen.FluxClean()
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
