package main

import (
	"github.com/iamneal/book_parser/server"
)

func main() {
	driveService, err := getDriveClient()
	if err != nil {
		panic(err)
	}
	rpcServer := server.NewRpcDriveServer(driveService)
	go func() {
		err := rpcServer.RunRpcServer(":9090")
		if err != nil {
			panic(err)
		}
	}()
	err = server.Run("localhost:9090", ":8080")
	if err != nil {
		panic(err)
	}
}
