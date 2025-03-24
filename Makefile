gen-oas:
	ogen --clean --package oas --target internal/oas _oas/openapi.yml 

gen-db:
	jet -source=sqlite -dsn="../database/database.db" -path=./internal/db
