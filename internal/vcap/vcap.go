package vcap

import (
	"bytes"
	"os"

	"golang.org/x/exp/slices"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gov.gsa.fac.cgov-util/internal/logging"
)

type CredentialsRDS struct {
	DB_Name  string `json:"db_name"`
	Host     string `json:"host"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Uri      string `json:"uri"`
}

type CredentialsS3 struct {
	Uri                string   `json:"uri"`
	InsecureSkipVerify bool     `json:"insecure_skip_verify"`
	AccessKeyId        string   `json:"access_key_id"`
	SecretAccessKey    string   `json:"secret_access_key"`
	Region             string   `json:"region"`
	Bucket             string   `json:"bucket"`
	Endpoint           string   `json:"endpoint"`
	FipsEndpoint       string   `json:"fips_endpoint"`
	AdditionalBuckets  []string `json:"additional_buckets"`
}

type InstanceS3 struct {
	Label        string        `json:"label"`
	Plan         string        `json:"plan"`
	Name         string        `json:"name"`
	Tags         []string      `json:"tags"`
	InstanceGuid string        `json:"instance_guid"`
	InstanceName string        `json:"instance_name"`
	BindingGuid  string        `json:"binding_guid"`
	BindingName  string        `json:"binding_name"`
	Credentials  CredentialsS3 `json:"credentials"`
}

type InstanceRDS struct {
	Label          string         `json:"label"`
	Provider       string         `json:"provider"`
	Plan           string         `json:"plan"`
	Name           string         `json:"name"`
	Tags           []string       `json:"tags"`
	InstanceGuid   string         `json:"instance_guid"`
	InstanceName   string         `json:"instance_name"`
	BindingGuid    string         `json:"binding_guid"`
	BindingName    string         `json:"binding_name"`
	Credentials    CredentialsRDS `json:"credentials"`
	SyslogDrainUrl string         `json:"syslog_drain_url"`
	VolumeMounts   string         `json:"volume_mounts"`
}

type UserProvided struct {
	Label        string            `json:"label"`
	Name         string            `json:"name"`
	Tags         []string          `json:"tags"`
	InstanceGuid string            `json:"instance_guid"`
	InstanceName string            `json:"instance_name"`
	BindingGuid  string            `json:"binding_guid"`
	BindingName  string            `json:"binding_name"`
	Credentials  map[string]string `json:"credentials"`
}

func GetRDSCredentials(label string) (*CredentialsRDS, error) {
	var instanceSlice []InstanceRDS
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
func GetLocalRDSCredentials(label string) (*CredentialsRDS, error) {
	var instanceSlice []InstanceRDS
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

type UserProvidedCredentials = map[string]string

// Returns a map, not a pointer to a structure
func GetUserProvidedCredentials(label string) (UserProvidedCredentials, error) {
	var instanceSlice []UserProvided
	err := viper.UnmarshalKey("user-provided", &instanceSlice)
	if err != nil {
		logging.Logger.Println("Could not unmarshal aws-rds from VCAP_SERVICES")
	}
	for _, instance := range instanceSlice {
		if instance.Label == label {
			return instance.Credentials, nil
		}
	}
	return nil, errors.Errorf("No credentials found for '%s'", label)
}

func GetS3Credentials(label string) (*CredentialsS3, error) {
	var instanceSlice []InstanceS3
	err := viper.UnmarshalKey("s3", &instanceSlice)
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

func GetRDSCreds(source_db string, dest_db string) (*CredentialsRDS, *CredentialsRDS) {
	var source *CredentialsRDS
	var dest *CredentialsRDS
	var err error

	if slices.Contains([]string{"LOCAL", "TESTING"}, os.Getenv("ENV")) {
		if source_db != "" {
			source, err = GetLocalRDSCredentials(source_db)
			if err != nil {
				logging.Logger.Println("BACKUPS Cannot get local source credentials")
				logging.Logger.Println(err)
				os.Exit(-1)
			}
		} else {
			source = nil
		}
		if dest_db != "" {
			dest, err = GetLocalRDSCredentials(dest_db)
			if err != nil {
				logging.Logger.Println("BACKUPS Cannot get local dest credentials")
				os.Exit(-1)
			}
		} else {
			dest = nil
		}

	} else {
		if source_db != "" {
			source, err = GetRDSCredentials(source_db)
			if err != nil {
				logging.Logger.Println("BACKUPS Cannot get RDS source credentials")
				os.Exit(-1)
			}
		} else {
			source = nil
		}
		if dest_db != "" {
			dest, err = GetRDSCredentials(dest_db)
			if err != nil {
				logging.Logger.Println("BACKUPS Cannot get RDS dest credentials")
				os.Exit(-1)
			}
		} else {
			dest = nil
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
