package util

import "github.com/bitfield/script"

func PG_dump_prep() {
	script.Exec("mkdir pg_dump_tables")
}

func PG_dump_cleanup() {
	script.Exec("rm -r ./pg_dump_tables")
}
