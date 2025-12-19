package ledger

import "sync"

type Repository struct {
    mu      sync.RWMutex
    entries []Entry
}

func NewRepository() *Repository {
    r := &Repository{
        entries: make([]Entry, 0),
    }

    // r.entries = append(r.entries, Entry{
    //     AccountNumber: "alice",
    //     Amount: 100,
    // })

    return r
}


func (r *Repository) Append(entry Entry) {
    // r.mu.Lock()
    // defer r.mu.Unlock()

    r.entries = append(r.entries, entry)
}

func (r *Repository) EntriesForAccount(accountNumber string) []Entry {
    // r.mu.RLock()
    // defer r.mu.RUnlock()

    var result []Entry
    for _, e := range r.entries {
        if e.AccountNumber == accountNumber {
            result = append(result, e)
        }
    }
    return result
}
