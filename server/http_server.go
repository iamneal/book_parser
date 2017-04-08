package server

import (
	"fmt"
	"time"
	"io/ioutil"
  "net/http"
	"encoding/json"
	"encoding/base64"
  "golang.org/x/net/context"
  "github.com/grpc-ecosystem/grpc-gateway/runtime"
  "google.golang.org/grpc"
	"golang.org/x/oauth2"
  gw "github.com/iamneal/book_parser/server/proto"
	mydrive "github.com/iamneal/book_parser/mydrive"
)



type MyHttpServer struct {
	HttpMux *http.ServeMux
	RpcMux *runtime.ServeMux
	HttpAddr string
	RpcAddr string
	state string
	Config *oauth2.Config
	cancel context.CancelFunc
	Clients map[*oauth2.Token]*http.Client
}

func NewMyHttpServer(rpcAddr, httpAddr string) (*MyHttpServer, error) {
	conf, err := mydrive.GetGoogleDriveConfig()
	if err != nil {
		return nil, err
	}

	mhs := &MyHttpServer{
		HttpAddr: httpAddr,
		RpcAddr: rpcAddr,
		Config: conf,
		Clients: make(map[*oauth2.Token]*http.Client),
	}

  ctx := context.Background()
  ctx, cancel := context.WithCancel(ctx)

	mhs.cancel = cancel
	mhs.state = "random"

  mux := runtime.NewServeMux()
  opts := []grpc.DialOption{grpc.WithInsecure()}
  err = gw.RegisterBookParserHandlerFromEndpoint(ctx, mux, mhs.RpcAddr, opts)
  if err != nil {
    return nil, err
  }

	httpMux := http.NewServeMux()

	mhs.HttpMux = httpMux
	mhs.RpcMux = mux

	mhs.HttpMux.Handle("/", http.FileServer(http.Dir(HTML_LOCATION)))
	mhs.HttpMux.Handle("/api/",http.Handler(mhs.RpcMux))
	mhs.HttpMux.HandleFunc("/auth", mhs.HandleAuth)
	mhs.HttpMux.HandleFunc("/login", mhs.HandleLogin)

	return mhs, nil
}

func (mhs *MyHttpServer) Run() error {
  return http.ListenAndServe(mhs.HttpAddr, mhs.HttpMux)
}

func (mhs *MyHttpServer) Shutdown() error {
	return nil
}

func (mhs *MyHttpServer) CreateTokenCookie(tok *oauth2.Token) (*http.Cookie, error) {
	expires := tok.Expiry.Add(time.Second * -1)
	maxAge := int(expires.Sub(time.Now()).Seconds())
	valBytes, err := json.Marshal(tok)
	if err != nil {
		return nil, err
	}
	val := string(valBytes[:])
	fmt.Printf("the cookie value: %+v\n", val)
	fmt.Printf("the max age: %+v\n", maxAge)
	fmt.Printf("expires: %+v\n", expires)
	encoded := base64.StdEncoding.EncodeToString(valBytes)

	return &http.Cookie{
		Name: COOKIE_NAME,
		Path: "/",
		Domain: "localhost",
		Value: encoded,
		Expires: expires,
		MaxAge: maxAge,
		Raw: val,
	}, nil
}

func (mhs *MyHttpServer) HandleAuth(res http.ResponseWriter, req *http.Request) {
	state := req.FormValue("state")
	if mhs.state != state {
		fmt.Printf("INVALID STATE %+v\n", state)
		http.Redirect(res, req, "/", http.StatusUnauthorized)
		return
	}
	code := req.FormValue("code")
	tok, err := mhs.Config.Exchange(context.Background(), code)

	if err != nil {
		fmt.Printf("could not aquire token\n")
		http.Redirect(res, req, "/", http.StatusUnauthorized)
	}
	mhs.Clients[tok] = mhs.Config.Client(context.Background(), tok)
	// set the jwt with the token
	cookie, err := mhs.CreateTokenCookie(tok)
	if err != nil {
		fmt.Printf("error creating cookie: %+v\n", err)
		http.Redirect(res, req, "/", http.StatusUnauthorized)
	}
	http.SetCookie(res, cookie)
	resp, err := mhs.Clients[tok].Get("https://www.googleapis.com/userinfo/v2/me")
	if err != nil {
		fmt.Printf("could not get profile, %s\n", err)
		http.Redirect(res, req, "/", http.StatusUnauthorized)
	}
	fmt.Printf("THE RESPONSE:: %#v\n%+v\n", resp, resp)
	credBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading the profile response body %s\n", err)
		http.Redirect(res, req, "/", http.StatusUnauthorized)
	}
	fmt.Printf("\ncreds? %s\n", string(credBytes[:]))
	http.Redirect(res, req, fmt.Sprintf("/?%s=%s",COOKIE_NAME, cookie.Raw), http.StatusFound)
}

func (mhs *MyHttpServer) HandleLogin(res http.ResponseWriter, req *http.Request) {
	url := mhs.Config.AuthCodeURL(mhs.state)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}
