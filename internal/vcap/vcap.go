package vcap

import (
	"bytes"
	"fmt"
	"os"

	"golang.org/x/exp/slices"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gov.gsa.fac.cgov-util/internal/logging"
	"gov.gsa.fac.cgov-util/internal/structs"

	"gov.gsa.fac.cgov-util/internal/util"
)

func GetRDSCredentials(label string) (*structs.CredentialsRDS, error) {
	var instanceSlice []structs.InstanceRDS
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
func GetLocalRDSCredentials(label string) (*structs.CredentialsRDS, error) {
	var instanceSlice []structs.InstanceRDS
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

// Returns a map, not a pointer to a structure
func GetUserProvidedCredentials(label string) (structs.UserProvidedCredentials, error) {
	var instanceSlice []structs.UserProvided
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

func GetS3Credentials(name string) (map[string]string, error) {
	var instanceSlice []structs.InstanceS3
	err := viper.UnmarshalKey("s3", &instanceSlice)
	if err != nil {
		logging.Logger.Println("Could not unmarshal s3 from VCAP_SERVICES")
	}
	for _, instance := range instanceSlice {
		if instance.Name == name {
			fmt.Println("INST", instance)
			fmt.Println("AKI", instance.Credentials["access_key_id"])
			fmt.Println("SAK", instance.Credentials["secret_access_key"])
			fmt.Println("REG", instance.Credentials["region"])

			return instance.Credentials, nil
		}
	}

	return nil, errors.Errorf("No credentials found for '%s'", name)
}

func GetRDSCreds(source_db string, dest_db string) (*structs.CredentialsRDS, *structs.CredentialsRDS) {
	var source *structs.CredentialsRDS
	var dest *structs.CredentialsRDS
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
	if util.IsDebugLevel("DEBUG") {
		logging.Logger.Printf("---- VCAP ----\n%s\n---- END VCAP ----\n", vcap)
	}

	viper.ReadConfig(bytes.NewBufferString(vcap))
}
