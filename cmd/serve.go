package cmd

import (
	"github.com/khushmeeet/flux/fluxgen"
	"log"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Go HTTP server",
	Long:  `Serve command will run a local HTTP server to serve _site/ content`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			log.Fatalf("%v", err)
		}

		watch, err := cmd.Flags().GetBool("watch")
		if err != nil {
			log.Fatalf("%v", err)
		}

		fluxgen.WatchAndServe(port, watch)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP("port", "p", "5050", "port on which go http server is running")
	serveCmd.Flags().BoolP("watch", "w", false, "detect file changes and restart the server")
}
