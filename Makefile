.PHONY: generate-oas

generate-oas:
	ogen --clean --package oas --target internal/oas _oas/openapi.yml 