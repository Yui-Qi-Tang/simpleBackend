package pianogame

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
)

func strConcate(s ...string) string {
	result := ""
	for _, v := range s {
		result += v
	}
	return result
}
func strConcateF(str string, user string, pwd string, db string) string {
	return fmt.Sprintf(str, user, pwd, db)
}

func getURLInfo(c *gin.Context) *url.URL {
	return location.Get(c)
}

// BindIPPort a tool to bind
func BindIPPort(IP string, PORT int) string {
	return fmt.Sprintf("%s:%d", IP, PORT)
}

// GenerateToken generate JWT token
func GenerateToken(username, password string) (string, error) {
	tokenExpireTimestamp := time.Now().Add(3 * time.Hour).Unix()

	claims := jwtClaim{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: tokenExpireTimestamp,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(SysConfig.JwtSec))

	return token, err
}

// IsJwtValid JWT Validation
func IsJwtValid(tokenString string) bool {

	claims := jwtClaim{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SysConfig.JwtSec), nil
		},
	)

	if err != nil || !token.Valid {
		return false
	}
	return true
}

// IsJwtExpired check if JWT expired
func IsJwtExpired(tokenString string) bool {

	claims := jwtClaim{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SysConfig.JwtSec), nil
		},
	)

	if err != nil || !token.Valid {
		panic("IsJwtExpires => parse or token is invalid")
	}

	if !(time.Now().Unix() > claims.ExpiresAt) {
		return false
	}
	return true
}

func readFile(filePath string) []byte {
	fileBytes, err := ioutil.ReadFile(filePath) // open file and read
	errorCheck(err, "readFile Error")
	return fileBytes
}

func errorCheck(e error, msg ...string) {
	// TO-DO: better to logging error
	if e != nil {
		errorMsg := ""
		for _, v := range msg {
			errorMsg += v
		}
		log.Panicf("%s => %v", errorMsg, e)
	} // fi
}

func shutDownGraceful(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Printf("Server %s graceful exiting...", server.Addr)
}

func waitQuitSignal(hint string) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println(hint)
}

// StartServers HINT
func StartServers(handler *gin.Engine) {
	servers := make([]*http.Server, len(SysConfig.Ports))
	for i, v := range SysConfig.Ports {
		servers[i] = &http.Server{
			Addr:    BindIPPort(SysConfig.IP, v),
			Handler: handler,
		}
		log.Println("Start server", servers[i].Addr)
		go runserverTLS(servers[i], SysConfig.Ssl.Cert, SysConfig.Ssl.Key)
	}

	waitQuitSignal("Receive Quit server Signal") // block until receive quit signal from system

	// stop servers
	for _, v := range servers {
		shutDownGraceful(v) // terminate each server
	} // for
}

func runserverTLS(server *http.Server, cert string, key string) {
	// Start HTTPS server by net/http
	if err := server.ListenAndServeTLS(cert, key); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
