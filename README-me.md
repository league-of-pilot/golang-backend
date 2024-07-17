dbdiagram.io -> export to PostgreSQL

# Install course tools

Setup sqlc theo docs docker style
https://docs.sqlc.dev/en/latest/overview/install.html#docker
docker run --rm -v $(pwd):/src -w /src sqlc/sqlc generate

install TablePlus

https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

```bash
docker run --rm -v $(pwd):/src -w /src sqlc/sqlc generate
# one time run with -rm to remove container after run

docker pull postgres:12-alpine
docker run --name gostgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
docker exec -it gostgres12 psql -U root

```

# Exec inside container

```bash
select now();
\q
# exit

# enter docker shell
# inside gostgres also tool of postgres
docker exec -it gostgres12 /bin/sh
createdb --username=root --owner=root simple_bank
psql simple_bank
dropdb simple_bank

# or exec from outside
docker exec -it gostgres12 createdb --username=root --owner=root simple_bank
docker exec -it gostgres12 psql -U root simple_bank

```

## Create Makefile to store command for reuse

to run, window ko thấy build được

```bash
make [cmd-defined-in-makefile]

```

# Some more docker basic command

```bash
docker logs gostgres12
```

# Golang-migrate

```bash
migrate -version
migrate create -ext sql -dir db/migration -seq init_schema
migrate -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" up
```
