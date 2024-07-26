package middleware

import (
	"fmt"
	connection "go-gin-book/database"
	"go-gin-book/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authorize(c *gin.Context) {

	//Get cookie off from request
	tokenStr, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"data": "Token is missing I", "status ": http.StatusUnauthorized})
		return
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"data": err.Error(), "status ": http.StatusUnauthorized})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		//Check the exp time
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"data": "Token has been expired", "status ": http.StatusUnauthorized})
			return
		}

		//Find the user with userId
		var user models.User
		connection.DB.Where("id=?", claims["userId"]).First(&user)

		if user.ID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"data": "Token is missing IV", "status ": http.StatusUnauthorized})
			return
		}

		//Attach to request
		c.Set("user", user)

	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"data": "Token is missing V", "status ": http.StatusUnauthorized})
		return
	}

	c.Next()
}
