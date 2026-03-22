include .env
export

migrate-up:
	migrate -path migrations -database "$(DB_CONNECTION)" up 

migrate-down:
	migrate -path migrations -database "$(DB_CONNECTION)" down