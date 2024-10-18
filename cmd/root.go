package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd is the base command for the CLI
var RootCmd = &cobra.Command{
	Use:   "todo",
	Short: "Todo CLI Application",
	Long: `Welcome to the Todo CLI!

This CLI application helps you manage your todo tasks through a REST API server. 
You can start the server, fetch todos, create tasks, and much more.`,
	// Custom logic for what to do when no subcommands are provided
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("===================================")
		fmt.Println(" Todo CLI - Command Line Interface")
		fmt.Println("===================================")
		fmt.Println("Available commands:")
		fmt.Println("  - serve: Start the API server")
		fmt.Println("  - get: Fetch all todos")
		fmt.Println("  - help: Display help information")
		fmt.Println("\nRun 'todo [command]' to use the CLI.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	// RootCmd.AddCommand(createCmd)
	// RootCmd.AddCommand(getCmd)
	// // RootCmd.AddCommand(getCmdAll) not required its attached from separate file
	// RootCmd.AddCommand(updateCmd)
	// RootCmd.AddCommand(deleteCmd)
	// RootCmd.AddCommand(registerCmd) // Add register command
	// RootCmd.AddCommand(loginCmd)    // Add login command
	// RootCmd.AddCommand(logoutCmd)   // Add logout command

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
