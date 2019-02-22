package pianogame

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"simpleBackend/ann-service/pianogame/datastructure"
	"simpleBackend/ann-service/pianogame/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
	"golang.org/x/crypto/bcrypt"

	gameMsg "simpleBackend/ann-service/pianogame/msg"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

type msg struct {
	Text     string
	PianoKey interface{}
	MyID     interface{}
	To       interface{}
	From     interface{}
}

var clients = make(map[*datastructure.WebSocketUser]bool)

// UserLogin api for user login request
func UserLogin(c *gin.Context) {
	var json Login
	// json decode
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// get db collection
	collection := Mongodb.Database("testing").Collection("user")
	// prepare filter to query
	filter := bson.M{
		"name":     json.User,
		"password": json.Password,
	}
	r := Login{}
	// query
	err := collection.FindOne(context.Background(), filter).Decode(&r)
	// log.Println(r)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "user not found"})
		return
	}
	if token, jwtErr := GenerateToken(json.User, json.Password); err != nil {
		log.Println("JWT error: ", jwtErr)
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error for auth, please retry"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "you are logged in",
			"token":  token,
		})
		return
	}
}

// UserRegister api for user register request
func UserRegister(c *gin.Context) {

	var registerData Login

	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	collection := gaCollection("testing", "user")
	filter := bson.M{
		"name": registerData.User,
	}
	r := Login{}
	err := collection.FindOne(context.Background(), filter).Decode(&r)
	if err != nil {
		newUserData := bson.M{
			"name":     registerData.User,
			"password": registerData.Password,
		}
		_, err := collection.InsertOne(context.Background(), newUserData)

		if err != nil {
			log.Println("Insert one failed: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "create new user Failed!"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "register ok!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "This account has been registed!"})
	}
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

// AddUser api for adding user into Mysql User table
func AddUser(c *gin.Context) {
	var userData struct {
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binding:"required"`
		// profile
		Dob    string   `json:"birthday"`
		Emails []string `json:"emails"`
		Name   string   `json:"name"`
	}
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var accountToDB sql.NullString
	var pwdToDB sql.NullString
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userData.Password), 8)
	errorCheck(accountToDB.Scan(userData.Account), "account for signup is Failed")
	errorCheck(pwdToDB.Scan(string(hashedPassword)), "password for signup is Failed")

	// create new user
	var user User
	newUser := MysqlDB.FirstOrCreate(
		&user,
		User{
			Account:  accountToDB,
			Password: pwdToDB,
		},
	)
	// check if acount exists
	if newUser.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account has be registed"})
		return
	}
	// bind profile to user
	user.Profile = UserProfile{
		Birthday: userData.Dob,
		Name:     userData.Name,
		Emails: func(emailData []string) []Email {
			var results []Email
			for _, v := range emailData {
				var email Email
				email.Email = v
				results = append(results, email)
			}
			return results
		}(userData.Emails),
	}
	MysqlDB.Save(&user)

	c.JSON(http.StatusCreated, gin.H{"msg": "success"})
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

func webPusher(p http.Pusher, resource string) {
	if p != nil {
		if err := p.Push(resource, nil); err != nil {
			log.Printf("Failed to push: %v", err)
		}
	}
}

// GameWebSocketHandler web socket handler for pian game
func GameWebSocketHandler(c *gin.Context) {
	// establish web socket
	websocketUpgrader := utils.NewWSocketUpgrader(1024, 1024)
	conn, err := websocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	errorCheck(err, "Web socket connection failed")

	// allocate user id
	newUserID := uuid.New().String()
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

func broadcastToClient(msg interface{}) {
	for client := range clients {
		client.SendMsg(msg)
	}
}

// close socket connection if error or leave
func connErrorOrExit(connErr error, u *datastructure.WebSocketUser) bool {
	if connErr == nil {
		return false
	}
	delete(clients, u) // remove socket user from shared model
	u.Close()          // close socket connectioin
	var userLeaveMsg interface{}
	if c := connErr.Error(); c == "client disconnects" {
		userLeaveMsg = &gameMsg.Exit{
			Base: gameMsg.Base{
				To:     "all",
				From:   u.GetID(),
				Action: gameMsg.ExitConn,
			},
			Text: strConcate("Good bye!", u.GetID()),
		}
	}
	broadcastToClient(userLeaveMsg)
	return true
}

func gameHandle(gamer *datastructure.WebSocketUser) {
	var recMsg interface{}
	for {
		err := gamer.GetConn().ReadJSON(&recMsg) // load data from client as inteface
		if connErrorOrExit(err, gamer) {
			return
		}
		recMsgMap := recMsg.(map[string]interface{}) // decode interface{}
		// check Action of receivced msg and do something
		switch act := recMsgMap["Action"]; act.(float64) { // decode map["Action"] interface as float64 and switch
		case gameMsg.SendPianoKey:
			var pianoMsg gameMsg.PianoKey
			mapstructure.Decode(recMsgMap, &pianoMsg)
			pianoMsg.From = gamer.GetID()
			if pianoMsg.To == nil {
				pianoMsg.To = "all"
			}
			broadcastToClient(pianoMsg)
		} // switch
	} // for
}
