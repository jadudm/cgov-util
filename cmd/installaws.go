/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"slices"

	"github.com/spf13/cobra"
	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/util"
)

// installawsCmd represents the installaws command
var installawsCmd = &cobra.Command{
	Use:   "installaws",
	Short: "Install AWS CLI on cloud.gov instances",
	Long:  `Install AWS CLI on cloud.gov instances`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Installing aws cli")
	},
}

func init() {
	rootCmd.AddCommand(installawsCmd)
	if slices.Contains([]string{"DEV", "PREVIEW", "STAGING", "PRODUCTION"}, os.Getenv("ENV")) {
		logging.Logger.Printf("ENV detected to be a cloud.gov environment. Installing AWS CLI.")
		util.InstallAWS()
	} else {
		logging.Logger.Printf("ENV set to local or testing, aws not necessary to install.")
	}
}
