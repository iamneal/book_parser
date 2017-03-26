package main

import (
	"fmt"
	"golang.org/x/net/context"
	drive "google.golang.org/api/drive/v3"
)

func main() {
	driveService, err := getDriveClient()
	if err != nil {
		panic(err)
	}
	fmt.Println("got drive service")

	fileService := drive.NewFilesService(driveService)
	if err != nil {
		panic(err)
	}
	list, err := fileService.List().Corpora("user").Context(context.Background()).
		Spaces("drive").Do()
	fmt.Printf("listed files: %+v\n", list)
	for _, f := range list.Files {
		fmt.Printf("fileId: %s\n Name: %s\n\n" ,f.Id, f.Name)
	}
}
