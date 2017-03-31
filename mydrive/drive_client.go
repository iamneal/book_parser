package mydrive

import (
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	drive "google.golang.org/api/drive/v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	SECRET_LOCATION_NAME = "SECRET_LOCATION"
)

func GetDriveClient(config *oauth2.Config, token *oauth2.Token) (*drive.Service, error) {
	url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("web browser should prompt you at this url: %s\n\ncode: ", url)
	var code string

	//read typed code from stdin
	_, err := fmt.Scan(&code)
	if err != nil {
		return nil, err
	}

	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	tokenSrc := oauth2.ReuseTokenSource(tok, nil)
	oauthClient := oauth2.NewClient(context.Background(), tokenSrc)
	if err != nil {
		return nil, err
	}
	return drive.New(oauthClient)
}

func GetGoogleDriveConfig() (*oauth2.Config, error) {
	secretLoc := "./client_secret.json"

	if setLoc := os.Getenv(SECRET_LOCATION_NAME); setLoc != "" {
		secretLoc = setLoc
	}
	absSecretLoc, err := filepath.Abs(secretLoc)
	if err != nil {
		return nil, err
	}

	file, err := ioutil.ReadFile(absSecretLoc)
	if err != nil {
		return nil, err
	}

	return google.ConfigFromJSON(file, drive.DriveScope)
}
