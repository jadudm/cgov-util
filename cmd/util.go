package cmd

import (
	"net/url"
	"os"
	"path/filepath"

	"gov.gsa.fac.cgov-util/internal/logging"
)

func parseS3Path() {
	u, err := url.Parse(s3path)
	if err != nil {
		logging.Logger.Printf("S3toDB could not parse s3 path: %s", s3path)
		os.Exit(logging.S3_PATH_PARSE_ERROR)
	}
	if u.Scheme != "s3" {
		logging.Logger.Printf("S3toDB does not look like an S3 path (e.g. `s3://`): %s", s3path)
		os.Exit(logging.S3_PATH_PARSE_ERROR)
	}
	bucket = filepath.Clean(u.Host)
	key = filepath.Clean(u.Path)

}
