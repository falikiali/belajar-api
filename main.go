package main

import (
	"belajar_api/controllers/authcontroller"
	"belajar_api/controllers/cartcontroller"
	"belajar_api/controllers/productcontroller"
	"belajar_api/controllers/usercontroller"
	"belajar_api/initializers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	initializers.DatabaseConnection()
	defer initializers.DB.Close()

	r.POST("/api/register", usercontroller.RegisterUser)
	r.POST("/api/login", authcontroller.Login)
	r.DELETE("/api/logout", authcontroller.Logout)

	r.GET("/api/product", productcontroller.ShowAllProducts)
	r.GET("/api/category/:idCategory/product", productcontroller.ShowProductByCategory)

	r.POST("/api/cart", cartcontroller.AddProductToCart)
	r.GET("/api/cart", cartcontroller.ShowProductInCart)
	r.DELETE("/api/cart", cartcontroller.RemoveProductInCart)

	r.GET("/api/testing", authcontroller.Check)

	r.Run()
}
