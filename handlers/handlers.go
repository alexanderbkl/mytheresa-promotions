// handlers/handlers.go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/alexanderbkl/mytheresa-promotions/models"
	"github.com/alexanderbkl/mytheresa-promotions/store"
	"github.com/alexanderbkl/mytheresa-promotions/utils"
)

func ProductsHandler(store store.ProductStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Parse query parameters
		category := r.URL.Query().Get("category")
		priceLessThanStr := r.URL.Query().Get("priceLessThan")
		var priceLessThan int
		if priceLessThanStr != "" {
			var err error
			priceLessThan, err = strconv.Atoi(priceLessThanStr)
			if err != nil {
				http.Error(w, "Invalid priceLessThan value", http.StatusBadRequest)
				return
			}
		}

		// Fetch products from the store
		products, err := store.GetProducts(ctx)
		if err != nil {
			http.Error(w, "Error fetching products", http.StatusInternalServerError)
			return
		}

		log.Printf("Fetched %d products", len(products))

		var responses []models.ProductResponse
		for _, product := range products {
			// Apply filters from query parameters
			if category != "" && product.Category != category {
				continue
			}
			if priceLessThanStr != "" && product.Price > priceLessThan {
				continue
			}

			// Apply discounts
			discountPercentage, finalPrice := utils.CalculateDiscount(product)

			priceDetail := models.PriceDetail{
				Original: product.Price,
				Final:    finalPrice,
				Currency: "EUR",
			}
			if discountPercentage != nil {
				priceDetail.DiscountPercentage = *discountPercentage
			}

			responses = append(responses, models.ProductResponse{
				SKU:      product.SKU,
				Name:     product.Name,
				Category: product.Category,
				Price:    priceDetail,
			})

			// Limit to 5 products
			if len(responses) == 5 {
				break
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string][]models.ProductResponse{
			"products": responses,
		})
	}
}
