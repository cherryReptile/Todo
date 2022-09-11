include ${PWD}/.env
export

USER:=$(shell echo $USER)
GROUP:=$(shell id -g)

build.app:
	docker compose run --rm app sh -c "go build -o bin/todo_bot cmd/main.go"
deploy.server:
	make build.app
	ansible-playbook -i deploy/hosts.yml deploy/server.yml -t deploy -e @deploy/vars/prod.yml -e "USER=$(USER)" -e "GROUP=$(GROUP)" --ask-vault-pass
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