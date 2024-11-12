package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alexanderbkl/mytheresa-promotions/handlers"
	"github.com/alexanderbkl/mytheresa-promotions/models"
	"github.com/alexanderbkl/mytheresa-promotions/store"

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
	http.HandleFunc("/products",
		handlers.ProductsHandler(productsStore),
	)

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
