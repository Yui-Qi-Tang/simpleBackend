package clientapi

import (
	"context"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"

	authenticationPb "simpleBackend/ann-service/pianogame/protocol-buffer/authentication"

	"github.com/gin-gonic/gin"
)

// Login a client api to call gRPC service
func Login(c *gin.Context) {
	const (
		address  = "localhost:9001" // gRPC server that is set in ann-servie/main.go now
		account  = "tester"
		password = "testpassword"
	)
	var user struct {
		Account  string `form:"account" json:"account" xml:"account"  binding:"required"`
		Password string `form:"password" json:"password" xml:"password"  binding:"required"`
	}

	if c.ShouldBind(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "data marshal FATAL error",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	/* Set connection */
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	/* set client */
	grpcClient := authenticationPb.NewAuthenticationGreeterClient(conn)

	r, err := grpcClient.Login(ctx, &authenticationPb.LoginRequest{Account: user.Account, Password: user.Password})
	if err != nil {
		log.Printf("could not greet: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":   r.Msg,
		"token": r.Token,
	})
}
