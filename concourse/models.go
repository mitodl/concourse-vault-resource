package concourse

import (
	"encoding/json"
	"io"
	"log"
)

// custom type structs
type Secrets struct {
	Engine string   `json:"engine"`
	Paths  []string `json:"paths"`
}

// key is secret "<mount>-<path>", and value is secret keys and values
// key-value pairs would be arbitrary for kv1 and kv2, but are standardized schema for credential generators
type SecretValue map[string]interface{}

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

// concourse standard custom type structs
type Source struct {
	AuthEngine   string `json:"auth_engine,omitempty"`
	Address      string `json:"address,omitempty"`
	AWSMountPath string `json:"aws_mount_path,omitempty"`
	AWSVaultRole string `json:"aws_vault_role,omitempty"`
	Token        string `json:"token,omitempty"`
	Insecure     bool   `json:"insecure"`
}

type Metadata []SecretValue

type Version struct {
	Version string `json:"version"`
}

// check/in custom type structs for inputs and outputs
type CheckResponse []Version

type inRequest struct {
	// key is secret mount, and nested map is paths-[<path>, <path>] and engine-<engine>
	// cannot use nested structs because mount keys are arbitrary
	Params  map[string]Secrets `json:"params"`
	Source  Source             `json:"source"`
	Version Version            `json:"version"`
}

type inResponse struct {
	Metadata Metadata `json:"metadata"`
	Version  Version  `json:"version"`
}

// inRequest constructor with pipeline param as io.Reader but typically os.Stdin *os.File input because concourse
func NewInRequest(pipelineJSON io.Reader) *inRequest {
	// read, decode, and unmarshal the pipeline json io.Reader, and assign to the inRequest pointer
	var inRequest inRequest
	if err := json.NewDecoder(pipelineJSON).Decode(&inRequest); err != nil {
		log.Print("error decoding pipline input from JSON")
		log.Fatal(err)
	}

	// return reference
	return &inRequest
}

// inResponse constructor
func NewInResponse(version Version) *inResponse {
	// return reference to initialized struct
	return &inResponse{
		Version:  version,
		Metadata: Metadata([]SecretValue{}),
	}
}
