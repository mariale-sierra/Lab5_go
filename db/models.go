package db

import "database/sql"

type Series struct {
	ID      int
	Name    string
	Current int
	Total   int
    Rating  sql.NullInt64
}