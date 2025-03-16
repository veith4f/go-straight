package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"{{.ModuleName}}/pkg"
)

// Define the root command
var rootCmd = &cobra.Command{
	Use:   "$project_name",
	Short: "A simple CLI app",
	Long:  "A longer description of your app.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Fib(10): %d\n", pkg.Fib(10))
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
