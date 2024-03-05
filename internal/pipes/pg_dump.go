package pipes

import (
	"fmt"
	"strings"

	"github.com/bitfield/script"
	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/vcap"
)

func PG_Dump_Table(creds *vcap.CredentialsRDS, schema string, table string, debug bool) *script.Pipe {
	// Compose the command as a slice
	cmd := []string{
		"pg_dump",
		"--clean",
		"--no-password",
		"--if-exists",
		"--no-privileges",
		"--no-owner",
		"--format plain",
		"--table",
		fmt.Sprintf("%s.%s", schema, table),
		"--dbname",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			creds.Username,
			creds.Password,
			creds.Host,
			creds.Port,
			creds.DB_Name,
		),
	}
	// Combine the slice for printing and execution.
	combined := strings.Join(cmd[:], " ")
	logging.Logger.Printf("BACKUPS pg_dump targeting %s\n", creds.DB_Name)
	if debug {
		fmt.Printf("command: %s\n", combined)
	}
	return script.Exec(combined)
}

// https://bitfieldconsulting.com/golang/scripting
func PG_Dump(creds *vcap.CredentialsRDS, debug bool) *script.Pipe {
	// Compose the command as a slice
	cmd := []string{
		"pg_dump",
		"--clean",
		"--no-password",
		"--if-exists",
		"--no-privileges",
		"--no-owner",
		"--format plain",
		"--dbname",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			creds.Username,
			creds.Password,
			creds.Host,
			creds.Port,
			creds.DB_Name,
		),
	}
	// Combine the slice for printing and execution.
	combined := strings.Join(cmd[:], " ")
	logging.Logger.Printf("BACKUPS pg_dump targeting %s\n", creds.DB_Name)
	if debug {
		fmt.Printf("command: %s\n", combined)
	}
	return script.Exec(combined)
}
