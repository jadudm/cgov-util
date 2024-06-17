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

func drop_tables(db_creds vcap.Credentials, tables []string) {
	if len(tables) > 0 {
		for _, table := range tables {
			logging.Logger.Printf("DROP dropping table %s\n", table)
			drop_pipe := pipes.Psql(
				script.Echo(fmt.Sprintf("DROP TABLE %s", table)),
				db_creds)
			// If things failed completely, we'll see an error
			err := drop_pipe.Error()
			// If the DB had a recoverable error -- e.g. the table doesn't exist,
			// it will come back via the pipe's stdout/stderr, which are combined.
			if err != nil {
				logging.Logger.Printf("DROP drop failed with error\n")
				logging.Logger.Println(drop_pipe.Error())
				os.Exit(logging.DB_DROP_ERROR)
			}

			stdout, _ := drop_pipe.String()
			if strings.Contains(stdout, "ERROR") {
				logging.Logger.Printf("DROP drop table reported an error\n")
				logging.Logger.Println(stdout)
				os.Exit(logging.DB_DROP_ERROR)
			}
		}
	}
}

// truncateCmd represents the truncate command
var dropCmd = &cobra.Command{
	Use:   "drop",
	Args:  cobra.ArbitraryArgs,
	Short: "drops one or more tables",
	Long:  `drops one or more tables`,
	Run: func(cmd *cobra.Command, args []string) {

		db_creds, err := vcap.VCS.GetCredentials("aws-rds", target_db)
		if err != nil {
			logging.Logger.Printf("DROP could not get DB credentials for %s", target_db)
			os.Exit(logging.COULD_NOT_FIND_CREDENTIALS)
		}
		truncate_tables(db_creds, args)
	},
}

var (
	drop      []string
	target_db string
)

func init() {
	rootCmd.AddCommand(dropCmd)
	parseFlags("drop", dropCmd)
}
