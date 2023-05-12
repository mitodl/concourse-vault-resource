.PHONY: build

fmt:
	@go fmt ./...

tidy:
	@go mod tidy

get:
	@go get github.com/mitodl/concourse-vault-resource

build: tidy
	@go build -o check cmd/check/main.go
	@go build -o in cmd/in/main.go
	@go build -o out cmd/out/main.go

release: tidy
	@go build -o check -ldflags="-s -w"  cmd/check/main.go
	@go build -o in -ldflags="-s -w"  cmd/in/main.go
	@go build -o out -ldflags="-s -w"  cmd/out/main.go

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
