package domain

// Product represents a product in the inventory.
// swagger:model
type Product struct {
	// The ID of the product.
	//
	// required: true
	// example: 1
	ID int `json:"id"`

	// The name of the product.
	//
	// required: true
	// example: "Milk"
	Name string `json:"name"`

	// The quantity of the product.
	//
	// required: true
	// example: 10
	Quantity int `json:"quantity"`

	// The code value of the product.
	//
	// required: true
	// example: "ABC123"
	CodeValue string `json:"code_value"`

	// Whether the product is published.
	//
	// required: false
	// example: true
	IsPublished bool `json:"is_published"`

	// The expiration date of the product.
	//
	// required: true
	// example: "2023-12-31"
	Expiration string `json:"expiration"`

	// The price of the product.
	//
	// required: true
	// example: 2.99
	Price float64 `json:"price"`
}
