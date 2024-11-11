// store/store.go
package store

import (
	"context"

	"github.com/alexanderbkl/mytheresa-promotions/models"
)

type ProductStore interface {
	GetProducts(ctx context.Context) ([]models.Product, error)
}
