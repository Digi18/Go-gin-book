package main

import (
	"fmt"
	"go-gin-book/controllers"
	connection "go-gin-book/database"
	"go-gin-book/middleware"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	loadDatabase()
	runServer()
}

func loadEnv() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully loaded env file...")
	}
}

func loadDatabase() {
	connection.Connect()
	// connection.db.AutoMigrate
}

func runServer() {

	r := gin.Default()

	public := r.Group("/api")

	public.GET("/healthCheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello world"})
	})

	r.GET("/books", middleware.Authorize, controllers.FetchBooks)
	r.POST("/createBook", controllers.AddBook)
	r.DELETE("/deleteBook/:id", controllers.DeleteBook)
	r.GET("/getSpecificBook/:id", controllers.GetSpecificBook)
	r.PATCH("/updateBookTitle/:id", controllers.UpdateBookTitle)
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// r.Run() //By default it runs on 8080 port
	r.Run(":2000")
}
