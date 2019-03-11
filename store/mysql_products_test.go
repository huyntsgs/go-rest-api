package store

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/huyntsgs/go-rest-api/models"
	"github.com/joho/godotenv"
)

func InsertProducts(n int) *MySqlDB {
	mysql := InitDB()

	for i := 0; i < n; i++ {
		p := models.Product{ProductName: "Samsaung S" + string(i), Image: string(i) + "ss.jpg", Rate: 4}
		mysql.CreateProduct(&p)
	}
	return mysql
}

// test getproducts funcs
func TestGetProducts(t *testing.T) {
	mysql := InsertProducts(30)
	ps, err := mysql.GetProducts(10, 0, 0)
	assert.Equal(t, nil, err)
	assert.Equal(t, 10, len(ps))

	// Get next page as paging
	ps1, err := mysql.GetProducts(10, 0, ps[9].ProductId)
	assert.Equal(t, nil, err)
	assert.Equal(t, 10, len(ps1))

	assert.Equal(t, ps1[0].ProductId, ps[9].ProductId+1)
}

// test createproducts funcs
func TestCreateProduct(t *testing.T) {
	mysql := InitDB()
	mysql.DeleteProducts()
	p := models.Product{ProductName: "Samsung S", Image: "ss.jpg", Rate: 4}
	_, err := mysql.CreateProduct(&p)
	assert.Equal(t, nil, err)
}
func TestCreateProductWithInvalidData(t *testing.T) {
	mysql := InitDB()
	mysql.DeleteProducts()
	p := models.Product{ProductName: "Samsung S", Image: "ss.jpg", Rate: 7}
	_, err := mysql.CreateProduct(&p)
	assert.NotEqual(t, nil, err)
}
func TestCreateProductWithNoDBConnection(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fail()
	}
	mysql := new(MySqlDB)
	p := models.Product{ProductName: "Samsung S", Image: "ss.jpg", Rate: 7}
	_, err = mysql.CreateProduct(&p)
	assert.NotEqual(t, nil, err)
}

// test deleteproduct funcs
func TestDeleteProducts(t *testing.T) {
	mysql := InitDB()
	mysql.DeleteProducts()
	p := models.Product{ProductName: "Samsung S", Image: "ss.jpg", Rate: 4}
	_, err := mysql.CreateProduct(&p)
	assert.Equal(t, nil, err)

	_, err = mysql.DeleteProduct(p.ProductId)
	assert.Equal(t, nil, err)
}

// test getsingleproduct funcs
func TestGetSingleProduct(t *testing.T) {
	mysql := InitDB()
	mysql.DeleteProducts()
	p := models.Product{ProductName: "Samsung S", Image: "ss.jpg", Rate: 4}
	_, err := mysql.CreateProduct(&p)
	assert.Equal(t, nil, err)

	pr, err := mysql.GetSingleProduct(p.ProductId)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, pr)
}

// test updateproduct funcs
func TestUpdateProduct(t *testing.T) {
	mysql := InitDB()
	mysql.DeleteProducts()
	p := models.Product{ProductName: "Samsung S", Image: "ss.jpg", Rate: 4}
	pid, err := mysql.CreateProduct(&p)
	assert.Equal(t, nil, err)
	p.ProductName = "SamsungS10"
	p.ProductId = pid
	_, err = mysql.UpdateProduct(&p)
	assert.Equal(t, nil, err)

	pn, err := mysql.GetSingleProduct(pid)
	assert.NotEqual(t, nil, pn)
	assert.Equal(t, nil, err)
	assert.Equal(t, "SamsungS10", pn.ProductName)
}

func TestUpdateProductWithInvalidName(t *testing.T) {
	mysql := InitDB()
	mysql.DeleteProducts()
	p := models.Product{ProductName: "Samsung S", Image: "ss.jpg", Rate: 4}
	pid, err := mysql.CreateProduct(&p)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, 0, pid)
	p.ProductName = ""
	p.ProductId = pid
	_, err = mysql.UpdateProduct(&p)
	assert.NotEqual(t, nil, err)

}
func TestUpdateProductWithInvalidId(t *testing.T) {
	mysql := InitDB()
	mysql.DeleteProducts()
	p := models.Product{ProductName: "Samsung S", Image: "ss.jpg", Rate: 4}
	pid, err := mysql.CreateProduct(&p)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, 0, pid)

	p.ProductId = 0
	_, err = mysql.UpdateProduct(&p)
	assert.NotEqual(t, nil, err)
}
func TestUpdateProductWithInvalidRate(t *testing.T) {
	mysql := InitDB()
	mysql.DeleteProducts()
	p := models.Product{ProductName: "Samsung S", Image: "ss.jpg", Rate: 4}
	pid, err := mysql.CreateProduct(&p)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, 0, pid)

	p.Rate = 6
	_, err = mysql.UpdateProduct(&p)
	assert.NotEqual(t, nil, err)
}
