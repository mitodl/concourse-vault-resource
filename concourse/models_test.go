package concourse

import "testing"

// test inRequest constructor
func TestNewInRequest(test *testing.T) {
	//inRequest := NewInRequest()
}

// test inResponse constructor
func TestNewInResponse(test *testing.T) {
	version := Version{Version: 2}
	inResponse := NewInResponse(version)

	if len(inResponse.Metadata.Values) != 0 || inResponse.Version != version {
		test.Error("the in response constructor returned unexpected values")
		test.Errorf("expected Metadata Values field to be empty map, actual: %v", inResponse.Metadata)
		test.Errorf("expected Version: %d, actual: %d", version, inResponse.Version)
	}
}
