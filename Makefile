# Makefile

# Target to run tests
test:
    @go test ./... -v

# Default target
.PHONY: test