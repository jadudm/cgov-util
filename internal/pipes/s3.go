package pipes

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitfield/script"
	"gov.gsa.fac.backups/internal/logging"
	"gov.gsa.fac.backups/internal/vcap"
)

// export AWS_PRIVATE_ACCESS_KEY_ID=longtest
// export AWS_PRIVATE_SECRET_ACCESS_KEY=longtest
// export AWS_S3_PRIVATE_ENDPOINT="http://minio:9000"
// mc alias set myminio "${AWS_S3_PRIVATE_ENDPOINT}" minioadmin minioadmin
// # Do nothing if the bucket already exists.
// # https: //min.io/docs/minio/linux/reference/minio-mc/mc-mb.html
// mc mb --ignore-existing myminio/gsa-fac-private-s3

// https://bitfieldconsulting.com/golang/scripting
func S3(in_pipe *script.Pipe, up *vcap.CredentialsS3, prefix string, source_db string) *script.Pipe {
	os.Setenv("ACCESS_KEY_ID", up.AccessKeyId)
	os.Setenv("SECRET_ACCESS_KEY", up.SecretAccessKey)

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
