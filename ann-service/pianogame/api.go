package pianogame

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"simpleBackend/ann-service/pianogame/datastructure"
	"simpleBackend/ann-service/pianogame/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	gameMsg "simpleBackend/ann-service/pianogame/msg"

	"github.com/google/uuid"
)

type msg struct {
	Text     string
	PianoKey interface{}
	MyID     interface{}
	To       interface{}
	From     interface{}
}

var clients = make(map[*datastructure.WebSocketUser]bool) // bad idea to stroe all user in websocket service

// Available for test
func Available(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "yes")
}

// MysqlCheckTable api for checking mysql table exist or not
func MysqlCheckTable(c *gin.Context) {
	var t struct {
		Name string `form:"name" json:"name" xml:"name"  binding:"required"`
	}
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := MysqlDB.HasTable(t.Name); err == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Table name does not exist"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": " ok!"})
	}
}

// UserValidation verify username/password, if valid, response JWT, otherwise respone failed
func UserValidation(c *gin.Context) {
	var userData struct {
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var accountToDB sql.NullString
	errorCheck(accountToDB.Scan(userData.Account), "account for sign in is Failed")

	var user User
	queryResult := MysqlDB.Where(&User{Account: accountToDB}).First(&user)

	if queryResult.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到使用者"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(userData.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		c.JSON(http.StatusUnauthorized, gin.H{"error": "帳號或是密碼錯誤"})
		return
	}

	// JWT
	if tokenStr, err := GenerateMemberToken(user.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"msg": "Welcome", "token": tokenStr})
	}
}

// GetUserInfoByToken just put this as demo
func GetUserInfoByToken(c *gin.Context) {
	tokenStr := c.Param("token")
	if IsMemberJWTValid(tokenStr) != true {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Token is invalid",
		})
		return
	}

	if IsMemberJWTExpired(tokenStr) {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Token is expired",
		})
		return
	}

	// A simple way to get user info by orm relation
	// TO-DO: try to use Specify Foreign Key & Association Key in table definition
	// refer: http://doc.gorm.io/associations.html#has-one
	userID := getUserIDByToken(tokenStr)
	var user User
	var userProfile UserProfile
	MysqlDB.First(&user, userID)
	MysqlDB.Model(&user).Related(&userProfile)
	// log.Println(userID, userProfile)

	c.JSON(http.StatusOK, gin.H{
		"msg":      "Token is valid",
		"username": userProfile.Name,
	})
	// c.String(http.StatusOK, "Hello %s", name)

}

// GetUsers api for getting all of users info or single user by id
func GetUsers(c *gin.Context) {
	userIDStr := c.Query("id")
	if userIDStr == "" {
		// Get all users
		var users []User
		MysqlDB.Find(&users)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    users,
		})
	} else {
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"msg":     err.Error(),
			})
			return
		}

		var user User
		user.ID = 0
		MysqlDB.Find(&user, userID)
		if user.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"msg":     "Can not find user",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    user,
		})
	}

}

// DeleteUser api for delete user by id
func DeleteUser(c *gin.Context) {
	userIDStr := c.Query("id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"msg":     "未帶使用者ID",
		})
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"msg":     err.Error(),
		})
		return
	}

	var uUserID uint
	uUserID = uint(userID)
	var user User
	user.ID = uUserID
	MysqlDB.Delete(&user)
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"msg":     "Can not find user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// UploadFileSample just put this as demo
func UploadFileSample(c *gin.Context) {
	// single file
	savePlace := "/tmp"

	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusOK, "File upload error!!")
	}

	// log.Println(file.Filename)

	// Upload the file to specific dst.
	c.SaveUploadedFile(file, fmt.Sprintf("%s/%s", savePlace, file.Filename))

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

// DecodeJwt just put this as demo
func DecodeJwt(c *gin.Context) {
	var json authData
	// json decode
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !IsJwtValid(json.Token) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data":    "invalid",
		})
		return
	}

	if IsJwtExpired(json.Token) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data":    "expired",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    "ok!",
	})
}

// DecodeJwtFromCookie just put this as demo
func DecodeJwtFromCookie(c *gin.Context) {
	tokenStr, err := c.Cookie("token")
	errorCheck(err, "token is unset in cookie")
	if !IsJwtValid(tokenStr) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data":    "invalid",
		})
		return
	}

	if IsJwtExpired(tokenStr) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data":    "expired",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    "Some data",
	})
}

// GameWebSocketHandler web socket handler for pian game
func GameWebSocketHandler(c *gin.Context) {
	/*
	   TODO: bind token and UUID to userplaytable(Mysql)?
	*/
	// establish web socket
	websocketUpgrader := utils.NewWSocketUpgrader(1024, 1024)
	conn, err := websocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	errorCheck(err, "Web socket connection failed")

	// allocate user id
	newUserID := uuid.New().String() // use redis to store the uuid?if (token, uuid) does not in redis then set it, otherwise skip;(This data hust exist a time we set(3 hrs?or 24hrs?))
	newUser := &datastructure.WebSocketUser{}
	newUser.SetID(newUserID)
	newUser.SetWsConn(conn)
	// send first msg for new user
	newUser.SendMsg(&gameMsg.Welcome{ID: newUserID, Text: strConcate("Welcom!", "userName!", "Your game ID is: ", newUserID)})
	// the behavior of handler: receive 'close' then board msg to all client
	newUser.GetConn().SetCloseHandler(
		func(code int, text string) error {
			// loggin code?
			log.Println(newUser.GetID(), "in close handler, just mark!")
			return errors.New("client disconnects")
		},
	)

	clients[newUser] = true // store new user for system

	go gameHandle(newUser)
}
