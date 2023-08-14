package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	grpcserver "testing-backend/server"
	"testing-backend/service"
)

func main() {
	protoDir := "proto"
	generatedDir := "generated"

	if err := os.RemoveAll(generatedDir); err != nil {
		log.Fatal(err)
	}

	if err := os.Mkdir(generatedDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command(
		"protoc",
		fmt.Sprintf("--go_out=%s", generatedDir),
		fmt.Sprintf("--go-grpc_out=%s", generatedDir),
		fmt.Sprintf("--proto_path=%s", protoDir),
		fmt.Sprintf("%s/*.proto", protoDir),
	)

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	testService := service.New()

	hostname, port := "127.0.0.1", 8000

	server := grpcserver.New(hostname, uint32(port))

	server.BindService(testService)

	defer server.Stop()

	server.Launch()
}
