include .env
export


start-project:
	DB_CONNECTION="$(DB_CONNECTION)" \
	ADMIN_PASS="$(ADMIN_PASS)" \
	go run cmd/main.go

migrate-up:
	migrate -path migrations -database "$(DB_CONNECTION)" up 

migrate-down:
	migrate -path migrations -database "$(DB_CONNECTION)" down