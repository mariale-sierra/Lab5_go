package main

import (
	
	"log"
	"net"
	"database/sql"
    _ "github.com/mattn/go-sqlite3"
	"Lab5_go/handlers"
    "fmt"
    "path/filepath"
)

func main() {

    db, err := sql.Open("sqlite3", "./db/series.db")
    if err != nil {
        log.Fatal(err)
    }
    absPath, _ := filepath.Abs("./db/series.db")
    fmt.Println("Using DB at:", absPath)
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

