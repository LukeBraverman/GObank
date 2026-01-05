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

func (r *InMemoryRepository) EntriesForAccount(accountNumber string) ([]Entry, error) {
	var result []Entry
	for _, e := range r.entries {
		if e.AccountNumber == accountNumber {
			result = append(result, e)
		}
	}
	return result, nil
}

func (r *InMemoryRepository) AllEntries() ([]Entry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// return a copy to avoid mutation
	entries := make([]Entry, len(r.entries))
	copy(entries, r.entries)
	return entries, nil
}
