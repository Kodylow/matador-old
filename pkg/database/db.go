package database

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

var DBCon *sql.DB

// InitDatabase creates a new database file with a tokens table
func InitDatabase() error {
	db, err := sql.Open("sqlite3", "./tokens.db") // Open the database
	if err != nil {
		return err
	}

	DBCon = db

	sqlTable := `CREATE TABLE IF NOT EXISTS tokens(
		Token TEXT NOT NULL PRIMARY KEY,
		Spent BOOLEAN NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_tokens ON tokens (Token);`

	_, err = DBCon.Exec(sqlTable) // Execute the SQL Statement
	if err != nil {
		return err
	}

	return nil
}

// AddToken adds a new token to the database and marks it as unspent
func AddToken(token string) error {
	sqlAdd := `INSERT INTO tokens(Token, Spent) VALUES (?, ?)`
	statement, err := DBCon.Prepare(sqlAdd) // Prepare SQL Statement
	if err != nil {
		return err
	}
	_, err = statement.Exec(token, 0)
	if err != nil {
		return err
	}

	return nil
}

// GetToken retrieves the spent status of a token
func GetToken(token string) (bool, error) {
	sqlGet := `SELECT Spent FROM tokens WHERE Token = ?`
	row := DBCon.QueryRow(sqlGet, token)

	var spent bool
	err := row.Scan(&spent)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no rows are returned, return a custom error message
			return false, errors.New("Token not found")
		}
		// If any other error is returned, return it directly
		return false, err
	}

	return spent, nil
}

// UpdateToken updates the spent status of a token in the database
func UpdateToken(token string, spent bool) error {
	sqlUpdate := `UPDATE tokens SET Spent = ? WHERE Token = ?`
	statement, err := DBCon.Prepare(sqlUpdate) // Prepare SQL Statement
	if err != nil {
		return err
	}
	_, err = statement.Exec(spent, token)
	if err != nil {
		return err
	}

	return nil
}
