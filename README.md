# todo-list-app

## Installation
- [GoLang](https://go.dev/doc/install)
- [PostgreSql](https://www.postgresql.org)
- [Docker](https://docs.docker.com/engine/install/)

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

### Docker
run `sudo docker build -t todo-app .` to build an image container

run `sudo docker run --network host todo-app` to launch app


### Urls
 - [Postman](https://api.postman.com/collections/459354-d9a68bfc-5acf-4755-9ae3-22b6b106b1d8?access_key=PMAT-01HJ64NV55Q2R8ZF3C8R8RR1MG)
 - `GET "/rest"` RootIndex
 - `GET "/rest/task/offset"` GetTaskOffset
 - `GET "/rest/task/status"` GetTaskStatusList
 - `POST "/rest/task"` CreateTask
 - `GET "/rest/task/:id"` GetTask
 - `PUT "/rest/task/:id"` EditTask
 - `PUT "/rest/task/:id/start_progress"` StartTaskProgress
 - `PUT "/rest/task/:id/pause"` PauseTask
 - `PUT "/rest/task/:id/done"` DoneTask
 - `DELETE "/rest/task/:id"` DeleteTask // only changes status to 'deleted'
 - `PUT "/rest/task/:id/restore"` RestoreTask
 - `DELETE "/rest/task/:id/completely"` DeleteTaskCompletely
 - `DELETE "/rest/task/free_trash"` FreeTaskTrash
 
