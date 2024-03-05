package pipes

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitfield/script"
	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/util"
)

// For reasons that are unclear, the access key id and secret key
// are coming through from VCAP empty. But, the endpoint is not.
// This makes no sense.
func S3(in_pipe *script.Pipe,
	up map[string]string,
	prefix string,
	source_db string,
	schema string, table string) *script.Pipe {
	os.Setenv("AWS_SECRET_ACCESS_KEY", up["secret_access_key"])
	os.Setenv("AWS_ACCESS_KEY_ID", up["access_key_id"])
	os.Setenv("AWS_DEFAULT_REGION", up["region"])
	// https://serverfault.com/questions/886562/streaming-postgresql-pg-dump-to-s3
	cmd := []string{
		"aws",
		"s3",
		"cp",
		"-",
		fmt.Sprintf("s3://%s/backups/%s-%s_%s.dump",
			up["bucket"],
			prefix,
			schema, table),
	}

	// Combine the slice for printing and execution.
	combined := strings.Join(cmd[:], " ")
	logging.Logger.Printf("BACKUPS s3 targeting %s\n", prefix)
	if util.IsDebugLevel("DEBUG") {
		fmt.Printf("command: %s\n", combined)
	}
	return in_pipe.Exec(combined)
}
