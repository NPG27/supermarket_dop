package store

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/NPG27/supermarket_dop/internal/domain"
)

type Store interface {
	GetAllProducts() ([]domain.Product, error)
	GetProductByID(id int) (domain.Product, error)
	CreateProduct(product *domain.Product) (*domain.Product, error)
	UpdateProduct(product domain.Product) error
	DeleteProduct(id int) error
	saveProducts(products []domain.Product) error
	loadProducts() ([]domain.Product, error)
}

type jsonStore struct {
	pathToFile string
}

func (s *jsonStore) loadProducts() ([]domain.Product, error) {
	var products []domain.Product
	file, err := os.Open(s.pathToFile)
	if err != nil {
		return nil, errors.New("File not found")
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.New("Cannot read file")
	}

	err = json.Unmarshal(byteValue, &products)
	if err != nil {
		return nil, errors.New("Cannot load json")
	}
	return products, nil
}

func (s *jsonStore) saveProducts(products []domain.Product) error {
	bytes, err := json.Marshal(products)
	if err != nil {
		return err
	}
	return os.WriteFile(s.pathToFile, bytes, 0644)
}

func NewStore(path string) Store {
	return &jsonStore{
		pathToFile: path,
	}
}

func (s *jsonStore) GetAllProducts() ([]domain.Product, error) {
	products, err := s.loadProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *jsonStore) GetProductByID(id int) (domain.Product, error) {
	products, err := s.loadProducts()
	if err != nil {
		return domain.Product{}, err
	}
	for _, product := range products {
		if product.ID == id {
			return product, nil
		}
	}
	return domain.Product{}, errors.New("product not found")
}

func (s *jsonStore) CreateProduct(product *domain.Product) (*domain.Product, error) {
	products, err := s.loadProducts()
	if err != nil {
		return &domain.Product{}, err
	}
	product.ID = len(products) + 1
	products = append(products, *product)
	s.saveProducts(products)
	return product, nil
}

func (s *jsonStore) UpdateProduct(product domain.Product) error {
	products, err := s.loadProducts()
	if err != nil {
		return err
	}
	for i, p := range products {
		if p.ID == product.ID {
			products[i] = product
			return s.saveProducts(products)
		}
	}
	return errors.New("Product not found")
}

func (s *jsonStore) DeleteProduct(id int) error {
	products, err := s.loadProducts()
	if err != nil {
		return err
	}
	for i, p := range products {
		if p.ID == id {
			products = append(products[:i], products[i+1:]...)
			return s.saveProducts(products)
		}
	}
	return errors.New("product not found")
}
