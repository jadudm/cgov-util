package util

import (
	"database/sql"
	"os"

	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/vcap"
)

func Get_table_and_schema_names(source_creds *vcap.CredentialsRDS) map[string]string {
	// Do this table-by-table for RAM reasons.
	db, err := sql.Open("postgres", source_creds.Uri)
	if err != nil {
		logging.Logger.Println("BACKUPS could not connect to DB for table-by-table dump")
		logging.Logger.Printf("BACKUPS %s\n", err)
		os.Exit(-1)
	}

	tables, err := db.Query("SELECT schemaname, tablename FROM pg_tables WHERE schemaname = 'public'")
	if err != nil {
		logging.Logger.Println("BACKUPS could not get table names for table-by-table dump")
		logging.Logger.Printf("BACKUPS %s\n", err)
		os.Exit(-1)
	}

	table_names := make(map[string]string, 0)

	for tables.Next() {
		var table string
		var schema string
		if err := tables.Scan(&schema, &table); err != nil {
			logging.Logger.Println("BACKUPS could not scan table names in SELECT")
			os.Exit(-1)
		}
		table_names[table] = schema
	}

	return table_names
}
