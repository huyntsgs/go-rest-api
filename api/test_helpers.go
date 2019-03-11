package api

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/huyntsgs/go-rest-api/store"
	"github.com/huyntsgs/go-rest-api/utils"
)

func GinAbort(c *gin.Context, status, errCode int, msg string) {
	c.JSON(status, gin.H{
		"status": "fail", "errCode": errCode, "msg": msg,
	})
	c.Abort()
}

func SetupUserRouter(store store.UserStore) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	userHandler := NewUserHandler(store)

	// Always has versioning for api
	// Default(initial) is v1
	v1 := router.Group("/api/v1")
	{
		v1.POST("/users/register", userHandler.Register())
		v1.POST("/users/login", userHandler.Login())
	}

	return router
}
func SetupProductRouter(store store.ProductStore) *gin.Engine {
	router := gin.New()
	productHandler := NewProductHandle(store)
	// Always has versioning for api
	// Default(initial) is v1
	v1 := router.Group("/api/v1")
	{
		// Get products with params limit offset likes /products?limit=50&offset=0
		v1.GET("/products", productHandler.GetProducts())
		v1.GET("/products/:productId", productHandler.GetSingleProduct())
	}

	v1.Use(Authorize())
	{
		v1.POST("/products", productHandler.CreateProduct())
		v1.PUT("/products", productHandler.UpdateProduct())
		v1.DELETE("/products/:productId", productHandler.DeleteProduct())
	}
	return router
}

func MakeRequest(method, apiUrl string, body []byte, setAuthHeader bool) (*http.Request, error) {
	req, err := http.NewRequest(method, apiUrl, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error ", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if setAuthHeader {
		claims := map[string]interface{}{
			"userId": 1,
			"name":   "huyntsgs",
		}
		token, err := utils.GenToken(claims, []byte(os.Getenv("TOKEN_KEY")), 1000)
		if err != nil {
			fmt.Println("Error ", err)
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req, nil
}
