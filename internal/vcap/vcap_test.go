package vcap

import (
	"io/ioutil"
	"os"
	"testing"
)

var db_vcap = `{
    "s3": [],
    "user-provided": [
        {
            "label": "mc",
            "name": "backups",
            "tags": [],
            "instance_guid": "UUIDALPHA1",
            "instance_name": "backups",
            "binding_guid": "UUIDALPHA2",
            "binding_name": null,
            "credentials": {
                "access_key_id": "longtest",
                "secret_access_key": "longtest",
                "bucket": "gsa-fac-private-s3",
                "endpoint": "localhost",
                "admin_username": "minioadmin",
                "admin_password": "minioadmin"
            }
        }
    ],
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
            "instance_name": "fac-snapshot-db",
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

func TestReadVCAP(t *testing.T) {
	// Load a test string into the env.
	os.Setenv("VCAP_SERVICES", db_vcap)
	// Read the VCAP config.
	ReadVCAPConfig()
	// Check to see if we can find the source DB
	creds, err := GetRDSCredentials("fac-db")
	if err != nil {
		t.Error("Could not read fac-db credentials from env.")
	}
	if creds.DB_Name != "the-source-db-name" {
		t.Error("Did not get fac-db db_name")
	}
	// How about the dest DB?
	creds, err = GetRDSCredentials("fac-snapshot-db")
	if err != nil {
		t.Error("Could not read fac-db credentials from env.")
	}
	if creds.DB_Name != "the-dest-db-name" {
		t.Error("Did not get fac-db db_name")
	}
}

func TestReadUserProvided(t *testing.T) {
	// Load a test string into the env.
	os.Setenv("VCAP_SERVICES", db_vcap)
	// Read the VCAP config.
	ReadVCAPConfig()
	creds, err := GetUserProvidedCredentials("mc")
	if err != nil {
		t.Error("Could not read user-provided credentials from env.")
	}
	_, ok := creds["admin_username"]
	if !ok {
		t.Error("Could not find a username")
	}
}

func TestReadS3(t *testing.T) {
	buffer, err := ioutil.ReadFile("example.json")
	if err != nil {
		t.Error("Could not read example.json")
	}
	os.Setenv("VCAP_SERVICES", string(buffer))
	ReadVCAPConfig()

	creds, err := GetS3Credentials("backups")
	if err != nil {
		t.Error("Could not read backups credentials from s3.")
	}
	if creds["access_key_id"] != "ACCESSKEYIDALPHA" {
		t.Error("Did not get s3 access key ACCESSKEYIDALPHA")
	}
	if creds["secret_access_key"] != "SECRETACCESSKEY+ALPHA" {
		t.Error("Did not get s3 secret key SECRETACCESSKEY+ALPHA")
	}
}
