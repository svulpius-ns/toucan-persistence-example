package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

func main() {
	os.Remove("sqlite-database.db") // I delete the file to avoid duplicated records.
	// SQLite is a file based database.

	log.Println("Creating sqlite-database.db...")
	file, err := os.Create("sqlite-database.db") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite-database.db created")

	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	defer sqliteDatabase.Close()                                     // Defer Closing the database
	createTable(sqliteDatabase)                                      // Create Database Tables

	// INSERT RECORDS
	insertTreasure(sqliteDatabase, "gold", "3 doubloons")
	insertTreasure(sqliteDatabase, "gems", "1 emerald")
	insertTreasure(sqliteDatabase, "gems", "1 ruby")
	insertTreasure(sqliteDatabase, "jewelry", "1 pearl necklace")

	// DISPLAY INSERTED RECORDS
	displayTreasures(sqliteDatabase)
}

func createTable(db *sql.DB) {
	createTreasureTableSQL := `CREATE TABLE treasure (
		"treasureID" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"type" TEXT,
		"description" TEXT
	  );` // SQL Statement for Create Table

	log.Println("Create treasure table...")
	statement, err := db.Prepare(createTreasureTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("Treasure table created")
}

// We are passing db reference connection from main to our method with other parameters
func insertTreasure(db *sql.DB, treasureType string, description string) {
	log.Println("Inserting treasure record ...")
	insertTreasureSQL := `INSERT INTO treasure(type, description) VALUES (?, ?)`
	statement, err := db.Prepare(insertTreasureSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(treasureType, description)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func displayTreasures(db *sql.DB) {
	row, err := db.Query("SELECT * FROM Treasure ORDER BY type")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var tType string
		var description string
		row.Scan(&id, &tType, &description)
		log.Println("Treasure: ", tType, " ", description)
	}
}
