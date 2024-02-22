package vcap

import (
	"fmt"
	"testing"
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

func TestOne(t *testing.T) {
	fmt.Println("Hi")
}
