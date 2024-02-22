package vcap

import (
	"os"

	"golang.org/x/exp/slices"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"gov.gsa.fac.backups/internal/logging"
)

func get_vcap_services() (string, error) {
	// var data map[string]interface{}
	// err := json.Unmarshal([]byte(test_json), &data)
	// if err != nil {
	// 	return nil, errors.New("Could not unmarshal vcap services")
	// }
	return "", nil
}

type RDSCreds struct {
	DBName   string
	Host     string
	Name     string
	Password string
	Port     string
	Uri      string
	Username string
}

func GetRDSCredentials(label string) (*RDSCreds, error) {
	data := os.Getenv("VCAP_SERVICES")
	instances := gjson.Get(data, "aws-rds")
	for _, instance := range instances.Array() {
		if instance.Get("name").String() == label {
			creds := instance.Get("credentials")
			return &RDSCreds{
				DBName:   creds.Get("db_name").String(),
				Host:     creds.Get("host").String(),
				Name:     creds.Get("name").String(),
				Password: creds.Get("password").String(),
				Port:     creds.Get("port").String(),
				Uri:      creds.Get("uri").String(),
				Username: creds.Get("username").String(),
			}, nil
		}
	}
	return nil, errors.Errorf("No credentials found for '%s'", label)
}

// These are hardcoded to match the FAC stack.
func GetLocalCredentials(label string) (*RDSCreds, error) {

	return &RDSCreds{
		DBName:   viper.GetString(label + ".db_name"),
		Name:     viper.GetString(label + ".db_name"),
		Host:     viper.GetString(label + ".host"),
		Port:     viper.GetString(label + ".port"),
		Username: viper.GetString(label + ".username"),
		Password: viper.GetString(label + ".password"),
		Uri:      viper.GetString(label + ".uri"),
	}, nil

}

func GetCreds() (*RDSCreds, *RDSCreds) {
	var source *RDSCreds
	var dest *RDSCreds
	var err error

	if slices.Contains([]string{"LOCAL", "TESTING"}, os.Getenv("ENV")) {
		source, err = GetLocalCredentials("fac-db")
		if err != nil {
			logging.Logger.Println("BACKUPS Cannot get local source credentials")
			os.Exit(-1)
		}
		dest, err = GetLocalCredentials("fac-snapshot-db")
		if err != nil {
			logging.Logger.Println("BACKUPS Cannot get local dest credentials")
			os.Exit(-1)
		}

	} else {
		source, _ = GetRDSCredentials("fac-db")
		if err != nil {
			logging.Logger.Println("BACKUPS Cannot get RDS source credentials")
			os.Exit(-1)
		}
		dest, _ = GetRDSCredentials("fac-snapshot-db")
		if err != nil {
			logging.Logger.Println("BACKUPS Cannot get RDS dest credentials")
			os.Exit(-1)
		}
	}

	return source, dest
}
