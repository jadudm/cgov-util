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
	"gov.gsa.fac.cgov-util/internal/structs"
	"gov.gsa.fac.cgov-util/internal/util"
	"gov.gsa.fac.cgov-util/internal/vcap"
)

func bucket_to_local_tables(
	db_creds vcap.Credentials,
	bucket_creds vcap.Credentials,
	s3path *structs.S3Path,
) {
	table_to_schema := get_table_and_schema_names(db_creds)
	//fmt.Sprintf("%s%s/%s-%s.dump", s3path.Bucket, s3path.Key, schema, table)
	for table, schema := range table_to_schema {
		dump_file_name := fmt.Sprintf("%s-%s.dump", schema, table)
		// mc_pipe := pipes.McRead(
		// 	bucket_creds,
		// 	fmt.Sprintf("%s%s/%s-%s.dump", s3path.Bucket, s3path.Key, schema, table))
		// mc_pipe := pipes.McRead(
		// 	bucket_creds,
		// 	fmt.Sprintf("%s%s/%s-%s.dump", s3path.Bucket, s3path.Key, schema, table),
		// ).FilterLine(func(s string) string {
		// 	if strings.Contains(s, "CREATE") {
		// 		fmt.Printf("REPLACING IN %s\n", s)
		// 	}
		// 	if strings.Contains(s, "CREATE TABLE") {
		// 		return strings.Replace(s, "CREATE TABLE", "CREATE TABLE IF NOT EXISTS", -1)
		// 	} else if strings.Contains(s, "CREATE INDEX") {
		// 		return strings.Replace(s, "CREATE INDEX", "CREATE INDEX IF NOT EXISTS", -1)
		// 	} else {
		// 		return s
		// 	}
		// })
		// psql_pipe := pipes.Psql(mc_pipe, db_creds)
		// pg_restore_schema_pipe := pipes.PG_Restore_Schema(mc_pipe, db_creds, schema, table)

		exit_code := 0
		mc_copy := pipes.McCopy(bucket_creds, fmt.Sprintf("%s%s/%s", s3path.Bucket, s3path.Key, dump_file_name))
		stdout, _ := mc_copy.String()
		if strings.Contains(stdout, "ERR") {
			logging.Logger.Printf("PGCOPY reported an error\n")
			logging.Logger.Println(stdout)
			exit_code = logging.PIPE_FAILURE
		}

		truncate_tables(db_creds, []string{table})

		pg_restore := pipes.PG_Restore(db_creds, schema, table)
		restoreOut, restoreError := pg_restore.String()
		util.ErrorCheck(restoreOut, restoreError)

		os.Remove(fmt.Sprintf("./pg_dump_tables/%s", dump_file_name))
		logging.Logger.Printf("REMOVING FILE: %s", dump_file_name)
		// func PG_dump_cleanup() {
		// 	script.Exec("rm -r ./pg_dump_tables")
		// }

		// if strings.Contains(stdout, "ERR") {
		// 	logging.Logger.Printf("S3TODB `mc` reported an error\n")
		// 	logging.Logger.Println(stdout)
		// 	exit_code = logging.PIPE_FAILURE
		// }
		// util.ErrorCheck(stdout, stderr)

		// if mc_pipe.Error() != nil {
		// 	logging.Logger.Println("S3TODB `dump | mc` pipe failed")
		// 	exit_code = logging.PIPE_FAILURE
		// }

		//stdout, _ = psql_pipe.String()
		// stdout, _ = pg_restore_schema_pipe.String()
		// if strings.Contains(stdout, "ERR") {
		// 	logging.Logger.Printf("PGRESTORESCHEMA reported an error\n")
		// 	logging.Logger.Println(stdout)
		// 	exit_code = logging.PIPE_FAILURE
		// }

		//pg_restore_data_pipe := pipes.PG_Restore_Data(mc_pipe, db_creds, schema, table)

		// stdout, _ = pg_restore_data_pipe.String()
		// if strings.Contains(stdout, "ERR") {
		// 	logging.Logger.Printf("PGRESTOREDATA reported an error\n")
		// 	logging.Logger.Println(stdout)
		// 	exit_code = logging.PIPE_FAILURE
		// }

		if exit_code != 0 {
			os.Exit(exit_code)
		}
	}

}

// FIXME: need s3read...
func bucket_to_cgov_tables(
	s3_creds vcap.Credentials,
	db_creds vcap.Credentials,
	s3path *structs.S3Path,
) {
	table_to_schema := get_table_and_schema_names(db_creds)
	//fmt.Sprintf("%s%s/%s-%s.dump", s3path.Bucket, s3path.Key, schema, table)
	for table, schema := range table_to_schema {
		dump_file_name := fmt.Sprintf("%s-%s.dump", schema, table)

		exit_code := 0
		s3_copy := pipes.S3Copy(s3_creds, fmt.Sprintf("%s%s/%s", s3path.Bucket, s3path.Key, dump_file_name))
		stdout, _ := s3_copy.String()
		if strings.Contains(stdout, "ERR") {
			logging.Logger.Printf("S3COPY reported an error\n")
			logging.Logger.Println(stdout)
			exit_code = logging.PIPE_FAILURE
		}

		truncate_tables(db_creds, []string{table})

		pg_restore := pipes.PG_Restore(db_creds, schema, table)
		restoreOut, restoreError := pg_restore.String()
		util.ErrorCheck(restoreOut, restoreError)

		os.Remove(fmt.Sprintf("./pg_dump_tables/%s", dump_file_name))
		logging.Logger.Printf("REMOVING FILE: %s", dump_file_name)
		if exit_code != 0 {
			os.Exit(exit_code)
		}
	}
	// s3_pipe := pipes.S3Read(
	// 	s3_creds,
	// 	fmt.Sprintf("%s%s", s3path.Bucket, s3path.Key),
	// )
	// psql_pipe := pipes.Psql(s3_pipe, dest_db_creds)

	// psql_pipe.Wait()
	// if err := psql_pipe.Error(); err != nil {
	// 	logging.Logger.Println("DUMPDBTOS3 `dump | mc` pipe failed")
	// 	os.Exit(logging.PIPE_FAILURE)
	// }
}

// S3toDBCmd represents the S3toDB command
var S3toDBCmd = &cobra.Command{
	Use:   "s3_to_db",
	Args:  cobra.ArbitraryArgs,
	Short: "Restore pg_dump file to database.",
	Long: `This command takes database and s3 path input, determining
	if this is being run locally with minio or on cloud.gov, copies the .dump
	files from the dedicated s3 storage to disk, truncates the target table,
	and then performs a pg_restore on the designated table.`,
	Run: func(cmd *cobra.Command, args []string) {
		util.UnsetProxy()
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
			os.Remove("pg_dump_tables")
			logging.Logger.Println("Finished Restore and cleaning residual files/folders.")
		case "DEV":
			fallthrough
		case "STAGING":
			fallthrough
		case "PREVIEW":
			fallthrough
		case "PRODUCTION":
			bucket_creds, err := vcap.VCS.GetCredentials("s3", path_struct.Bucket)
			if err != nil {
				logging.Logger.Printf("S3toDB could not get s3 credentials")
				os.Exit(logging.COULD_NOT_FIND_CREDENTIALS)
			}
			bucket_to_cgov_tables(bucket_creds, db_creds, path_struct)
			os.Remove("pg_dump_tables")
			logging.Logger.Println("Finished Restore and cleaning residual files/folders.")
		}
	},
}

var (
	s3_to_db_s3path string
	s3_to_db_db     string
)

func init() {
	PG_dump_prep()
	rootCmd.AddCommand(S3toDBCmd)
	parseFlags("s3_to_db", S3toDBCmd)
}

func PG_dump_prep() {
	script.Exec("mkdir -p pg_dump_tables")
}
