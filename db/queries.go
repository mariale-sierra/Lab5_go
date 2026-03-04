package db

import "database/sql"

func GetAllSeries(database *sql.DB) ([]Series, error) {

	rows, err := database.Query(
		"SELECT id, name, current_episode, total_episodes FROM series",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Series

	for rows.Next() {
		var s Series
		err := rows.Scan(&s.ID, &s.Name, &s.Current, &s.Total)
		if err != nil {
			continue
		}
		list = append(list, s)
	}

	return list, nil
}