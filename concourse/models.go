package concourse

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
type InRequest struct {
	// key is secret mount, and nested map is paths-[<path>, <path>] and engine-<engine>
	// cannot use nested structs because mount keys are arbitrary
	Secrets map[string]map[string]any `json:"secrets"`
	Source  Source                    `json:"source"`
	Version int                       `json:"version"`
}

type InResponse struct {
	Metadata Metadata `json:"metadata"`
	Version  int      `json:"version"`
}
