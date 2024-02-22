/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/bitfield/script"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gov.gsa.fac.backups/internal/logging"
	vcap "gov.gsa.fac.backups/internal/vcap"

	_ "github.com/lib/pq"
)

// https://bitfieldconsulting.com/golang/scripting
func pg_dump(creds *vcap.RDSCreds) *script.Pipe {
	// Compose the command as a slice
	cmd := []string{
		"pg_dump",
		"--clean",
		"--no-password",
		"--if-exists",
		"--no-privileges",
		"--no-owner",
		"--format plain",
		fmt.Sprintf("--host %s", creds.Host),
		fmt.Sprintf("--port %s", creds.Port),
		fmt.Sprintf("--username %s", creds.Username),
		creds.DBName,
	}
	// Combine the slice for printing and execution.
	combined := strings.Join(cmd[:], " ")
	logging.Logger.Printf("BACKUPS Running `%s`\n", combined)
	return script.Exec(combined)
}

func psql(in_pipe *script.Pipe, creds *vcap.RDSCreds) *script.Pipe {
	cmd := []string{
		"psql",
		"--no-password",
		fmt.Sprintf("--host %s", creds.Host),
		fmt.Sprintf("--port %s", creds.Port),
		fmt.Sprintf("--username %s", creds.Username),
		creds.DBName,
	}
	combined := strings.Join(cmd[:], " ")
	logging.Logger.Printf("BACKUPS Running `%s`\n", combined)
	return in_pipe.Exec(combined)
}

func clone(source *vcap.RDSCreds, dest *vcap.RDSCreds) {
	psql_pipe := psql(pg_dump(source), dest)
	psql_pipe.Wait()
	if err := psql_pipe.Error(); err != nil {
		logging.Logger.Println("BACKUPS Pipe failed")
		os.Exit(-1)
	}
}

func get_row_count(creds *vcap.RDSCreds, table string) int {
	var count int
	// FIXME: Not sure if `disable` is correct for RDS sslmode.
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		creds.Username,
		creds.Password,
		creds.Host,
		creds.DBName)
	db, _ := sql.Open("postgres", connStr)
	defer db.Close()
	row := db.QueryRow(fmt.Sprintf("SELECT count(*) FROM %s", table))
	if err := row.Scan(&count); err != nil {
		logging.Logger.Printf("BACKUPS Could not get count of %s", table)
	}
	return count
}

// snapshotDbToDbCmd represents the snapshotDbToDb command
var cloneDBToDB = &cobra.Command{
	Use:   "clone-db-to-db",
	Short: "Clones one database to another, destructively",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		source, dest := vcap.GetCreds()
		clone(source, dest)
		// FIXME: These won't exist in the VCAP_SERVICES version
		// of the config. We'll have to always... load both?
		// There needs to be a way to configure this in the remote env.
		for _, table := range viper.GetStringSlice("check-counts") {
			source_row_count := get_row_count(source, table)
			dest_row_count := get_row_count(dest, table)
			logging.Logger.Printf("BACKUPS table %s source %d dest %d",
				table, source_row_count, dest_row_count)
		}
	},
}

func init() {
	rootCmd.AddCommand(cloneDBToDB)
}
