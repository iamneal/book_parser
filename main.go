package main

import (
	"github.com/iamneal/book_parser/server"
)

func main() {
	rpcServer := server.NewRpcDriveServer(nil)
	go func() {
		err := rpcServer.RunRpcServer(":9090")
		if err != nil {
			panic(err)
		}
	}()
	webServer, err := server.NewMyHttpServer("locallhost:9090", "localhost:8080")
	if err != nil {
		panic(err)
	}
	err = webServer.Run()
	if err != nil {
		panic(err)
	}
}
