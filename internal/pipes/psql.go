package pipes

import (
	"fmt"
	"strings"

	"github.com/bitfield/script"
	"gov.gsa.fac.backups/internal/logging"
	"gov.gsa.fac.backups/internal/vcap"
)

func Psql(in_pipe *script.Pipe, creds *vcap.CredentialsRDS) *script.Pipe {
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
	logging.Logger.Printf("BACKUPS psql targeting %s", creds.DB_Name)
	return in_pipe.Exec(combined)
}
