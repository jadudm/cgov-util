package pipes

import (
	"fmt"
	"strings"

	"github.com/bitfield/script"
	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/util"
	"gov.gsa.fac.cgov-util/internal/vcap"
)

func Psql(in_pipe *script.Pipe, creds vcap.Credentials) *script.Pipe {
	cmd := []string{
		util.PSQL_path,
		"--no-password",
		"--dbname",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			creds.Get("username").String(),
			creds.Get("password").String(),
			creds.Get("host").String(),
			creds.Get("port").String(),
			creds.Get("db_name"),
		),
	}
	combined := strings.Join(cmd[:], " ")
	if util.IsDebugLevel("DEBUG") {
		logging.Logger.Printf("command: %s\n", combined)
	}
	logging.Logger.Printf(util.PSQL_path + " running\n")
	return in_pipe.Exec(combined)
}
