/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"gov.gsa.fac.backups/internal/logging"
	"gov.gsa.fac.backups/internal/pipes"
	vcap "gov.gsa.fac.backups/internal/vcap"

	_ "github.com/lib/pq"
)

func clone(source *vcap.CredentialsRDS, dest *vcap.CredentialsRDS) {
	psql_pipe := pipes.Psql(pipes.PG_Dump(source), dest)
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
		source_creds, dest_creds := vcap.GetRDSCreds(SourceDB, DestinationDB)
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
