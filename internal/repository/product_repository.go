package repository

import (
	"errors"

	"github.com/NPG27/supermarket_dop/internal/domain"
	"github.com/NPG27/supermarket_dop/pkg/store"
)

type productRepository struct {
	storage       store.Store
	productByCode map[string]domain.Product
}

func NewProductRepository(storage store.Store) (ProductRepository, error) {
	repo := &productRepository{storage: storage}
	repo.productByCode = make(map[string]domain.Product)
	products, err := storage.GetAllProducts()
	if err != nil {
		return nil, err
	}
	for _, product := range products {
		repo.productByCode[product.CodeValue] = product
	}
	return repo, nil
}

func (r *productRepository) GetProductByCode() map[string]domain.Product {
	return r.productByCode
}

func (r *productRepository) GetAllProducts() ([]domain.Product, error) {
	products, err := r.storage.GetAllProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) GetProductByID(id int) (domain.Product, error) {
	products, err := r.storage.GetAllProducts()
	if err != nil {
		return domain.Product{}, err
	}
	for _, p := range products {
		if p.ID == id {
			return p, nil
		}
	}
	return domain.Product{}, errors.New("Product not found")
}

func (r *productRepository) GetProductByPriceGreaterThan(price float64) []domain.Product {
	products, _ := r.storage.GetAllProducts()
	var filteredProducts []domain.Product
	for _, product := range products {
		if product.Price > price {
			filteredProducts = append(filteredProducts, product)
		}
	}
	return filteredProducts
}

func (r *productRepository) CreateProduct(product *domain.Product) (*domain.Product, error) {
	product, err := r.storage.CreateProduct(product)
	if err != nil {
		return &domain.Product{}, err
	}
	r.productByCode[product.CodeValue] = *product
	return product, nil
}

func (r *productRepository) UpdateProduct(id int, product *domain.Product) error {
	product.ID = id
	err := r.storage.UpdateProduct(*product)
	if err != nil {
		return err
	}
	r.productByCode[product.CodeValue] = *product
	return nil
}

func (r *productRepository) PatchProduct(id int, product *domain.Product) error {
	products, err := r.storage.GetAllProducts()
	if err != nil {
		return err
	}
	for i, p := range products {
		if p.ID == id {
			product.ID = id
			if product.Name == "" {
				product.Name = products[i].Name
			}
			if product.Quantity == 0 {
				product.Quantity = products[i].Quantity
			}
			if product.CodeValue == "" {
				product.CodeValue = products[i].CodeValue
			}
			if product.IsPublished == false {
				product.IsPublished = products[i].IsPublished
			}
			if product.Expiration == "" {
				product.Expiration = products[i].Expiration
			}
			if product.Price == 0 {
				product.Price = products[i].Price
			}
			r.storage.UpdateProduct(*product)
			r.productByCode[product.CodeValue] = *product
			return nil
		}
	}
	return errors.New("Product not found")
}

func (r *productRepository) DeleteProduct(id int) error {
	products, err := r.storage.GetAllProducts()
	if err != nil {
		return err
	}
	for _, p := range products {
		if p.ID == id {
			r.storage.DeleteProduct(p.ID)
			delete(r.productByCode, p.CodeValue)
			return nil
		}
	}
	return errors.New("Product not found")
}
