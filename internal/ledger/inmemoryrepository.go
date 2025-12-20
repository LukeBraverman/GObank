package ledger

import "sync"

type InMemoryRepository struct {
    mu      sync.RWMutex
    entries []Entry
}

func NewRepository() *InMemoryRepository {
    r := &InMemoryRepository{
        entries: make([]Entry, 0),
    }

    return r
}


func (r *InMemoryRepository) Append(entry Entry) {
    r.entries = append(r.entries, entry)
}

func (r *InMemoryRepository) EntriesForAccount(accountNumber string) []Entry {
    var result []Entry
    for _, e := range r.entries {
        if e.AccountNumber == accountNumber {
            result = append(result, e)
        }
    }
    return result
}
