package util

import (
	"os"

	"gov.gsa.fac.cgov-util/internal/logging"
)

var PSQL_path = "NO PATH SET"
var PGDUMP_path = "NO PATH SET"
var AWS_path = "NO PATH SET"

func SetPaths(env string) {
	switch env {
	case "LOCAL":
		fallthrough
	case "TESTING":
		PSQL_path = "psql"
		PGDUMP_path = "pg_dump"
		AWS_path = "aws"
	case "DEV":
		fallthrough
	case "STAGING":
		fallthrough
	case "PREVIEW":
		fallthrough
	case "PRODUCTION":
		PSQL_path = "/home/vcap/deps/0/apt/usr/lib/postgresql/15/bin/psql"
		PGDUMP_path = "/home/vcap/deps/0/apt/usr/lib/postgresql/15/bin/pg_dump"
		AWS_path = "/home/vcap/app/bin/aws"
	default:
		logging.Logger.Println("No environment set, paths for executables were not set.")
		os.Exit(-1)
	}
}
