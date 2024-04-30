package cmd

import (
	"net/url"
	"os"
	"path/filepath"

	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/structs"
	"gov.gsa.fac.cgov-util/internal/vcap"
)

func parseS3Path(path string) *structs.S3Path {
	u, err := url.Parse(path)
	if err != nil {
		logging.Logger.Printf("S3toDB could not parse s3 path: %s", path)
		os.Exit(logging.S3_PATH_PARSE_ERROR)
	}
	if u.Scheme != "s3" {
		logging.Logger.Printf("S3toDB does not look like an S3 path (e.g. `s3://`): %s", path)
		os.Exit(logging.S3_PATH_PARSE_ERROR)
	}

	return &structs.S3Path{
		Bucket: filepath.Clean(u.Host),
		Key:    filepath.Clean(u.Path),
	}

}

func runLocalOrRemote(funs structs.Choice) {
	switch os.Getenv("ENV") {
	case "LOCAL":
		fallthrough
	case "TESTING":
		funs.Local()
	case "DEV":
		fallthrough
	case "STAGING":
		fallthrough
	case "PRODUCTION":
		funs.Remote()
	default:
		logging.Logger.Printf("LOCALORREMOTE impossible condition")
		os.Exit(-1)
	}
}

func getDBCredentials(db_name string) vcap.Credentials {
	// Check that we can get credentials.
	db_creds, err := vcap.VCS.GetCredentials("aws-rds", db_name)
	if err != nil {
		logging.Logger.Printf("S3toDB could not get DB credentials for %s", db_name)
		os.Exit(logging.COULD_NOT_FIND_CREDENTIALS)
	}
	return db_creds
}

func getBucketCredentials(s3path *structs.S3Path) vcap.Credentials {
	switch os.Getenv("ENV") {
	case "LOCAL":
		fallthrough
	case "TESTING":
		bucket_creds, err := vcap.VCS.GetCredentials("user-provided", s3path.Bucket)
		if err != nil {
			logging.Logger.Printf("DBTOS3 could not get minio credentials")
			os.Exit(logging.COULD_NOT_FIND_CREDENTIALS)
		}
		return bucket_creds
	case "DEV":
		fallthrough
	case "STAGING":
		fallthrough
	case "PREVIEW":
		fallthrough
	case "PRODUCTION":
		bucket_creds, err := vcap.VCS.GetCredentials("s3", s3path.Bucket)
		if err != nil {
			logging.Logger.Printf("DBTOS3 could not get s3 credentials")
			os.Exit(logging.COULD_NOT_FIND_CREDENTIALS)
		}
		return bucket_creds
	default:
		logging.Logger.Printf("DBTOS3 could not get env for bucket credentials")
		os.Exit(logging.COULD_NOT_FIND_CREDENTIALS)
	}

	// We should never get here.
	logging.Logger.Printf("DBTOS3 impossible condition")
	os.Exit(-1)
	return vcap.Credentials{}
}
