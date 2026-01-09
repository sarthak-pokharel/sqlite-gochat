.PHONY: build run test test-v clean

APP_NAME=app
BINARY=dist/$(APP_NAME)

build:
	@mkdir -p dist
	CGO_ENABLED=0 go build -o $(BINARY) src/main.go

run: build
	./$(BINARY)

test:
	go test ./src/... --count=1

test-v:
	go test ./src/... --count=1 -v

clean:
	rm -rf dist
