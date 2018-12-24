package pianogame

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserLogin test
func UserLogin(c *gin.Context) {
	gaCollection("a", "b")
	// log.Println(gaCollection) // This Mongodb is not set because, this variable does init at the same package
	c.JSON(http.StatusBadRequest, gin.H{"status": "Invalid username and password!"})

	/*
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
		}*/
}
