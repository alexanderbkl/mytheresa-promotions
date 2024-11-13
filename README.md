# Mytheresa Promotions API

## Description

This application provides a REST API endpoint to retrieve a list of products with discounts applied according to specific rules. The application is designed to be performant even with a large number of products.

## Features

- Fetch products with applicable discounts.
- Filter products by category and price.
- Designed to handle over 20,000 products efficiently.
- Includes unit tests for critical functionalities.

## Getting Started

### Prerequisites

- **Docker** and **Docker Compose** installed on your machine.

### Running the Application

Clone the repository:

```shell
git clone https://github.com/alexanderbkl/mytheresa-promotions.git
cd mytheresa-promotions
```

Build and start the application using Docker Compose:
```shell
docker-compose up --build
```

The API will be available at http://localhost:8080.

## API Endpoint

>GET /products

Retrieve a list of products with discounts applied.

## Query Parameters:

- `category` _(optional)_: Filter products by category.

- `priceLessThan` (_optional_): Filter products with prices less than or equal to this value (before discounts).

### Example Request:

> curl "http://localhost:8080/products?category=boots&priceLessThan=80000"

### Response:

```json
{
  "products": [
    {
      "sku": "000003",
      "name": "Ashlington leather ankle boots",
      "category": "boots",
      "price": {
        "original": 71000,
        "final": 49700,
        "discount_percentage": "30%",
        "currency": "EUR"
      }
    },
    // ... up to 5 products
  ]
}
```

## Running Tests
To run the unit tests:
```bash
# If you have Linux/MAC OS installed
make test

# Or directly using Go
go test ./... -v
```

## Project Structure

- `main.go`o: Entry point of the application.
- `handlers/`: Contains HTTP handlers.
- `models/`: Defines data models.
- `store/`: Contains interfaces and implementations for data storage.
- `utils`: Utility functions.
- `tests/`: Contains test files.
- `Dockerfile`: Docker configuration for the application.
- `docker-compose.yml`: Docker Compose configuration.
- `products.json`: Initial dataset.

## Design Decisions
- __Interfaces and Dependency Injection__: To make the code testable without external dependencies, interfaces are used to abstract data stores, and dependency injection is applied.
- __Redis as Data Store__: Redis is used for its in-memory data storage capabilities, ensuring high performance.
- __Modular Code Structure__: The code is organized into packages (handlers, models, store, utils) for better maintainability and testability.
- __Error Handling__: Errors are properly handled and logged to ensure robustness.
- __Unit Tests__: Critical functions are covered with unit tests using the testing package and testify for assertions.