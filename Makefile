build:
	@rm -rf .verso;\
	go build -o bin/verso cmd/verso/main.go;

test:
	go test ./...
