package main

import (
	"database/sql"
	"log"
	"uber-replica/util"

	db "uber-replica/db/sqlc"

	"uber-replica/api"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Print("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Print("cannot connect to db:", err)
	}
	defer conn.Close()
	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Print("create server fail:", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Print("cannot start server:", err)
	}
}
