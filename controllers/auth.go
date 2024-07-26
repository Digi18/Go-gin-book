package controllers

import (
	connection "go-gin-book/database"
	"go-gin-book/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {

	var body struct {
		Email    string
		Password string
	}

	err := c.Bind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": "Invalid body"})
		return
	}

	//Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": "Failed to hash password"})
		return
	}

	saveData := models.User{Email: body.Email, Password: string(hash)}
	result := connection.DB.Save(&saveData)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "User registerd successfully"})

}

func Login(c *gin.Context) {

	var user models.User

	var body struct {
		Email    string
		Password string
	}

	c.Bind(&body)

	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"data": "Data validation failed"})
	// 	return
	// }

	if body.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"data": "Enter email"})
		return
	}

	if body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"data": "Enter password"})
		return
	}

	// fmt.Printf("Email is %s", body.Email)
	// fmt.Println("")
	connection.DB.Where("email=?", body.Email).First(&user)

	if user.ID == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"data": "User not found", "status": http.StatusNotFound})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"data": "Password is incorrect", "status": http.StatusBadRequest})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"nbf":    time.Now().Unix(),
		"exp":    time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"data": "Failed to generate token", "status ": http.StatusBadRequest})
		return
	}

	// c.JSON(http.StatusOK, gin.H{"token": tokenString, "status": http.StatusOK})

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

}
