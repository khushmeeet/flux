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
}
