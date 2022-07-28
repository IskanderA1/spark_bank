postgres:
	docker run --name spark-bank-pg -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root250700 -d postgres:12

createdb:
	docker exec -it spark-bank-pg createdb --username=root --owner=root spark_bank

dropdb:
	docker exec -it spark-bank-pg dropdb spark_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:root250700@localhost:5432/spark_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:root250700@localhost:5432/spark_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown sqlc