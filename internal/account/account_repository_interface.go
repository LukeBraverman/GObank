package account

type Repository interface {
	CreateAccount(account *TokenAccount) error
	GetAccount(accountNumber string) (*TokenAccount, bool, error)
}
