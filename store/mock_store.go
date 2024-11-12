// store/mock_store.go
package store

import (
	"context"

	"github.com/alexanderbkl/mytheresa-promotions/models"
)

// MockStore is a mock implementation of the ProductStore interface
type MockStore struct {
	Products []models.Product
}

// NewMockStore creates a new MockStore instance
func NewMockStore(products []models.Product) *MockStore {
	return &MockStore{Products: products}
}

// GetProducts returns the list of products
func (s *MockStore) GetProducts(ctx context.Context) ([]models.Product, error) {
	return s.Products, nil
}
