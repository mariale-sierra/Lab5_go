package main

import (
	
	"log"
	"net"
	_ "modernc.org/sqlite"
	"database/sql"
	"Lab5_go/handlers"
)

func main() {

    db, err := sql.Open("sqlite", "file:db/series.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatal(err)
    }
    defer listener.Close()

    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
		go handlers.Handle(conn, db)
    }
}

