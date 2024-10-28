start:
	docker compose up --build

migrate:
	#migrate -source file://db/migrations -database postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE) up 1
	migrate -source file://db/migrations -database "postgres://postgres:12345@localhost:5432/2024_2_42?sslmode=disable" up 1

