package main

import (
	"log"

	"merge-guardian/cmd/merge-guardian/cmd" // Import the cmd package
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "merge-guardian [command]",
	Short: "CLI for Merge Guardian AI, an intelligent merge management solution",
	Long:  `Merge Guardian AI CLI is an enterprise solution that combines GitHub's native merge queue capabilities with custom AI-powered conflict prediction and resolution.`,
}

func Execute() {
	// Add the analyze command to the root command
	rootCmd.AddCommand(cmd.NewAnalyzeCmd())  
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	Execute()
}
