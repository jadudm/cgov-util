package cmd

import (
	"net/url"
	"os"
	"path/filepath"

	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/structs"
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
