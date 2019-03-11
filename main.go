package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/huyntsgs/go-rest-api/api"

	"github.com/huyntsgs/go-rest-api/store"
	"github.com/joho/godotenv"
)

func main() {
	// Load config data from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var mySqlDB = new(store.MySqlDB)
	mySqlDB.Connect()

	router := SetupRouter(mySqlDB)
	router.Run(":8081")
}

func SetupRouter(store *store.MySqlDB) *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	//router.Use(gin.Recovery())

	productHandler := api.NewProductHandle(store)
	userHandler := api.NewUserHandler(store)

	// Always has versioning for api
	// Default(initial) is v1
	v1 := router.Group("/api/v1")
	{
		// Get products with params limit offset likes /products?limit=50&offset=0
		v1.GET("/products", productHandler.GetProducts())
		v1.GET("/products/:productId", productHandler.GetSingleProduct())
	}

	v1.POST("/users/register", userHandler.Register())
	v1.POST("/users/login", userHandler.Login())

	// apis needs to authenticate
	v1.Use(api.Authorize())
	{
		v1.POST("/products", productHandler.CreateProduct())
		v1.PUT("/products", productHandler.UpdateProduct())
		v1.DELETE("/products/:productId", productHandler.DeleteProduct())
	}
	return router
}
