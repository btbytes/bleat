BINARY_NAME := bleat

.PHONY: build clean

build:
	@go build -o $(BINARY_NAME) .
	@echo "Built $(BINARY_NAME)"

clean:
	@rm -f $(BINARY_NAME)
	@echo "Cleaned $(BINARY_NAME)"
