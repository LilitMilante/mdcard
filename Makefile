local-up: local-down
	docker run --name mdcard_db -p 8080:5432 --network mdcard -e POSTGRES_PASSWORD=dev \
 	-e POSTGRES_USER=dev -e POSTGRES_DB=mdcard -d postgres:15.2-alpine

local-down:
	docker rm -f mdcard_db

migrate-up:
	docker run -v /Users/lilit/GolandProjects/medical-card/migrations:/migrations --network mdcard migrate/migrate \
        -path=/migrations/ -database postgres://dev:dev@mdcard_db:5432/mdcard?sslmode=disable up

migrate-new:
	docker run -v /Users/lilit/GolandProjects/medical-card/migrations:/migrations --network mdcard migrate/migrate create -ext sql -dir /migrations "$(name)"

up-test-db: down-test-db
	docker run --name mdcard_db_test -p 8181:5432 --network mdcard -e POSTGRES_PASSWORD=dev \
     	-e POSTGRES_USER=dev -e POSTGRES_DB=mdcard -d postgres:15.2-alpine

down-test-db:
	docker rm -f mdcard_db_test
