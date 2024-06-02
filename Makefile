download:
	go mod download

local_build_redis:
	go build -o woc .
	./woc -db "redis"

local_build_postgres:
	go build -o woc .
	./woc -db "postgres"

local_run_redis:
	go run server.go -db "redis"

local_run_postgres:
	go run server.go -db "postgres"

docker_build:
	docker-compose build app $(DB_TYPE)

docker:
	chmod +x start.sh
	./start.sh "$(DB_TYPE)"
