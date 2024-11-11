package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"

    "github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type Product struct {
    SKU      string `json:"sku"`
    Name     string `json:"name"`
    Category string `json:"category"`
    Price    int    `json:"price"`
}

type ProductResponse struct {
    SKU               string       `json:"sku"`
    Name              string       `json:"name"`
    Category          string       `json:"category"`
    Price             PriceDetails `json:"price"`
}

type PriceDetails struct {
    Original           int     `json:"original"`
    Final              int     `json:"final"`
    DiscountPercentage *string `json:"discount_percentage"`
    Currency           string  `json:"currency"`
}

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
    loadProducts(rdb)

    // Set up HTTP handlers
    http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
        handleProducts(w, r, rdb)
    })

    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadProducts(rdb *redis.Client) {
    // Check if products are already loaded
    exists, err := rdb.Exists(ctx, "products_loaded").Result()
    if err != nil || exists == 1 {
        return
    }

    file, err := os.ReadFile("products.json")
    if err != nil {
        log.Fatal(err)
    }

    var data struct {
        Products []Product `json:"products"`
    }

    err = json.Unmarshal(file, &data)
    if err != nil {
        log.Fatal(err)
    }

    for _, product := range data.Products {
        productJSON, err := json.Marshal(product)
        if err != nil {
            log.Println(err)
            continue
        }
        err = rdb.Set(ctx, product.SKU, productJSON, 0).Err()
        if err != nil {
            log.Println(err)
            continue
        }
    }

    // Set a flag indicating products are loaded
    rdb.Set(ctx, "products_loaded", "true", 0)
}

func handleProducts(w http.ResponseWriter, r *http.Request, rdb *redis.Client) {
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
    keys, err := rdb.Keys(ctx, "00000*").Result()
    if err != nil {
        http.Error(w, "Error fetching products", http.StatusInternalServerError)
        return
    }

    var products []ProductResponse
    for _, key := range keys {
        val, err := rdb.Get(ctx, key).Result()
        if err != nil {
            continue
        }

        var product Product
        err = json.Unmarshal([]byte(val), &product)
        if err != nil {
            continue
        }

        // Apply filters
        if category != "" && product.Category != category {
            continue
        }
        if priceLessThanStr != "" && product.Price > priceLessThan {
            continue
        }

        // Apply discounts
        discountPercentage, finalPrice := calculateDiscount(product)

        priceDetails := PriceDetails{
            Original: product.Price,
            Final:    finalPrice,
            Currency: "EUR",
        }
        if discountPercentage != nil {
            priceDetails.DiscountPercentage = discountPercentage
        }

        products = append(products, ProductResponse{
            SKU:      product.SKU,
            Name:     product.Name,
            Category: product.Category,
            Price:    priceDetails,
        })

        // Limit to 5 products
        if len(products) == 5 {
            break
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string][]ProductResponse{
        "products": products,
    })
}

func calculateDiscount(product Product) (*string, int) {
    var discount int

    if product.SKU == "000003" {
        discount = 15
    }
    if product.Category == "boots" && discount < 30 {
        discount = 30
    }

    if discount > 0 {
        discountPercentage := fmt.Sprintf("%d%%", discount)
        finalPrice := product.Price * (100 - discount) / 100
        return &discountPercentage, finalPrice
    }

    return nil, product.Price
}
