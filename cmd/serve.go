package cmd

import (
	"todo-cli/api"

	"github.com/spf13/cobra"
)

// serveCmd starts the Gin server
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",
	Run: func(cmd *cobra.Command, args []string) {
		api.StartServer() // Start the Gin server
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
