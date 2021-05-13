package cmd

import (
	"github.com/khushmeeet/flux/fluxgen"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds your website",
	Long: `Build subcommand will gather Markdown, HTML/CSS/JS and Assets files and 
generate complete website in _site/ folder, ready to be deployed to a cloud provider.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fluxgen.FluxBuild()
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
