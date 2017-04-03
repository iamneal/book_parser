package server

import (
	"fmt"
	"time"
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


var (
	SECRET_LOCATION_NAME = "SECRET_LOCATION"
)

type MyHttpServer struct {
	HttpMux *http.ServeMux
	RpcMux *runtime.ServeMux
	HttpAddr string
	RpcAddr string
	state string
	Config *oauth2.Config
	cancel context.CancelFunc
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
		fmt.Printf("could not aquire token")
		http.Redirect(res, req, "/", http.StatusUnauthorized)
	}
	// set the jwt with the token
	cookie, err := mhs.CreateTokenCookie(tok)
	if err != nil {
		fmt.Printf("error creating cookie: %+v", err)
		http.Redirect(res, req, "/", http.StatusUnauthorized)
	}
	http.SetCookie(res, cookie)
	http.Redirect(res, req, fmt.Sprintf("/?%s=%s",COOKIE_NAME, cookie.Raw), http.StatusFound)
}

func (mhs *MyHttpServer) HandleLogin(res http.ResponseWriter, req *http.Request) {
	url := mhs.Config.AuthCodeURL(mhs.state)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}
