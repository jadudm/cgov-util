package util

import "gov.gsa.fac.cgov-util/internal/logging"

func ErrorCheck(output string, err error) {
	// https://stackoverflow.com/a/7786922
	if err != nil {
		logging.Logger.Println(err.Error())
		return
	} else {
		logging.Logger.Println(output)
	}
}
