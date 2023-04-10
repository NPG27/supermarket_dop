package service

import (
	"errors"
	"time"

	"github.com/NPG27/supermarket_dop/internal/domain"
	"github.com/NPG27/supermarket_dop/internal/repository"
)

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *productService {
	return &productService{productRepo: repo}
}

func (s *productService) GetAllProducts() ([]domain.Product, error) {
	products, err := s.productRepo.GetAllProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *productService) GetProductByID(id int) (domain.Product, error) {
	product, err := s.productRepo.GetProductByID(id)
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

func (s *productService) GetProductByPriceGreaterThan(price float64) []domain.Product {
	return s.productRepo.GetProductByPriceGreaterThan(price)
}

func validateProduct(product domain.Product) bool {
	isCorrect := true
	if product.Name == "" || product.Quantity == 0 || product.CodeValue == "" || product.Expiration == "" || product.Price == float64(0) {
		isCorrect = false
	}
	return isCorrect
}

func (s *productService) CreateProduct(product *domain.Product) (*domain.Product, error) {
	if !validateProduct(*product) {
		return nil, errors.New("Product is missing required values")
	}
	_, errProductExpiration := time.Parse("02/01/2006", product.Expiration)
	if errProductExpiration != nil {
		return &domain.Product{}, errors.New("Invalid expiration date format")
	}
	productByCode := s.productRepo.GetProductByCode()
	if _, codeValueExists := productByCode[product.CodeValue]; codeValueExists {
		return &domain.Product{}, errors.New("Code value already exists")
	}
	pr, err := s.productRepo.CreateProduct(product)
	if err != nil {
		return &domain.Product{}, err
	}
	return pr, nil
}

func (s *productService) UpdateProduct(id int, product *domain.Product) error {
	if !validateProduct(*product) {
		return errors.New("Product is missing required values")
	}
	productByCode := s.productRepo.GetProductByCode()
	if productMap, codeValueExists := productByCode[product.CodeValue]; productMap.ID != id && codeValueExists {
		return errors.New("Code value already exists")
	}
	_, errProductExpiration := time.Parse("02/01/2006", product.Expiration)
	if errProductExpiration != nil {
		return errors.New("Invalid expiration date format")
	}

	err := s.productRepo.UpdateProduct(id, product)
	if err != nil {
		return err
	}
	return nil
}

func (s *productService) PatchProduct(id int, product *domain.Product) error {
	productByCode := s.productRepo.GetProductByCode()
	if productMap, codeValueExists := productByCode[product.CodeValue]; product.CodeValue != "" && productMap.ID != id && codeValueExists {
		return errors.New("Code value already exists")
	}
	if _, errProductExpiration := time.Parse("02/01/2006", product.Expiration); product.Expiration != "" && errProductExpiration != nil {
		return errors.New("Invalid expiration date format")
	}
	err := s.productRepo.PatchProduct(id, product)
	if err != nil {
		return err
	}
	return nil
}

func (s *productService) DeleteProduct(id int) error {
	err := s.productRepo.DeleteProduct(id)
	if err != nil {
		return err
	}
	return nil
}
