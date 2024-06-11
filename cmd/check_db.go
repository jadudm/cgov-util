/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/vcap"
)

var (
	source_database string
)

func check_if_table_exists(source_creds vcap.Credentials) {
	//SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'schema_name' AND tablename  = 'table_name');
	db, err := sql.Open("postgres", source_creds.Get("uri").String())
	if err != nil {
		logging.Logger.Println("TABLECHECK could not connect to DB for checking table existance")
		logging.Logger.Printf("DBTOS3 %s\n", err)
		os.Exit(logging.DB_SCHEMA_SCAN_FAILURE)
	}

	file, err := os.Open("db_tables.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var not_existing []string
	for scanner.Scan() {
		//scanner.Text()
		query := fmt.Sprintf("select * from %s ;", scanner.Text())
		//query := fmt.Sprintf("SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename  = '%s')", scanner.Text())
		rows, table_check := db.Query(query)

		if table_check == nil {
			//fmt.Println(i + " exists")
			rows.Next()
		} else {
			//logging.Logger.Println(scanner.Text() + " does not exist")
			// store all scanner.Text() into a map
			// if map != nil
			// hard exit
			not_existing = append(not_existing, scanner.Text())
		}
	}
	if len(not_existing) > 0 {
		logging.Logger.Println("An array of tables that does not exist in the database, but does exist in a manifest has been returned.")
		logging.Logger.Println("System exiting...")
		joined_tables := strings.Join(not_existing[:], " ")
		logging.Logger.Printf(joined_tables)
		os.Exit(3)
	} else {
		logging.Logger.Printf("Manifest and Database tables appear to be in sync.")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	//for table := range table_to_schema {
	// for table := range list_of_tables {
	// 	//"SELECT schemaname, tablename FROM pg_tables WHERE schemaname = 'public'"
	// 	query := fmt.Sprintf("SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename  = '%s')", table)
	// 	rows, table_check := db.Query(query)

	// 	if table_check == nil {
	// 		fmt.Printf(table + " exists")
	// 		rows.Close()
	// 	} else {
	// 		fmt.Println(table + " does not exist")
	// 	}
	// }
	// exists, err := db.Query(query)
	// if err != nil {
	// 	logging.Logger.Println("DBTOS3 could not get table names to check if it exists.")
	// 	logging.Logger.Printf("DBTOS3 %s\n", err)
	// 	os.Exit(logging.DB_SCHEMA_SCAN_FAILURE)
	// }
}

// https://stackoverflow.com/a/18479916
// func readLines(path string) string {
// 	// file, err := os.Open(path)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// defer file.Close()

// 	// var list_of_tables []string
// 	// scanner := bufio.NewScanner(file)
// 	// for scanner.Scan() {
// 	// 	list_of_tables = append(list_of_tables, scanner.Text())
// 	// }
// 	// return list_of_tables, scanner.Err()

// 	file, err := os.Open(path)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		I := scanner.Text()
// 	}

// 	if err := scanner.Err(); err != nil {
// 		log.Fatal(err)
// 	}

// }

// https://stackoverflow.com/a/15323988
// func stringInSlice(table string, list_of_tables []string) bool {
// 	logging.Logger.Println(reflect.DeepEqual(table_to_schema, list_of_tables))
// 	eq := reflect.DeepEqual(table, list_of_tables)
// 	if eq {
// 		logging.Logger.Println("Database and Manifest appear to be in sync.")
// 		return true
// 	} else {
// 		logging.Logger.Println("Database and Manifest appear to differ.")
// 		return false
// 	}
// 	for _, i := range list_of_tables {
// 		if table == i {
// 			//logging.Logger.Printf(table + " exists in manifest and database.")
// 			//logging.Logger.Printf("table: " + table + " appears to be missing.")
// 			return true
// 		}
// 	}
// 	logging.Logger.Printf("table: " + table + " appears to be missing.")
// 	return false
// }

// checkDbCmd represents the checkDb command
var checkDbCmd = &cobra.Command{
	Use:   "check_db",
	Short: "A brief description of your command",
	Long:  `A`,
	Run: func(cmd *cobra.Command, args []string) {
		db_creds := getDBCredentials(source_database)
		//stringInSlice(table, list_of_tables)
		check_if_table_exists(db_creds)

	},
}

func init() {
	rootCmd.AddCommand(checkDbCmd)
	parseFlags("check_tables", checkDbCmd)
}
