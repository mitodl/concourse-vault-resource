package vault

import (
	"context"
	"log"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/aws"
)

// authentication engine with pseudo-enum
type AuthEngine string

const (
	awsIam AuthEngine = "aws iam"
	token  AuthEngine = "token"
)

// VaultConfig defines vault api interface config
type VaultConfig struct {
	Engine       AuthEngine
	Address      string
	AWSMountPath string
	Token        string
	Insecure     bool
}

// VaultConfig constructor
func (config *VaultConfig) New() {
	// validate authentication inputs; TODO: engine now specified by inputs and not constructor logic, so update validation for issues e.g. aws and token co-specification; also validate on missing inputs
	if len(config.Token) > 0 && len(config.AWSMountPath) > 0 {
		log.Fatal("Token and AWS authentication were simultaneously selected; these are mutually exclusive options")
	}
	if len(config.Token) > 0 && len(config.Token) != 28 {
		log.Fatal("the specified Vault Token is invalid")
	}
	/*if len(config.Token) == 0 {
		log.Print("AWS IAM authentication will be utilized with the Vault client")
		config.Engine = awsIam
	} else {
		log.Print("Token authentication will be utilized with the Vault client")
		config.Engine = token
	}*/

	// vault address
	if len(config.Address) == 0 {
		config.Address = "http://127.0.0.1:8200"
	}

	// aws mount path
	if len(config.AWSMountPath) == 0 {
		config.AWSMountPath = "aws"
	}
}

// instantiate authenticated vault client with aws-iam auth
func (config *VaultConfig) AuthClient() *vault.Client {
	// initialize config
	VaultConfig := &vault.Config{Address: config.Address}
	err := VaultConfig.ConfigureTLS(&vault.TLSConfig{Insecure: config.Insecure})
	if err != nil {
		log.Print("Vault TLS configuration failed to initialize")
		log.Fatal(err)
	}

	// initialize client
	client, err := vault.NewClient(VaultConfig)
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
	switch config.Engine {
	case token:
		client.SetToken(config.Token)
	case awsIam:
		// authenticate with aws iam
		awsAuth, err := auth.NewAWSAuth(auth.WithIAMAuth())
		if err != nil {
			log.Fatal("unable to initialize AWS IAM authentication")
		}

		authInfo, err := client.Auth().Login(context.Background(), awsAuth)
		if err != nil {
			log.Print("unable to login to AWS IAM auth method")
			log.Fatal(err)
		}
		if authInfo == nil {
			log.Fatal("no auth info was returned after login")
		}
	}

	// return authenticated vault client
	return client
}
