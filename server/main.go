package server;

import (
	pb "github.com/iamneal/book_parser/server/proto"
	drive "google.golang.org/api/drive/v3"
	"golang.org/x/net/context"
	"fmt"
	"net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/codes"
)

type Server struct {
	Cache *OAuth2TokenCache
}

type MsgWithToken interface {
	GetToken() string
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

func (s *Server) DebugPrintCache(ctx context.Context, em *pb.Empty) (*pb.DebugMsg, error) {
	fmt.Println("Debug Print")
	fmt.Printf("cache: %s", s.Cache)
	return &pb.DebugMsg{Msg: fmt.Sprintf("%s", s.Cache)}, nil
}

func (s *Server) ListBooks(ctx context.Context, file *pb.Token) (*pb.BookList, error) {
	userCache, err := s.getUserCache(ctx, file)
	if err != nil {
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

func (s *Server) PullBook(ctx context.Context, file *pb.File) (*pb.DebugMsg, error) {
	fmt.Printf("Pull book recieved: %#v\n", file)
	if file.Id == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "no file id found")
	}
	userCache, err := s.getUserCache(ctx, file)
	if err != nil {
		return nil, err
	}
	if userCache.User == nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "unable to determin user")
	}
	filename := userCache.User.Id + file.Id

	//resp, err := drive.NewFilesService(userCache.Drive).Get(file.Id).Download()

	return &pb.DebugMsg{ Msg: filename}, nil

}

func (s *Server) getToken(ctx context.Context, msg MsgWithToken) (string, error) {
	token := msg.GetToken()
	if token != "" {
		return token, nil
	}
	metadata, ok := metadata.FromContext(ctx)

	if ok {
		tokenArr := metadata[GRPC_GATEWAY_TOKEN]
		//fmt.Println("printing out keys on the metadata")
		//for key := range metadata {
		//	fmt.Printf("%s\n", key)
		//}
		if len(tokenArr) == 1 {
			token = tokenArr[0]
			if token != "" {
				return token, nil
			}
		} else if len(tokenArr) > 1 {
			return "", fmt.Errorf("token array? %+v", tokenArr)
		}
	}
	fmt.Printf("could not find token on ctx")
	return "", grpc.Errorf(codes.Unauthenticated, "no token found on metadata")
}

func (s *Server) getUserCache(ctx context.Context, msg MsgWithToken) (*UserCache, error) {
	token, err := s.getToken(ctx, msg)
	if err != nil {
		return nil, err
	}
	c, err := s.Cache.Get(token)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, err.Error())
	}
	return c, nil
}
