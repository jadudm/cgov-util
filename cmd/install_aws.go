/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/util"
)

func ErrorCheck(output []byte, err error) {
	// https://stackoverflow.com/a/7786922
	if err != nil {
		logging.Logger.Println(err.Error())
		return
	}
	logging.Logger.Println(output)
}

// func CheckInstallation() {
// 	checkinstallation := []string{
// 		util.AWS_path,
// 		"--version",
// 	}
// 	command := strings.Join(checkinstallation[:], " ")
// 	check := exec.Command("bash", "-c", command)
// 	checkOutput, checkError := check.Output()
// 	ErrorCheck(checkOutput, checkError)
// }

// installAwsCmd represents the installAws command
var installAwsCmd = &cobra.Command{
	Use:   "install_aws",
	Short: "Install AWS-CLI",
	Long:  `This command will curl the necessary aws-cli package and install it`,
	Run: func(cmd *cobra.Command, args []string) {
		if slices.Contains([]string{"DEV", "PREVIEW", "STAGING", "PRODUCTION"}, os.Getenv("ENV")) {
			logging.Logger.Printf("ENV detected to be a cloud.gov environment. Installing AWS CLI.")
			// curl -x $https_proxy -L "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
			getaws := []string{
				"curl",
				"-x",
				os.Getenv("https_proxy"),
				"-L",
				"https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip",
				"-o",
				"awscliv2.zip",
			}
			curlCommand := strings.Join(getaws[:], " ")
			logging.Logger.Printf("Fetching aws cli via curl...")
			logging.Logger.Printf("Running command: " + curlCommand)
			curl := exec.Command("bash", "-c", curlCommand)
			curlOutput, curlError := curl.Output()
			ErrorCheck(curlOutput, curlError)

			// unzip awscliv2.zip && rm awscliv2.zip
			unzipaws := []string{
				"unzip",
				"awscliv2.zip",
				"&&",
				"rm",
				"awscliv2.zip",
			}
			unzipCommand := strings.Join(unzipaws[:], " ")
			logging.Logger.Printf("Unzipping aws cli...")
			logging.Logger.Printf("Running command: " + unzipCommand)
			extract := exec.Command("bash", "-c", unzipCommand)
			unzipOutput, unzipError := extract.Output()
			ErrorCheck(unzipOutput, unzipError)

			// ./aws/install -i ~/usr -b ~/bin
			installaws := []string{
				"./aws/install",
				"-i",
				"~/usr",
				"-b",
				"~/bin",
			}
			installCommand := strings.Join(installaws[:], " ")
			logging.Logger.Printf("Installing aws to bin...")
			logging.Logger.Printf("Running command: " + installCommand)
			install := exec.Command("bash", "-c", installCommand)
			installOutput, installError := install.Output()
			ErrorCheck(installOutput, installError)
			// Regardless of the case, check to see if AWS-CLI is installed or not.
			CheckInstallation()
		} else {
			logging.Logger.Printf("ENV set to LOCAL or TESTING, aws-cli is not necessary to install.")
			// Regardless of the case, check to see if AWS-CLI is installed or not.
			CheckInstallation()
		}
	},
}

func init() {
	rootCmd.AddCommand(installAwsCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installAwsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installAwsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
