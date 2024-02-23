package vcap

import (
	"os"
	"testing"
)

var db_vcap = `{
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
