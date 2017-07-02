# secretctl

secretctl is a secret management tool that uses GPG and Vault as storage. This allows you to push secrets to storage like Git, or pull secrets from storage.

## Install

Go to the [releases page](https://github.com/summerwind/secretctl/releases), find the version you want, and download the tarball file.

## Getting started

You need to create a YAML file describing the secret managed by storage as follows.

```
$ vim .secret.yml
```
```
storage:
  vault:
    addr: https://127.0.0.1:8200

files:
  secrets/token.txt:
    vault:
      path: secret/token.txt

  config/database.yml:
    vault:
      path: secret/config/database.yml

env_vars:
  AWS_ACCESS_KEY_ID:
    vault:
      path: secret/aws/access_key_id

  AWS_SECRET_ACCESS_KEY:
    vault:
      path: secret/aws/secret_access_key
```

Push your secrets based on the configuration file.

```
$ secretctl push --vault-token ${VAULT_TOKEN}
[File] Pushed: secrets/token.txt
[File] Pushed: config/database.yml
[EnvVar] Pushed: AWS_ACCESS_KEY_ID
[EnvVar] Pushed: AWS_SECRET_ACCESS_KEY
```

After pushing, you can get your secret files with the `pull` command.

```
$ secretctl pull --vault-token ${VAULT_TOKEN}
[File] Pulled: secrets/token.txt
[File] Pulled: config/database.yml
```

You can use the `exec` command to execute any command with the pushed secret environment variables.

```
$ secretctl exec --vault-token ${VAULT_TOKEN} -- aws ec2 describe-instances
```

