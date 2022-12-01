package main

// import (
// 	"context"
// 	"fmt"
// 	pb "icarus/message/grpc/icarus"
// 	"log"
// 	"net"

// 	"google.golang.org/grpc"
// )

// const (
// 	port = ":50051"
// )

// // server is used to implement helloworld.GreeterServer.
// type server struct{}

// // SayHello implements helloworld.GreeterServer
// func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
// 	log.Printf("Received: %v", in.Name)
// 	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
// }

// func main() {
// 	lis, err := net.Listen("tcp", port)
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}
// 	s := grpc.NewServer()
// 	pb.RegisterGreeterServer(s, &server{})
// 	fmt.Println("启动服务", port)
// 	if err := s.Serve(lis); err != nil {
// 		log.Fatalf("failed to serve: %v", err)
// 	}
// }
