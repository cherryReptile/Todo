include ${PWD}/.env
export

up:
	docker-compose up -d && make log
down:
	docker-compose down
exec:
	docker-compose exec app bash
migrate:
	docker-compose exec app migrate create -ext sql -dir internal/database/migrations ${name}
migrate.up:
	docker-compose exec app migrate -database "sqlite3://tmp/db.db?query" -path internal/database/migrations up
migrate.down:
	docker-compose exec app migrate -database "sqlite3://tmp/db.db?query" -path internal/database/migrations down
exec.root:
	docker-compose exec -u root app bash
log:
	docker-compose logs -f app