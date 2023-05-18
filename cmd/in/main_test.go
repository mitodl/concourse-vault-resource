package main

import (
	"os"
	_ "testing"
)

// TODO: add second kv2 pair at same secret mount
func ExampleMain() {
	// deliver test pipeline file content as stdin to "in" the same as actual pipeline execution
	os.Stdin, _ = os.OpenFile("fixtures/token_kv.json", os.O_RDONLY, 0o644)
	defer os.Stdin.Close()

	// invoke main and validate stdout
	main()
  // Output: {"metadata":[{"secret-foo/bar":{"password":"supersecret"}},{"kv-foo/bar":{"password":"supersecret"}}],"version":{"version":""}}
}
