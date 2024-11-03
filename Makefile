start:
	docker compose up --build

stop:
	docker compose down -v

create-env:
	powershell.exe -noprofile -command "Get-Content .env | ForEach-Object {if ($_ -match '^(.*)=(.*)$') {[System.Environment]::SetEnvironmentVariable($matches[1], $matches[2], 'Process')}}"

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

migrate-windows:
	# Загружаем переменные из .env и устанавливаем их для PowerShell
	powershell -Command "Get-Content .env | ForEach-Object {if ($$_ -match '^(.*)=(.*)$$') {[System.Environment]::SetEnvironmentVariable($$matches[1], $$matches[2])}}; \
	migrate -source file://db/migrations -database postgres://$env:DB_USER:$env:DB_PASSWORD@localhost:$env:DB_PORT/$env:DB_NAME?sslmode=$env:DB_SSLMODE up 1"


drop:
	migrate -source file://db/migrations -database postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE) down