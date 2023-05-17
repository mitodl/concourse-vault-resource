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
	Token        string `json:"token,omitempty"`
	Insecure     bool   `json:"insecure"`
}

type Metadata struct {
	// key is secret "<mount>-<path>" and value is secret keys and values
	// key-value pairs would be arbitrary for kv1 and kv2, but are standardized schema for credential generators
	Values map[string]map[string]interface{} `json:"values"`
}

// in custom type struct for inputs and outputs TODO: concourse may be passing `params` value merged with value of `source` and retain `source` as key for post-merge according to comcast resource, but that sounds ridiculous; if params key is being passed then restructure so redundant `secrets` key is removed
type inRequest struct {
	// key is secret mount, and nested map is paths-[<path>, <path>] and engine-<engine>
	// cannot use nested structs because mount keys are arbitrary
	Params  map[string]map[string]any `json:"params"`
	Source  Source                    `json:"source"`
	Version int                       `json:"version"`
}

type inResponse struct {
	Metadata Metadata `json:"metadata"`
	Version  int      `json:"version"`
}

// inRequest constructor with pipeline param as io.Reader but typically os.Stdin *os.File input because concourse
func NewInRequest(pipelineJSON io.Reader) *inRequest {
	// read, decode, and unmarshal the pipeline json io.Reader, and assign to the inRequest pointer
	var inRequest inRequest
	if err := json.NewDecoder(pipelineJSON).Decode(&inRequest); err != nil {
		log.Print("error decoding stdin from JSON")
		log.Fatal(err)
	}

  // return reference
	return &inRequest
}

// inResponse constructor
func NewInResponse(version int) *inResponse {
	// return reference to initialized struct
	return &inResponse{
		Version:  version,
		Metadata: Metadata{Values: map[string]map[string]interface{}{}},
	}
}
