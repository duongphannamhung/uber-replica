createdb:
	docker compose exec -it uber-replica-db createdb --username=root --owner=root uber_replica
dropdb:
	docker compose exec -it uber-replica-db dropdb uber_replica

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/uber_replica?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/uber_replica?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/uber_replica?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/uber_replica?sslmode=disable" -verbose down 1

migrateup2:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/uber_replica?sslmode=disable" -verbose up 2

migratedown2:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/uber_replica?sslmode=disable" -verbose down 2

test_osrm:
	curl "http://localhost:8282/route/v1/driving/108.2126533032648,16.060370207645025;108.23284753838242,16.061554708431885"