package pianogame

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
)

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

// InsertUserToMysql api for adding user into Mysql User table
func InsertUserToMysql(c *gin.Context) {
	timeFormat := "2006-01-02"
	var t struct {
		Name string `json:"name" binding:"required"`
		Age  int    `json:"age" binding:"required"`
		Dob  string `json:"dob" binding:"required"`
	}
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(t)
	dobTime, dobErr := time.Parse(timeFormat, t.Dob)
	if dobErr != nil {
		log.Println("Parse time error!", dobErr.Error())
	}
	user := User{
		Name:     t.Name,
		Age:      t.Age,
		Birthday: dobTime,
	}
	if f := MysqlDB.NewRecord(user); f == true {
		MysqlDB.Create(&user) // bind base Model data into 'user' and create
		c.JSON(http.StatusCreated, gin.H{"status": "user created"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": " create failed!!"})
	}
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
