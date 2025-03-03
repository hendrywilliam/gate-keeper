DB_URL=postgresql://root:secret@localhost:5432/gatekeeper?sslmode=disable

pg:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17.4-alpine3.21

start_pg:
	docker start postgres

stop_pg:
	docker stop postgres

redis:
	docker run --name redis -p 6379:6379 -d redis redis-server --save 60 1 --loglevel warning

start_redis:
	docker start redis

stop_redis:
	docker stop redis

create_db:
	docker exec -it postgres createdb --username=root --owner=root gatekeeper

create_migrate:
	migrate create -ext sql -dir db/migrations -seq $(name)

migrate_up:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migrate_down:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

.PHONY: pg start_pg stop_pg redis start_redis stop_redis create_migrate migrate_up migrate_down