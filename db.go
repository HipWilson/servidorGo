package main

import (
	"database/sql"
	"log"
)

type Series struct {
	ID             int
	Name           string
	CurrentEpisode int
	TotalEpisodes  int
}

func initDB(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS series (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		current_episode INTEGER NOT NULL DEFAULT 1,
		total_episodes INTEGER NOT NULL
	)`)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}
}

func getAllSeries(db *sql.DB) ([]Series, error) {
	rows, err := db.Query("SELECT id, name, current_episode, total_episodes FROM series")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Series
	for rows.Next() {
		var s Series
		err := rows.Scan(&s.ID, &s.Name, &s.CurrentEpisode, &s.TotalEpisodes)
		if err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, nil
}

func insertSeries(db *sql.DB, name string, current int, total int) error {
	_, err := db.Exec(
		"INSERT INTO series (name, current_episode, total_episodes) VALUES (?, ?, ?)",
		name, current, total,
	)
	return err
}

func incrementEpisode(db *sql.DB, id string) error {
	_, err := db.Exec(
		"UPDATE series SET current_episode = current_episode + 1 WHERE id = ? AND current_episode < total_episodes",
		id,
	)
	return err
}

func decrementEpisode(db *sql.DB, id string) error {
	_, err := db.Exec(
		"UPDATE series SET current_episode = current_episode - 1 WHERE id = ? AND current_episode > 1",
		id,
	)
	return err
}

func deleteSeries(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM series WHERE id = ?", id)
	return err
}
