package concourse

import (
	"encoding/json"
	"io"
	"log"
)

// TODO https://itnext.io/how-to-use-golang-generics-with-structs-8cabc9353d75

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

type SecretSource struct {
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
	AuthEngine   string       `json:"auth_engine,omitempty"`
	Address      string       `json:"address,omitempty"`
	AWSMountPath string       `json:"aws_mount_path,omitempty"`
	AWSVaultRole string       `json:"aws_vault_role,omitempty"`
	Token        string       `json:"token,omitempty"`
	Insecure     bool         `json:"insecure"`
	Secret       SecretSource `json:"secret"`
}

// key is "<mount>-<path>" and value is version of secret
type Version map[string]string

// check/in/out custom type structs for inputs and outputs
type checkRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

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

// inRequest constructor with pipeline param as io.Reader but typically os.Stdin *os.File input because concourse
func NewCheckRequest(pipelineJSON io.Reader) *checkRequest {
	// read, decode, and unmarshal the pipeline json io.Reader, and assign to the inRequest pointer
	var checkRequest checkRequest
	if err := json.NewDecoder(pipelineJSON).Decode(&checkRequest); err != nil {
		log.Print("error decoding pipline input from JSON")
		log.Fatal(err)
	}

	// initialize empty version if unspecified
	if checkRequest.Version == nil {
		checkRequest.Version = map[string]string{}
	} else if checkRequest.Source.Secret.Engine == "kv1" && checkRequest.Version != nil {
		// validate version not specified for kv1
		log.Fatal("version cannot be specified in conjunction with a kv version 1 engine secret")
	}

	return &checkRequest
}

// checkResponse constructor NOW oops this should be empty and request should be populated with version; or possible use as both?
func NewCheckResponse() *checkResponse {
	// return reference to slice of version
	return &checkResponse{}
}

// inRequest constructor with pipeline param as io.Reader but typically os.Stdin *os.File input because concourse
func NewInRequest(pipelineJSON io.Reader) *inRequest {
	// read, decode, and unmarshal the pipeline json io.Reader, and assign to the inRequest pointer
	var inRequest inRequest
	if err := json.NewDecoder(pipelineJSON).Decode(&inRequest); err != nil {
		log.Print("error decoding pipline input from JSON")
		log.Fatal(err)
	}
	// initialize request version with empty map
	if inRequest.Version != nil {
		log.Print("version is currently ignored in the get step as it must be tied to a specific secret path")
	}
	inRequest.Version = map[string]string{}
	// validate params versus source.secret
	if inRequest.Source.Secret != (SecretSource{}) && inRequest.Params != nil {
		log.Fatal("secrets cannot be simultaneously specified in both source and params")
	} else if inRequest.Source.Secret == (SecretSource{}) && inRequest.Params == nil {
		log.Fatal("one secret must be specified in source, or one or more secrets in params, and neither was specified")
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
	// validate
	if outRequest.Source.Secret != (SecretSource{}) {
		log.Print("specifying a secret in source for a put step has no effect, and that value will be ignored during this step execution")
	}
	if outRequest.Params == nil {
		log.Fatal("no secret parameters were specified for this put step")
	}

	// return reference
	return &outRequest
}
