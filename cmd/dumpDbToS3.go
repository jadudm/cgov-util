/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"slices"

	_ "github.com/lib/pq"

	"github.com/spf13/cobra"
	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/pipes"

	"gov.gsa.fac.cgov-util/internal/vcap"
)

var (
	db     string
	bucket string
	s3path string
	key    string
)

func get_table_and_schema_names(source_creds vcap.Credentials) map[string]string {
	// Do this table-by-table for RAM reasons.
	db, err := sql.Open("postgres", source_creds.Get("uri").String())
	if err != nil {
		logging.Logger.Println("DUMPDBTOS3 could not connect to DB for table-by-table dump")
		logging.Logger.Printf("DUMPDBTOS3 %s\n", err)
		os.Exit(logging.DB_SCHEMA_SCAN_FAILURE)
	}

	tables, err := db.Query("SELECT schemaname, tablename FROM pg_tables WHERE schemaname = 'public'")
	if err != nil {
		logging.Logger.Println("DUMPDBTOS3 could not get table names for table-by-table dump")
		logging.Logger.Printf("DUMPDBTOS3 %s\n", err)
		os.Exit(logging.DB_SCHEMA_SCAN_FAILURE)
	}

	table_names := make(map[string]string, 0)

	for tables.Next() {
		var table string
		var schema string
		if err := tables.Scan(&schema, &table); err != nil {
			logging.Logger.Println("DUMPDBTOS3 could not scan table names in SELECT")
			os.Exit(logging.DB_SCHEMA_SCAN_FAILURE)
		}
		table_names[table] = schema
	}

	return table_names
}
func bucket_local_tables(source_creds vcap.Credentials, up vcap.Credentials, table_names []string) {
	var BACKUP_ALL = len(table_names) == 0

	logging.Logger.Printf("DUMPDBTOS3 backing up from %s to %s\n",
		source_creds.Get("name").String(),
		up.Get("name").String(),
	)
	table_to_schema := get_table_and_schema_names(source_creds)

	for table, schema := range table_to_schema {
		// Back up tables under two conditions:
		// 1. When it is in a list of names we want backed up, or
		// 2. When there are no names in the list (backup all).
		if slices.Contains(table_names, table) || BACKUP_ALL {
			mc_pipe := pipes.Mc(
				pipes.PG_Dump_Table(source_creds, schema, table),
				up,
				fmt.Sprintf("%s%s/%s-%s.dump", bucket, key, schema, table),
			)
			mc_pipe.Wait()
			if err := mc_pipe.Error(); err != nil {
				logging.Logger.Println("DUMPDBTOS3 `dump | mc` pipe failed")
				os.Exit(logging.PIPE_FAILURE)
			}
		}
	}

}

func bucket_cgov_tables(source_creds vcap.Credentials, up vcap.Credentials, table_names []string) {
	var BACKUP_ALL = len(table_names) == 0

	table_to_schema := get_table_and_schema_names(source_creds)
	for table, schema := range table_to_schema {
		if slices.Contains(table_names, table) || BACKUP_ALL {
			s3_pipe := pipes.S3(
				pipes.PG_Dump_Table(source_creds, schema, table),
				up,
				fmt.Sprintf("%s%s/%s-%s.dump", bucket, key, schema, table),
			)
			s3_pipe.Wait()
			if err := s3_pipe.Error(); err != nil {
				logging.Logger.Println("DUMPDBTOS3 `dump | s3` pipe failed")
				os.Exit(logging.PIPE_FAILURE)
			}
		}
	}
}

// dumpDbToS3Cmd represents the dumpDbToS3 command
var dumpDbToS3Cmd = &cobra.Command{
	Use:   "dumpDbToS3",
	Args:  cobra.ArbitraryArgs,
	Short: "Dumps a full database to a file in S3",
	Long: `Dumps a full database to a file in S3
Takes 0 or more table names as arguments. If no arguments are
provided, all tables are backed up.
	`,
	Run: func(cmd *cobra.Command, table_names []string) {
		u, err := url.Parse(s3path)
		if err != nil {
			logging.Logger.Printf("DUMPDBTOS3 could not parse s3 path: %s", s3path)
			os.Exit(logging.S3_PATH_PARSE_ERROR)
		}
		if u.Scheme != "s3" {
			logging.Logger.Printf("DUMPDBTOS3 does not look like an S3 path (e.g. `s3://`): %s", s3path)
			os.Exit(logging.S3_PATH_PARSE_ERROR)
		}
		bucket = filepath.Clean(u.Host)
		key = filepath.Clean(u.Path)

		// Check that we can get credentials.
		db_creds, err := vcap.VCS.GetCredentials("aws-rds", db)
		if err != nil {
			logging.Logger.Printf("DUMPDBTOS3 could not get DB credentials for %s", db)
			os.Exit(logging.COULD_NOT_FIND_CREDENTIALS)
		}

		switch os.Getenv("ENV") {
		case "LOCAL":
			fallthrough
		case "TESTING":
			up, err := vcap.VCS.GetCredentials("user-provided", bucket)
			if err != nil {
				logging.Logger.Printf("DUMPDBTOS3 could not get minio credentials")
				os.Exit(logging.COULD_NOT_FIND_CREDENTIALS)
			}
			bucket_local_tables(db_creds, up, table_names)
		case "DEV":
			fallthrough
		case "STAGING":
			fallthrough
		case "PRODUCTION":
			up, err := vcap.VCS.GetCredentials("aws-rds", bucket)
			if err != nil {
				logging.Logger.Printf("DUMPDBTOS3 could not get s3 credentials")
				os.Exit(logging.COULD_NOT_FIND_CREDENTIALS)
			}
			bucket_cgov_tables(db_creds, up, table_names)

		}
	},
}

func init() {
	rootCmd.AddCommand(dumpDbToS3Cmd)
	dumpDbToS3Cmd.Flags().StringVarP(&db, "db", "", "", "source database label")
	dumpDbToS3Cmd.Flags().StringVarP(&s3path, "s3path", "", "", "destination path")

	dumpDbToS3Cmd.MarkFlagRequired("db")
	dumpDbToS3Cmd.MarkFlagRequired("bucket")
	dumpDbToS3Cmd.MarkFlagRequired("s3path")

}
