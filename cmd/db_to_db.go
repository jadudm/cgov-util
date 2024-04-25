/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// db2dbCmd represents the db2db command
var DbToDb = &cobra.Command{
	Use:   "db_to_db",
	Short: "Copies tables from one database to another",
	Long:  `Copies tables from one database to another`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("db2db called")
	},
}

func init() {
	rootCmd.AddCommand(DbToDb)

}
