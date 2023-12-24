build:
	docker compose build todo-app

run:
	docker compose up todo-app

test:
	go test -v

migrate:
	docker exec -i todo_app_db psql -U postgres -W postgres -d todo < init.sql

swag:
	swag init

todo:
	go run main.go

