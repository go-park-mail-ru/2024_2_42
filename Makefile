start:
	docker compose up --build

migrate:
	migrate -source file://db/migrations -database postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE) up 1
	#если не видит переменные из .env используй
	#sudo -E $(cat .env | xargs) make migrate 
