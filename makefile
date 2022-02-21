build:
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/app

run:
	PORT=8090 go run main.go