/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	db       string
	bucket   string
	s3path   string
	key      string
	truncate string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "backups",
		Short: "A tool for backing up, testing, and restoring",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("sha", "s", false, "Print short build SHA")
}
