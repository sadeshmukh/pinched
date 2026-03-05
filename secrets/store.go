package secrets

import (
	"database/sql"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

func GetDB() *sql.DB {
	db, err := sql.Open("sqlite", "secrets.db")
	if err != nil {
		log.Fatal("failed to open secrets database:", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS secrets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE,
		value TEXT
	);
	`)
	if err != nil {
		log.Fatal("failed to create secrets table:", err)
	}

	return db
}

func StoreSecret(name, value string) error {
	db := GetDB()
	defer db.Close()

	_, err := db.Exec("INSERT OR REPLACE INTO secrets (name, value) VALUES (?, ?)", name, value)
	return err
}

func GetSecret(name string) (string, error) {
	db := GetDB()
	defer db.Close()

	var value string
	err := db.QueryRow("SELECT value FROM secrets WHERE name = ?", name).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}
