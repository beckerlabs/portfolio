package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsnUser := flag.String("dsnUser", "web", "MySQL data source user")
	dsnPass := flag.String("dsnPass", "password", "MySQL data source password")
	flag.Parse()

	dsn := fmt.Sprintf("%s:%s@/beckerlabs?parseTime=true", *dsnUser, *dsnPass)
	db, err := openDb(dsn)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	mux := routes()

	err = http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}

func openDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
