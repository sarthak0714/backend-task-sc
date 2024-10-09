.PHONY: run build clean

run: build
	@./bin/server

build:
	@mkdir -p bin
	@go build -o bin/server cmd/main.go

clean:
	@rm -rf ./bin/server