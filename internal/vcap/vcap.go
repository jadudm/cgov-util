package vcap

import (
	"fmt"
	"os"

	"golang.org/x/exp/slices"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"gov.gsa.fac.backups/internal/logging"
)

var test_json = `{
    "s3": [],
    "user-provided": [],
    "aws-rds": [
        {
            "label": "aws-rds",
            "provider": null,
            "plan": "medium-gp-psql",
            "name": "fac-db",
            "tags": [
                "database",
                "RDS"
            ],
            "instance_guid": "source-guid",
            "instance_name": "fac-db",
            "binding_guid": "source-binding-guid",
            "binding_name": null,
            "credentials": {
                "db_name": "the-source-db-name",
                "host": "the-source-db.us-gov-west-1.rds.amazonaws.com",
                "name": "the-source-name",
                "password": "the-source-password",
                "port": "54321",
                "uri": "the-source-uri",
                "username": "source-username"
            },
            "syslog_drain_url": null,
            "volume_mounts": []
        },
        {
            "label": "aws-rds",
            "provider": null,
            "plan": "medium-gp-psql",
            "name": "fac-snapshot-db",
            "tags": [
                "database",
                "RDS"
            ],
            "instance_guid": "dest-instance-guid",
            "instance_name": "fac-db",
            "binding_guid": "dest-binding-guid",
            "binding_name": null,
            "credentials": {
                "db_name": "the-dest-db-name",
                "host": "the-dest-db.us-gov-west-1.rds.amazonaws.com",
                "name": "the-dest-name",
                "password": "the-dest-password",
                "port": "65432",
                "uri": "the-dest-uri",
                "username": "dest-username"
            },
            "syslog_drain_url": null,
            "volume_mounts": []
        }
    ]
}`

func get_vcap_services() (string, error) {
	// var data map[string]interface{}
	// err := json.Unmarshal([]byte(test_json), &data)
	// if err != nil {
	// 	return nil, errors.New("Could not unmarshal vcap services")
	// }
	return test_json, nil
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
	data, _ := get_vcap_services()
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
	var port string
	if label == "fac-db" {
		port = "5432"
	} else if label == "fac-snapshot-db" {
		port = "6543"
	}

	return &RDSCreds{
		DBName:   "postgres",
		Host:     "127.0.0.1",
		Name:     "postgres",
		Password: "",
		Port:     port,
		Uri:      "",
		Username: "postgres",
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
	fmt.Printf("%s\n%s\n", source, dest)
	return source, dest
}
