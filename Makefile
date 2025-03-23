.PHONY: build test clean install

# Build the application
build:
	@echo "Building FixFiles..."
	go build -o fixfiles

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	rm -f fixfiles
	rm -f error-context-*.txt

# Install globally
install: build
	@echo "Installing FixFiles..."
	cp ./fixfiles /usr/local/bin/fixfiles
