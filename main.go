package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

var (
	CLIENT_SECRET_NAME   = "CLIENT_SECRET"
	CLIENT_ID_NAME       = "CLIENT_ID"
	SECRET_LOCATION_NAME = "SECRET_LOCATION"
)

func main() {
	creds, err := GetGoogleDriveCreds()
	if err != nil {
		panic(err)
	}
	fmt.Printf("credentials %+v", creds)
}

type Credentials struct {
	ClientSecret string
	ClientId     string
}

func GetGoogleDriveCreds() (*Credentials, error) {
	clientSecret := os.Getenv("CLIENT_SECRET")
	clientId := os.Getenv("CLIENT_ID")

	if clientSecret != "" && clientId != "" {
		return &Credentials{
			ClientSecret: clientSecret,
			ClientId:     clientId,
		}, nil
	}
	secretLoc := "~/.book_parser/client_secret.env"

	if setLoc := os.Getenv(SECRET_LOCATION_NAME); setLoc != "" {
		secretLoc = setLoc
	}
	err := godotenv.Load(secretLoc)

	if err != nil {
		return nil, err
	}
	clientSecret = os.Getenv(CLIENT_SECRET_NAME)
	clientId = os.Getenv(CLIENT_ID_NAME)

	if clientSecret == "" || clientId == "" {
		return nil, fmt.Errorf(`Tried to read file "%s" for env vars "%s" and "%s" but one or both were not set`)
	}

	return &Credentials{
		ClientSecret: clientSecret,
		ClientId:     clientId,
	}, nil
}
