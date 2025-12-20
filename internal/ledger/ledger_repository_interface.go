package ledger

type Repository interface {
    Append(entry Entry)
    EntriesForAccount(accountNumber string) []Entry
}
