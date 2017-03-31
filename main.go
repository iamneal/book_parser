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
		err := rpcServer.RunRpcServer("localhost:9090")
		if err != nil {
			panic(err)
		}
	}()
	webServer, err := server.NewMyHttpServer("localhost:9090", "localhost:8080")
	if err != nil {
		panic(err)
	}
	err = webServer.Run()
	if err != nil {
		panic(err)
	}
}
