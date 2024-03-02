package pipes

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitfield/script"
	"github.com/google/uuid"
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
func Mc(in_pipe *script.Pipe, upc vcap.UserProvidedCredentials, prefix string, source_db string) *script.Pipe {
	// // mc pipe myminio/gsa-fac-private-s3/backups/${PREFIX}-${FROM_DATABASE}.dump
	// Always set the alias first.
	os.Setenv("AWS_PRIVATE_ACCESS_KEY_ID", upc["access_key_id"])
	os.Setenv("AWS_PRIVATE_SECRET_ACCESS_KEY", upc["secret_access_key"])

	minio_alias := fmt.Sprintf("minio_alias_%s", uuid.New())

	set_alias := []string{
		"mc", "alias", "set", minio_alias,
		upc["endpoint"],
		upc["admin_username"],
		upc["admin_password"],
	}
	sa_combined := strings.Join(set_alias[:], " ")
	logging.Logger.Printf("BACKUPS Running `%s`\n", sa_combined)
	script.Exec(sa_combined).Stdout()

	cmd := []string{
		"mc",
		"pipe",
		fmt.Sprintf("%s/%s/backups/%s-%s.dump",
			minio_alias,
			upc["bucket"],
			prefix,
			source_db),
	}
	// Combine the slice for printing and execution.
	combined := strings.Join(cmd[:], " ")
	// This will log the password...
	logging.Logger.Printf("BACKUPS Running `%s`\n", combined)
	logging.Logger.Printf("BACKUPS mc targeting %s", prefix)
	return in_pipe.Exec(combined)
}
