package pbserver

import (
	"log"
	"net"
	"simpleBackend/ann-service/pianogame"
	authenticationPb "simpleBackend/ann-service/pianogame/protocol-buffer/authentication"
	pianoPlayPb "simpleBackend/ann-service/pianogame/protocol-buffer/pianoplay"

	"google.golang.org/grpc"
)

// StartGrpcService start gRPC server
func StartGrpcService() {
	lis, err := net.Listen(pianogame.GrpcConfig.Protocol, pianogame.GrpcConfig.Server) // TODO: procotol of network that is set from config file
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	log.Printf("Start gRPC server at %v using %v", pianogame.GrpcConfig.Server, pianogame.GrpcConfig.Protocol)
	authenticationPb.RegisterAuthenticationGreeterServer(s, &authenticationService{})
	pianoPlayPb.RegisterPianoplayGreeterServer(s, &pianoplayService{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
