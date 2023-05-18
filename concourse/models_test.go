package concourse

import "testing"

// test inRequest constructor
func TestNewInRequest(test *testing.T) {
	//inRequest := NewInRequest()
}

// test inResponse constructor
func TestNewInResponse(test *testing.T) {
	version := Version{Version: "2"}
	inResponse := NewInResponse(version)

	if len(inResponse.Metadata) != 0 || inResponse.Version != version {
		test.Error("the in response constructor returned unexpected values")
		test.Errorf("expected Metadata field to be empty slice, actual: %v", inResponse.Metadata)
		test.Errorf("expected Version: %s, actual: %s", version, inResponse.Version)
	}
}
