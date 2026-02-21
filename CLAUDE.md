# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A minimal Go CLI tool that generates GitHub App installation access tokens. It reads GitHub App credentials from CLI flags (with environment variable fallback), creates a signed JWT, exchanges it for an installation token via the GitHub API, and prints the token to stdout.

## Build & Run

```bash
go build -o github-app-token    # build binary
go vet ./...                     # static analysis
go fmt ./...                     # format code
```

There are no tests currently. Standard `go test ./...` would run them if added.

## Configuration

Accepts CLI flags or environment variables (flags take precedence):

| Flag                | Env Var                          | Description                              |
|---------------------|----------------------------------|------------------------------------------|
| `-app-id`           | `GITHUB_APP_ID`                  | GitHub App identifier                    |
| `-key-path`         | `GITHUB_APP_PRIVATE_KEY_PATH`    | Filesystem path to PEM-encoded RSA private key (PKCS#1) |
| `-installation-id`  | `GITHUB_APP_INSTALLATION_ID`     | Installation ID for the target repo/org  |

## Architecture

Single-file application (`main.go`, ~80 lines) with three functions:

- `main()` — orchestrates the full flow: parse flags/env vars, load key, sign JWT (RS256, 10min expiry), POST to GitHub API, output token
- `loadPrivateKey(path)` — reads and parses a PKCS#1 PEM private key file
- `flagOrEnv(val, envKey)` — uses flag value if set, otherwise falls back to env var, fatals if both empty

The only external dependency is `github.com/golang-jwt/jwt/v5` for JWT creation. The GitHub API endpoint used is `POST /app/installations/{id}/access_tokens` (API version 2022-11-28).
