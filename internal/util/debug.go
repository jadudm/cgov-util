package util

import "os"

func getDebugLevel() string {
	return os.Getenv("CGOV_UTIL_DEBUG_LEVEL")
}

func IsDebugLevel(lvl string) bool {
	return getDebugLevel() == lvl
}
