package server

import(
	"time"
	"errors"
	"io/ioutil"
	"path/filepath"
	"golang.org/x/oauth2"
	"net/http"
	drive "google.golang.org/api/drive/v3"
	"golang.org/x/oauth2/google"
)

var BadCache = errors.New("the access token has expired")

type UpdateToken struct {
	Update func(*oauth2.Token, *User)
	User *User
	Token *oauth2.Token
	Parent oauth2.TokenSrc
}

func (ut *UpdateToken) Token() (*oauth2.Token, error) {
	if ut.User == nil {
		return ut.Token, fmt.Errorf("No user set in UpdateToken")
	}
	tok, err := ut.Parent.Token()
	if err != nil {
		return nil, err
	}
	if tok.AccessToken != ut.Token.AccessToken {
		ut.Update(tok, ut.User)
	}
	ut.Token = tok
	return ut.Token, nil
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

type OAuth2TokenCache struct {
	Tokens map[string] *UserCache
	Config *oauth2.Config
}

func NewOAuth2TokenCache() (*OAuth2TokenCache, error) {
	conf, err := GetGooleAuthConfig()
	if err != nil {
		return nil, err
	}
	// TODO instead of starting with a new map everytime,  read serialized tokens from a file
	return &OAuth2TokenCache{
		Tokens make(map[string]*UserCache)
		Config: conf,
	}
}

func (t *OAuth2TokenCache) Get(tok string) (*UserCache, error ) {
	cache := t.Tokens[tok]
	if cache == nil {
		return nil, BadCache
	}
	if time.Now().Before(cache.Token.Expiry) {
		return cache
	}
	return nil, BadCache
}

func (t *OAuth2TokenCache) NewToken(ctx context, code string) (*Token, error) {
	tok, err := t.Config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	tknSrc := UpdateToken{
		Parent: oauth2.ReuseTokenSource(tok, nil),
		Token: tok,
		Update: func(oldt, newt oauth2.Token) {
			if cache, exists := t.Tokens[oldt.AccessToken]; exists && cache != nil {
				delete(t.Tokens[olt.AccessToken])
				t.Tokens[newt] = cache
				cache.Token = token
			} else {
				fmt.Printf("USER DID NOT EXIST on UpdateToken request: %#v\n", u)
				t.Tokens[newt] = &UserCache{
					Token: token,
					RefreshToken: token.RefreshToken,
				}
			}
		},
	}
	tempCli := t.Config.NewClient(ctx, tknSrc)
	resp, err := tempCli.Get("https://www.googleapis.com/userinfo/v2/me")
	if err != nil {
		return nil, err
	}
	userBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	user, err := t.UserFromBytes(userBytes)
	if err != nil {
		return nil, err
	}
	driveCli, err := drive.New(httpCli)
	if err != nil {
		return nil
	}
	t.Tokens[tok] = user
	t.Clients[user] = &UserCache {
		AccessToken: tok,
		Http: httpCli,
		Drive: drive,
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
