package usercontroller

import (
	"belajar_api/initializers"
	"belajar_api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	id := uuid.New()
	name := user.Name
	email := user.Email
	password := user.Password

	if name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Nama tidak boleh kosong"})
		return
	}

	if email == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Email tidak boleh kosong"})
		return
	}

	if password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Password tidak boleh kosong"})
		return
	}

	var count int
	if err := initializers.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	if count >= 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Email sudah digunakan"})
		return
	}

	if _, err := initializers.DB.ExecContext(c, "INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?)", id, name, email, password); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Daftar berhasil"})
}
