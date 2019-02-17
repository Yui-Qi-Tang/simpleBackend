package clientapi

import (
	"context"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"

	usergrpc "simpleBackend/ann-service/pianogame/grpc"

	"github.com/gin-gonic/gin"
)

// ServiceLogin a client api to call gRPC service
func ServiceLogin(c *gin.Context) {
	const (
		address  = "localhost:9001" // gRPC server that is set in ann-servie/main.go now
		account  = "tester"
		password = "testpassword"
	)
	var user struct {
		Account  string `form:"name" json:"name" xml:"name"  binding:"required"`
		Password string `form:"pwd" json:"pwd" xml:"pwd"  binding:"required"`
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
	grpcClient := usergrpc.NewUserGreetingClient(conn)

	r, err := grpcClient.Login(ctx, &usergrpc.LoginRequest{Account: user.Account, Password: user.Password})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "grpcClient Fatal error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": r.Msg,
	})
}
