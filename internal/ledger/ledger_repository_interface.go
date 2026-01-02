package ledger

type Repository interface {
    WithTransaction(fn func(TxRepository) error) error
    // Append(entry Entry)
    EntriesForAccount(accountNumber string) []Entry
    AllEntries() []Entry 
}

type TxRepository interface { 
    HasIdempotencyKey(key string) (bool, error) 
    RecordIdempotencyKey(key string) error 
    AppendEntry(entry Entry) error
    Balance(accountNumber string) (float64, error)
 }