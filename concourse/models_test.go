package concourse

import "testing"

const versionKey = "secret-foo/bar"

var version = map[string]string{versionKey: "1"}

// test checkrequest constructor

// test checkresponse constructor
func TestCheckResponse(test *testing.T) {
	checkResponse := NewCheckResponse([]Version{version})

	if len(checkResponse) != 1 || checkResponse[0] == nil {
		test.Error("the check response constructor returned an unexpected value")
		test.Errorf("expected value: &[], actual: %v", checkResponse)
	}
}

// test inRequest constructor
func TestNewInRequest(test *testing.T) {
	//inRequest := NewInRequest()
}

// test inResponse constructor
func TestNewInResponse(test *testing.T) {
	inResponse := NewResponse(version)

	if len(inResponse.Metadata) != 0 || inResponse.Version[versionKey] != version[versionKey] {
		test.Error("the in response constructor returned unexpected values")
		test.Errorf("expected Metadata field to be empty slice, actual: %v", inResponse.Metadata)
		test.Errorf("expected Version: %v, actual: %v", version, inResponse.Version)
	}
}

// test outRequest constructor

// test outResponse constructor
func TestOutResponse(test *testing.T) {
	outResponse := NewResponse(nil)
	outResponse.Version = version

	if len(outResponse.Metadata) != 0 || outResponse.Version[versionKey] != version[versionKey] {
		test.Error("the in response constructor returned unexpected values")
		test.Errorf("expected Metadata field to be slice of one element, actual: %v", outResponse.Metadata)
		test.Errorf("expected Metadata field only element to be empty map, actual: %v", outResponse.Metadata[0])
		test.Errorf("expected Version: %v, actual: %v", version, outResponse.Version)
	}
}
