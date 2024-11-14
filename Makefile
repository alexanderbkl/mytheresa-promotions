# Makefile

# Target to run the application
run:
    @docker compose up app --build

# Target to run tests
test:
    @docker compose up test --build

# Default target
.PHONY: test