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

        CREATE TABLE IF NOT EXISTS idempotency_keys (
            key TEXT PRIMARY KEY,
            created_at TEXT NOT NULL
        );
    `

    if _, err := db.Exec(schema); err != nil {
        return nil, err
    }

    return &SQLiteRepository{db: db}, nil
}
// TODO: What does this function signiture mean
func (r *SQLiteRepository) WithTransaction(
    fn func(TxRepository) error,
) error {

    tx, err := r.db.Begin()
    if err != nil {
        return err
    }

    txRepo := &sqliteTxRepository{tx: tx}

    if err := fn(txRepo); err != nil {
        _ = tx.Rollback()
        return err
    }

    return tx.Commit()
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


type sqliteTxRepository struct {
    tx *sql.Tx
}


func (r *sqliteTxRepository) AppendEntry(entry Entry) error {
    _, err := r.tx.Exec(
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
    return err
}

func (r *sqliteTxRepository) HasIdempotencyKey(key string) (bool, error) {
    row := r.tx.QueryRow(
        `SELECT 1 FROM idempotency_keys WHERE key = ?`,
        key,
    )

    var dummy int
    err := row.Scan(&dummy)
    if err == sql.ErrNoRows {
        return false, nil
    }
    if err != nil {
        return false, err
    }

    return true, nil
}

func (r *sqliteTxRepository) RecordIdempotencyKey(key string) error {
    _, err := r.tx.Exec(
        `
        INSERT INTO idempotency_keys (key, created_at)
        VALUES (?, ?)
        `,
        key,
        time.Now().UTC().Format(time.RFC3339Nano),
    )
    return err
}



