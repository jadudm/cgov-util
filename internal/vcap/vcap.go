package vcap

import (
	"bytes"
	"os"

	"golang.org/x/exp/slices"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gov.gsa.fac.backups/internal/logging"
)

type RDSCreds struct {
	DB_Name  string `json:"db_name"`
	Host     string `json:"host"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Uri      string `json:"uri"`
}

type RDSInstance struct {
	Plan        string   `json:"plan`
	Name        string   `json:"name"`
	Credentials RDSCreds `json:"credentials"`
}

func GetRDSCredentials(label string) (*RDSCreds, error) {
	var instanceSlice []RDSInstance
	err := viper.UnmarshalKey("aws-rds", &instanceSlice)
	if err != nil {
		logging.Logger.Println("Could not unmarshal aws-rds from VCAP_SERVICES")
	}
	for _, instance := range instanceSlice {
		if instance.Name == label {
			return &instance.Credentials, nil
		}
	}
	return nil, errors.Errorf("No credentials found for '%s'", label)
}

// These are hardcoded to match the FAC stack.
func GetLocalCredentials(label string) (*RDSCreds, error) {
	return &RDSCreds{
		DB_Name:  viper.GetString(label + ".db_name"),
		Name:     viper.GetString(label + ".name"),
		Host:     viper.GetString(label + ".host"),
		Port:     viper.GetString(label + ".port"),
		Username: viper.GetString(label + ".username"),
		Password: viper.GetString(label + ".password"),
		Uri:      viper.GetString(label + ".uri"),
	}, nil

}

func GetCreds(source_db string, dest_db string) (*RDSCreds, *RDSCreds) {
	var source *RDSCreds
	var dest *RDSCreds
	var err error

	if slices.Contains([]string{"LOCAL", "TESTING"}, os.Getenv("ENV")) {
		source, err = GetLocalCredentials(source_db)
		if err != nil {
			logging.Logger.Println("BACKUPS Cannot get local source credentials")
			os.Exit(-1)
		}
		dest, err = GetLocalCredentials(dest_db)
		if err != nil {
			logging.Logger.Println("BACKUPS Cannot get local dest credentials")
			os.Exit(-1)
		}

	} else {
		source, err = GetRDSCredentials(source_db)
		if err != nil {
			logging.Logger.Println("BACKUPS Cannot get RDS source credentials")
			os.Exit(-1)
		}
		dest, err = GetRDSCredentials(dest_db)
		if err != nil {
			logging.Logger.Println("BACKUPS Cannot get RDS dest credentials")
			os.Exit(-1)
		}
	}

	return source, dest
}

func ReadVCAPConfig() {
	// Remotely, read it in from the VCAP_SERVICES env var, which will
	// provide a large JSON structure.
	viper.SetConfigType("json")
	vcap := os.Getenv("VCAP_SERVICES")
	viper.ReadConfig(bytes.NewBufferString(vcap))
}
