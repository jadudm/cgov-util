/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/pipes"
	"gov.gsa.fac.cgov-util/internal/structs"
	"gov.gsa.fac.cgov-util/internal/vcap"
)

func bucket_to_local_tables(
	db_creds vcap.Credentials,
	bucket_creds vcap.Credentials,
	s3path *structs.S3Path,
) {

	mc_pipe := pipes.McRead(
		bucket_creds,
		fmt.Sprintf("%s%s", s3path.Bucket, s3path.Key),
	).FilterLine(func(s string) string {
		if strings.Contains(s, "CREATE") {
			fmt.Printf("REPLACING IN %s\n", s)
		}
		if strings.Contains(s, "CREATE TABLE") {
			return strings.Replace(s, "CREATE TABLE", "CREATE TABLE IF NOT EXISTS", -1)
		} else if strings.Contains(s, "CREATE INDEX") {
			return strings.Replace(s, "CREATE INDEX", "CREATE INDEX IF NOT EXISTS", -1)
		} else {
			return s
		}
	})
	psql_pipe := pipes.Psql(mc_pipe, db_creds)

	exit_code := 0
	stdout, _ := mc_pipe.String()
	if strings.Contains(stdout, "ERR") {
		logging.Logger.Printf("S3TODB `mc` reported an error\n")
		logging.Logger.Println(stdout)
		exit_code = logging.PIPE_FAILURE
	}

	if mc_pipe.Error() != nil {
		logging.Logger.Println("S3TODB `dump | mc` pipe failed")
		exit_code = logging.PIPE_FAILURE
	}

	stdout, _ = psql_pipe.String()
	if strings.Contains(stdout, "ERR") {
		logging.Logger.Printf("S3TODB database reported an error\n")
		logging.Logger.Println(stdout)
		exit_code = logging.PIPE_FAILURE
	}

	if exit_code != 0 {
		os.Exit(exit_code)
	}

}

// FIXME: need s3read...
func bucket_to_cgov_tables(
	s3_creds vcap.Credentials,
	dest_db_creds vcap.Credentials,
	s3path *structs.S3Path,
) {
	s3_pipe := pipes.S3Read(
		s3_creds,
		fmt.Sprintf("%s%s", s3path.Bucket, s3path.Key),
	)
	psql_pipe := pipes.Psql(s3_pipe, dest_db_creds)

	psql_pipe.Wait()
	if err := psql_pipe.Error(); err != nil {
		logging.Logger.Println("DUMPDBTOS3 `dump | mc` pipe failed")
		os.Exit(logging.PIPE_FAILURE)
	}
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

		path_struct := parseS3Path(s3_to_db_s3path)
		// Check that we can get credentials.
		db_creds, err := vcap.VCS.GetCredentials("aws-rds", s3_to_db_db)
		if err != nil {
			logging.Logger.Printf("S3toDB could not get DB credentials for %s", s3_to_db_db)
			os.Exit(logging.COULD_NOT_FIND_CREDENTIALS)
		}

		switch os.Getenv("ENV") {
		case "LOCAL":
			fallthrough
		case "TESTING":
			bucket_creds, err := vcap.VCS.GetCredentials("user-provided", path_struct.Bucket)
			if err != nil {
				logging.Logger.Printf("S3TODB could not get minio credentials")
				os.Exit(logging.COULD_NOT_FIND_CREDENTIALS)
			}
			bucket_to_local_tables(db_creds, bucket_creds, path_struct)
		case "DEV":
			fallthrough
		case "STAGING":
			fallthrough
		case "PRODUCTION":
			bucket_creds, err := vcap.VCS.GetCredentials("aws-rds", path_struct.Bucket)
			if err != nil {
				logging.Logger.Printf("S3toDB could not get s3 credentials")
				os.Exit(logging.COULD_NOT_FIND_CREDENTIALS)
			}
			bucket_to_cgov_tables(bucket_creds, db_creds, path_struct)
		}
	},
}

var (
	s3_to_db_s3path string
	s3_to_db_db     string
)

func init() {
	rootCmd.AddCommand(S3toDBCmd)
	parseFlags("s3_to_db", S3toDBCmd)

}
