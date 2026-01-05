package ledger

import (
	"fmt"
	"sync"
	"testing"
)

func mustAllEntries(t *testing.T, repo Repository) []Entry {
	t.Helper()
	entries, err := repo.AllEntries()
	if err != nil {
		t.Fatalf("failed to fetch entries: %v", err)
	}
	return entries
}

func mustBalance(t *testing.T, svc *Service, account string) float64 {
	t.Helper()
	balance, err := svc.Balance(account)
	if err != nil {
		t.Fatalf("failed to fetch balance for %s: %v", account, err)
	}
	return balance
}

func TestMoneyIsConserved(t *testing.T) {
	repo := NewTestSQLiteRepository(t)
	svc := NewService(repo)

	if err := svc.Credit("init-1", "alice", 1000, "initial funding"); err != nil {
		t.Fatal(err)
	}

	_ = svc.Transfer("tx-1", "alice", "bob", 200)
	_ = svc.Transfer("tx-2", "bob", "carol", 50)
	_ = svc.Transfer("tx-3", "alice", "carol", 100)

	var total float64
	for _, e := range mustAllEntries(t, repo) {
		total += e.Amount
	}

	if total != 1000 {
		t.Fatalf("money invariant violated: expected 1000, got %f", total)
	}
}

func TestIdempotencyKeyAppliedOnce(t *testing.T) {
	repo := NewTestSQLiteRepository(t)
	svc := NewService(repo)

	_ = svc.Credit("init-1", "alice", 100, "initial")

	key := "dup-123"
	_ = svc.Transfer(key, "alice", "bob", 50)
	_ = svc.Transfer(key, "alice", "bob", 50) // replay

	if mustBalance(t, svc, "alice") != 50 {
		t.Fatalf("expected alice=50, got %f", mustBalance(t, svc, "alice"))
	}

	if mustBalance(t, svc, "bob") != 50 {
		t.Fatalf("expected bob=50, got %f", mustBalance(t, svc, "bob"))
	}
}

func TestNoNegativeBalances(t *testing.T) {
	repo := NewTestSQLiteRepository(t)
	svc := NewService(repo)

	_ = svc.Credit("init-1", "alice", 100, "initial")

	err := svc.Transfer("tx-1", "alice", "bob", 200)
	if err == nil {
		t.Fatal("expected insufficient funds error")
	}

	if mustBalance(t, svc, "alice") < 0 {
		t.Fatal("negative balance invariant violated")
	}
}

func TestTransferIsZeroSum(t *testing.T) {
	repo := NewTestSQLiteRepository(t)
	svc := NewService(repo)

	_ = svc.Credit("init-1", "alice", 100, "initial")
	_ = svc.Transfer("tx-1", "alice", "bob", 40)

	var sum float64
	for _, e := range mustAllEntries(t, repo) {
		if e.IdempotencyKey == "tx-1" {
			sum += e.Amount
		}
	}

	if sum != 0 {
		t.Fatalf("transfer not zero-sum: got %f", sum)
	}
}

func TestDifferentIdempotencyKeysBothApply(t *testing.T) {
	repo := NewTestSQLiteRepository(t)
	svc := NewService(repo)

	_ = svc.Credit("init-1", "alice", 100, "initial")
	_ = svc.Transfer("tx-1", "alice", "bob", 30)
	_ = svc.Transfer("tx-2", "alice", "bob", 30)

	if mustBalance(t, svc, "alice") != 40 {
		t.Fatalf("expected alice=40, got %f", mustBalance(t, svc, "alice"))
	}

	if mustBalance(t, svc, "bob") != 60 {
		t.Fatalf("expected bob=60, got %f", mustBalance(t, svc, "bob"))
	}
}

func TestTransferCreatesExactlyTwoEntries(t *testing.T) {
	repo := NewTestSQLiteRepository(t)
	svc := NewService(repo)

	_ = svc.Credit("init-1", "alice", 100, "initial")
	_ = svc.Transfer("tx-1", "alice", "bob", 40)

	count := 0
	for _, e := range mustAllEntries(t, repo) {
		if e.IdempotencyKey == "tx-1" {
			count++
		}
	}

	if count != 2 {
		t.Fatalf("expected 2 entries, got %d", count)
	}
}

func TestCreditIncreasesTotalMoney(t *testing.T) {
	repo := NewTestSQLiteRepository(t)
	svc := NewService(repo)

	_ = svc.Credit("c-1", "alice", 100, "credit")

	var total float64
	for _, e := range mustAllEntries(t, repo) {
		total += e.Amount
	}

	if total != 100 {
		t.Fatalf("expected total=100, got %f", total)
	}
}

func TestConcurrentTransfers(t *testing.T) {
	repo := NewTestSQLiteRepository(t)
	svc := NewService(repo)

	_ = svc.Credit("init-1", "alice", 1000, "initial")

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = svc.Transfer(
				fmt.Sprintf("tx-%d", i),
				"alice",
				"bob",
				10,
			)
		}(i)
	}
	wg.Wait()

	if mustBalance(t, svc, "alice") != 500 {
		t.Fatalf("expected alice=500, got %f", mustBalance(t, svc, "alice"))
	}

	if mustBalance(t, svc, "bob") != 500 {
		t.Fatalf("expected bob=500, got %f", mustBalance(t, svc, "bob"))
	}
}
