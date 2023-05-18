# Concourse Vault Resource

A [concourse-ci](https://concourse-ci.org) resource for interacting with secrets via [Vault](https://www.vault.io).

## Behavior

### `source`: designates the Vault server and authentication engine information

**parameters**
- `auth_engine`: _optional_ The authentication engine for use with Vault. Allowed values are `aws` or `token`. If unspecified will default to `aws` with no `token` parameter specified, or `token` if `token` parameter is specified.

- `address`: _optional_ The address for the Vault server in format of `URL:PORT`. default: `http://127.0.0.1:8200`

- `aws_mount_path`: _optional_ The mount path for the AWS authentication engine. Parameter is ignored if authentication engine is not `aws`. default: `aws`

- `aws_iam_role`: _optional_ The AWS IAM Role for the AWS authentication login to Vault. Parameter is ignored if authentication engine is not `aws`. default: (default IAM Role according to utilized AWS config)

- `token`: _optional_ The token for the token authentication engine. Required if `auth_engine` parameter is `token`.

- `insecure`: _optional_ Whether to utilize an insecure connection with Vault (e.g. no HTTP or HTTPS with self-signed cert). default: `false`

### `check`: not implemented

### `in`: interacts with the supported Vault secrets engines

**parameters**

- `<secret_mount path>`: _required_ One or more map/hash/dictionary of the following YAML schema for specifying the secrets to retrieve and/or generate.

```yaml
<secret_mount_path>:
  paths:
  - <path/to/secret>
  - <path/to/other_secret>
  engine: <secret engine> # supported values: database, aws, kv1, kv2
```

### `out`: not implemented

## Example

```yaml
resource_types:
- name: vault
  type: docker-image
  source:
    repository: mitodl/concourse-vault-resource
    tag: latest

resources:
- name: vault
  type: vault
  source:
    address: https://mitodl.vault.com:8200
    auth_engine: aws

jobs:
- name: do something
  plan:
  - get: my-code
  - get: vault
    params:
      postres-mitxonline:
        paths:
        - readonly
        engine: database
      secret:
        paths:
        - path/to/secret
        engine: kv2
      kv:
        paths:
        - path/to/secret
        engine: kv1
```
