package controllers

import (
	"fmt"
	connection "go-gin-book/database"
	"go-gin-book/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FetchBooks(c *gin.Context) {

	var books []models.Book
	connection.DB.Find(&books)
	c.JSON(http.StatusOK, gin.H{"data": books})
}

func AddBook(c *gin.Context) {

	var input models.CreateBookInput
	// if err := c.ShouldBindJSON(&input); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": "input validation failed"})
		return
	}

	newBook := models.Book{Title: input.Title, Author: input.Author}
	connection.DB.Save(&newBook)
	fmt.Printf("New book %s %s was created successfully!\\n", newBook.Title, newBook.Author)
	c.JSON(http.StatusOK, gin.H{"data": newBook})
}

func DeleteBook(c *gin.Context) {

	bookId := c.Param("id")
	fmt.Println(bookId)
	var book models.Book

	findBook := connection.DB.Where("id=?", bookId).First(&book)

	if findBook.Error != nil {
		panic("Failed to retrieve book:" + findBook.Error.Error())
	}

	connection.DB.Delete(&book)

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

func GetSpecificBook(c *gin.Context) {

	bookId := c.Param("id")
	var newBook models.Book

	findBook := connection.DB.Where("id=?", bookId).First(&newBook)

	if findBook.Error != nil {
		// panic("Error in fetching book:" + findBook.Error.Error())
		c.JSON(http.StatusNotFound, gin.H{"data": findBook.Error.Error()})
		return
	}
	fmt.Printf(newBook.Title)
	c.JSON(http.StatusOK, gin.H{"data": newBook})
}

func UpdateBookTitle(c *gin.Context) {

	bookId := c.Param(("id"))
	var newBook models.Book
	var input models.UpdateBook

	//First find book
	findBook := connection.DB.Where("id=?", bookId).First(&newBook)

	if findBook.Error != nil {
		// panic("Error in fetching book" + findBook.Error.Error())
		c.JSON(http.StatusNotFound, gin.H{"data": findBook.Error.Error()})
		return
	}

	//Validate input
	err := c.ShouldBindJSON(&input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": err.Error()})
		return
	}

	connection.DB.Model(&newBook).Updates(&input)

	c.JSON(http.StatusOK, gin.H{"data": "Book has been updated successfully"})
}
