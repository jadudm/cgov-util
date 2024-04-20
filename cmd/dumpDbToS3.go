/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/spf13/cobra"
	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/pipes"

	"gov.gsa.fac.cgov-util/internal/vcap"
)

var (
	db     string
	s3path string
)

func get_table_and_schema_names(source_creds vcap.Credentials) map[string]string {
	// Do this table-by-table for RAM reasons.
	db, err := sql.Open("postgres", source_creds.Get("uri").String())
	if err != nil {
		logging.Logger.Println("DUMPDBTOS3 could not connect to DB for table-by-table dump")
		logging.Logger.Printf("DUMPDBTOS3 %s\n", err)
		os.Exit(-1)
	}

	tables, err := db.Query("SELECT schemaname, tablename FROM pg_tables WHERE schemaname = 'public'")
	if err != nil {
		logging.Logger.Println("DUMPDBTOS3 could not get table names for table-by-table dump")
		logging.Logger.Printf("DUMPDBTOS3 %s\n", err)
		os.Exit(-1)
	}

	table_names := make(map[string]string, 0)

	for tables.Next() {
		var table string
		var schema string
		if err := tables.Scan(&schema, &table); err != nil {
			logging.Logger.Println("DUMPDBTOS3 could not scan table names in SELECT")
			os.Exit(-1)
		}
		table_names[table] = schema
	}

	return table_names
}
func bucket_local_tables(source_creds vcap.Credentials, up vcap.Credentials) {
	logging.Logger.Printf("DUMPDBTOS3 backing up from %s to %s\n",
		source_creds.Get("name").String(),
		up.Get("name").String(),
	)
	table_to_schema := get_table_and_schema_names(source_creds)
	for table, schema := range table_to_schema {
		mc_pipe := pipes.Mc(
			pipes.PG_Dump_Table(source_creds, schema, table),
			up,
			fmt.Sprintf("%s/%s-%s.dump", s3path, schema, table),
		)
		mc_pipe.Wait()
		if err := mc_pipe.Error(); err != nil {
			logging.Logger.Println("DUMPDBTOS3 `dump | mc` pipe failed")
			os.Exit(-1)
		}
	}
}

func bucket_cgov_tables(source_creds vcap.Credentials, up vcap.Credentials) {
	table_to_schema := get_table_and_schema_names(source_creds)
	for table, schema := range table_to_schema {
		s3_pipe := pipes.S3(
			pipes.PG_Dump_Table(source_creds, schema, table),
			up,
			fmt.Sprintf("%s/%s-%s.dump", s3path, schema, table),
		)
		s3_pipe.Wait()
		if err := s3_pipe.Error(); err != nil {
			logging.Logger.Println("DUMPDBTOS3 `dump | s3` pipe failed")
			os.Exit(-1)
		}
	}
}

// dumpDbToS3Cmd represents the dumpDbToS3 command
var dumpDbToS3Cmd = &cobra.Command{
	Use:   "dumpDbToS3",
	Short: "Dumps a full database to a file in S3",
	Long:  `Dumps a full database to a file in S3`,
	Run: func(cmd *cobra.Command, args []string) {

		// Check that we can get credentials.
		db_creds, err := vcap.VCS.GetCredentials("aws-rds", db)
		if err != nil {
			logging.Logger.Printf("DUMPDBTOS3 could not get DB credentials for %s", db)
			os.Exit(-1)
		}

		switch os.Getenv("ENV") {
		case "LOCAL":
			fallthrough
		case "TESTING":
			up, err := vcap.VCS.GetCredentials("user-provided", "backups")
			if err != nil {
				logging.Logger.Printf("DUMPDBTOS3 could not get minio credentials")
				os.Exit(-1)
			}
			bucket_local_tables(db_creds, up)
		case "DEV":
			fallthrough
		case "STAGING":
			fallthrough
		case "PRODUCTION":
			up, err := vcap.VCS.GetCredentials("aws-rds", s3path)
			if err != nil {
				logging.Logger.Printf("DUMPDBTOS3 could not get s3 credentials")
				os.Exit(-1)
			}
			bucket_cgov_tables(db_creds, up)

		}
	},
}

func init() {
	rootCmd.AddCommand(dumpDbToS3Cmd)
	dumpDbToS3Cmd.Flags().StringVarP(&db, "db", "", "", "source database label")
	dumpDbToS3Cmd.Flags().StringVarP(&s3path, "s3path", "", "", "destination path")

	dumpDbToS3Cmd.MarkFlagRequired("db")
	dumpDbToS3Cmd.MarkFlagRequired("s3path")

}
