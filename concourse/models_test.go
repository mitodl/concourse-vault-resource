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

// test outRequest constructor

// test outResponse constructor
func TestOutResponse(test *testing.T) {
	version := Version{Version: "2"}
	outResponse := NewOutResponse(version)

	if len(outResponse.Metadata) != 1 || len(outResponse.Metadata[0]) != 0 || outResponse.Version != version {
		test.Error("the in response constructor returned unexpected values")
		test.Errorf("expected Metadata field to be slice of one element, actual: %v", outResponse.Metadata)
		test.Errorf("expected Metadata field only element to be empty map, actual: %v", outResponse.Metadata[0])
		test.Errorf("expected Version: %s, actual: %s", version, outResponse.Version)
	}
}
