package main

import (
	"log"

	// "github.com/duongphannamhung/go-bank/api"
	// db "github.com/duongphannamhung/go-bank/db/sqlc"
	"github.com/duongphannamhung/uber-replica/backend/util"
	// "github.com/techschool/simplebank/api"
	// db "github.com/techschool/simplebank/db/sqlc"
	// _ "github.com/lib/pq"
)

func main() {
	_, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	// conn, err := sql.Open(config.DBDriver, config.DBSource)
	// if err != nil {
	// 	log.Fatal("cannot connect to db:", err)
	// }
	// store := db.NewStore(conn)
	// server := api.NewServer(store)
	// err = server.Start(config.ServerAddress)
	// if err != nil {
	// 	log.Fatal("cannot start server:", err)
	// }
}
