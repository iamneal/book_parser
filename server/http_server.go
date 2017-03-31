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

  mux := runtime.NewServeMux()
  opts := []grpc.DialOption{grpc.WithInsecure()}
  err = gw.RegisterBookParserHandlerFromEndpoint(ctx, mux, mhs.RpcAddr, opts)
  if err != nil {
    return nil, err
  }

	httpMux := http.NewServeMux()

	mhs.HttpMux = httpMux
	mhs.RpcMux = mux

	mhs.HttpMux.Handle("/api/",http.Handler(mhs.RpcMux))
	mhs.HttpMux.HandleFunc("/auth", mhs.HandleAuth)
	mhs.HttpMux.HandleFunc("/login", mhs.HandleLogin)
	mhs.HttpMux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		fmt.Printf("headers %+v\n", req.Header)
		jwt := req.Header.Get(COOKIE_NAME)
		if jwt != "" {
			res.Write([]byte("already authenticated"))
		} else {
			res.Write([]byte(`
				<html><body>
					<button href="/login"> Login with google </button>
				</body></html>
			`))
		}
	})
	mhs.state = "random"

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
	fmt.Printf("the cookie value: %+v", val)
	encoded := base64.StdEncoding.EncodeToString(valBytes)

	return &http.Cookie{
		Name: COOKIE_NAME,
		Value: encoded,
		Expires: expires,
		MaxAge: maxAge,
		Secure: true,
		Raw: val,
		Unparsed: []string{COOKIE_NAME, val},
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
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Success!"))
}

func (mhs *MyHttpServer) HandleLogin(res http.ResponseWriter, req *http.Request) {
	url := mhs.Config.AuthCodeURL(mhs.state)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}
