package ogenexample

//go:generate go run github.com/ogen-go/ogen/cmd/ogen --target oas/oasgen --package oasgen --clean oas/petstore.yml
//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc generate -f db/sqlc.yml
