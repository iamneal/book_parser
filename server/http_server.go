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
)



type MyHttpServer struct {
	HttpMux *http.ServeMux
	RpcMux *runtime.ServeMux
	HttpAddr string
	RpcAddr string
	state string
	cancel context.CancelFunc
	Cache *OAuth2TokenCache
}

func NewMyHttpServer(rpcAddr, httpAddr string, cache *OAuth2TokenCache) (*MyHttpServer, error) {
	mhs := &MyHttpServer{
		HttpAddr: httpAddr,
		RpcAddr: rpcAddr,
		Cache: cache,
	}

  ctx := context.Background()
  ctx, cancel := context.WithCancel(ctx)

	mhs.cancel = cancel
	mhs.state = "random"

  mux := runtime.NewServeMux()
  opts := []grpc.DialOption{grpc.WithInsecure()}
  err := gw.RegisterBookParserHandlerFromEndpoint(ctx, mux, mhs.RpcAddr, opts)
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
	mhs.HttpMux.HandleFunc("/update/token", mhs.HandleUpdateToken)

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

func (mhs *MyHttpServer) HandleUpdateToken(res http.ResponseWriter, req *http.Request) {
	oldTok := req.FormValue("token")
	fmt.Printf("handleUpdateToken called with: %s", oldTok)

	newTok, err := mhs.Cache.UpdateOldToken(oldTok)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(newTok))
}

func (mhs *MyHttpServer) HandleAuth(res http.ResponseWriter, req *http.Request) {
	state := req.FormValue("state")
	if mhs.state != state {
		fmt.Printf("INVALID STATE %+v\n", state)
		http.Redirect(res, req, "/", http.StatusUnauthorized)
		return
	}
	code := req.FormValue("code")
	tok, err := mhs.Cache.NewToken(context.Background(), code)
	if err != nil {
		fmt.Printf("could not aquire token\n")
		http.Redirect(res, req, "/", http.StatusUnauthorized)
	}
	// set the jwt with the token
	cookie, err := mhs.CreateTokenCookie(tok)
	if err != nil {
		fmt.Printf("error creating cookie: %+v\n", err)
		http.Redirect(res, req, "/", http.StatusUnauthorized)
	}
	http.SetCookie(res, cookie)
	http.Redirect(res, req, fmt.Sprintf("/?%s=%s",COOKIE_NAME, cookie.Raw), http.StatusFound)
}

func (mhs *MyHttpServer) HandleLogin(res http.ResponseWriter, req *http.Request) {
	url := mhs.Cache.Config.AuthCodeURL(mhs.state)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}
