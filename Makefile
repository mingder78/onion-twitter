
.PHONY: 

cleandb:
	rm -r *.db

migrate:
	./onion migratedb

run:
	./onion serve

