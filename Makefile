dc:
	docker-compose up  --remove-orphans --build

build:
	go build -o app_port cmd/main.go

buildLinux:
	GOOS=linux go build -o app_port cmd/main.go

run:
	go build -o app_port cmd/main.go && \
	HTTP_ADDR=:8080 \
	DSN="file::memory:?cache=shared" \
	MIGRATIONS_PATH="file://./internal/app/migrations" \
	./app_port

test:
	go test -v ./...

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8

lint:
	golangci-lint run ./...