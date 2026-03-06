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
	<th>Next Episode</th>
	<th>Rating</th>
	</tr>`

	for _, s := range series {
	var ratingDisplay string

	if s.Rating.Valid {
		ratingDisplay = fmt.Sprintf("%d", s.Rating.Int64)
	} else {
		ratingDisplay = fmt.Sprintf(`
			<form method="POST" action="/rate" class="rating-form">
				<input type="hidden" name="id" value="%d">
				<input type="number" name="rating" min="0" max="10" placeholder="0-10" required>
				<button type="submit">Rate</button>
			</form>
		`, s.ID)
	}

	body += fmt.Sprintf(
		"<tr>" +
			"<td>%d</td>" +
			"<td>%s</td>" +
			"<td>%d</td>" +
			"<td>%d</td>" +
			"<td><button onclick=\"nextEpisode(%d)\">+1</button></td>" +
			"<td>%s</td>" +
		"</tr>",
		s.ID, s.Name, s.Current, s.Total, s.ID, ratingDisplay,
	)
}

	body += `<div class="add-container">
    <a href="/create">
        <button class="add-btn">Add Series</button>
    </a>
	</div>`

	body += `
	</table>
	<script src="/static/script.js"></script>
	</body>
	</html>`

	return body
}