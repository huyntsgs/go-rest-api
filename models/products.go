package models

import (
	//"log"
	"time"

	"github.com/asaskevich/govalidator"
)

type (
	Product struct {
		ProductId   int64     `json:"productId"`
		ProductName string    `json:"productName"`
		Image       string    `json:"image"`
		Rate        int       `json:"rate"`
		CreatedAt   time.Time `json:"createAt"`
	}

	ProductTime struct {
		products []Product
		expire   int64
	}
)

func (p Product) Validate() bool {
	if govalidator.IsNull(p.ProductName) || govalidator.IsNull(p.Image) || !govalidator.InRange(p.Rate, 1, 5) {
		return false
	}
	return true
}

func (p Product) ValidateUpdate() bool {
	if p.ProductId == 0 {
		return false
	}
	if govalidator.IsNull(p.ProductName) || govalidator.IsNull(p.Image) || !govalidator.InRange(p.Rate, 1, 5) {
		return false
	}

	return true
}
