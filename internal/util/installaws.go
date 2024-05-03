package util

import (
	"os"
	"strings"

	"github.com/bitfield/script"
	"gov.gsa.fac.cgov-util/internal/logging"
)

func InstallAWS() {
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
	curlaws := strings.Join(getaws[:], " ")
	logging.Logger.Printf("Fetching aws cli via curl...")
	script.Exec(curlaws)

	// unzip awscliv2.zip && rm awscliv2.zip
	unzipaws := []string{
		"unzip",
		"awscliv2.zip",
		"&&",
		"rm",
		"awscliv2.zip",
	}
	logging.Logger.Printf("Unzipping aws cli...")
	unzip := strings.Join(unzipaws[:], " ")
	script.Exec(unzip)

	// ./aws/install -i ~/usr -b ~/bin
	installaws := []string{
		"./aws/install",
		"-i",
		"~/usr",
		"-b",
		"~/bin",
	}
	logging.Logger.Printf("Installing aws to bin...")
	install := strings.Join(installaws[:], " ")
	script.Exec(install)
}
