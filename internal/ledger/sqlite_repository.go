package ledger

import (
    "database/sql"
    "time"

    _ "github.com/mattn/go-sqlite3"
)


type SQLiteRepository struct {
    db *sql.DB
}

func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }

    schema := `
    CREATE TABLE IF NOT EXISTS ledger_entries (
        transaction_id   TEXT NOT NULL,
        idempotency_key  TEXT NOT NULL,
        account_number   TEXT NOT NULL,
        amount           REAL NOT NULL,
        description      TEXT NOT NULL,
        created_at       TEXT NOT NULL
    );
    `

    if _, err := db.Exec(schema); err != nil {
        return nil, err
    }

    return &SQLiteRepository{db: db}, nil
}

func (r *SQLiteRepository) Append(entry Entry) {
    _, err := r.db.Exec(
        `
        INSERT INTO ledger_entries (
            transaction_id,
            idempotency_key,
            account_number,
            amount,
            description,
            created_at
        ) VALUES (?, ?, ?, ?, ?, ?)
        `,
        entry.TransactionID,
        entry.IdempotencyKey,
        entry.AccountNumber,
        entry.Amount,
        entry.Description,
        entry.CreatedAt.UTC().Format(time.RFC3339Nano),
    )
    if err != nil {
        panic(err) // fine for now; later return error
    }
}

func (r *SQLiteRepository) EntriesForAccount(accountNumber string) []Entry {
    rows, err := r.db.Query(
        `
        SELECT transaction_id, idempotency_key, account_number,
               amount, description, created_at
        FROM ledger_entries
        WHERE account_number = ?
        `,
        accountNumber,
    )
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    var entries []Entry
    for rows.Next() {
        var e Entry
        var createdAt string

        if err := rows.Scan(
            &e.TransactionID,
            &e.IdempotencyKey,
            &e.AccountNumber,
            &e.Amount,
            &e.Description,
            &createdAt,
        ); err != nil {
            panic(err)
        }

        e.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdAt)
        entries = append(entries, e)
    }

    return entries
}

func (r *SQLiteRepository) AllEntries() []Entry {
    rows, err := r.db.Query(`
        SELECT transaction_id, idempotency_key, account_number,
               amount, description, created_at
        FROM ledger_entries
    `)
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    var entries []Entry
    for rows.Next() {
        var e Entry
        var createdAt string

        rows.Scan(
            &e.TransactionID,
            &e.IdempotencyKey,
            &e.AccountNumber,
            &e.Amount,
            &e.Description,
            &createdAt,
        )
        e.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdAt)
        entries = append(entries, e)
    }

    return entries
}
