createdb:
	docker compose exec -it uber-replica-db createdb --username=root --owner=root uber_replica
dropdb:
	docker compose exec -it uber-replica-db dropdb uber_replica

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/uber_replica?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/uber_replica?sslmode=disable" -verbose down