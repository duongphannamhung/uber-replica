package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"uber-replica/util"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Print("cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Print("cannot connect to db:", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
