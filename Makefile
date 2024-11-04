start:
	docker compose up --build

stop:
	docker compose down -v

migrate:
	migrate -source file://db/migrations -database postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT_OUTER)/$(DB_NAME)?sslmode=$(DB_SSLMODE) up
	
	#если не видит переменные из .env используй
	#Для linux
	#sudo -E $(cat .env | xargs) make migrate
	#Для Windows
	#Get-Content .env | ForEach-Object {
    #	if ($_ -match '^(.*)=(.*)$') {
    #		[System.Environment]::SetEnvironmentVariable($matches[1], $matches[2])
    #	}
	#}
	#make migrate