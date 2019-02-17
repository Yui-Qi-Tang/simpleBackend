package pbserver

import (
	"context"
	"log"
	authenticationPb "simpleBackend/ann-service/pianogame/protocol-buffer/authentication"
)

type authenticationService struct{}

// Login implements authenticationPb.AuthenticationGreeter
func (s *authenticationService) Login(ctx context.Context, in *authenticationPb.LoginRequest) (*authenticationPb.LoginResponse, error) {
	log.Printf("Received account/password: %v/%v", in.Account, in.Password)
	return &authenticationPb.LoginResponse{Msg: "Hello, got login req and response token for you", Token: "token value"}, nil
}

// Logout implements authenticationPb.AuthenticationGreeter
func (s *authenticationService) Logout(ctx context.Context, in *authenticationPb.LogoutRequest) (*authenticationPb.LogoutResponse, error) {
	log.Printf("Received token: %v", in.Token)
	return &authenticationPb.LogoutResponse{Msg: "Goodbye"}, nil
}
