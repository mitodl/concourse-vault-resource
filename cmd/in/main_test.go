package main

import (
	"os"
	"testing"
)

func TestE2ERetrieveKVSecrets(test *testing.T) {
	// deliver test pipeline file content as stdin to "in" the same as actual pipeline execution
	os.Stdin, _ = os.OpenFile("fixtures/token_kv.json", os.O_RDONLY, 0o644)
	defer os.Stdin.Close()

	// invoke main
	main()

	// test stdout TODO: decode from json to map and test entries
	// var stdout bytes.Buffer
}
