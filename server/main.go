package server;

import (
	pb "github.com/iamneal/book_parser/server/proto"
	mydrive "github.com/iamneal/book_parser/mydrive"
	drive "google.golang.org/api/drive/v3"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"time"
	"fmt"
	"net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"encoding/json"
)

type Server struct {
	Config *oauth2.Config
	Tokens map[string] *drive.Service
}

func NewRpcDriveServer() (*Server, error) {
	conf, err := mydrive.GetGoogleDriveConfig()
	if err != nil {
		return nil, err
	}
	return &Server{Config: conf}, nil
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

func (s *Server) FindDriveService(token string) (*drive.Service, error) {
	tok := new(oauth2.Token)
	err := json.Unmarshal([]byte(token), tok)
	if err != nil {
		return nil, err
	}
	// check if the token is expired
	now := time.Now()
	if tok.Expiry.After(now) {
		return nil, fmt.Errorf("token expired")
	}
	//check if there is a drive service in the Tokens map
	serv := s.Tokens[token]
	if serv != nil {
		return serv, nil
	}
	serv, err = mydrive.GetDriveClient(s.Config, tok)
	if err != nil {
		return nil, err
	}
	s.Tokens[token] = serv
	return serv, nil
}

func (s *Server) getTokenFromCtx(ctx context.Context) (string, error) {
	metadata, ok := metadata.FromContext(ctx)
	var token string

	if ok {
		fmt.Printf("metadata on the request: %+v", metadata)
		tokenArr := metadata["token"]
		if len(tokenArr) == 1 {
			token = tokenArr[0]
			if token != "" {
				return token, nil
			}
		} else if len(tokenArr) > 1 {
			return "", fmt.Errorf("token array? %+v", tokenArr)
		}
	}
	return "", fmt.Errorf("no token found on metadata")
}

func (s *Server) PullBook(ctx context.Context, file *pb.File) (*pb.Empty, error) {
	fmt.Printf("rpc Server recieved: %+v", file)
	token := file.Token

	if token == "" {
		tok, err := s.getTokenFromCtx(ctx)
		if err != nil {
			return nil, err
		}
		token = tok
	}
	serv, err := s.FindDriveService(token)
	if err != nil {
		return nil, fmt.Errorf("bad token")
	}
	fileService := drive.NewFilesService(serv)
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
