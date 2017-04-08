package mydrive

import (
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	drive "google.golang.org/api/drive/v3"
	"os"
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

