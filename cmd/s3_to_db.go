/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/bitfield/script"
	"github.com/spf13/cobra"
	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/pipes"
	"gov.gsa.fac.cgov-util/internal/vcap"
)

func bucket_to_local_tables(db_creds vcap.Credentials, bucket_creds vcap.Credentials) {

	if truncate != "" {
		logging.Logger.Printf("S3TODB truncating table %s\n", truncate)

		truncate_pipe := pipes.Psql(script.Echo(fmt.Sprintf("TRUNCATE TABLE %s", truncate)), db_creds)
		truncate_pipe.Wait()
	}

	mc_pipe := pipes.McRead(
		bucket_creds,
		fmt.Sprintf("%s%s", bucket, key),
	)
	psql_pipe := pipes.Psql(mc_pipe, db_creds)

	psql_pipe.Wait()
	if err := mc_pipe.Error(); err != nil {
		logging.Logger.Println("DUMPDBTOS3 `dump | mc` pipe failed")
		os.Exit(logging.PIPE_FAILURE)

	}
}

func bucket_to_cgov_tables(source_creds vcap.Credentials, up vcap.Credentials) {

}

// S3toDBCmd represents the S3toDB command
var S3toDBCmd = &cobra.Command{
	Use:   "s3_to_db",
	Args:  cobra.ArbitraryArgs,
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		parseS3Path()
		// Check that we can get credentials.
		db_creds, err := vcap.VCS.GetCredentials("aws-rds", db)
		if err != nil {
			logging.Logger.Printf("S3toDB could not get DB credentials for %s", db)
			os.Exit(logging.COULD_NOT_FIND_CREDENTIALS)
		}

		switch os.Getenv("ENV") {
		case "LOCAL":
			fallthrough
		case "TESTING":
			bucket_creds, err := vcap.VCS.GetCredentials("user-provided", bucket)
			if err != nil {
				logging.Logger.Printf("S3TODB could not get minio credentials")
				os.Exit(logging.COULD_NOT_FIND_CREDENTIALS)
			}
			bucket_to_local_tables(db_creds, bucket_creds)
		case "DEV":
			fallthrough
		case "STAGING":
			fallthrough
		case "PRODUCTION":
			bucket_creds, err := vcap.VCS.GetCredentials("aws-rds", bucket)
			if err != nil {
				logging.Logger.Printf("S3toDB could not get s3 credentials")
				os.Exit(logging.COULD_NOT_FIND_CREDENTIALS)
			}
			bucket_to_cgov_tables(db_creds, bucket_creds)

		}
	},
}

func init() {
	rootCmd.AddCommand(S3toDBCmd)
	S3toDBCmd.Flags().StringVarP(&s3path, "s3path", "", "", "destination path")
	S3toDBCmd.Flags().StringVarP(&db, "db", "", "", "source database label")
	S3toDBCmd.Flags().StringVarP(&truncate, "truncate", "", "", "table to truncate before load")

	S3toDBCmd.MarkFlagRequired("db")
	S3toDBCmd.MarkFlagRequired("s3path")

}
