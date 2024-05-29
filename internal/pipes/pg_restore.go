package pipes

import (
	"fmt"
	"strings"

	"github.com/bitfield/script"
	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/util"
	"gov.gsa.fac.cgov-util/internal/vcap"
)

func PG_Restore(creds vcap.Credentials, schema string, table string) *script.Pipe {
	// Compose the command as a slice
	// pg_restore -U postgres --schema-only -d new_db /directory/path/db-dump-name.dump
	cmd := []string{
		util.PGRESTORE_path,
		"-v",
		"--no-password",
		"--no-privileges",
		"--no-owner",
		"--data-only",
		"--dbname",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			creds.Get("username").String(),
			creds.Get("password").String(),
			creds.Get("host").String(),
			creds.Get("port").String(),
			creds.Get("db_name").String(),
		),
		fmt.Sprintf("./pg_dump_tables/%s-%s.dump", schema, table),
	}
	// Combine the slice for printing and execution.
	combined := strings.Join(cmd[:], " ")
	logging.Logger.Printf("RESTORE "+util.PGRESTORE_path+" targeting %s.%s\n", schema, table)
	//logging.Logger.Printf("CALLING COMMAND: " + combined)
	if util.IsDebugLevel("DEBUG") {
		fmt.Printf("command: %s\n", combined)
	}
	return script.Exec(combined)
}
