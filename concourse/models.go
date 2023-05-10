package concourse

// concourse standard custom type structs
type Source struct {
	// key is secret mount, and nested map is path-[<path>] and engine-<engine>
	// maybe concourse and go would be fine with a slice of custom type structs here, but not interested in fighting that battle
	Secrets      map[string]map[string]any `json:"secrets"`
	Engine       string                    `json:"auth_engine"`
	Address      string                    `json:"address"`
	AWSMountPath string                    `json:"aws_mount_path"`
	Token        string                    `json:"token"`
	Insecure     bool                      `json:"insecure"`
}

type Version struct {
	Version int `json:"version"`
}

type Metadata struct {
	// key is secret "<mount>-<path>" and value is secret keys and values
	// key-value pairs would be arbitrary for kv1 and kv2, but are standardized schema for credential generators
	Values map[string]map[string]interface{} `json:"values"`
}

// in custom type struct for inputs and outputs
type InRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type InResponse struct {
	Metadata Metadata `json:"metadata"`
	Version  Version  `json:"version"`
}
