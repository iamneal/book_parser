package main

import (
	"fmt"
	"golang.org/x/net/context"
	"github.com/iamneal/book_parser/server"
	drive "google.golang.org/api/drive/v3"
)

func main() {
	driveService, err := getDriveClient()
	if err != nil {
		panic(err)
	}
	rpcServer := NewRpcDriveServer(driveService)
	go func() {
		err := rpcServer.RunRpcServer(":9090")
		if err != nil {
			panic(err)
		}
	}()
	err := server.Run("localhost:9090", ":8080")
	if err != nil {
		panic(err)
}
