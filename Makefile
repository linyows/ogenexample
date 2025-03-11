default: gen

server:
	go run cmd/server/main.go

client:
	go run cmd/client/main.go

gen:
	go generate ./...

mysql:
	@mysql -uroot -e "SHOW DATABASES LIKE 'ogenexample';" | grep 'ogenexample' > /dev/null || mysql -uroot -e "CREATE DATABASE ogenexample;"
	mysqldef ogenexample -uroot < ./db/schema.sql
