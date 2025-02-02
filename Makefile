DB_URL:=postgres://root:secret@localhost:5432/simple_bank?sslmode=disable
DOC_NAME:=gostgres12

postgres:
	docker run --name $DOC_NAME -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it $DOC_NAME createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it $DOC_NAME dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

# ko setup động với window powershell được, 
# đã thử Get-Location, %%cd%%, $(PWD) thì không được
DIR := D:\myCode\goback
# SQLC = docker run --rm -v $(DIR):/src -w /src sqlc/sqlc
# gosq:
# 	$(SQLC) $(ARGS)

gosql:
	docker run --rm -v $(DIR):/src -w /src sqlc/sqlc generate

gotest:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown gosql gotest