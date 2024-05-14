package util

import (
	"os"
	"os/exec"
	"slices"
	"strings"

	"gov.gsa.fac.cgov-util/internal/logging"
)

func Unset_Proxy() {
	if slices.Contains([]string{"LOCAL", "DEV", "PREVIEW", "STAGING", "PRODUCTION"}, os.Getenv("ENV")) {
		unset := []string{
			"unset",
			"https_proxy",
		}
		command := strings.Join(unset[:], " ")
		logging.Logger.Println("Proxy:", os.Getenv("https_proxy"))
		logging.Logger.Printf("Unsetting https_proxy variable...")
		logging.Logger.Printf("Running command: " + command)
		unset_proxy := exec.Command("bash", "-c", command)
		unsetOutput, unsetError := unset_proxy.Output()
		ErrorCheck(string(unsetOutput), unsetError)
	}
}
