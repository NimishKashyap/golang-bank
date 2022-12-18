package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/NimishKashyap/simplebank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../..")

	if err != nil {
		log.Fatal("Couln't load cofig:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Couldn't connect to database", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())

}
