package main

import (
	"github.com/iamneal/book_parser/server"
)

func main() {
	cache, err := server.NewOAuth2TokenCache()
	if err != nil {
		panic(err)
	}
	rpcServer, err := server.NewRpcDriveServer(cache)
	if err != nil {
		panic(err)
	}
	go func() {
		err := rpcServer.RunRpcServer("0.0.0.0:9090")
		if err != nil {
			panic(err)
		}
	}()
	webServer, err := server.NewMyHttpServer("0.0.0.0:9090", "0.0.0.0:8080", cache)
	if err != nil {
		panic(err)
	}
	err = webServer.Run()
	if err != nil {
		panic(err)
	}
}
func thing() string{
	var ohello string
	ohello = "what the"
	return ohello
}
