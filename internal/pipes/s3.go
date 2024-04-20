package pipes

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitfield/script"
	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/util"
	"gov.gsa.fac.cgov-util/internal/vcap"
)

// For reasons that are unclear, the access key id and secret key
// are coming through from VCAP empty. But, the endpoint is not.
// This makes no sense.
func S3(in_pipe *script.Pipe,
	up vcap.Credentials,
	path string) *script.Pipe {

	os.Setenv("AWS_SECRET_ACCESS_KEY", up.Get("secret_access_key").String())
	os.Setenv("AWS_ACCESS_KEY_ID", up.Get("access_key_id").String())
	os.Setenv("AWS_DEFAULT_REGION", up.Get("region").String())
	// https://serverfault.com/questions/886562/streaming-postgresql-pg-dump-to-s3
	cmd := []string{
		"aws",
		"s3",
		"cp",
		"-",
		fmt.Sprintf("s3://%s/%s",
			up.Get("bucket").String(),
			path),
	}

	// Combine the slice for printing and execution.
	combined := strings.Join(cmd[:], " ")
	logging.Logger.Printf("S3 s3 targeting %s\n", path)
	if util.IsDebugLevel("DEBUG") {
		fmt.Printf("command: %s\n", combined)
	}
	return in_pipe.Exec(combined)
}
