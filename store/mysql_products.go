package store

import (
	"github.com/huyntsgs/go-rest-api/models"
)

// GetProducts gets products from table products with number is limit and product id greater than lastId.
// Query database with paging base on limit, offset has bad performance when the number of record is big
// Using lastId will has better performance. lastId is retrieved from previous page query.
func (mysql *MySqlDB) GetProducts(limit, offset int, lastId int64) ([]*models.Product, error) {

	// productsT := cache[lastId]
	// if productsT != nil && productsT.expire > time.Now().Unix() {
	// 	return productsT.products, nil
	// }
	//db := GetDB()
	// if db == nil {
	// 	return nil, errors.New("Can not connect to database")
	// }

	query, err := mysql.DB.Query("SELECT * FROM products where product_id > ? ORDER BY product_id ASC limit ?", lastId, limit)
	if err != nil {
		return nil, models.NewError("Internal server error", ERR_INTERNAL)
	}
	res := []*models.Product{}
	for query.Next() {
		product := models.Product{}
		err = query.Scan(&product.ProductId, &product.ProductName, &product.Image, &product.Rate, &product.CreatedAt)
		if err != nil {
			return nil, models.NewError("Internal server error", ERR_INTERNAL)
		}
		res = append(res, &product)
	}

	return res, nil
}

// GetSingleProduct returns product from database with provided productId
func (mysql *MySqlDB) GetSingleProduct(productId int64) (*models.Product, error) {

	query, err := mysql.DB.Query("SELECT * FROM products where product_id = ?", productId)
	if err != nil {
		return nil, models.NewError("Internal server error", ERR_INTERNAL)
	}
	defer query.Close()
	var product models.Product
	for query.Next() {
		err = query.Scan(&product.ProductId, &product.ProductName, &product.Image, &product.Rate, &product.CreatedAt)
		if err != nil {
			return nil, models.NewError("Internal server error", ERR_INTERNAL)
		}
		return &product, nil
	}

	return nil, nil
}

// DeleteProduct deletes from products table with provided productId
func (mysql *MySqlDB) DeleteProduct(productId int64) (int64, error) {
	res, err := mysql.DB.Exec("DELETE FROM products where product_id = ?", productId)
	if err != nil {
		return int64(0), models.NewError("Internal server error", ERR_INTERNAL)
	}
	row, _ := res.RowsAffected()
	return row, nil
}

// UpdateProduct updates product data to database.
// Function returns number of row is updated.
func (mysql *MySqlDB) UpdateProduct(p *models.Product) (int64, error) {

	if !p.ValidateUpdate() {
		return 0, models.NewError("Invalid product data", INVALID_DATA)
	}

	query, err := mysql.DB.Exec("UPDATE products set product_name =?, image =?, rate =? where product_id = ?", p.ProductName, p.Image, p.Rate, p.ProductId)
	if err != nil {
		return 0, models.NewError("Internal server error", ERR_INTERNAL)
	}
	row, _ := query.RowsAffected()
	return row, nil
}

// CreateProduct add new product to database.
// Function check duplicated product name before insert.
// Returns the product id just inserted.
func (mysql *MySqlDB) CreateProduct(p *models.Product) (int64, error) {

	if !p.Validate() {
		return 0, models.NewError("Invalid product data", INVALID_DATA)
	}

	sql := "INSERT INTO products (product_name, image, rate) " +
		"SELECT ?, ?, ? FROM DUAL WHERE NOT EXISTS " +
		"(SELECT product_name FROM products WHERE product_name = ?)"

	query, err := mysql.DB.Exec(sql, p.ProductName, p.Image, p.Rate, p.ProductName)
	if err != nil {
		return 0, models.NewError("Internal server error", ERR_INTERNAL)
	}
	lastId, _ := query.LastInsertId()
	return lastId, nil
}
