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
	"gov.gsa.fac.backups/internal/logging"
	vcap "gov.gsa.fac.backups/internal/vcap"
)

// https://bitfieldconsulting.com/golang/scripting

func pg_dump(creds *vcap.RDSCreds) *script.Pipe {
	cmd := []string{
		"pg_dump",
		"--clean",
		"--no-password",
		"--if-exists",
		"--no-privileges",
		"--no-owner",
		"--format plain",
		fmt.Sprintf("--host %s", creds.Host),
		fmt.Sprintf("--port %s", creds.Port),
		fmt.Sprintf("--username %s", creds.Username),
		creds.DBName,
	}
	combined := strings.Join(cmd[:], " ")
	logging.Logger.Printf("BACKUPS Running `%s`\n", combined)
	return script.Exec(combined)
}

func psql(in_pipe *script.Pipe, creds *vcap.RDSCreds) *script.Pipe {
	cmd := []string{
		"psql",
		"--no-password",
		fmt.Sprintf("--host %s", creds.Host),
		fmt.Sprintf("--port %s", creds.Port),
		fmt.Sprintf("--username %s", creds.Username),
		creds.DBName,
	}
	combined := strings.Join(cmd[:], " ")
	logging.Logger.Printf("BACKUPS Running `%s`\n", combined)
	return in_pipe.Exec(combined)
}

func clone() {
	source, dest := vcap.GetCreds()
	psql_pipe := psql(pg_dump(source), dest)
	psql_pipe.Wait()
	if err := psql_pipe.Error(); err != nil {
		logging.Logger.Printf("BACKUPS Pipe failed: %w", err)
		os.Exit(-1)
	}

}

// snapshotDbToDbCmd represents the snapshotDbToDb command
var cloneDBToDB = &cobra.Command{
	Use:   "clone-db-to-db",
	Short: "Clones one database to another, destructively",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("clone-db-to-db called")
		clone()
	},
}

func init() {
	rootCmd.AddCommand(cloneDBToDB)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// snapshotDbToDbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// snapshotDbToDbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
