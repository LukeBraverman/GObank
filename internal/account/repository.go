package account

type repository struct {
    // later this will hold a DB or map
}

func NewRepository() *repository {
    return &repository{}
}

// Temporary test method
func (r *repository) GetTestAccount() *BankAccount {
    return &BankAccount{
        AccountNumber: "12345678",
        Balance:       1000.50,
        Name:          "Test User",
    }
}
