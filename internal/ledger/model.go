package ledger

import "time"


type Entry struct {
    TransactionID string
    IdempotencyKey  string
    AccountNumber string
    Amount        float64
    Description   string
    CreatedAt     time.Time
}
