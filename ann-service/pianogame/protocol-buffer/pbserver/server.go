package pbserver

import (
	"log"
	"net"
	authenticationPb "simpleBackend/ann-service/pianogame/protocol-buffer/authentication"
	pianoPlayPb "simpleBackend/ann-service/pianogame/protocol-buffer/pianoplay"

	"google.golang.org/grpc"
)

// StartAuthenticationService start gRPC server
func StartAuthenticationService() {
	const port = ":9001" // TODO: config in file

	lis, err := net.Listen("tcp", port) // TODO: procotol of network that is set from config file
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	log.Printf("Start gRPC server at %v", port)
	authenticationPb.RegisterAuthenticationGreeterServer(s, &authenticationService{})
	pianoPlayPb.RegisterPianoplayGreeterServer(s, &pianoplayService{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
