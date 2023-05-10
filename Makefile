.PHONY: build

fmt:
	@go fmt ./...

tidy:
	@go mod tidy

build: tidy
	@go build -o foo

release: tidy
	@go build -s -w -o foo

unit:
	@go test -v ./...

#TODO: use api module
bootstrap:
	#@vault server -dev
	@vault auth enable aws
	@vault secrets enable database
	@vault secrets enable aws
	@vault secrets enable -version=1 kv
	@vault kv put -mount=kv foo/bar password=supersecret
	@vault kv put -mount=secret foo/bar password=supersecret
