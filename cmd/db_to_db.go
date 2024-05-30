/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/pipes"
	"gov.gsa.fac.cgov-util/internal/structs"
	"gov.gsa.fac.cgov-util/internal/util"
	"gov.gsa.fac.cgov-util/internal/vcap"
)

var (
	source_db string
	dest_db   string
	operation string
)

func get_table_and_schema_names_db(source_creds vcap.Credentials) map[string]string {
	// Do this table-by-table for RAM reasons.
	db, err := sql.Open("postgres", source_creds.Get("uri").String())
	if err != nil {
		logging.Logger.Println("DBTODB could not connect to DB for table-by-table dump")
		logging.Logger.Printf("DBTODB %s\n", err)
		os.Exit(logging.DB_SCHEMA_SCAN_FAILURE)
	}

	tables, err := db.Query("SELECT schemaname, tablename FROM pg_tables WHERE schemaname = 'public'")
	if err != nil {
		logging.Logger.Println("DBTODB could not get table names for table-by-table dump")
		logging.Logger.Printf("DBTODB %s\n", err)
		os.Exit(logging.DB_SCHEMA_SCAN_FAILURE)
	}

	table_names := make(map[string]string, 0)

	for tables.Next() {
		var table string
		var schema string
		if err := tables.Scan(&schema, &table); err != nil {
			logging.Logger.Println("DBTODB could not scan table names in SELECT")
			os.Exit(logging.DB_SCHEMA_SCAN_FAILURE)
		}
		table_names[table] = schema
	}

	return table_names
}

func LocalDatabaseSync(
	source_db_creds vcap.Credentials,
	dest_db_creds vcap.Credentials,
	table_names []string,
	operation string) {
	var BACKUP_ALL = len(table_names) == 0

	logging.Logger.Println("DBTODB " + source_db_creds.Get("name").String() + " to " + dest_db_creds.Get("name").String() + " starting")
	table_to_schema := get_table_and_schema_names_db(source_db_creds)
	//pg_dump -t table_to_copy source_db | psql target_db
	for table, schema := range table_to_schema {
		if slices.Contains(table_names, table) || BACKUP_ALL {
			switch operation {
			case "initial":
				logging.Logger.Printf("Initial db2db operation, truncate not required")
			case "backup":
				fallthrough
			case "restore":
				truncate_tables(dest_db_creds, []string{table})
			default:
				logging.Logger.Printf("Correct operation not supplied. Please supply initial, backup, or restore")
				os.Exit(-1)
			}

			psql_write := pipes.Psql(
				pipes.PG_Dump_Table(source_db_creds, schema, table),
				dest_db_creds,
			)
			psql_write.Wait()
			stdout, _ := psql_write.String()
			if strings.Contains(stdout, "ERR") {
				logging.Logger.Println("DBTODB " + source_db_creds.Get("name").String() + " to " + dest_db_creds.Get("name").String() + " pipe failed")
				os.Exit(logging.PIPE_FAILURE)
			}
		}
	}
}

func CgovDatabaseSync(
	source_db_creds vcap.Credentials,
	dest_db_creds vcap.Credentials,
	table_names []string,
	operation string) {
	var BACKUP_ALL = len(table_names) == 0

	logging.Logger.Println("DBTODB " + source_db_creds.Get("name").String() + " to " + dest_db_creds.Get("name").String() + " starting")
	table_to_schema := get_table_and_schema_names_db(source_db_creds)
	//pg_dump -t table_to_copy source_db | psql target_db
	for table, schema := range table_to_schema {
		if slices.Contains(table_names, table) || BACKUP_ALL {
			switch operation {
			case "initial":
				logging.Logger.Printf("Initial db2db operation, truncate not required")
			case "backup":
				fallthrough
			case "restore":
				truncate_tables(dest_db_creds, []string{table})
			default:
				logging.Logger.Printf("Correct operation not supplied. Please supply initial, backup, or restore")
				os.Exit(-1)
			}

			psql_write := pipes.Psql(
				pipes.PG_Dump_Table(source_db_creds, schema, table),
				dest_db_creds,
			)
			psql_write.Wait()
			stdout, _ := psql_write.String()
			if strings.Contains(stdout, "ERR") {
				logging.Logger.Println("DBTODB " + source_db_creds.Get("name").String() + " to " + dest_db_creds.Get("name").String() + " pipe failed")
				os.Exit(logging.PIPE_FAILURE)
			}
		}
	}
}

// db2dbCmd represents the db2db command
var DbToDb = &cobra.Command{
	Use:   "db_to_db",
	Short: "Copies tables from one database to another",
	Long:  `Copies tables from one database to another`,
	Run: func(cmd *cobra.Command, table_names []string) {
		fmt.Println("db2db called")
		util.UnsetProxy()
		source_db_creds := getDBCredentials(source_db)
		dest_db_creds := getDBCredentials(dest_db)

		ch := structs.Choice{
			Local: func() {
				LocalDatabaseSync(source_db_creds, dest_db_creds, table_names, operation)
			},
			Remote: func() {
				CgovDatabaseSync(source_db_creds, dest_db_creds, table_names, operation)
			}}
		runLocalOrRemote(ch)

	},
}

func init() {
	rootCmd.AddCommand(DbToDb)
	parseFlags("db_to_db", DbToDb)
	// ./gov.gsa.fac.cgov-util.exe db_to_db --src_db fac-db --dest_db fac-snapshot-db
}
