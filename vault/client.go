package vault

import (
  "log"
  "context"

  vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/aws"

  "github.com/mitodl/concourse-vault-resource/concourse"
)

// vaultConfig defines vault api interface config
type vaultConfig struct {
	vaultAddr    string
	token        string
	awsMountPath string
	insecure     bool
}

// vaultConfig constructor
func NewVaultConfig(params *concourse.Params) *vaultConfig {
  // initialize vault config return from concourse params
  newVaultConfig := &vaultConfig{
    vaultAddr:    params.VaultAddr,
    token:        params.Token,
    awsMountPath: params.AWSMountPath,
    insecure:     params.Insecure,
  }

  return newVaultConfig
}

// instantiate authenticated vault client with aws-iam auth
func (config *vaultConfig) authClient() *vault.Client {
  // initialize config
	vaultConfig := &vault.Config{Address: config.vaultAddr}
	err := vaultConfig.ConfigureTLS(&vault.TLSConfig{Insecure: config.insecure})
	if err != nil {
		log.Print("Vault TLS configuration failed to initialize")
		log.Fatal(err)
	}

	// initialize client
	client, err := vault.NewClient(vaultConfig)
	if err != nil {
		log.Print("Vault client failed to initialize")
		log.Fatal(err)
	}

  // verify vault is unsealed
  sealStatus, err := client.Sys().SealStatus()
  if err != nil {
    log.Print("unable to verify that the Vault cluster is unsealed")
    log.Fatal(err)
  }
  if sealStatus.Sealed {
    log.Fatal("the Vault cluster is sealed and no operations can be executed")
  }

	// determine authentication method
	if len(config.token) > 0 {
		client.SetToken(config.token)
	} else {
    // authenticate with aws iam
		awsAuth, err := auth.NewAWSAuth(auth.WithIAMAuth())
		if err != nil {
			log.Fatal("Unable to initialize AWS IAM authentication")
		}

		authInfo, err := client.Auth().Login(context.TODO(), awsAuth)
		if err != nil {
			log.Print("Unable to login to AWS IAM auth method")
      log.Fatal(err)
		}
		if authInfo == nil {
			log.Fatal("No auth info was returned after login")
		}
	}

	// return authenticated vault client
  return client
}
