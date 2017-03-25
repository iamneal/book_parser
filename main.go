package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

var (
	CLIENT_SECRET_NAME   = "CLIENT_SECRET"
	CLIENT_ID_NAME       = "CLIENT_ID"
	SECRET_LOCATION_NAME = "SECRET_LOCATION"
)

type Credentials struct {
	ClientSecret string
	ClientId     string
}

func main() {
	creds, err := GetGoogleDriveCreds()
	if err != nil {
		panic(err)
	}
	fmt.Printf("credentials %+v\n", creds)
}

func GetGoogleDriveCreds() (*Credentials, error) {
	secretLoc := "./client_secret.env"

	if setLoc := os.Getenv(SECRET_LOCATION_NAME); setLoc != "" {
		secretLoc = setLoc
	}
	absSecretLoc, err := filepath.Abs(secretLoc)
	fmt.Printf("secretloc: %s\n abs: %s\n", secretLoc, absSecretLoc)

	if err != nil {
		return nil, err
	}
	err = godotenv.Load(absSecretLoc)

	if err != nil {
		return nil, err
	}
	clientSecret := os.Getenv(CLIENT_SECRET_NAME)
	clientId := os.Getenv(CLIENT_ID_NAME)

	if clientSecret == "" || clientId == "" {
		return nil, fmt.Errorf(`Tried to read file "%s" for env vars "%s" and "%s" but one or both were not set`)
	}

	return &Credentials{
		ClientSecret: clientSecret,
		ClientId:     clientId,
	}, nil
}
