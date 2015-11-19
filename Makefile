
.PHONY: all clean

all: onion.go 
	@go build
	./onion migratedb
	./onion serve

clean:
	@go clean
	rm -f *_resource.go
	rm -f main.go
	rm -f web_service.go
	rm -f config.yaml
	rm -f Makefile
	rm -f Dockerfile
	rm -f dockerize.sh
	rm -f -r Godeps

cleandb:
	rm -r *.db

migrate:
	./onion migratedb

run:
	./onion serve

