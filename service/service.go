package service

import (
	"context"
	"io"
	"sync/atomic"
	pb "testing-backend/generated"
)

type Impl struct {
	pb.UnimplementedTestServiceServer
}

func New() pb.TestServiceServer {
	return &Impl{}
}

func (s *Impl) UnaryCall(_ context.Context, request *pb.UnaryCallRequest) (*pb.UnaryCallResponse, error) {
	enum := pb.TestEnum_FIRST
	innerMessage := &pb.TestInnerMessage{
		DoubleValue: request.DoubleValue,
		FloatValue:  request.FloatValue,
		IntValue:    request.IntValue,
		BoolValue:   request.BoolValue,
		StringValue: request.StringValue,
		BytesValue:  request.BytesValue,
	}
	message := &pb.TestMessage{
		EnumValue:    enum,
		MessageValue: innerMessage,
	}
	return &pb.UnaryCallResponse{Message: message}, nil
}

func (s *Impl) ClientStreaming(server pb.TestService_ClientStreamingServer) error {
	var accumulator int64

	done := make(chan bool)
	go func() {
		defer close(done)

		for {
			req, err := server.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				return
			}

			atomic.AddInt64(&accumulator, req.Value)
		}
	}()

	<-done

	return server.SendAndClose(&pb.ClientStreamingResponse{Sum: accumulator})
}

func (s *Impl) ServerStreaming(request *pb.ServerStreamingRequest, server pb.TestService_ServerStreamingServer) error {
	done := make(chan bool)
	go func() {
		defer close(done)

		word := request.GetWord()

		for i := 0; i < len(word); i++ {
			if err := server.Send(&pb.ServerStreamingResponse{Letter: string(word[i])}); err != nil {
				break
			}
		}
	}()

	<-done

	return nil
}

func (s *Impl) BidiStreaming(server pb.TestService_BidiStreamingServer) error {
	done := make(chan bool)
	go func() {
		defer close(done)

		var accumulator int64

		for {
			req, err := server.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				return
			}

			accumulator += req.Value

			err = server.Send(&pb.BidiStreamingResponse{Accumulator: accumulator})
			if err != nil {
				return
			}
		}
	}()

	<-done

	return nil
}
