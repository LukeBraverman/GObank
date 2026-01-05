package account

import (
	"fmt"

	"github.com/LukeBraverman/GObank/internal/ledger"
)

type Service struct {
	repo   Repository
	ledger *ledger.Service
}

func NewService(
	repo Repository,
	ledgerSvc *ledger.Service,
) *Service {
	return &Service{
		repo:   repo,
		ledger: ledgerSvc,
	}
}

func (s *Service) OpenAccount(
	accountNumber string,
	name string,
) (*TokenAccount, error) {

	if accountNumber == "" {
		return nil, fmt.Errorf("account number required")
	}

	acc := &TokenAccount{
		AccountNumber: accountNumber,
		Name:          name,
	}

	if err := s.repo.CreateAccount(acc); err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *Service) GetAccount(accountNumber string) (*TokenAccount, error) {
	acc, ok, err := s.repo.GetAccount(accountNumber)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("account not found")
	}
	return acc, nil
}

func (s *Service) Transfer(
	idempotencyKey string,
	from string,
	to string,
	amount float64,
) error {

	// Account-level validation
	if _, ok, err := s.repo.GetAccount(from); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("from account not found")
	}
	if _, ok, err := s.repo.GetAccount(to); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("to account not found")
	}

	// Delegate money movement to ledger
	return s.ledger.Transfer(idempotencyKey, from, to, amount)
}

func (s *Service) GetLedgerEntries(
	accountNumber string,
) ([]ledger.Entry, error) {

	if _, ok, err := s.repo.GetAccount(accountNumber); err != nil {
		return nil, err
	} else if !ok {
		return nil, fmt.Errorf("account not found")
	}

	entries, err := s.ledger.EntriesForAccount(accountNumber)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (s *Service) Withdraw(
	idempotencyKey string,
	accountNumber string,
	amount float64,
) error {

	// Account-level validation
	if _, ok, err := s.repo.GetAccount(accountNumber); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("from account not found")
	}

	// Delegate money movement to ledger
	return s.ledger.Debit(idempotencyKey, accountNumber, amount, "Withdraw")
}

func (s *Service) Deposit(
	idempotencyKey string,
	accountNumber string,
	amount float64,
) error {

	// Account-level validation
	if _, ok, err := s.repo.GetAccount(accountNumber); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("from account not found")
	}

	// Delegate money movement to ledger
	return s.ledger.Credit(idempotencyKey, accountNumber, amount, "Deposit")
}
