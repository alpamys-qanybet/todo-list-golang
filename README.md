# todo-list-app

## Installation
- [GoLang](https://go.dev/doc/install)
- [PostgreSql](https://www.postgresql.org)

### Setup

Clone this repository into `$GOPATH/src` directory, by default it is located in `~/go`. Or set `$GOPATH` in your bash profile.

install dependencies:

`go mod tidy`

`go mod vendor`

run `go run main.go`.

build `go build main.go`.

you can change environment variables in `.env`.
env sample:
```
SERVER_HOST=localhost
# optional, default: "", empty string is localhost
SERVER_PORT=9292
# optional, default: 9292

DATABASE_URL=postgresql://postgres:postgres@localhost:5432/todo
# optional, default: postgresql://postgres:postgres@localhost:5432/todo

APP_SECRET=kj3mSJbsw4lpFWUsHasQZf9r
# required

DEBUG=true
# optional, default false
```