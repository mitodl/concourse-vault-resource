package main

import (
	"os"
	_ "testing"
)

func ExampleMain() {
	// deliver test pipeline file content as stdin to "in" the same as actual pipeline execution
	os.Stdin, _ = os.OpenFile("fixtures/token_kv.json", os.O_RDONLY, 0o644)
	defer os.Stdin.Close()

	// invoke main and validate stdout
	main()
  // Output: {"metadata":{"values":{"kv-foo/bar":{"password":"supersecret"},"secret-foo/bar":{"password":"supersecret"}}},"version":{"version":""}}
}
