package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	_ "modernc.org/sqlite"
	"database/sql"
)

func main() {

    db, err := sql.Open("sqlite", "file:series.db")
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

        go handle(conn, db)
    }
}

func handle(conn net.Conn, db *sql.DB) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	fmt.Println("Request:", requestLine)

	parts := strings.Fields(requestLine)
	if len(parts) < 2 {
		return
	}

	path := parts[1]

	for {
		line, err := reader.ReadString('\n')
		if err != nil || line == "\r\n" {
			break
		}
	}

	var body string
	statusLine := "HTTP/1.1 200 OK\r\n"

	if path == "/" {

		rows, err := db.Query("SELECT id, name, current_episode, total_episodes FROM series")
		if err != nil {
			log.Println(err)
			return
		}
		defer rows.Close()

		body = `
		<html>
		<head>
			<title>Series Tracker</title>
			<style>
				body { font-family: Arial, sans-serif; }
				table { border-collapse: collapse; width: 60%; margin: auto; }
				th, td { border: 1px solid black; padding: 8px; text-align: center; }
				th { background-color: #fba6bd; }
				h1 { text-align: center; }
			</style>
		</head>
		<body>
		<h1>Series Tracker</h1>
		<table>
		<tr>
			<th>#</th>
			<th>Name</th>
			<th>Current</th>
			<th>Total</th>
			<th>Completa</th>
		</tr>
		`

		for rows.Next() {
			var id int
			var name string
			var current int
			var total int

			err := rows.Scan(&id, &name, &current, &total)
			if err != nil {
				log.Println(err)
				continue
			}

			body += fmt.Sprintf(
				"<tr><td>%d</td><td>%s</td><td>%d</td><td>%d</td><td><button onclick=\"markComplete(this)\">Complete</button></td></tr>",
				id, name, current, total,
			)
		}

		body += `
		</table>
		<script>
		function markComplete(button) {
			const row = button.parentElement.parentElement;
			row.style.backgroundColor = "#c8f7c5";
		}
		</script>
		</body>
		</html>
		`

	} else {
		statusLine = "HTTP/1.1 404 Not Found\r\n"
		body = "<html><body><h1>404 Not Found</h1></body></html>"
	}

	response := fmt.Sprintf(
		"%sContent-Type: text/html\r\nContent-Length: %d\r\n\r\n%s",
		statusLine,
		len(body),
		body,
	)

	conn.Write([]byte(response))
}