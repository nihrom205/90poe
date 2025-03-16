build:
	go build -o app cmd/main.go

run:
	go build -o app cmd/main.go && \
	HTTP_ADDR=:8080 \
	DEBUG_ERRORS=1 \
	DSN="host=localhost user=postgres password=postgres dbname=link port=5432 sslmode=disable search_path=go_home_work" \
	MIGRATIONS_PATH="file://./internal/app/migrations" \
	./app