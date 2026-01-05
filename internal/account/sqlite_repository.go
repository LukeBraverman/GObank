package account

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_txlock=immediate")
	if err != nil {
		return nil, err
	}
	db.Exec("PRAGMA journal_mode=WAL;")
	db.Exec("PRAGMA busy_timeout = 5000;")

	schema := `
		CREATE TABLE IF NOT EXISTS accounts (
			account_number TEXT PRIMARY KEY,
			name TEXT NOT NULL
		);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return &SQLiteRepository{db: db}, nil
}

func (r *SQLiteRepository) CreateAccount(account *TokenAccount) error {
	_, err := r.db.Exec(
		`INSERT INTO accounts (account_number, name) VALUES (?, ?)`,
		account.AccountNumber,
		account.Name,
	)
	return err
}

func (r *SQLiteRepository) GetAccount(accountNumber string) (*TokenAccount, bool, error) {
	row := r.db.QueryRow(
		`SELECT account_number, name FROM accounts WHERE account_number = ?`,
		accountNumber,
	)

	var acc TokenAccount
	if err := row.Scan(&acc.AccountNumber, &acc.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &acc, true, nil
}
