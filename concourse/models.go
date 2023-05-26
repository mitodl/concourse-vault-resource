package concourse

import (
	"encoding/json"
	"io"
	"log"
)

// custom type structs
// key-value pairs would be arbitrary for kv1 and kv2, but are standardized schema for credential generators
type SecretValue map[string]interface{}

// key is secret "<mount>-<path>", and value is secret keys and values
type SecretValues map[string]SecretValue

type MetadataEntry struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Secrets struct {
	Engine string   `json:"engine"`
	Paths  []string `json:"paths"`
}

// TODO potentially combine both below with above by converting Paths to any (also probably rename) and doing a bunch of type checks BUT wow that seems like not great cost/benefit
type SecretsPut struct {
	Engine string `json:"engine"`
	Patch  bool   `json:"patch"`
	// key is secret path
	Secrets SecretValues `json:"secrets"`
}

type SecretCheck struct {
	Engine string `json:"engine"`
	Mount  string `json:"mount"`
	Path   string `json:"path"`
}

// TODO: for future fine-tuning of secret value (enum?)
type DBSecretValue struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type KVSecretValue map[string]interface{}
type AWSSecretValue struct {
	AccessKey     string `json:"access_key"`
	SecretKey     string `json:"secret_key"`
	SecurityToken string `json:"security_token,omitempty"`
	ARN           string `json:"arn"`
}

// concourse standard type structs
type Source struct {
	AuthEngine   string      `json:"auth_engine,omitempty"`
	Address      string      `json:"address,omitempty"`
	AWSMountPath string      `json:"aws_mount_path,omitempty"`
	AWSVaultRole string      `json:"aws_vault_role,omitempty"`
	Token        string      `json:"token,omitempty"`
	Insecure     bool        `json:"insecure"`
	Secret       SecretCheck `json:"secret,omitempty"`
}

// key is "<mount>-<path>" and value is version of secret
type Version map[string]string

// check/in custom type structs for inputs and outputs
type checkResponse []Version

// TODO use version for specific secret version retrieval
type inRequest struct {
	// key is secret mount
	Params  map[string]Secrets `json:"params"`
	Source  Source             `json:"source"`
	Version Version            `json:"version"`
}

type outRequest struct {
	// key is secret mount
	Params map[string]SecretsPut `json:"params"`
	Source Source                `json:"source"`
}

type response struct {
	Metadata []MetadataEntry `json:"metadata"`
	Version  Version         `json:"version"`
}

// checkResponse constructor NOW oops this should be empty and request should be populated with version; or possible use as both?
func NewCheckResponse(version Version) *checkResponse {
	// return reference to slice of version
	return &checkResponse{version}
}

// inRequest constructor with pipeline param as io.Reader but typically os.Stdin *os.File input because concourse
func NewInRequest(pipelineJSON io.Reader) *inRequest {
	// read, decode, and unmarshal the pipeline json io.Reader, and assign to the inRequest pointer
	var inRequest inRequest
	if err := json.NewDecoder(pipelineJSON).Decode(&inRequest); err != nil {
		log.Print("error decoding pipline input from JSON")
		log.Fatal(err)
	}
	// initialize request version with empty map if unspecified
	if inRequest.Version == nil {
		inRequest.Version = map[string]string{}
	}

	// return reference
	return &inRequest
}

// response constructor
func NewResponse(version Version) *response {
	// default empty version for out
	responseVersion := map[string]string{}

	if version != nil {
		// use input version for in
		responseVersion = version
	}

	// return initialized reference
	return &response{
		Version:  responseVersion,
		Metadata: []MetadataEntry{},
	}
}

// outRequest constructor with pipeline param as io.Reader but typically os.Stdin *os.File input because concourse
func NewOutRequest(pipelineJSON io.Reader) *outRequest {
	// read, decode, and unmarshal the pipeline json io.Reader, and assign to the outRequest pointer
	var outRequest outRequest
	if err := json.NewDecoder(pipelineJSON).Decode(&outRequest); err != nil {
		log.Print("error decoding pipline input from JSON")
		log.Fatal(err)
	}

	// return reference
	return &outRequest
}
