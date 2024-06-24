package pipes

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/bitfield/script"
	"github.com/google/uuid"
	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/util"
	"gov.gsa.fac.cgov-util/internal/vcap"
)

func getExecutable() string {
	if runtime.GOOS == "windows" {
		return "./mc.exe"
	} else {
		return "mc"
	}
}

// https://bitfieldconsulting.com/golang/scripting
func McWrite(in_pipe *script.Pipe,
	creds vcap.Credentials,
	path string) *script.Pipe {
	// // mc pipe myminio/gsa-fac-private-s3/backups/${PREFIX}-${FROM_DATABASE}.dump
	// Always set the alias first.
	os.Setenv("AWS_PRIVATE_ACCESS_KEY_ID", creds.Get("access_key_id").String())
	os.Setenv("AWS_PRIVATE_SECRET_ACCESS_KEY", creds.Get("secret_access_key").String())

	minio_alias := fmt.Sprintf("minio_alias_%s", uuid.New())

	set_alias := []string{
		getExecutable(), "alias", "set", minio_alias,
		creds.Get("endpoint").String(),
		creds.Get("admin_username").String(),
		creds.Get("admin_password").String(),
	}
	sa_combined := strings.Join(set_alias[:], " ")
	logging.Logger.Printf("MC Running `%s`\n", sa_combined)
	script.Exec(sa_combined).Stdout()

	cmd := []string{
		getExecutable(),
		"pipe",
		fmt.Sprintf("%s/%s/%s",
			minio_alias,
			creds.Get("bucket").String(),
			path),
	}
	// Combine the slice for printing and execution.
	combined := strings.Join(cmd[:], " ")
	if util.IsDebugLevel("DEBUG") {
		fmt.Printf("command: %s\n", combined)
	}
	logging.Logger.Printf("MC "+getExecutable()+" targeting %s", path)
	return in_pipe.Exec(combined)
}

func McRead(
	creds vcap.Credentials,
	path string) *script.Pipe {
	// // mc pipe myminio/gsa-fac-private-s3/backups/${PREFIX}-${FROM_DATABASE}.dump
	// Always set the alias first.
	os.Setenv("AWS_PRIVATE_ACCESS_KEY_ID", creds.Get("access_key_id").String())
	os.Setenv("AWS_PRIVATE_SECRET_ACCESS_KEY", creds.Get("secret_access_key").String())

	minio_alias := fmt.Sprintf("minio_alias_%s", uuid.New())

	set_alias := []string{
		getExecutable(), "alias", "set", minio_alias,
		creds.Get("endpoint").String(),
		creds.Get("admin_username").String(),
		creds.Get("admin_password").String(),
	}
	sa_combined := strings.Join(set_alias[:], " ")
	logging.Logger.Printf("MC Running `%s`\n", sa_combined)
	script.Exec(sa_combined).Stdout()

	cmd := []string{
		getExecutable(),
		"cat",
		fmt.Sprintf("%s/%s/%s",
			minio_alias,
			creds.Get("bucket").String(),
			path),
	}
	// Combine the slice for printing and execution.
	combined := strings.Join(cmd[:], " ")
	if util.IsDebugLevel("DEBUG") {
		fmt.Printf("command: %s\n", combined)
	}
	logging.Logger.Printf("MC mc cat'ing %s", path)
	return script.Exec(combined)
}

func McCopy(
	creds vcap.Credentials,
	path string) *script.Pipe {
	// // mc pipe myminio/gsa-fac-private-s3/backups/${PREFIX}-${FROM_DATABASE}.dump
	// Always set the alias first.
	os.Setenv("AWS_PRIVATE_ACCESS_KEY_ID", creds.Get("access_key_id").String())
	os.Setenv("AWS_PRIVATE_SECRET_ACCESS_KEY", creds.Get("secret_access_key").String())

	minio_alias := fmt.Sprintf("minio_alias_%s", uuid.New())

	set_alias := []string{
		getExecutable(), "alias", "set", minio_alias,
		creds.Get("endpoint").String(),
		creds.Get("admin_username").String(),
		creds.Get("admin_password").String(),
	}
	sa_combined := strings.Join(set_alias[:], " ")
	logging.Logger.Printf("MC Running `%s`\n", sa_combined)
	script.Exec(sa_combined).Stdout()

	cmd := []string{
		getExecutable(),
		"cp",
		fmt.Sprintf("%s/%s/%s",
			minio_alias,
			creds.Get("bucket").String(),
			path),
		"./pg_dump_tables/",
	}
	// Combine the slice for printing and execution.
	combined := strings.Join(cmd[:], " ")
	if util.IsDebugLevel("DEBUG") {
		fmt.Printf("command: %s\n", combined)
	}
	logging.Logger.Printf("MC copying %s", path)
	return script.Exec(combined)
}
