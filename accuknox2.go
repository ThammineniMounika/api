package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Note struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	Note   string `json:"note"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var db *gorm.DB

func main() {
	// Connect to MySQL database
	dsn := "root:Mouni@9797@tcp(127.0.0.1:3306)/balaji?charset=utf8&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&User{}, &Note{})

	r := gin.Default()

	r.POST("/signup", CreateUser)
	r.POST("/login", UserLogin)
	r.GET("/notes", ListNotes)
	r.POST("/notes", CreateNote)
	r.DELETE("/notes", DeleteNote)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func CreateUser(c *gin.Context) {
	var user User
	c.Bind(&user)
	if len(user.Name)<3{
		c.JSON(201,"invalid name")
	}
	if !checkEmailFormat(user.Email) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "please make sure email contains @ and .com"})
		return
	}

	/*if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request format"})
		return
	}

	result := db.Create(&user)
	if result.Error != nil {
		log.Println(result.Error)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create user"})
		return
	}*/

	c.JSON(http.StatusOK, "user created sussesfully")
}
func checkEmailFormat(email string) bool {
	// Implement your email format validation logic here
	// This is a simple example to check for '@' and '.com'
	return strings.Contains(email, "@") && strings.Contains(email, ".com")
}

func UserLogin(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request format"})
		return
	}

	result := db.Where("email = ? AND password = ?", user.Email, user.Password).First(&user)
	if result.Error != nil {
		log.Println(result.Error)
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sid": user.ID})
}

func ListNotes(c *gin.Context) {
	userID := c.Query("sid")
	if userID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "User ID missing"})
		return
	}

	var notes []Note
	result := db.Where("user_id = ?", userID).Find(&notes)
	if result.Error != nil {
		log.Println(result.Error)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve notes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notes": notes})
}

func CreateNote(c *gin.Context) {
	var note Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request format"})
		return
	}

	result := db.Create(&note)
	if result.Error != nil {
		log.Println(result.Error)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": note.ID})
}

func DeleteNote(c *gin.Context) {
	noteID := c.PostForm("id")
	if noteID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Note ID missing"})
		return
	}

	var note Note
	result := db.First(&note, noteID)
	if result.Error != nil {
		log.Println(result.Error)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to find note"})
		return
	}

	result = db.Delete(&note)
	if result.Error != nil {
		log.Println(result.Error)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}

