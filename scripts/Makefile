.PHONY: run help all gqlgen gqlgen-safe

help:
	@ECHO +------------------------------------------------------+
	@ECHO ^|            Gamler DeFi Project - Make Commands       ^|
	@ECHO +------------------------------------------------------+
	@ECHO   ^> make run        - Run the application server
	@ECHO                       (executes: go run ../core/main.go)
	@ECHO   ^> make gqlgen     - Generate GraphQL code from schema
	@ECHO                       (executes: go run github.com/99designs/gqlgen generate)
	@ECHO   ^> make gqlgen-safe - Generate GraphQL code with dependency check
	@ECHO                       (ensures gqlgen is available first)
	@ECHO                       Use this if 'make gqlgen' fails with missing dependency
	@ECHO   ^> make help       - Display this help information
	@ECHO   ^> make all        - Same as 'make run'

run:
	go run ../core/main.go

gqlgen:
	go run -mod=mod github.com/99designs/gqlgen generate --config ../gqlgen.yml

gqlgen-safe:
	go get github.com/99designs/gqlgen
	go run github.com/99designs/gqlgen generate --config ../gqlgen.yml

all: run