package server;

import (
	pb "github.com/iamneal/book_parser/server/proto"
	drive "google.golang.org/api/drive/v3"
	"golang.org/x/net/context"
	"fmt"
	"net"
	"google.golang.org/grpc"
)

type Server struct {
	driveService *drive.Service
}

func NewRpcDriveServer(d *drive.Service) *Server {
	return &Server{driveService: d}
}

func (s *Server) RunRpcServer(conn string) error {
	listener, err := net.Listen("tcp", conn)
	if err != nil {
		return err
	}
	server := grpc.NewServer()

	pb.RegisterBookParserServer(server, s)
	return server.Serve(listener)
}

func (s *Server) PullBook(ctx context.Context, file *pb.File) (*pb.Empty, error) {
	fmt.Printf("rpc Server recieved: %+v", file)

	fileService := drive.NewFilesService(s.driveService)
	list, err := fileService.List().Corpora("user").Context(context.Background()).
		Spaces("drive").Do()
	if err != nil {
		return nil, err
	}

	fmt.Printf("listed files: %+v\n", list)
	for _, f := range list.Files {
		fmt.Printf("fileId: %s\n Name: %s\n\n" ,f.Id, f.Name)
	}
	return &pb.Empty{}, nil
}
