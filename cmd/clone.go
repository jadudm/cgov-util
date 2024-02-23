/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitfield/script"
	"github.com/spf13/cobra"
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
		"--dbname",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			creds.Username,
			creds.Password,
			creds.Host,
			creds.Port,
			creds.DB_Name,
		),
	}
	// Combine the slice for printing and execution.
	combined := strings.Join(cmd[:], " ")
	// This will log the password...
	// logging.Logger.Printf("BACKUPS Running `%s`\n", combined)
	logging.Logger.Printf("BACKUPS pg_dump targeting %s", creds.DB_Name)
	return script.Exec(combined)
}

func psql(in_pipe *script.Pipe, creds *vcap.RDSCreds) *script.Pipe {
	cmd := []string{
		"psql",
		"--no-password",
		"--dbname",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			creds.Username,
			creds.Password,
			creds.Host,
			creds.Port,
			creds.DB_Name,
		),
	}
	combined := strings.Join(cmd[:], " ")
	logging.Logger.Printf("BACKUPS psql targeting %s", creds.DB_Name)
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

// snapshotDbToDbCmd represents the snapshotDbToDb command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "pg_dump | psql, source to destination",
	Long: `
An imperfect, point-in-time snapshot.
	
This command copies the source database to the destination 
database by streaming STDOUT of 'pg_dump' piped into the STNDIN 
of 'psql'. The former reads from the FAC production database, and 
writes to a snapshot clone DB.
`,
	Run: func(cmd *cobra.Command, args []string) {
		source_creds, dest_creds := vcap.GetCreds(SourceDB, DestinationDB)
		clone(source_creds, dest_creds)
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().StringVarP(&SourceDB, "source-db", "", "", "source database (req)")
	cloneCmd.Flags().StringVarP(&DestinationDB, "destination-db", "", "", "destination database (req)")
	cloneCmd.MarkFlagRequired("source-db")
	cloneCmd.MarkFlagRequired("destination-db")

}
