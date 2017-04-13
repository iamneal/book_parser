package server

import(
	"os"
	"fmt"
	"context"
	"time"
	"errors"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"golang.org/x/oauth2"
	"net/http"
	drive "google.golang.org/api/drive/v3"
	"golang.org/x/oauth2/google"
)

var BadCache = errors.New("the access token has expired")

type UpdateToken struct {
	Update func(*oauth2.Token, *oauth2.Token)
	User *User
	CurToken *oauth2.Token
	Parent oauth2.TokenSource
}

func (ut *UpdateToken) Token() (*oauth2.Token, error) {
	if ut.User == nil {
		return ut.CurToken, nil
	}
	token, err := ut.Parent.Token()
	if err != nil {
		return nil, err
	}
	if token.AccessToken != ut.CurToken.AccessToken {
		ut.Update(ut.CurToken, token)
	}
	ut.CurToken = token
	return ut.CurToken, nil
}



type User struct {
	Id string `json:"id,omitempty"`
	FirstName string `json:"given_name,omitempty"`
	LastName string `json:"family_name,omitempty"`
}

type UserCache struct {
	Token *oauth2.Token
	User *User
	Http *http.Client
	Drive *drive.Service
}

func (uc *UserCache) String() string {
	return fmt.Sprintf("UserCache:\n\t\tUser: %#v", uc.User)
}

func NewUserFromBytes(bytes []byte) (*User, error) {
	fmt.Printf("userbytes: %s\n", bytes)
	user := new(User)
	err := json.Unmarshal(bytes, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

type OAuth2TokenCache struct {
	Tokens map[string] *UserCache
	Config *oauth2.Config
}

func (oa *OAuth2TokenCache) String() string {
	pr := "OAuth2TokenCache\n\tTokens:\n\t"
	for k, v := range oa.Tokens {
		pr += fmt.Sprintf("%s: %s\n", k, v)
	}
	return pr
}

func NewOAuth2TokenCache() (*OAuth2TokenCache, error) {
	conf, err := GetGoogleAuthConfig()
	if err != nil {
		return nil, err
	}
	// TODO instead of starting with a new map everytime,  read serialized tokens from a file
	return &OAuth2TokenCache{
		Tokens: make(map[string]*UserCache),
		Config: conf,
	}, nil
}

func (t *OAuth2TokenCache) Get(tok string) (*UserCache, error ) {
	cache := t.Tokens[tok]
	if cache == nil {
		return nil, BadCache
	}
	if time.Now().Before(cache.Token.Expiry) {
		return cache, nil
	}
	return nil, BadCache
}

func (t *OAuth2TokenCache) UpdateOldToken(tok string) (string, error) {
	cache := t.Tokens[tok]
	if cache == nil || cache.Token == nil || cache.Token.AccessToken == "" {
		return "", fmt.Errorf("cache for token not found or updated. Token: %s", tok)
	}
	if cache.Token.AccessToken == tok {
		fmt.Println("tried to update a token that is up to date...")
	}
	t.Tokens[cache.Token.AccessToken] = cache
	delete(t.Tokens, tok)

	return cache.Token.AccessToken, nil
}

func (t *OAuth2TokenCache) NewToken(ctx context.Context, code string) (*oauth2.Token, error) {
	tok, err := t.Config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	tknSrc := &UpdateToken{
		Parent: oauth2.ReuseTokenSource(tok, nil),
		CurToken: tok,
		Update: func(oldt, newt *oauth2.Token) {
			if cache, exists := t.Tokens[oldt.AccessToken]; exists && cache != nil {
				// delete(t.Tokens[olt.AccessToken])
				// t.Tokens[newt] = cache
				cache.Token = newt
			} else {
				fmt.Printf("USER DID NOT EXIST on UpdateToken request: %#v\n", oldt)
				t.Tokens[newt.AccessToken] = &UserCache{Token: newt}
			}
		},
	}
	httpCli := oauth2.NewClient(ctx, tknSrc)
	resp, err := httpCli.Get("https://www.googleapis.com/userinfo/v2/me")
	if err != nil {
		return nil, err
	}
	userBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	user, err := NewUserFromBytes(userBytes)
	if err != nil {
		return nil, err
	}
	driveCli, err := drive.New(httpCli)
	if err != nil {
		return nil, err
	}
	t.Tokens[tok.AccessToken] = &UserCache {
		Token: tok,
		User: user,
		Http: httpCli,
		Drive: driveCli,
	}

	return tok, nil
}

func GetGoogleAuthConfig() (*oauth2.Config, error) {
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
	scopes := []string{
		"https://www.googleapis.com/auth/userinfo.profile",
		drive.DriveScope,
	}
	return google.ConfigFromJSON(file, scopes...)
}
