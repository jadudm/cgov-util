package pipes

import (
	"fmt"
	"strings"

	"github.com/bitfield/script"
	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/vcap"
)

func Psql(in_pipe *script.Pipe, creds *vcap.CredentialsRDS, debug bool) *script.Pipe {
	cmd := []string{
		"psql",
		"--no-password",
		"--dbname",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			creds.Username,
			creds.Password,
			creds.Host,
			creds.Port,
			creds.DB_Name,
		),
	}
	combined := strings.Join(cmd[:], " ")
	if debug {
		logging.Logger.Printf("command: %s\n", combined)
	}
	logging.Logger.Printf("BACKUPS psql targeting %s\n", creds.DB_Name)
	return in_pipe.Exec(combined)
}
