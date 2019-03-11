package api

import (
	"bytes"
	"encoding/json"

	"fmt"
	"os"
	"time"

	"github.com/huyntsgs/go-rest-api/utils"
	"github.com/joho/godotenv"

	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/huyntsgs/go-rest-api/models"
)

var id int64 = 13
var products = []*models.Product{
	{1, "Iphone 3", "iphone3.jpg", 3, time.Now()},
	{2, "Iphone 4", "iphone4.jpg", 4, time.Now()},
	{3, "Iphone 5", "iphone5.jpg", 4, time.Now()},
	{4, "Iphone 5s", "iphone5s.jpg", 4, time.Now()},
	{5, "Iphone 6", "iphone6.jpg", 4, time.Now()},
	{6, "Iphone 6s", "iphone6s.jpg", 4, time.Now()},
	{7, "Iphone 7", "iphone7.jpg", 4, time.Now()},
	{8, "Iphone 7", "iphone7.jpg", 4, time.Now()},
	{9, "Iphone 7s", "iphone7s.jpg", 4, time.Now()},
	{10, "Iphone 8", "iphone8.jpg", 5, time.Now()},
	{11, "Iphone 8s", "iphone8s.jpg", 5, time.Now()},
	{12, "Iphone Xs", "iphonexs.jpg", 5, time.Now()},
	{13, "Iphone Xss", "iphonexss.jpg", 5, time.Now()},
}

type Mock struct{}

func (mock *Mock) GetProducts(limit, offset int, lastId int64) ([]*models.Product, error) {
	for i, p := range products {
		if p.ProductId > lastId {
			return products[i:], nil
		}
	}
	return nil, nil
}

func (mock *Mock) GetSingleProduct(productId int64) (*models.Product, error) {
	for _, p := range products {
		if p.ProductId == productId {
			return p, nil
		}
	}
	return nil, nil
}
func (mock *Mock) DeleteProduct(productId int64) (int64, error) {
	for i, p := range products {
		if p.ProductId == productId {
			products = append(products[0:i], products[i+1:]...)
			return productId, nil
		}
	}
	return 0, nil
}
func (mock *Mock) UpdateProduct(p *models.Product) (int64, error) {
	for i, p := range products {
		if p.ProductId == p.ProductId {
			products[i] = p
			return p.ProductId, nil
		}
	}
	return 0, nil
}
func (mock *Mock) CreateProduct(pr *models.Product) (int64, error) {

	for _, p := range products {
		if p.ProductName == pr.ProductName {
			return 0, nil
		}
	}
	id++
	pr.ProductId = id
	products = append(products, pr)
	return pr.ProductId, nil
}

type ErrorMock struct{}

func (mock *ErrorMock) GetProducts(limit, offset int, lastId int64) ([]*models.Product, error) {
	return nil, models.NewError("Error get products", ERR_INTERNAL)
}

func (mock *ErrorMock) GetSingleProduct(productId int64) (*models.Product, error) {
	return nil, models.NewError("Error get product", ERR_INTERNAL)
}
func (mock *ErrorMock) DeleteProduct(productId int64) (int64, error) {
	return 0, models.NewError("Error delete product", ERR_INTERNAL)
}
func (mock *ErrorMock) UpdateProduct(p *models.Product) (int64, error) {
	return 0, models.NewError("Error update product", ERR_INTERNAL)
}
func (mock *ErrorMock) CreateProduct(p *models.Product) (int64, error) {
	return 0, models.NewError("Error add new product", ERR_INTERNAL)
}

func MakeRequestWithInvalidToken(method, apiUrl, token string, body []byte, setAuthHeader bool) (*http.Request, error) {
	req, err := http.NewRequest(method, apiUrl, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error ", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if setAuthHeader {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req, nil
}

// test middleware authorize
func TestAuthInvalidToken(t *testing.T) {
	r := SetupProductRouter(&Mock{})

	req, err := MakeRequestWithInvalidToken("DELETE", "/api/v1/products/1", "invalidtokey", []byte{}, true)
	if err != nil {
		fmt.Println("Error ", err)
	}

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusUnauthorized)
}

func TestAuthTokenExpired(t *testing.T) {
	r := SetupProductRouter(&Mock{})

	claims := map[string]interface{}{
		"userId": 1,
		"name":   "huyntsgs",
	}
	token, err := utils.GenToken(claims, []byte(os.Getenv("TOKEN_KEY")), 1)
	if err != nil {
		fmt.Println("Error ", err)
		t.Fail()
	}

	time.Sleep(2 * time.Minute)
	req, err := MakeRequestWithInvalidToken("DELETE", "/api/v1/products/1", token, []byte{}, true)
	if err != nil {
		fmt.Println("Error ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusUnauthorized)
}

func TestAuthTokenWithWrongKey(t *testing.T) {
	r := SetupProductRouter(&Mock{})

	claims := map[string]interface{}{
		"userId": 1,
		"name":   "huyntsgs",
	}
	token, err := utils.GenToken(claims, []byte("wrongkey"), 10)
	if err != nil {
		fmt.Println("Error ", err)
		t.Fail()
	}

	req, err := MakeRequestWithInvalidToken("DELETE", "/api/v1/products/1", token, []byte{}, true)
	if err != nil {
		fmt.Println("Error ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusUnauthorized)
}

// test funcs for GetProducts
func TestGetProducts(t *testing.T) {

	r := SetupProductRouter(&Mock{})
	req, err := MakeRequest("GET", "/api/v1/products", []byte{}, false)
	if err != nil {
		fmt.Println("Error ", err)
	}

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
}
func TestGetProductsWithStoreError(t *testing.T) {

	r := SetupProductRouter(&ErrorMock{})
	req, err := MakeRequest("GET", "/api/v1/products", []byte{}, false)
	if err != nil {
		fmt.Println("Error ", err)
	}

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusInternalServerError)
}
func TestGetProductsWithLimit(t *testing.T) {

	r := SetupProductRouter(&Mock{})
	req, err := MakeRequest("GET", "/api/v1/products?limit=5&lastId=0", []byte{}, false)
	if err != nil {
		fmt.Println("Error ", err)
	}

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
}
func TestGetProductsWithInvalidLimit(t *testing.T) {

	r := SetupProductRouter(&Mock{})
	req, err := MakeRequest("GET", "/api/v1/products?limit=aaa&lastId=0", []byte{}, false)
	if err != nil {
		fmt.Println("Error ", err)
	}

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusBadRequest)
}
func TestGetProductsWithInvalidLastId(t *testing.T) {

	r := SetupProductRouter(&Mock{})
	req, err := MakeRequest("GET", "/api/v1/products?limit=5&lastId=fff", []byte{}, false)
	if err != nil {
		fmt.Println("Error ", err)
	}

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusBadRequest)
}

// test funcs for GetSingleProduct
func TestGetSingleProduct(t *testing.T) {
	r := SetupProductRouter(&Mock{})

	req, err := MakeRequest("GET", "/api/v1/products/5", []byte{}, false)
	if err != nil {
		fmt.Println("Error ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
}
func TestGetSingleProductWithStoreError(t *testing.T) {
	r := SetupProductRouter(&ErrorMock{})

	req, err := MakeRequest("GET", "/api/v1/products/5", []byte{}, false)
	if err != nil {
		fmt.Println("Error ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusInternalServerError)
}
func TestGetSingleProductNonExist(t *testing.T) {
	r := SetupProductRouter(&Mock{})

	req, err := MakeRequest("GET", "/api/v1/products/555", []byte{}, false)
	if err != nil {
		fmt.Println("Error ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusNoContent)
}
func TestGetSingleProductWithInvalidId(t *testing.T) {
	r := SetupProductRouter(&Mock{})

	req, err := MakeRequest("GET", "/api/v1/products/abc", []byte{}, false)
	if err != nil {
		fmt.Println("Error ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusBadRequest)
}

// test funcs for DeleteProduct
func TestDeleteProduct(t *testing.T) {
	r := SetupProductRouter(&Mock{})

	req, err := MakeRequest("DELETE", "/api/v1/products/1", []byte{}, true)
	if err != nil {
		fmt.Println("Error ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
}
func TestDeleteProductWithStoreError(t *testing.T) {
	r := SetupProductRouter(&ErrorMock{})

	req, err := MakeRequest("DELETE", "/api/v1/products/1", []byte{}, true)
	if err != nil {
		fmt.Println("Error ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusInternalServerError)
}
func TestDeleteProductWithNonExistId(t *testing.T) {
	r := SetupProductRouter(&Mock{})

	req, err := MakeRequest("DELETE", "/api/v1/products/111", []byte{}, true)
	if err != nil {
		fmt.Println("Error ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusNoContent)
}
func TestDeleteProductWithInvalidId(t *testing.T) {
	r := SetupProductRouter(&Mock{})

	req, err := MakeRequest("DELETE", "/api/v1/products/abc", []byte{}, true)
	if err != nil {
		fmt.Println("Error ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusBadRequest)
}

// test funcs for CreateProduct
func TestCreateProduct(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := SetupProductRouter(&Mock{})
	p := models.Product{ProductName: "Samsung GA", Image: "ss.jpg", Rate: 4}
	json, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Error ", err)
		t.Fail()
	}

	req, err := MakeRequest("POST", "/api/v1/products", json, true)
	if err != nil {
		fmt.Println("Error ", err)
	}

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 201)
}
func TestCreateProductWithStoreError(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := SetupProductRouter(&ErrorMock{})
	p := models.Product{ProductName: "Samsung GA", Image: "ss.jpg", Rate: 4}
	json, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Error ", err)
		t.Fail()
	}

	req, err := MakeRequest("POST", "/api/v1/products", json, true)
	if err != nil {
		fmt.Println("Error ", err)
	}

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusInternalServerError)
}
func TestCreateProductWithDupName(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := SetupProductRouter(&Mock{})
	p := models.Product{ProductName: "Iphone 3", Image: "ss.jpg", Rate: 4}
	json, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Error ", err)
		t.Fail()
	}

	req, err := MakeRequest("POST", "/api/v1/products", json, true)
	if err != nil {
		fmt.Println("Error ", err)
	}

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 201)
}
func TestCreateProductWithInvalidInput(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := SetupProductRouter(&Mock{})
	p := models.Product{ProductName: "", Image: "ss.jpg", Rate: 4}
	json, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Error ", err)
		t.Fail()
	}

	req, err := MakeRequest("POST", "/api/v1/products", json, true)
	if err != nil {
		fmt.Println("Error ", err)
	}

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusBadRequest)
}
func TestCreateProductWithBindJsonErr(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := SetupProductRouter(&Mock{})
	json := []byte(`{ProductNameInvalid: "Samsung", Image: "ss.jpg", Rate: 4}`)

	req, err := MakeRequest("POST", "/api/v1/products", json, true)
	if err != nil {
		fmt.Println("Error ", err)
	}

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusBadRequest)
}

// test funcs for UpdateProduct
func TestUpdateProduct(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := SetupProductRouter(&Mock{})
	p := models.Product{ProductId: 1, ProductName: "Iphone 3 New", Image: "ss.jpg", Rate: 4}
	json, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Error ", err)
		t.Fail()
	}

	req, err := MakeRequest("PUT", "/api/v1/products", json, true)
	if err != nil {
		fmt.Println("Error ", err)
	}

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
}
func TestUpdateProductWithStoreError(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := SetupProductRouter(&ErrorMock{})
	p := models.Product{ProductId: 2, ProductName: "Samsung GA", Image: "ss.jpg", Rate: 4}
	json, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Error ", err)
		t.Fail()
	}

	req, err := MakeRequest("PUT", "/api/v1/products", json, true)
	if err != nil {
		fmt.Println("Error ", err)
	}

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusInternalServerError)
}

// func TestUpdateProductWithDupName(t *testing.T) {
// 	err := godotenv.Load("../.env")
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}
// 	r := SetupProductRouter(&Mock{})
// 	p := models.Product{ProductId: 1, ProductName: "Iphone 3", Image: "ss.jpg", Rate: 4}
// 	json, err := json.Marshal(p)
// 	if err != nil {
// 		fmt.Println("Error ", err)
// 		t.Fail()
// 	}

// 	req, err := MakeRequest("POST", "/api/v1/products", json, true)
// 	if err != nil {
// 		fmt.Println("Error ", err)
// 	}

// 	resp := httptest.NewRecorder()

// 	r.ServeHTTP(resp, req)
// 	assert.Equal(t, resp.Code, 201)
// }
func TestUpdateProductWithInvalidInput(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := SetupProductRouter(&Mock{})
	p := models.Product{ProductId: 1, ProductName: "", Image: "", Rate: 6}
	json, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Error ", err)
		t.Fail()
	}

	req, err := MakeRequest("PUT", "/api/v1/products", json, true)
	if err != nil {
		fmt.Println("Error ", err)
	}

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusBadRequest)
}
func TestUpdateProductWithBindJsonErr(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := SetupProductRouter(&Mock{})
	json := []byte(`{productId2: 1, ProductNameInvalid: "Samsung", Image: "ss.jpg", Rate: 4}`)

	req, err := MakeRequest("PUT", "/api/v1/products", json, true)
	if err != nil {
		fmt.Println("Error ", err)
	}

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusBadRequest)
}
