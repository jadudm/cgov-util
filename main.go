/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"

	"golang.org/x/exp/slices"

	"gov.gsa.fac.cgov-util/cmd"
	"gov.gsa.fac.cgov-util/internal/vcap"
)

// Useful documentation for people new to Go, and
// related to modules in this command.
// https://www.digitalocean.com/community/tutorials/how-to-use-the-cobra-package-in-go
// https://github.com/tidwall/gjson
// https://martengartner.medium.com/my-favourite-go-project-setup-479563662834
// https://github.com/spf13/cobra
// https://github.com/spf13/viper
// https://go.dev/doc/tutorial/handle-errors

// Looks for config.yaml in the same directory as the app.
// Optionally, the config can be in `$HOME/.fac/config.yaml`
func readConfig() {
	// Do the right thing in the right env.
	if slices.Contains([]string{"LOCAL", "TESTING"}, os.Getenv("ENV")) {
		// Locally, load the file from one of two places.
		vcap.ReadVCAPConfigFile("config.json")
	} else {
		vcap.ReadVCAPConfig()
	}
}

func main() {

	readConfig()
	cmd.Execute()
}
