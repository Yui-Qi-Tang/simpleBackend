package pianogame

import (
	"net/http"
	"context"
	"log"
	
	
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
)

// UserLogin test
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
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Invalid username and password!"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	}
}

// UserRegister user register
func UserRegister(c *gin.Context) {
	var registerData Login
	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	collection := gaCollection("testing", "user")
	filter := bson.M{
		"name":     registerData.User,
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
			log.Fatalf("Insert one failed: %v", err)
		}
		c.JSON(http.StatusOK, gin.H{"status": "register ok!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "This account has been registed!"})
	}
}