gen-oas:
	ogen --clean --package api --target ./.gen/api openapi/openapi.yml 

gen-db:
	jet -source=sqlite -dsn="../database/database.db" -path=./.gen/db
