.PHONY: docker-up docker-down sqlc-product run-product build-product tidy

docker-up:
	docker compose -f deploy/docker-compose.yml up -d

docker-down:
	docker compose -f deploy/docker-compose.yml down

sqlc-product:
	cd services/product && go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate

run-product:
	cd services/product && go run ./cmd/api

build-product:
	cd services/product && go build -o bin/product-api ./cmd/api

tidy:
	cd pkg && go mod tidy
	cd services/product && go mod tidy
	go work sync
