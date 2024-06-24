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
	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/pipes"
	"gov.gsa.fac.cgov-util/internal/vcap"
)

func truncate_tables(db_creds vcap.Credentials, tables []string) {
	if len(tables) > 0 {
		for _, table := range tables {
			logging.Logger.Printf("TRUNCATE truncating table %s\n", table)
			truncate_pipe := pipes.Psql(
				script.Echo(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)),
				db_creds)
			// If things failed completely, we'll see an error
			err := truncate_pipe.Error()
			// If the DB had a recoverable error -- e.g. the table doesn't exist,
			// it will come back via the pipe's stdout/stderr, which are combined.
			if err != nil {
				logging.Logger.Printf("TRUNCATE failed with error\n")
				logging.Logger.Println(truncate_pipe.Error())
				os.Exit(logging.DB_TRUNCATE_ERROR)
			}

			stdout, _ := truncate_pipe.String()
			if strings.Contains(stdout, "ERROR") {
				logging.Logger.Printf("TRUNCATE database reported an error\n")
				logging.Logger.Println(stdout)
				os.Exit(logging.DB_TRUNCATE_ERROR)
			}
		}
	}
}

// truncateCmd represents the truncate command
var truncateCmd = &cobra.Command{
	Use:   "truncate",
	Args:  cobra.ArbitraryArgs,
	Short: "Truncates one or more tables",
	Long:  `Truncates one or more tables`,
	Run: func(cmd *cobra.Command, args []string) {

		db_creds, err := vcap.VCS.GetCredentials("aws-rds", truncate_db)
		if err != nil {
			logging.Logger.Printf("TRUNCATE could not get DB credentials for %s", truncate_db)
			os.Exit(logging.COULD_NOT_FIND_CREDENTIALS)
		}
		truncate_tables(db_creds, args)
	},
}

var (
	truncate_truncate []string
	truncate_db       string
)

func init() {
	rootCmd.AddCommand(truncateCmd)
	parseFlags("truncate", truncateCmd)

}
