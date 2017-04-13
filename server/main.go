package server;

import (
	pb "github.com/iamneal/book_parser/server/proto"
	drive "google.golang.org/api/drive/v3"
	"golang.org/x/net/context"
	"fmt"
	"net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	Cache *OAuth2TokenCache
}

func NewRpcDriveServer(cache *OAuth2TokenCache) (*Server) {
	return &Server{Cache: cache}
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

func (s *Server) getTokenFromCtx(ctx context.Context) (string, error) {
	metadata, ok := metadata.FromContext(ctx)
	var token string

	if ok {
		fmt.Printf("metadata on the request: %+v\n", metadata)
		tokenArr := metadata[COOKIE_NAME]
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

func (s *Server) DebugPrintCache(ctx context.Context, em *pb.Empty) (*pb.DebugMsg, error) {
	fmt.Println("Debug Print")
	fmt.Printf("cache: %s", s.Cache)
	return &pb.DebugMsg{Msg: fmt.Sprintf("%s", s.Cache)}, nil
}

func (s *Server) ListBooks(ctx context.Context, file *pb.Token) (*pb.BookList, error) {
	fmt.Printf("rpc Server recieved: %+v", file)
	tokenStr := file.Token

	if tokenStr == "" {
		tok, err := s.getTokenFromCtx(ctx)
		if err != nil {
			fmt.Printf("error when getting token from context %s", err)
			return nil, err
		}
		tokenStr = tok
	}
	fmt.Printf("got token: %#v", tokenStr)
	userCache, err := s.Cache.Get(tokenStr)
	if err != nil {
		fmt.Printf("error when grabbing from cache: %s", err)
		return nil, err
	}
	fs := drive.NewFilesService(userCache.Drive)
	list, err := fs.List().Corpora("user").Context(context.Background()).
		Spaces("drive").Do()
	if err != nil {
		fmt.Printf("error when listing files from drive: %s", err)
		return nil, err
	}

	books := make([]*pb.Book, 0)
	fmt.Printf("listed files: %+v\n", list)
	for _, f := range list.Files {
		fmt.Printf("fileId: %s\n Name: %s\n\n" ,f.Id, f.Name)
		books = append(books, &pb.Book{Id: f.Id, Name: f.Name})
	}
	return &pb.BookList{Books: books}, nil
}
