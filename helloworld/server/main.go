package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/bignyap/helloworld/service"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloResponse{Message: "Hello " + in.GetName()}, nil
}

func loggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	fmt.Printf("Calling the Logging Interceptor: %v\n", info.FullMethod)
	fmt.Printf("User %v is requesting\n", ctx.Value("user"))
	resp, err := handler(ctx, req)
	if err != nil {
		return "", err
	}
	return resp, nil
}

type MyKeyType string

func authInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	fmt.Printf("Calling the Auth Interceptor: %v\n", info.FullMethod)
	var userKey MyKeyType = "user"
	resp, err := handler(context.WithValue(ctx, userKey, "Bignya"), req)
	if err != nil {
		return "", err
	}
	return resp, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			authInterceptor,
			loggingInterceptor,
		),
	)
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("Sever listening at %v", lis.Addr().String())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}

}
