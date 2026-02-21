package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	appID := flag.String("app-id", "", "GitHub App ID (env: GITHUB_APP_ID)")
	keyPath := flag.String("key-path", "", "path to PEM private key (env: GITHUB_APP_PRIVATE_KEY_PATH)")
	installationID := flag.String("installation-id", "", "GitHub App installation ID (env: GITHUB_APP_INSTALLATION_ID)")
	flag.Parse()

	flagOrEnv(appID, "GITHUB_APP_ID")
	flagOrEnv(keyPath, "GITHUB_APP_PRIVATE_KEY_PATH")
	flagOrEnv(installationID, "GITHUB_APP_INSTALLATION_ID")

	key := loadPrivateKey(*keyPath)

	now := time.Now()
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": now.Add(-60 * time.Second).Unix(),
		"exp": now.Add(10 * time.Minute).Unix(),
		"iss": *appID,
	}).SignedString(key)
	if err != nil {
		log.Fatalf("sign jwt: %v", err)
	}

	url := fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", *installationID)
	req, _ := http.NewRequest(http.MethodPost, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("decode response: %v", err)
	}
	if result.Token == "" {
		log.Fatal("empty token in response")
	}

	fmt.Print(result.Token)
}

func loadPrivateKey(path string) *rsa.PrivateKey {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("read key: %v", err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		log.Fatal("failed to decode PEM block")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("parse key: %v", err)
	}
	return key
}

func flagOrEnv(val *string, envKey string) {
	if *val == "" {
		*val = os.Getenv(envKey)
	}
	if *val == "" {
		log.Fatalf("missing required value: set %s env var or see -help for flags", envKey)
	}
}
