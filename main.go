package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/alexanderbkl/mytheresa-promotions/models"
	"github.com/alexanderbkl/mytheresa-promotions/store"
	"github.com/alexanderbkl/mytheresa-promotions/utils"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	// Set up Redis connection
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Load products into Redis
	err := loadProducts(rdb)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize product store
	productsStore := store.NewRedisStore(rdb)

	// Set up HTTP handlers
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		handleProducts(w, r, productsStore)
	})

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadProducts(rdb *redis.Client) error {
	// Check if products are already loaded
	exists, err := rdb.Exists(ctx, "products_loaded").Result()
	if err != nil {
		return err
	}

	if exists == 1 {
		return nil
	}

	file, err := os.ReadFile("products.json")
	if err != nil {
		return err
	}

	var data struct {
		Products []models.Product `json:"products"`
	}

	err = json.Unmarshal(file, &data)
	if err != nil {
		return err
	}

	for _, product := range data.Products {
		productJSON, err := json.Marshal(product)
		if err != nil {
			continue
		}
		err = rdb.Set(ctx, product.SKU, productJSON, 0).Err()
		if err != nil {
			continue
		}
	}

	// Set a flag indicating products are loaded
	return rdb.Set(ctx, "products_loaded", "true", 0).Err()
}

func handleProducts(w http.ResponseWriter, r *http.Request, store store.ProductStore) {
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

	// Fetch all products
	products, err := store.GetProducts(ctx)
	if err != nil {
		http.Error(w, "Error fetching products", http.StatusInternalServerError)
		return
	}

	var responses []models.ProductResponse
	for _, product := range products {
		// Apply filters
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
			priceDetail.DiscountPercentage = discountPercentage
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
