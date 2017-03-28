package server

import (
  "net/http"

  "golang.org/x/net/context"
  "github.com/grpc-ecosystem/grpc-gateway/runtime"
  "google.golang.org/grpc"
  gw "github.com/iamneal/book_parser/server/proto"
)


func Run(rpcLoc, httpLoc string) error {
  ctx := context.Background()
  ctx, cancel := context.WithCancel(ctx)
  defer cancel()

  mux := runtime.NewServeMux()
  opts := []grpc.DialOption{grpc.WithInsecure()}
  err := gw.RegisterBookParserHandlerFromEndpoint(ctx, mux, rpcLoc, opts)
  if err != nil {
    return err
  }
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("hello world"))
	})
  return http.ListenAndServe(httpLoc, mux)
}

