package main

import (
	"log"

	"merge-guardian/cmd/merge-guardian/cmd" // Import the cmd package
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "merge-guardian [command]",
	Short: "CLI for Merge Guardian AI, an intelligent merge management solution",
	Long:  `Merge Guardian AI CLI is an enterprise solution that combines GitHub's native merge queue capabilities with custom AI-powered conflict prediction and resolution.`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Merge Guardian",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Merge Guardian version %s (commit: %s, date: %s)", version, commit, date)
	},
}

func Execute() {
	// Add commands
	rootCmd.AddCommand(cmd.NewAnalyzeCmd())
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	Execute()
}
