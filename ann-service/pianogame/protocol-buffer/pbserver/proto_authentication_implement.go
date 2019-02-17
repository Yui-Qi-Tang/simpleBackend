package pbserver

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"simpleBackend/ann-service/pianogame"
	authenticationPb "simpleBackend/ann-service/pianogame/protocol-buffer/authentication"

	"golang.org/x/crypto/bcrypt"
)

type authenticationService struct{}

// Login implements authenticationPb.AuthenticationGreeter
func (s *authenticationService) Login(ctx context.Context, in *authenticationPb.LoginRequest) (*authenticationPb.LoginResponse, error) {
	log.Printf("Received account/password: %v/%v", in.Account, in.Password)
	var accountToDB sql.NullString
	pianogame.ErrorCheck(accountToDB.Scan(in.Account), "account for sign in is Failed")
	var user pianogame.User
	queryResult := pianogame.MysqlDB.Where(&pianogame.User{Account: accountToDB}).First(&user)

	// check if account exists
	if queryResult.RowsAffected == 0 {
		return &authenticationPb.LoginResponse{Msg: "找不到使用者", Token: ""}, errors.New("can not found user")
	}

	// check if password valid
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(in.Password)); err != nil {
		// Msg just shows a fuzzy tip for user; in this phase it's 'Password'
		return &authenticationPb.LoginResponse{Msg: "帳號或是密碼錯誤", Token: ""}, errors.New("Password error")
	}

	// gen jwt token
	if tokenStr, err := pianogame.GenerateMemberToken(user.ID); err != nil {
		return &authenticationPb.LoginResponse{Msg: "System error, please contact system administrator", Token: ""}, err

	} else {
		return &authenticationPb.LoginResponse{Msg: "Welcome", Token: tokenStr}, nil
	}
}

// Logout implements authenticationPb.AuthenticationGreeter
func (s *authenticationService) Logout(ctx context.Context, in *authenticationPb.LogoutRequest) (*authenticationPb.LogoutResponse, error) {
	log.Printf("Received token: %v", in.Token)
	return &authenticationPb.LogoutResponse{Msg: "Goodbye"}, nil
}
