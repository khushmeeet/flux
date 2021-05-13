package cmd

import (
	"github.com/khushmeeet/flux/fluxgen"
	"log"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve command will run http server on _site folder",
	Long:  `Serve command will run http server on _site folder`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			log.Fatalf("%v", err)
		}
		fluxgen.FluxServe(port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	serveCmd.Flags().StringP("port", "p", "5050", "Port on which go http server is running")
}
