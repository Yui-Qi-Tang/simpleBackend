package pianogame

import (
	"fmt"
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
	type claimData struct {
		username string
		password string
		jwt.StandardClaims
	}
	var jwtSecret = []byte("secret")
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := claimData{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}
