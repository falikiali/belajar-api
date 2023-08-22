package productcontroller

import (
	"belajar_api/initializers"
	"belajar_api/models"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ShowAllProducts(c *gin.Context) {
	perPage := 10
	page := 1
	pageStr := c.DefaultQuery("page", "1")

	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}

	if page < 1 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Page not found"})
		return
	}

	offset := (page - 1) * perPage

	var products []models.Product

	query := "SELECT id_product, name_product, qty_product, price_product, name_category FROM products INNER JOIN categories ON category_product = id_category LIMIT ? OFFSET ?"
	rows, err := initializers.DB.QueryContext(c, query, perPage, offset)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var p models.Product

		err = rows.Scan(&p.IdProduct, &p.NamaProduct, &p.QtyProduct, &p.PriceProduct, &p.CategoryProduct)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}

		products = append(products, p)
	}

	var totalData int
	query = "SELECT COUNT(*) FROM products"
	if err := initializers.DB.QueryRowContext(c, query).Scan(&totalData); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"products": products,
		},
		"page":      page,
		"perPage":   perPage,
		"totalData": totalData,
		"totalPage": math.Ceil(float64(totalData) / float64(perPage)),
	})
}

func ShowProductByCategory(c *gin.Context) {
	perPage := 10
	page := 1
	pageStr := c.DefaultQuery("page", "1")

	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}

	if page < 1 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Page not found"})
		return
	}

	offset := (page - 1) * perPage

	var products []models.Product
	idCategory := c.Param("idCategory")

	query := "SELECT id_product, name_product, qty_product, price_product, name_category FROM products INNER JOIN categories ON category_product = id_category WHERE category_product = ? LIMIT ? OFFSET ?"
	rows, err := initializers.DB.QueryContext(c, query, idCategory, perPage, offset)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var p models.Product

		err = rows.Scan(&p.IdProduct, &p.NamaProduct, &p.QtyProduct, &p.PriceProduct, &p.CategoryProduct)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}

		products = append(products, p)
	}

	var totalData int
	query = "SELECT COUNT(*) FROM products WHERE category_product = ?"
	if err := initializers.DB.QueryRowContext(c, query, idCategory).Scan(&totalData); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"products": products,
		},
		"page":      page,
		"perPage":   perPage,
		"totalData": totalData,
		"totalPage": math.Ceil(float64(totalData) / float64(perPage)),
	})
}
