package main

import (
	"github.com/iamneal/book_parser/server"
)

func main() {
	rpcServer, err := server.NewRpcDriveServer()
	if err != nil {
		panic(err)
	}
	go func() {
		err := rpcServer.RunRpcServer("0.0.0.0:9090")
		if err != nil {
			panic(err)
		}
	}()
	webServer, err := server.NewMyHttpServer("0.0.0.0:9090", "0.0.0.0:8080")
	if err != nil {
		panic(err)
	}
	err = webServer.Run()
	if err != nil {
		panic(err)
	}
}
