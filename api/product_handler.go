package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/huyntsgs/go-rest-api/models"
	//"github.com/huyntsgs/go-rest-api/store"
)

type ProductHandle struct {
	productRepo ProductStore
}

// Creates new ProductHandle.
// ProductHandle accepts ProductStore interface.
// Any data store implements ProductStore could be input of the handle.
func NewProductHandle(productStore ProductStore) ProductHandle {
	return ProductHandle{productStore}
}

// GetProducts handles get product router.
// Function validates parameters and call GetProducts from ProductStore.
func (h ProductHandle) GetProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		limit := DEFAULT_LIMIT
		offset := DEFAULT_OFFSET
		lastId := int64(0)
		if len(c.Query("limit")) != 0 {
			limit, err = strconv.Atoi(c.Query("limit"))
			if err != nil {
				GinAbort(c, http.StatusBadRequest, INVALID_PARAMS, "")
				return
			}
		}

		if len(c.Query("lastId")) > 0 {
			lastId, err = strconv.ParseInt(c.Query("lastId"), 10, 64)
			if err != nil {
				GinAbort(c, http.StatusBadRequest, INVALID_PARAMS, "")
				return
			}
		}

		//products := []*models.Product{}
		products, errc := h.productRepo.GetProducts(limit, offset, lastId)
		if errc != nil {
			log.Println(errc)
			GinAbort(c, http.StatusInternalServerError, errc.(*models.ErrorC).ErrCode, errc.(*models.ErrorC).ErrMsg)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"res":    products,
			})
		}
	}
}

// GetSingleProduct handles get single product router.
// Function validates the parameters and call GetSingleProduct from ProductStore.
func (h ProductHandle) GetSingleProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		//var productId int64
		//var product *models.Product
		//var err error

		productId, err := strconv.ParseInt(c.Param("productId"), 10, 64)
		if err != nil {
			log.Printf("GetSingleProduct.Error %v\n", err)
			GinAbort(c, http.StatusBadRequest, INVALID_PARAMS, "")
			return
		}

		product, errc := h.productRepo.GetSingleProduct(productId)
		if errc != nil {
			GinAbort(c, http.StatusInternalServerError, errc.(*models.ErrorC).ErrCode, errc.(*models.ErrorC).ErrMsg)
			return
		} else {
			if product == nil {
				c.JSON(http.StatusNoContent, gin.H{
					"status": "success", "res": "{}",
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"status": "success", "res": product,
				})
			}
		}
	}
}

// CreateProduct handles create product router.
// Function validates the parameters and call CreateProduct from ProductStore.
func (h ProductHandle) CreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var product models.Product
		err := c.BindJSON(&product)
		if err != nil {
			log.Printf("CreateProduct.Error binding json %v\n", err)
			GinAbort(c, http.StatusBadRequest, INVALID_PARAMS, "")
			return
		}

		if !product.Validate() {
			GinAbort(c, http.StatusBadRequest, INVALID_PARAMS, "")
			return
		}

		lastId, errc := h.productRepo.CreateProduct(&product)
		if errc != nil {
			log.Printf("CreateProduct.Error %v\n", errc)
			GinAbort(c, http.StatusInternalServerError, errc.(*models.ErrorC).ErrCode, errc.(*models.ErrorC).ErrMsg)
			return
		}

		if lastId == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail", "errCode": DUPLICATE_NAME,
			})
		} else {
			c.JSON(http.StatusCreated, gin.H{
				"status": "success", "lastId": lastId,
			})
		}
	}
}

// UpdateProduct handles update product router.
// Function validates the parameters and call UpdateProduct from ProductStore.
func (h ProductHandle) UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var product models.Product
		err := c.BindJSON(&product)
		if err != nil {
			log.Printf("UpdateProduct.Error binding json %v\n", err)
			GinAbort(c, http.StatusBadRequest, INVALID_PARAMS, "")
			return
		}

		if !product.ValidateUpdate() {
			GinAbort(c, http.StatusBadRequest, INVALID_PARAMS, "")
			return
		}

		rows, errc := h.productRepo.UpdateProduct(&product)
		if errc != nil {
			log.Printf("UpdateProduct.Error %v\n", errc)
			GinAbort(c, http.StatusInternalServerError, errc.(*models.ErrorC).ErrCode, errc.(*models.ErrorC).ErrMsg)
		} else {
			if rows == 0 {
				c.JSON(http.StatusNoContent, gin.H{
					"status": "success",
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"status": "success",
				})
			}
		}
	}
}

// DeleteProduct handles delete product router.
// Function validates the parameters and call DeleteProduct from ProductStore.
func (h ProductHandle) DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productId int64
		var err error
		productId, err = strconv.ParseInt(c.Param("productId"), 10, 64)
		if err != nil {
			GinAbort(c, http.StatusBadRequest, INVALID_PARAMS, "")
			return
		}

		rows, errc := h.productRepo.DeleteProduct(productId)
		if errc != nil {
			GinAbort(c, http.StatusInternalServerError, errc.(*models.ErrorC).ErrCode, errc.(*models.ErrorC).ErrMsg)
		} else {
			if rows == 0 {
				c.JSON(http.StatusNoContent, gin.H{
					"status": "success",
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"status": "success",
				})
			}
		}
	}
}
