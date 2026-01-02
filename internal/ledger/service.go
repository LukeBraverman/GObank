package ledger

import (
    "fmt"
    "time"

    "github.com/google/uuid"
)

type Service struct {
    repo Repository
}

func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) Transfer(
    idempotencyKey string,
    from string,
    to string,
    amount float64,
) error {

    if amount <= 0 {
        return fmt.Errorf("amount must be positive")
    }
    if from == to {
        return fmt.Errorf("cannot transfer to same account")
    }

    // TODO Understnad this function signiture. Fix other methids
    return s.repo.WithTransaction(func(tx TxRepository) error {

        // 1️⃣ Idempotency check (DB is source of truth)
        seen, err := tx.HasIdempotencyKey(idempotencyKey)
        if err != nil {
            return err
        }
        if seen {
            return nil
        }

       balance, err := tx.Balance(from)
        if err != nil {
            return err
        }
        if balance < amount {
            return fmt.Errorf("insufficient funds")
        }

        if balance < amount {
            return fmt.Errorf("insufficient funds")
        }

        txID := uuid.NewString()
        now := time.Now().UTC()

        // 3️⃣ Ledger entries
        if err := tx.AppendEntry(Entry{
            TransactionID:  txID,
            IdempotencyKey: idempotencyKey,
            AccountNumber:  from,
            Amount:         -amount,
            Description:    "transfer to " + to,
            CreatedAt:      now,
        }); err != nil {
            return err
        }

        if err := tx.AppendEntry(Entry{
            TransactionID:  txID,
            IdempotencyKey: idempotencyKey,
            AccountNumber:  to,
            Amount:         amount,
            Description:    "transfer from " + from,
            CreatedAt:      now,
        }); err != nil {
            return err
        }

        // 4️⃣ Record idempotency key LAST
        return tx.RecordIdempotencyKey(idempotencyKey)
    })
}

func (s *Service) Debit(
    idempotencyKey string,
    accountNumber string,
    amount float64,
    description string,
) error {

    if accountNumber == "" {
        return fmt.Errorf("account number required")
    }
    if amount <= 0 {
        return fmt.Errorf("debit amount must be positive")
    }

    return s.repo.WithTransaction(func(tx TxRepository) error {

        // 1️⃣ Idempotency check
        seen, err := tx.HasIdempotencyKey(idempotencyKey)
        if err != nil {
            return err
        }
        if seen {
            return nil
        }

        // 2️⃣ Balance check (same caveat as Transfer)
     balance, err := tx.Balance(accountNumber)
        if err != nil {
            return err
        }
        if balance < amount {
            return fmt.Errorf("insufficient funds")
        }

        if balance < amount {
            return fmt.Errorf("insufficient funds")
        }

        // 3️⃣ Append ledger entry
        if err := tx.AppendEntry(Entry{
            TransactionID:  uuid.NewString(),
            IdempotencyKey: idempotencyKey,
            AccountNumber:  accountNumber,
            Amount:         -amount,
            Description:    description,
            CreatedAt:      time.Now().UTC(),
        }); err != nil {
            return err
        }

        // 4️⃣ Record idempotency key LAST
        return tx.RecordIdempotencyKey(idempotencyKey)
    })
}



func (s *Service) Credit(
    idempotencyKey string,
    accountNumber string,
    amount float64,
    description string,
) error {

    if accountNumber == "" {
        return fmt.Errorf("account number required")
    }
    if amount <= 0 {
        return fmt.Errorf("credit amount must be positive")
    }

    return s.repo.WithTransaction(func(tx TxRepository) error {

        // 1️⃣ Idempotency check
        seen, err := tx.HasIdempotencyKey(idempotencyKey)
        if err != nil {
            return err
        }
        if seen {
            return nil
        }

        // 2️⃣ Append ledger entry
        if err := tx.AppendEntry(Entry{
            TransactionID:  uuid.NewString(),
            IdempotencyKey: idempotencyKey,
            AccountNumber:  accountNumber,
            Amount:         amount,
            Description:    description,
            CreatedAt:      time.Now().UTC(),
        }); err != nil {
            return err
        }

        // 3️⃣ Record idempotency key LAST
        return tx.RecordIdempotencyKey(idempotencyKey)
    })
}

func (s *Service) Balance(accountNumber string) float64 {
    entries := s.repo.EntriesForAccount(accountNumber)

    var total float64
    for _, e := range entries {
        total += e.Amount
    }

    return total
}

func (s *Service) EntriesForAccount(accountNumber string) []Entry {
    return s.repo.EntriesForAccount(accountNumber)
}

