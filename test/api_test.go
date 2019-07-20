package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func main() {
	router := SetupRouter()
	router.Run()
}
func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"hello": "world",
		})
	})

	return router
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	req.Header.Add("Authorize") = "Bearer "
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
func TestHelloWorld(t *testing.T) {
	// Build our expected body
	body := gin.H{
		"hello": "world",
	}
	// Grab our router
	router := SetupRouter()
	// Perform a GET request with that handler.
	w := performRequest(router, "GET", "/")
	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)
	// Convert the JSON response to a map
	var response map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	// Grab the value & whether or not it exists
	value, exists := response["hello"]
	// Make some assertions on the correctness of the response.
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, body["hello"], value)
}

//https://medium.com/@craigchilds94/testing-gin-json-responses-1f258ce3b0b1
//https://golangcode.com/mysql-database-insert-get-last-insert-id/
//https://tutorialedge.net/golang/authenticating-golang-rest-api-with-jwts/
//https://medium.com/@noomerzx/the-easiest-way-to-benchmarking-your-server-with-bombardier-87cc9c6b8d6f
//https://github.com/tsenart/vegeta
//https://github.com/codesenberg/bombardier
//https://stackoverflow.com/questions/52456506/getting-too-many-open-files-during-load-test-with-gin-gonic
//https://semaphoreci.com/community/tutorials/test-driven-development-of-go-web-applications-with-gin
// func TestPingRoute(t *testing.T) {
// 	router := setupRouter()

// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("GET", "/ping", nil)
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, 200, w.Code)
// 	assert.Equal(t, "pong", w.Body.String())
// }
