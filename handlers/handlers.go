package handlers

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"database/sql"
	"Lab5_go/db"
	"Lab5_go/templates"
	"os"
)

func Handle(conn net.Conn, database *sql.DB) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	parts := strings.Fields(requestLine)
	if len(parts) < 2 {
		return
	}

	path := parts[1]

	for {
		line, _ := reader.ReadString('\n')
		if line == "\r\n" {
			break
		}
	}

	var body string
	statusLine := "HTTP/1.1 200 OK\r\n"

	if strings.HasPrefix(path, "/static/") {
    ServeStatic(conn, path)
	return
	}

	switch path {

	case "/":
		series, err := db.GetAllSeries(database)
		if err != nil {
			body = "<h1>Error interno</h1>"
			statusLine = "HTTP/1.1 500 Internal Server Error\r\n"
			break
		}

		body = templates.RenderHome(series)

	default:
		statusLine = "HTTP/1.1 404 Not Found\r\n"
		body = "<h1>404 Not Found</h1>"
	}

	response := fmt.Sprintf(
		"%sContent-Type: text/html\r\nContent-Length: %d\r\n\r\n%s",
		statusLine,
		len(body),
		body,
	)

	conn.Write([]byte(response))
}

func ServeStatic(conn net.Conn, path string) {

    filePath := "." + path 

    content, err := os.ReadFile(filePath)
    if err != nil {
        response := "HTTP/1.1 404 Not Found\r\n\r\n"
        conn.Write([]byte(response))
        return
    }

    contentType := "text/plain"

    if strings.HasSuffix(path, ".css") {
        contentType = "text/css"
    } else if strings.HasSuffix(path, ".js") {
        contentType = "application/javascript"
    }

    response := fmt.Sprintf(
        "HTTP/1.1 200 OK\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n",
        contentType,
        len(content),
    )
	fmt.Println("PATH:", path)

    conn.Write([]byte(response))
    conn.Write(content)
}