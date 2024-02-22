/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"bytes"
	"fmt"
	"os"

	"golang.org/x/exp/slices"

	"github.com/spf13/viper"
	"gov.gsa.fac.backups/cmd"
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

	if slices.Contains([]string{"LOCAL", "TESTING"}, os.Getenv("ENV")) {
		// Locally, load the file from one of two places.
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("$HOME/.fac")
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	} else {
		// Remotely, read it in from the VCAP_SERVICES env var, which will
		// provide a large JSON structure.
		viper.SetConfigType("json")
		vcap_services := []byte(os.Getenv("VCAP_SERVICES"))
		viper.ReadConfig(bytes.NewBuffer(vcap_services))
	}
}

func main() {
	readConfig()
	cmd.Execute()
}
