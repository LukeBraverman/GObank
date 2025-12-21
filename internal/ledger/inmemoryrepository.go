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

func (r *InMemoryRepository) AllEntries() []Entry {
    r.mu.RLock()
    defer r.mu.RUnlock()

    // return a copy to avoid mutation
    entries := make([]Entry, len(r.entries))
    copy(entries, r.entries)
    return entries
}

