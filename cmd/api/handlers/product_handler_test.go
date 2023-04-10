package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/NPG27/supermarket_dop/cmd/api/handlers"
	"github.com/NPG27/supermarket_dop/internal/domain"
	"github.com/NPG27/supermarket_dop/internal/middleware"
	"github.com/NPG27/supermarket_dop/internal/repository"
	"github.com/NPG27/supermarket_dop/internal/service"
	"github.com/NPG27/supermarket_dop/pkg/store"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type response struct {
	Data interface{} `json:"data"`
}

func createServer(token string) *gin.Engine {

	if token != "" {
		err := os.Setenv("TOKEN", token)
		if err != nil {
			panic(err)
		}
	}

	db := store.NewStore("./products_copy.json")
	repo, _ := repository.NewProductRepository(db)
	service := service.NewProductService(repo)
	productHandler := handlers.NewProductHandler(service)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	products := r.Group("/products")
	products.Use(middleware.VerifyToken())
	{
		products.GET("", productHandler.GetAllProducts)
		products.GET("/:id", productHandler.GetProductByID)
		products.GET("/filter", productHandler.GetProductByPriceGreaterThan)
		products.POST("", productHandler.CreateProduct)
		products.PATCH("/:id", productHandler.PatchProduct)
		products.PUT("/:id", productHandler.UpdateProduct)
		products.DELETE("/:id", productHandler.DeleteProduct)
	}
	return r
}

func createRequestTest(method string, url string, body string, token string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	if token != "" {
		req.Header.Add("TOKEN", token)
	}
	return req, httptest.NewRecorder()
}

func loadProducts(path string) ([]domain.Product, error) {
	var products []domain.Product
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(file), &products)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func writeProducts(path string, list []domain.Product) error {
	bytes, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return err
}

func Test_GetAllProducts_OK(t *testing.T) {
	var expected = response{Data: []domain.Product{}}

	r := createServer("my-secret-value")
	req, rr := createRequestTest(http.MethodGet, "/products", "", "my-secret-value")

	p, err := loadProducts("./products_copy.json")
	if err != nil {
		panic(err)
	}
	expected.Data = p
	actual := map[string][]domain.Product{}

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	err = json.Unmarshal(rr.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected.Data, actual["data"])
}

func Test_GetProductByID_OK(t *testing.T) {
	var expected = response{Data: domain.Product{}}

	r := createServer("my-secret-value")
	req, rr := createRequestTest(http.MethodGet, "/products/1", "", "my-secret-value")
	r.ServeHTTP(rr, req)

	p, err := loadProducts("./products_copy.json")
	if err != nil {
		panic(err)
	}
	expected.Data = p[0]
	actual := map[string]domain.Product{}

	assert.Equal(t, http.StatusOK, rr.Code)
	err = json.Unmarshal(rr.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected.Data, actual["data"])
}

func Test_CreateProduct_OK(t *testing.T) {
	var expected = response{Data: domain.Product{
		ID:          501,
		Name:        "New product",
		Quantity:    20,
		CodeValue:   "TEST1",
		IsPublished: false,
		Expiration:  "15/12/2023",
		Price:       100.50,
	}}

	product, _ := json.Marshal(expected.Data)
	r := createServer("my-secret-value")
	req, rr := createRequestTest(http.MethodPost, "/products", string(product), "my-secret-value")

	p, _ := loadProducts("./products_copy.json")

	r.ServeHTTP(rr, req)
	actual := map[string]domain.Product{}
	_ = json.Unmarshal(rr.Body.Bytes(), &actual)
	_ = writeProducts("./products_copy.json", p)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, expected.Data, actual["data"])
}

func Test_Delete_OK(t *testing.T) {

	r := createServer("my-secret-token")
	req, rr := createRequestTest(http.MethodDelete, "/products/1", "", "my-secret-token")

	p, err := loadProducts("./products_copy.json")
	if err != nil {
		panic(err)
	}

	r.ServeHTTP(rr, req)

	err = writeProducts("./products_copy.json", p)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, 204, rr.Code)
	assert.Nil(t, rr.Body.Bytes())
}

func Test_BadRequest(t *testing.T) {

	test := []string{http.MethodPut, http.MethodPatch, http.MethodDelete}

	r := createServer("my-secret-token")
	for _, method := range test {
		req, rr := createRequestTest(method, "/products/not_a_number", "", "my-secret-token")
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	}

}

func Test_BadRequest_GET(t *testing.T) {

	test := http.MethodGet
	r := createServer("my-secret-token")
	req, rr := createRequestTest(test, "/products/not_a_number", "", "my-secret-token")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

}

func Test_NotFound(t *testing.T) {

	test := []string{http.MethodGet, http.MethodPatch, http.MethodDelete}

	r := createServer("my-secret-token")
	for _, method := range test {
		req, rr := createRequestTest(method, "/products/800", "{}", "my-secret-token")
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	}
}

func Test_Unauthorized(t *testing.T) {

	test := []string{http.MethodPut, http.MethodPatch, http.MethodDelete}

	r := createServer("my-secret-token")
	for _, method := range test {
		req, rr := createRequestTest(method, "/products/10", "{}", "not-my-token")
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	}
}

func Test_Unauthorized_GET(t *testing.T) {
	r := createServer("my-secret-token")
	req, rr := createRequestTest(http.MethodGet, "/products", "", "not-my-token")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func Test_Unauthorized_POST(t *testing.T) {
	r := createServer("my-secret-token")
	req, rr := createRequestTest(http.MethodPost, "/products", "{}", "not-my-token")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
