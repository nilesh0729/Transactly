Container:
	docker run --name jansi -p 5433:5432 -e POSTGRES_PASSWORD=SituBen -e POSTGRES_USER=root -d postgres

Createdb:
	docker exec -it jansi createdb --username=root --owner=root Hiten

Dropdb:
	docker exec -it jansi dropdb -U root Hiten

MigrateUp:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up

MigrateDown:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down

Sqlc:
	sqlc generate

Test:
	go test -v -cover ./...

Server:
	go run main.go

Mock:
	mockgen -package mockDB -destination db/Mock/Store.go github.com/nilesh0729/OrdinaryBank/db/Result Store

.PHONY:	Container	Createdb	Dropdb	MigrateDown	MigrateUp	Sqlc	Test	Server	Mock