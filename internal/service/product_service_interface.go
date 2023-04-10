package service

import "github.com/NPG27/supermarket_dop/internal/domain"

type ProductService interface {
	GetAllProducts() ([]domain.Product, error)
	GetProductByID(id int) (domain.Product, error)
	GetProductByPriceGreaterThan(price float64) []domain.Product
	CreateProduct(product *domain.Product) (*domain.Product, error)
	UpdateProduct(id int, product *domain.Product) error
	PatchProduct(id int, product *domain.Product) error
	DeleteProduct(id int) error
}
