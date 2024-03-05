package pipes

import (
	"fmt"
	"strings"

	"github.com/bitfield/script"
	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/structs"
	"gov.gsa.fac.cgov-util/internal/util"
)

func Psql(in_pipe *script.Pipe, creds *structs.CredentialsRDS) *script.Pipe {
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
	if util.IsDebugLevel("DEBUG") {
		logging.Logger.Printf("command: %s\n", combined)
	}
	logging.Logger.Printf("BACKUPS psql targeting %s\n", creds.DB_Name)
	return in_pipe.Exec(combined)
}
