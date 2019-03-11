package store

import (
	"github.com/huyntsgs/go-rest-api/models"
)

// Defines ProductStore interface
// Product handle accepts the interface.
type ProductStore interface {
	GetProducts(limit, offset int, lastId int64) ([]*models.Product, error)
	GetSingleProduct(productId int64) (*models.Product, error)
	DeleteProduct(productId int64) (int64, error)
	UpdateProduct(p *models.Product) (int64, error)
	CreateProduct(p *models.Product) (int64, error)
}

// Defines UserStore interface
// User handle accepts the interface.
type UserStore interface {
	Register(user *models.User) error
	Login(user *models.User) (*models.User, error)
}
