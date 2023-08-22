package authcontroller

import (
	"belajar_api/initializers"
	"belajar_api/models"
	"belajar_api/utils"
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func Login(c *gin.Context) {
	var auth models.Auth

	if err := c.ShouldBindJSON(&auth); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if auth.Email == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Email tidak boleh kosong"})
		return
	}

	if auth.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Password tidak boleh kosong"})
		return
	}

	var emailValidation int
	if err := initializers.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", auth.Email).Scan(&emailValidation); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	if emailValidation == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Email tidak terdaftar"})
		return
	}

	var id string
	if err := initializers.DB.QueryRow("SELECT id FROM users WHERE email = ? AND password = ?", auth.Email, auth.Password).Scan(&id); err == sql.ErrNoRows {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Password salah"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	token := utils.GenerateJWT(id)

	if token == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	if _, err := initializers.DB.ExecContext(c, "INSERT INTO authentications (id, token) VALUES (?, ?)", id, token); err != nil {

		if sqlErr := err.(*mysql.MySQLError); sqlErr.Number == 1062 {
			if _, newErr := initializers.DB.ExecContext(c, "UPDATE authentications SET token = ? WHERE id = ?", token, id); newErr != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Internal Server Error"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Login berhasil",
				"data": gin.H{
					"accessToken": token,
				},
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"data": gin.H{
			"accessToken": token,
		},
	})
}

func Logout(c *gin.Context) {
	var token models.Token

	if err := c.ShouldBindJSON(&token); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var id string
	if err := initializers.DB.QueryRow("SELECT id FROM authentications WHERE token = ?", token.Token).Scan(&id); err == sql.ErrNoRows {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token is invalid"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	idFromToken, err := utils.CheckToken(token.Token)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	if id != idFromToken {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token is invalid"})
		return
	}

	if _, err := initializers.DB.ExecContext(c, "DELETE FROM authentications WHERE id = ?", id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout berhasil"})
}

func Check(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	// Memeriksa apakah header Authorization kosong atau tidak mengandung "Bearer"
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// Mendapatkan token dari header Authorization
	token := strings.TrimPrefix(authHeader, "Bearer ")

	id, err := utils.CheckToken(token)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": id})
}
