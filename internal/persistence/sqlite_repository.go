package persistence

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type SqliteRepository struct {
	db *sql.DB
}

func NewSqliteRepository(db string) *SqliteRepository {
	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	s := SqliteRepository{}
	s.db = sqliteDatabase
	return &s

}
func (s SqliteRepository) GetTreasureByID(treasureID int64) Treasure {
	query := "SELECT * FROM Treasure Where treasureID ='" + fmt.Sprint(treasureID) + "'"
	row, err := s.db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var id int64
	var tType string
	var description string

	for row.Next() { // Iterate and fetch the records from result cursor
		row.Scan(&id, &tType, &description)
	}
	return Treasure{
		TreasureID:  id,
		Type:        tType,
		Description: description,
	}
}

func (s SqliteRepository) ListTreasure(params map[string]string) []Treasure {
	query := "SELECT * FROM Treasure WHERE "
	for col, val := range params {
		query += col + "='" + val + "' AND "
	}
	query = query[:len(query)-5]
	row, err := s.db.Query(query)
	if err != nil {
		log.Fatal(err, query)
	}
	defer row.Close()
	treasures := make([]Treasure, 0)
	for row.Next() { // Iterate and fetch the records from result cursor
		treasure := Treasure{}
		row.Scan(&treasure.TreasureID, &treasure.Type, &treasure.Description)
		treasures = append(treasures, treasure)
	}
	return treasures
}
