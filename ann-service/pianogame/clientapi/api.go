package clientapi

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"

	"simpleBackend/ann-service/pianogame"
	authenticationPb "simpleBackend/ann-service/pianogame/protocol-buffer/authentication"

	// "github.com/gin-contrib/location" // bad design

	"github.com/gin-gonic/gin"
)

// Login a client api to call gRPC service
func Login(c *gin.Context) {
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
	conn, err := grpc.Dial(pianogame.GrpcConfig.Server, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	/* set client */
	grpcClient := authenticationPb.NewAuthenticationGreeterClient(conn)

	r, err := grpcClient.Login(ctx, &authenticationPb.LoginRequest{Account: user.Account, Password: user.Password})
	if err != nil {
		// log.Printf("could not greet: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	host, _, _ := net.SplitHostPort(c.Request.Host)

	c.SetCookie("token", r.Token, 86400, "/", host, true, true)

	c.JSON(http.StatusOK, gin.H{
		"msg": r.Msg,
	})
} // end of Login()

// Signup user register account
func Signup(c *gin.Context) {
	// Data check
	var registerData struct {
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binding:"required"`
		// profile
		Dob    string   `json:"birthday"`
		Emails []string `json:"emails"`
		Name   string   `json:"name"`
	}

	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// send data to backend server via gRPC
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.Dial(pianogame.GrpcConfig.Server, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	grpcClient := authenticationPb.NewAuthenticationGreeterClient(conn)

	r, err := grpcClient.Register(ctx, &authenticationPb.RegisterRequest{
		Account:  registerData.Account,
		Password: registerData.Password,
		Dob:      registerData.Dob,
		Emails:   registerData.Emails,
		Name:     registerData.Name,
	})
	if err != nil {
		// log.Printf("could not greet: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"msg": r.Msg,
	})

} // end of Signup
