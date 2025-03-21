.PHONY: generate-oas

generate-oas:
	ogen --clean --package oas --target internal/oas _oas/openapi.yml 

generate-db:
	jet -source=sqlite -dsn="../database/database.db" -path=./internal/db
