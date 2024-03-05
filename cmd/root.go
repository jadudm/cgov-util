/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	SourceDB          string
	DestinationDB     string
	DestinationBucket string
	SHA1              string
	Debug             bool

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "backups",
		Short: "A tool for backing up, testing, and restoring",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			ver, _ := cmd.Flags().GetBool("sha")
			if ver {
				fmt.Println(SHA1)
			}
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
