package account

import (
    "fmt"
    "sync"
)

type Repository struct {
    mu       sync.RWMutex
    accounts map[string]*TokenAccount
}

func NewRepository() *Repository {
    return &Repository{
        accounts: make(map[string]*TokenAccount),
    }
}

func (r *Repository) CreateAccount(account *TokenAccount) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, exists := r.accounts[account.AccountNumber]; exists {
        return fmt.Errorf("account %s already exists", account.AccountNumber)
    }

    r.accounts[account.AccountNumber] = account
    return nil
}

func (r *Repository) GetAccount(accountNumber string) (*TokenAccount, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    account, exists := r.accounts[accountNumber]
    return account, exists
}

// func (r *Repository) UpdateBalance(accountNumber string, delta float64) error {
//     r.mu.Lock()
//     defer r.mu.Unlock()

//     account, exists := r.accounts[accountNumber]
//     if !exists {
//         return fmt.Errorf("account not found")
//     }

//     account.Balance += delta
//     return nil
// }
