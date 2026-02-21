# github-app-token

A minimal CLI tool that generates GitHub App installation access tokens. It creates a signed JWT from your App credentials, exchanges it for an installation token via the GitHub API, and prints the token to stdout.

## Installation

### Go install

```bash
go install github.com/bntso/github-app-token@latest
```

### Download from releases

```bash
sudo curl -fsSL "https://github.com/bntso/github-app-token/releases/latest/download/github-app-token-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m | sed 's/x86_64/amd64/' | sed 's/aarch64/arm64/')" \
  -o /usr/local/bin/github-app-token && sudo chmod +x /usr/local/bin/github-app-token
```

Pre-built binaries for Linux and macOS (amd64/arm64) are also available on the [releases page](https://github.com/bntso/github-app-token/releases).

## Usage

The tool accepts configuration via CLI flags or environment variables. Flags take precedence over environment variables.

### With flags

```bash
github-app-token \
  -app-id 12345 \
  -key-path /path/to/private-key.pem \
  -installation-id 67890
```

### With environment variables

```bash
export GITHUB_APP_ID=12345
export GITHUB_APP_PRIVATE_KEY_PATH=/path/to/private-key.pem
export GITHUB_APP_INSTALLATION_ID=67890
github-app-token
```

### Parameters

| Flag                | Env Var                          | Description                              |
|---------------------|----------------------------------|------------------------------------------|
| `-app-id`           | `GITHUB_APP_ID`                  | GitHub App ID                            |
| `-key-path`         | `GITHUB_APP_PRIVATE_KEY_PATH`    | Path to PEM-encoded RSA private key      |
| `-installation-id`  | `GITHUB_APP_INSTALLATION_ID`     | Installation ID for the target repo/org  |

### Example: cloning a private repo

```bash
TOKEN=$(github-app-token -app-id 12345 -key-path key.pem -installation-id 67890)
git clone https://x-access-token:${TOKEN}@github.com/org/repo.git
```
