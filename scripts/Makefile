.PHONY: run help all gqlgen gqlgen-safe

help:
	@ECHO +------------------------------------------------------+
	@ECHO ^|        Gamler DeFi Project - Make Commands		      ^|
	@ECHO +------------------------------------------------------+
	@ECHO   ^> make run        - Run the application server
	@ECHO                       (executes: go run ../core/main.go)
	@ECHO   ^> make build      - Build the application server
	@ECHO                       (executes: go build -o ../bin/main ../core/main.go)
	@ECHO   ^> make run-build  - Run the built application server
	@ECHO                       (executes: ../bin/main)
	@ECHO   ^> make swagger    - Generate Swagger documentation
	@ECHO                       (executes: swag init -g core/main.go)
	@ECHO   ^> make gqlgen     - Generate GraphQL code from schema
	@ECHO                       (executes: go run github.com/99designs/gqlgen generate)
	@ECHO   ^> make gqlgen-safe - Generate GraphQL code with dependency check
	@ECHO                       (ensures gqlgen is available first)
	@ECHO                       Use this if 'make gqlgen' fails with missing dependency
	@ECHO   ^> make help       - Display this help information
	@ECHO   ^> make all        - Same as 'make run'

run:
	go run ../core/main.go

build:
	go build -o ../bin/main ../core/main.go

run-build:
	../bin/main

swagger:
	cd ../ && swag init -g core/main.go --parseDependency --parseInternal 
#--parseDepth 2 на всякий случай

gqlgen:
	go run -mod=mod github.com/99designs/gqlgen generate --config ../gqlgen.yml

gqlgen-safe:
	go get github.com/99designs/gqlgen
	go run github.com/99designs/gqlgen generate --config ../gqlgen.yml

all: run