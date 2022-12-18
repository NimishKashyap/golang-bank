package main

import (
	"database/sql"
	"log"

	"github.com/NimishKashyap/simplebank/api"
	db "github.com/NimishKashyap/simplebank/db/sqlc"
	"github.com/NimishKashyap/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Couldn't connect to database", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("CANNOT START SERVER", err)
	}
}
