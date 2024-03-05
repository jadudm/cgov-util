package structs

type UserProvidedCredentials = map[string]string

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
	SyslogDrainUrl     string   `json:"syslog_drain_url"`
	VolumeMounts       []string `json:"volume_mounts`
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
