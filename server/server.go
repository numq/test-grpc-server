package server

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "testing-backend/generated"
)

type Server struct {
	address string
	server  *grpc.Server
}

func New(hostname string, port uint32, opts ...grpc.ServerOption) Server {
	address := fmt.Sprintf("%s:%d", hostname, port)
	return Server{address, grpc.NewServer(opts...)}
}

func (s *Server) Launch() {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Server started: %s", s.address)
	}
	if err = s.server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) BindService(service pb.TestServiceServer) {
	pb.RegisterTestServiceServer(s.server, service)
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}
