/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/structs"

	vcap "gov.gsa.fac.cgov-util/internal/vcap"
)

func get_row_count(creds *structs.CredentialsRDS, table string) int {
	var count int
	// FIXME: Not sure if `disable` is correct for RDS sslmode.
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		creds.Username,
		creds.Password,
		creds.Host,
		creds.DB_Name)
	db, _ := sql.Open("postgres", connStr)
	defer db.Close()
	row := db.QueryRow(fmt.Sprintf("SELECT count(*) FROM %s", table))
	if err := row.Scan(&count); err != nil {
		logging.Logger.Printf("BACKUPS Could not get count of %s", table)
	}
	return count
}

func check_results(source *structs.CredentialsRDS, dest *structs.CredentialsRDS, tables []string) {
	// FIXME: These won't exist in the VCAP_SERVICES version
	// of the config. We'll have to always... load both?
	// There needs to be a way to configure this in the remote env.
	for _, table := range tables {
		source_row_count := get_row_count(source, table)
		dest_row_count := get_row_count(dest, table)
		logging.Logger.Printf("CHECK OK %s source %d dest %d",
			table, source_row_count, dest_row_count)
		if source_row_count < dest_row_count {
			logging.Logger.Printf("CHECK too many rows in '%s' source (%d < %d)",
				table, source_row_count, dest_row_count)
			os.Exit(-1)
		}
	}
}

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks table counts between source and destination",
	Long: `
When given a source and destination, this command returns 0 when the 
number of rows in the source are equal to or higher than the number 
of rows in the destination.

This is because the clone tool is used against live tables. It is likely
that the source will increase between the time of the clone and the check.

Expects a space-separated list of table names as arguments.
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		source_creds, dest_creds := vcap.GetRDSCreds(SourceDB, DestinationDB)
		check_results(source_creds, dest_creds, args)

	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringVarP(&SourceDB, "source-db", "", "", "source database (req)")
	checkCmd.Flags().StringVarP(&DestinationDB, "destination-db", "", "", "destination database (req)")
	checkCmd.MarkFlagRequired("source-db")
	checkCmd.MarkFlagRequired("destination-db")
}
