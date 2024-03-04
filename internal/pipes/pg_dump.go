package pipes

import (
	"fmt"
	"strings"

	"github.com/bitfield/script"
	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/vcap"
)

// https://bitfieldconsulting.com/golang/scripting
func PG_Dump(creds *vcap.CredentialsRDS) *script.Pipe {
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
	fmt.Printf("command: %s\n", combined)
	logging.Logger.Printf("BACKUPS pg_dump targeting %s", creds.DB_Name)
	return script.Exec(combined)
}
