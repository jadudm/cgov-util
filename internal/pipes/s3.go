package pipes

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bitfield/script"
	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/util"
	"gov.gsa.fac.cgov-util/internal/vcap"
)

// For reasons that are unclear, the access key id and secret key
// are coming through from VCAP empty. But, the endpoint is not.
// This makes no sense.
func S3Write(in_pipe *script.Pipe,
	s3_creds vcap.Credentials,
	path string) *script.Pipe {

	os.Setenv("AWS_SECRET_ACCESS_KEY", s3_creds.Get("secret_access_key").String())
	os.Setenv("AWS_ACCESS_KEY_ID", s3_creds.Get("access_key_id").String())
	os.Setenv("AWS_DEFAULT_REGION", s3_creds.Get("region").String())
	// https://serverfault.com/questions/886562/streaming-postgresql-pg-dump-to-s3
	cmd := []string{
		util.AWS_path,
		"s3",
		"cp",
		"-",
		fmt.Sprintf("s3://%s/%s",
			s3_creds.Get("bucket").String(),
			path),
	}

	// Combine the slice for printing and execution.
	combined := strings.Join(cmd[:], " ")
	logging.Logger.Printf("S3 "+util.AWS_path+" targeting %s\n", path)
	if util.IsDebugLevel("DEBUG") {
		fmt.Printf("command: %s\n", combined)
	}
	return in_pipe.Exec(combined)
}

func S3Read(s3_creds vcap.Credentials,
	path string) *script.Pipe {

	os.Setenv("AWS_SECRET_ACCESS_KEY", s3_creds.Get("secret_access_key").String())
	os.Setenv("AWS_ACCESS_KEY_ID", s3_creds.Get("access_key_id").String())
	os.Setenv("AWS_DEFAULT_REGION", s3_creds.Get("region").String())
	// NOTE: The entire change from "write" is that we pass "-" as the
	// path in a different location. Instead of reading from STDIN, we write
	// to STDOUT. This lets us pipe the read into another command.
	cmd := []string{
		util.AWS_path,
		"s3",
		"cp",
		fmt.Sprintf("s3://%s/%s",
			s3_creds.Get("bucket").String(),
			path),
		"-",
	}

	// Combine the slice for printing and execution.
	combined := strings.Join(cmd[:], " ")
	logging.Logger.Printf("S3 "+util.AWS_path+" targeting %s\n", path)
	if util.IsDebugLevel("DEBUG") {
		fmt.Printf("command: %s\n", combined)
	}
	return script.Exec(combined)
}

func S3Sync(source_creds vcap.Credentials,
	dest_creds vcap.Credentials) {

	os.Setenv("AWS_SECRET_ACCESS_KEY", source_creds.Get("secret_access_key").String())
	os.Setenv("AWS_ACCESS_KEY_ID", source_creds.Get("access_key_id").String())
	os.Setenv("AWS_DEFAULT_REGION", source_creds.Get("region").String())
	cmd := []string{
		util.AWS_path,
		"s3",
		"sync",
		fmt.Sprintf("s3://%s/",
			source_creds.Get("bucket").String(),
		),
		fmt.Sprintf("s3://%s/",
			dest_creds.Get("bucket").String(),
		),
	}
	combined := strings.Join(cmd[:], " ")
	logging.Logger.Printf("S3 Syncing " + source_creds.Get("bucket").String() + " to " + dest_creds.Get("bucket").String())
	logging.Logger.Printf("Running command: " + combined)
	sync := exec.Command("bash", "-c", combined)
	syncOutput, syncError := sync.Output()
	util.ErrorCheck(string(syncOutput), syncError)

}
