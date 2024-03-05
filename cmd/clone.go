/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/pipes"
	"gov.gsa.fac.cgov-util/internal/util"
	vcap "gov.gsa.fac.cgov-util/internal/vcap"

	_ "github.com/lib/pq"
)

func clone(source *vcap.CredentialsRDS, dest *vcap.CredentialsRDS) {
	psql_pipe := pipes.Psql(pipes.PG_Dump(source, Debug), dest, Debug)
	psql_pipe.Wait()
	if err := psql_pipe.Error(); err != nil {
		logging.Logger.Println("BACKUPS Pipe failed")
		os.Exit(-1)
	}
}

func clone_tables(source *vcap.CredentialsRDS, dest *vcap.CredentialsRDS) {
	table_to_schema := util.Get_table_and_schema_names(source)
	for table, schema := range table_to_schema {
		psql_pipe := pipes.Psql(pipes.PG_Dump_Table(source, schema, table, Debug), dest, Debug)
		psql_pipe.Wait()
		if err := psql_pipe.Error(); err != nil {
			logging.Logger.Printf("BACKUPS Pipe failed for %s, %s\n", schema, table)
			os.Exit(-1)
		}
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
		clone_tables(source_creds, dest_creds)
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().StringVarP(&SourceDB, "source-db", "", "", "source database (req)")
	cloneCmd.Flags().StringVarP(&DestinationDB, "destination-db", "", "", "destination database (req)")
	cloneCmd.Flags().BoolVarP(&Debug, "debug", "d", false, "Log debug statements")
	cloneCmd.MarkFlagRequired("source-db")
	cloneCmd.MarkFlagRequired("destination-db")

}
