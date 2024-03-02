package pipes

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitfield/script"
	"gov.gsa.fac.backups/internal/logging"
	"gov.gsa.fac.backups/internal/vcap"
)

// https://bitfieldconsulting.com/golang/scripting
func S3(in_pipe *script.Pipe, up *vcap.CredentialsS3, prefix string, source_db string) *script.Pipe {
	os.Setenv("ACCESS_KEY_ID", up.AccessKeyId)
	os.Setenv("SECRET_ACCESS_KEY", up.SecretAccessKey)
	// https://serverfault.com/questions/886562/streaming-postgresql-pg-dump-to-s3
	cmd := []string{
		"s3",
		"cp",
		fmt.Sprintf("s3://%s/backups/%s-%s.dump",
			up.Bucket,
			prefix,
			source_db),
	}

	// Combine the slice for printing and execution.
	combined := strings.Join(cmd[:], " ")
	// This will log the password...
	logging.Logger.Printf("BACKUPS Running `%s`\n", combined)
	logging.Logger.Printf("BACKUPS s3 targeting %s", prefix)
	return in_pipe.Exec(combined)
}
