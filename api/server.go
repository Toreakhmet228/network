package main

import (
	"context"
	"log"
	"net"
	"net/http"

	pb "chat/api/github.com/yourname/grpc-gateway-example/hello"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedHelloServiceServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Received request: %s", req.Name)
	return &pb.HelloResponse{Message: "Hello, " + req.Name + "!"}, nil
}

func main() {
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		s := grpc.NewServer()
		pb.RegisterHelloServiceServer(s, &server{})

		log.Println("gRPC server listening on :50051")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	mux := runtime.NewServeMux()
	ctx := context.Background()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("Failed to start HTTP Gateway: %v", err)
	}

	log.Println("HTTP-gRPC Gateway listening on :8080")
	http.ListenAndServe(":8080", mux)
}
