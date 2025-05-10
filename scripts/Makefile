.PHONY: run help all gqlgen gqlgen-safe

help:
	@ECHO +------------------------------------------------------+
	@ECHO ^|        Gamler DeFi Project - Make Commands           ^|
	@ECHO +------------------------------------------------------+
	@ECHO   ^> make run            - Start the application server
	@ECHO   ^> make build          - Compile the application into a binary
	@ECHO   ^> make run-build      - Run the compiled application binary
	@ECHO   ^> make swagger        - Generate Swagger API documentation
	@ECHO   ^> make mongo-up       - Launch MongoDB container using Docker
	@ECHO   ^> make mongo-down     - Stop and remove MongoDB container
	@ECHO   ^> make gqlgen         - Generate GraphQL code from schema
	@ECHO   ^> make gqlgen-safe    - Generate GraphQL code with dependency check
	@ECHO   ^> make help           - Display this help information
	@ECHO   ^> make all            - Start the application server (same as 'make run')
	
run:
	go run ../core/main.go

build:
	go build -o ../bin/main ../core/main.go

run-build:
	../bin/main

swagger:
	cd ../ && swag init -g core/main.go --parseDependency --parseInternal 
#--parseDepth 2 на всякий случай

mongo-up:
	docker-compose -f ../docker/docker.mongo.compose.yml up -d

mongo-down:
	docker-compose -f ../docker/docker.mongo.compose.yml down --rmi all -v
# docker system prune -f

gqlgen:
	go run -mod=mod github.com/99designs/gqlgen generate --config ../gqlgen.yml

gqlgen-safe:
	go get github.com/99designs/gqlgen
	go run github.com/99designs/gqlgen generate --config ../gqlgen.yml

all: run