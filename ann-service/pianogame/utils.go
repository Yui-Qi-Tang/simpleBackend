package pianogame

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
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

	//log.Println(token.Valid, claims.Username, claims.Password, time.Now().Unix() > claims.ExpiresAt)
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
