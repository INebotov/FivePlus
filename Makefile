download:
	go mod tidy
	go mod download

run:
	go run main.go ./resources/config.yaml

compile:
	go mod tidy
	go mod download
	GOOS=linux GOARCH=arm64 go build -o bin/linux-arm64 main.go
	GOOS=windows GOARCH=arm64 go build -o bin/windows-arm64 main.go
