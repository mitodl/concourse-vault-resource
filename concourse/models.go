package concourse

import (
	"encoding/json"
	"io"
	"log"
)

// concourse standard custom type structs
type Source struct {
	AuthEngine   string `json:"auth_engine,omitempty"`
	Address      string `json:"address,omitempty"`
	AWSMountPath string `json:"aws_mount_path,omitempty"`
	AWSVaultRole string `json:"aws_vault_role,omitempty"`
	Token        string `json:"token,omitempty"`
	Insecure     bool   `json:"insecure"`
}

type Metadata struct {
	// key is secret "<mount>-<path>" and value is secret keys and values
	// key-value pairs would be arbitrary for kv1 and kv2, but are standardized schema for credential generators
	Values map[string]map[string]interface{} `json:"values"`
}

type Version struct {
	Version string `json:"version"`
}

// custom type structs
type Secrets struct {
	Engine string   `json:"engine"`
	Paths  []string `json:"paths"`
}

// check/in custom type struct for inputs and outputs
type CheckResponse struct {
	Versions []Version
}

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
		Metadata: Metadata{Values: map[string]map[string]interface{}{}},
	}
}
