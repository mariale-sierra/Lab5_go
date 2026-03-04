package templates

import (
	"fmt"
	"Lab5_go/db"
)

func RenderHome(series []db.Series) string {
	body := `<html>
	<head>
	<meta charset="UTF-8">
	<title>Series Tracker</title>
	<link rel="stylesheet" href="/static/styles.css">
	</head>
	<body>
	<h1>Series Tracker</h1>
	<table>
	<tr>
	<th>#</th>
	<th>Name</th>
	<th>Current</th>
	<th>Total</th>
	<th>Complete</th>
	</tr>`

	for _, s := range series {
	body += fmt.Sprintf(
		"<tr>" +
			"<td>%d</td>" +
			"<td>%s</td>" +
			"<td>%d</td>" +
			"<td>%d</td>" +
			"<td><button onclick=\"markComplete(this)\">Complete</button></td>" +
		"</tr>",
		s.ID, s.Name, s.Current, s.Total,
	)
}

	body += `
	</table>
	<script src="/static/script.js"></script>
	</body>
	</html>`

	return body
}