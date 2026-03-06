package handlers

import (
	"bufio"
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"

	"Lab5_go/db"
	"Lab5_go/templates"
)

func Handle(conn net.Conn, database *sql.DB) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// leer primera línea (GET / HTTP/1.1)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	parts := strings.Fields(requestLine)
	if len(parts) < 2 {
		return
	}

	method := parts[0]
	path := parts[1]

	// leer headers y capturar Content-Length si existe
	var contentLength int

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		if strings.HasPrefix(line, "Content-Length:") {
			fmt.Sscanf(line, "Content-Length: %d", &contentLength)
		}

		if line == "\r\n" {
			break
		}
	}

	// static
	if strings.HasPrefix(path, "/static/") {
		ServeStatic(conn, path)
		return
	}

	// GET /create → mostrar formulario
	if path == "/create" && method == "GET" {
		handleCreateForm(conn)
		return
	}

	// POST /create → procesar formulario
	if path == "/create" && method == "POST" {
		handleCreatePost(conn, reader, contentLength, database)
		return
	}

	// POST /update?id=3 → actualizar episodio actual
	if strings.HasPrefix(path, "/update") && method == "POST" {
    handleUpdate(conn, path, database)
    return
	}

	// POST /rate?id=3&rating=5 → actualizar rating
	if path == "/rate" && method == "POST" {
    handleRate(conn, reader, contentLength, database)
    return
	}


	// body
	var body string
	statusLine := "HTTP/1.1 200 OK\r\n"

	switch path {

	case "/":
		series, err := db.GetAllSeries(database)
		if err != nil {
		fmt.Println("DB ERROR:", err)
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

func handleCreateForm(conn net.Conn) {

	body := `<html>
	<head>
	<meta charset="UTF-8">
	<title>Create Series</title>
	<link rel="stylesheet" href="/static/styles.css">
	</head>
	<body>

	<h1>Create Series</h1>

	<form method="POST" action="/create" class="create-form">

		<div class="form-group">
			<label>Title</label>
			<input type="text" name="series_name" required>
		</div>

		<div class="form-group">
			<label>Current Episode</label>
			<input type="number" name="current_episode" min="1" value="1" required>
		</div>

		<div class="form-group">
			<label>Total Episodes</label>
			<input type="number" name="total_episodes" min="1" required>
		</div>

		<button type="submit">Create</button>

	</form>

	<a href="/">Back</a>

	</body>
	</html>`

	response := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nContent-Length: %d\r\n\r\n%s",
		len(body),
		body,
	)

	conn.Write([]byte(response))
}

func handleCreatePost(conn net.Conn, reader *bufio.Reader, contentLength int, database *sql.DB) {

	// leer body
	bodyBytes := make([]byte, contentLength)
	_, err := reader.Read(bodyBytes)
	if err != nil {
		return
	}

	body := string(bodyBytes)

	// parsear body 
	values, err := url.ParseQuery(body)
	if err != nil {
		return
	}

	name := values.Get("series_name")
	currentStr := values.Get("current_episode")
	totalStr := values.Get("total_episodes")

	current, _ := strconv.Atoi(currentStr)
	total, _ := strconv.Atoi(totalStr)

	// insertar en bd
	_, err = database.Exec(
		"INSERT INTO series (name, current_episode, total_episodes) VALUES (?, ?, ?)",
		name, current, total,
	)
	if err != nil {
		return
	}

	// redirigir (POST/Redirect/GET)
	response := "HTTP/1.1 303 See Other\r\n"
	response += "Location: /\r\n\r\n"

	conn.Write([]byte(response))
}

func handleUpdate(conn net.Conn, path string, database *sql.DB) {

    // separar /update?id=3
    parts := strings.SplitN(path, "?", 2)

    if len(parts) < 2 {
        return
    }

    // parsear parámetros
    params, _ := url.ParseQuery(parts[1])
    id := params.Get("id")

    fmt.Println("Updating ID:", id)

    // ejecutar update
    _, err := database.Exec(
        `UPDATE series
         SET current_episode = current_episode + 1
         WHERE id = ? AND current_episode < total_episodes`,
        id,
    )

    if err != nil {
        fmt.Println("DB error:", err)
        return
    }

    // responder OK
    response := "HTTP/1.1 200 OK\r\n"
    response += "Content-Type: text/plain\r\n\r\n"
    response += "ok"

    conn.Write([]byte(response))
}

func handleRate(conn net.Conn, reader *bufio.Reader, contentLength int, database *sql.DB) {

    bodyBytes := make([]byte, contentLength)
    reader.Read(bodyBytes)
    body := string(bodyBytes)

    values, _ := url.ParseQuery(body)

    id := values.Get("id")
    rating := values.Get("rating")

    database.Exec(
        "INSERT OR REPLACE INTO ratings (series_id, rating) VALUES (?, ?)",
        id, rating,
    )

    response := "HTTP/1.1 303 See Other\r\n"
    response += "Location: /\r\n\r\n"

    conn.Write([]byte(response))
}

//static
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

	conn.Write([]byte(response))
	conn.Write(content)
}
