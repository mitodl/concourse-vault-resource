package vault

import (
  "log"
  "fmt"
  "context"

  vault "github.com/hashicorp/vault/api"

  "github.com/mitodl/concourse-vault-resource/concourse"
)

// secret engine with pseudo-enum
type SecretEngine string

const (
  database  SecretEngine = "DB"
  awsIAM    SecretEngine = "AWS IAM"
  keyvalue1 SecretEngine = "KV1"
  keyvalue2 SecretEngine = "KV2"
)

// secret defines a composite Vault secret configuration
type vaultSecret struct {
  engine SecretEngine
  path   string
  mount  string
}

func NewSecret(params *concourse.Params) *vaultSecret {
  // initialize vault secret return from concourse params
  newVaultSecret := &vaultSecret{
    engine: params.SecretEngine,
    path:   params.SecretPath,
    mount:  params.MountPath,
  }

  return newVaultSecret
}

// retrieve and return secrets; TODO: efficiency vs. compactness re-evaluate
func (secret *vaultSecret) retrieveSecret(client *vault.Client) map[string]interface{} {
  // declare func scope variables
  var secretValue map[string]interface{}
  var err error

  switch secret.engine {
  case database:
    // read db secret (in this case path signifies role name)
    response, err := client.Secrets.DatabaseGenerateCredentials(
      context.Background(),
      path,
      vaultMountPath,
    )
  case awsIAM:
    // read aws iam secret (in this case path signifies role name)
    response, err = client.Secrets.AwsGenerateCredentials(
      context.Background(),
      path,
      vaultMountPath,
    )
  }

  // verify secret read
  if err != nil {
    log.Printf("Failed to read secret %s from %s secrets engine", secret.path, secret.engine)
    log.Fatal(err)
  }

  return secretValue
}


func (secret *vaultSecret) retrieveKVSecret(client *vault.Client) map[string]interface{} {
  // declare func scope variable
  kvClient *vaut.Client.KVv1 //TODO or KVv2... generics?

  switch secret.engine {
  case keyvalue1:
    // initialize kv1 client
    kvClient = client.KVv1(secret.mount)
  case keyvalue2:
    // initialize kv2 client
    kvClient = client.KVv2(secret.mount)
  }

  // read kv secret
  kvSecret, err := kvClient.Get(
    context.Background(),
    secret.path,
  )

  // verify secret read
  if err != nil {
    log.Printf("Failed to read secret %s from %s secrets engine", secret.path, secret.engine)
    log.Fatal(err)
  }

  // return secret value as map[string]interface{}
  return kvSecret.Data
}
