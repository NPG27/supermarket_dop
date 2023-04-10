package handlers

import (
	"strconv"

	"github.com/NPG27/supermarket_dop/internal/domain"
	"github.com/NPG27/supermarket_dop/internal/service"
	"github.com/NPG27/supermarket_dop/pkg/web"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService}
}

func (h *ProductHandler) GetAllProducts(ctx *gin.Context) {
	products, err := h.productService.GetAllProducts()
	if err != nil {
		web.Failure(ctx, 500, err)
		return
	}
	web.Success(ctx, 200, products)
}

func (h *ProductHandler) GetProductByID(ctx *gin.Context) {
	id := ctx.Param("id")
	idConverted, errConverted := strconv.Atoi(id)
	if errConverted != nil {
		web.Failure(ctx, 400, errConverted)
		return
	}
	product, err := h.productService.GetProductByID(idConverted)
	if err != nil {
		web.Failure(ctx, 404, err)
		return
	}
	web.Success(ctx, 200, product)
}

func (h *ProductHandler) GetProductByPriceGreaterThan(ctx *gin.Context) {
	price := ctx.Query("price")
	priceConverted, errConverted := strconv.ParseFloat(price, 64)
	if errConverted != nil {
		web.Failure(ctx, 400, errConverted)
		return
	}
	products := h.productService.GetProductByPriceGreaterThan(priceConverted)
	web.Success(ctx, 200, products)
}

// Post godoc
// @Summary      Create a new product
// @Description  Create a new product in repository
// @Tags         products
// @Produce      json
// @Param        token header string true "token"
// @Param        product body domain.Product true "product"
// @Success      201 {object}  web.response
// @Failure      400 {object}  web.errorResponse
// @Failure      404 {object}  web.errorResponse
// @Router       /products [post]
func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	var newProduct domain.Product
	if err := ctx.ShouldBindJSON(&newProduct); err != nil {
		web.Failure(ctx, 400, err)
		return
	}
	product, err := h.productService.CreateProduct(&newProduct)
	if err != nil {
		web.Failure(ctx, 400, err)
		return
	}
	web.Success(ctx, 201, product)
}

func (h *ProductHandler) UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	idConverted, errConverted := strconv.Atoi(id)
	if errConverted != nil {
		web.Failure(ctx, 500, errConverted)
		return
	}
	var product domain.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		web.Failure(ctx, 400, err)
		return
	}
	err := h.productService.UpdateProduct(idConverted, &product)
	if err != nil && err.Error() == "Product not found" {
		web.Failure(ctx, 404, err)
		return
	} else if err != nil {
		web.Failure(ctx, 400, err)
		return
	}
	web.Success(ctx, 200, gin.H{"message": "product updated"})
}

func (h *ProductHandler) PatchProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	idConverted, errConverted := strconv.Atoi(id)
	if errConverted != nil {
		web.Failure(ctx, 500, errConverted)
	}
	var product domain.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		web.Failure(ctx, 400, err)
		return
	}
	err := h.productService.PatchProduct(idConverted, &product)
	if err != nil && err.Error() == "Product not found" {
		web.Failure(ctx, 404, err)
		return
	} else if err != nil {
		web.Failure(ctx, 400, err)
		return
	}
	web.Success(ctx, 200, gin.H{"message": "product updated"})
}

func (h *ProductHandler) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	idConverted, errConverted := strconv.Atoi(id)
	if errConverted != nil {
		web.Failure(ctx, 500, errConverted)
		return
	}
	err := h.productService.DeleteProduct(idConverted)
	if err != nil {
		web.Failure(ctx, 404, err)
		return
	}
	web.Success(ctx, 204, nil)
}
