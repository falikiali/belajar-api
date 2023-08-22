package cartcontroller

import (
	"belajar_api/initializers"
	"belajar_api/models"
	"belajar_api/utils"
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddProductToCart(c *gin.Context) {
	var req models.AddCartRequestBody
	var id string
	var idFromToken string
	var qty int

	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	idCart := uuid.New()

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	if err := initializers.DB.QueryRow("SELECT id FROM authentications WHERE token = ?", token).Scan(&id); err == sql.ErrNoRows {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token is invalid"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	if id, err := utils.CheckToken(token); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	} else {
		idFromToken = id
	}

	if id != idFromToken {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token is invalid"})
		return
	}

	if err := initializers.DB.QueryRow("SELECT qty FROM carts WHERE user = ? AND product = ?", idFromToken, req.Product).Scan(&qty); err == sql.ErrNoRows {
		if _, newErr := initializers.DB.ExecContext(c, "INSERT INTO carts (id, user, product, qty) VALUES (?, ?, ?, ?)", idCart, idFromToken, req.Product, req.Qty); newErr != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Berhasil menambahkan barang ke keranjang"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	if _, err := initializers.DB.ExecContext(c, "UPDATE carts SET qty = ? WHERE user = ? AND product = ?", req.Qty+qty, idFromToken, req.Product); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil menambahkan barang ke keranjang"})
}

func RemoveProductInCart(c *gin.Context) {
	var req models.RemoveCartRequestBody
	var id string
	var idFromToken string

	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	if err := initializers.DB.QueryRow("SELECT id FROM authentications WHERE token = ?", token).Scan(&id); err == sql.ErrNoRows {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token is invalid"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	if id, err := utils.CheckToken(token); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	} else {
		idFromToken = id
	}

	if id != idFromToken {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token is invalid"})
		return
	}

	if _, err := initializers.DB.ExecContext(c, "DELETE FROM carts WHERE id = ?", req.IdCart); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil menghapus barang dari keranjang"})
}

func ShowProductInCart(c *gin.Context) {
	var id string
	var idFromToken string
	var carts []models.Cart
	var totalData int

	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	perPage := 10
	page := 1
	pageStr := c.DefaultQuery("page", "1")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	if err := initializers.DB.QueryRow("SELECT id FROM authentications WHERE token = ?", token).Scan(&id); err == sql.ErrNoRows {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token is invalid"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	if id, err := utils.CheckToken(token); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	} else {
		idFromToken = id
	}

	if id != idFromToken {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token is invalid"})
		return
	}

	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}

	if page < 1 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Page not found"})
		return
	}

	offset := (page - 1) * perPage

	query := `
		SELECT carts.id, carts.product, carts.qty, name_product, qty_product, price_product * carts.qty FROM carts 
		INNER JOIN users ON carts.user = users.id 
		INNER JOIN products ON carts.product = products.id_product
		WHERE carts.user = ?
		LIMIT ? OFFSET ?
	`
	rows, err := initializers.DB.QueryContext(c, query, idFromToken, perPage, offset)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var cart models.Cart

		err = rows.Scan(&cart.IdCart, &cart.IdProduct, &cart.QtyPurchase, &cart.NamaProduct, &cart.QtyProduct, &cart.PricePurchase)
		if err != nil {
			fmt.Print(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}

		carts = append(carts, cart)
	}

	query = "SELECT COUNT(*) FROM carts WHERE carts.user = ?"
	if err := initializers.DB.QueryRowContext(c, query, idFromToken).Scan(&totalData); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"carts": carts,
		},
		"page":      page,
		"perPage":   perPage,
		"totalData": totalData,
		"totalPage": math.Ceil(float64(totalData) / float64(perPage)),
	})
}
