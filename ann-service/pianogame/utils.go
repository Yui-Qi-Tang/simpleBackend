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
	yaml "gopkg.in/yaml.v2"
)

func dumpStructData(data interface{}) {
	log.Printf("%+v\n", data)
}

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
	token, err := tokenClaims.SignedString([]byte(authSettings.Secret.Jwt))

	return token, err
}

// GenerateMemberToken generate JWT token
func GenerateMemberToken(ID uint) (string, error) {
	tokenExpireTimestamp := time.Now().Add(3 * time.Hour).Unix()

	claims := AuthMemberClaim{
		ID,
		jwt.StandardClaims{
			ExpiresAt: tokenExpireTimestamp,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(authSettings.Secret.Jwt))

	return token, err
}

// IsJwtValid JWT Validation
func IsJwtValid(tokenString string) bool {

	claims := jwtClaim{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(authSettings.Secret.Jwt), nil
		},
	)

	if err != nil || !token.Valid {
		return false
	}
	return true
}

// IsMemberJWTValid JWT Validation
func IsMemberJWTValid(tokenString string) bool {

	claims := AuthMemberClaim{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(authSettings.Secret.Jwt), nil
		},
	)
	if err != nil || !token.Valid {
		return false
	}
	return true
}

// getUserIDByToken JWT Validation
func getUserIDByToken(tokenString string) uint {

	claims := AuthMemberClaim{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(authSettings.Secret.Jwt), nil
		},
	)
	if err != nil || !token.Valid {
		panic("Token is invalid")
	}
	return claims.ID
}

// IsMemberJWTExpired check if JWT expired
func IsMemberJWTExpired(tokenString string) bool {

	claims := AuthMemberClaim{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(authSettings.Secret.Jwt), nil
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

// IsJwtExpired check if JWT expired
func IsJwtExpired(tokenString string) bool {

	claims := jwtClaim{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(authSettings.Secret.Jwt), nil
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

func ShutDownGraceful(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Printf("Server %s graceful exiting...", server.Addr)
}

func WaitQuitSignal(hint string) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println(hint)
}

// StartServers Server network setting
func StartServers(handler *gin.Engine, config []host, meta serviceMeta) []*http.Server {

	servers := make([]*http.Server, len(config))
	log.Println(meta.Name)
	for i, v := range config {
		servers[i] = &http.Server{
			Addr:    BindIPPort(v.Name, v.Port),
			Handler: handler,
		}
		log.Println("Start server", servers[i].Addr)
		go runserverTLS(servers[i], Ssl.Path.Cert, Ssl.Path.Key) // TODO: Ssl as function parameter?
	}
	return servers
}

func runserverTLS(server *http.Server, cert string, key string) {
	// Start HTTPS server by net/http
	if err := server.ListenAndServeTLS(cert, key); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func loadYAMLConfig(configFilePath, errMsg, successMsg string, configStructure interface{}) {
	// This is blank
	bytesData := readFile(configFilePath)
	configUnmarshalError := yaml.Unmarshal(bytesData, configStructure)
	errorCheck(configUnmarshalError, errMsg)
	log.Println(configFilePath, successMsg)
}
