package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/veith4f/go-straight/pkg/project"
)

func main() {

	var moduleName string
	var author string

	// Create the command
	var cmd = &cobra.Command{
		Use:   "go-straight -m module-name -c copyright-name <path/to/projectdir>",
		Short: "Create a straight project directory with tests, devcontainer docker-composea and makefile.",
		Args:  cobra.ExactArgs(1), // Ensures that exactly one argument (the directory path) is passed
		Run: func(cmd *cobra.Command, args []string) {
			// Get the directory path from arguments
			projectDir := args[0]

			// Validate the module name flag
			if moduleName == "" {
				fmt.Println("Error: Module name (-m) is required.")
				return
			}

			project, err := project.NewProject(projectDir, moduleName, author)
			if err != nil {
				fmt.Println("Error reading template files:", err)
				os.Exit(1)
			}

			err2 := project.GoStraight()
			if err2 != nil {
				fmt.Println("Error writing project files:", err2)
				os.Exit(1)
			}

			fmt.Printf("straight go-project %s created.\n", projectDir)
		},
	}

	cmd.Flags().StringVarP(&moduleName, "module-name", "m", "", "Go module name used for this project.")
	cmd.Flags().StringVarP(&author, "author", "a", "", "Project author's full name.")

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
