package game

import "database/sql"

func make_word(db *sql.DB, c chan<- string) {
	defer db.Close()
	var cur string
	query := "SELECT word FROM nouns ORDER BY rand() LIMIT 1"
	err := db.QueryRow(query).Scan(&cur)
	if err != nil {
		panic(err)
	}
	c <- cur
}
