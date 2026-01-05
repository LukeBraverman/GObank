# GObank — Single-Node Ledger Engine

## Overview

GObank is a **single-node, correctness-first ledger engine** designed to reliably record financial transactions in constrained or disaster-recovery scenarios.

It is **not** a full banking product.

Instead, it focuses on doing **one thing well**:

> **Recording a durable, auditable, append-only financial ledger with strict invariants, even under concurrency and failure.**

The project intentionally prioritises:
- correctness over throughput
- auditability over convenience
- explicit tradeoffs over hidden complexity

---

## Design Goals

- **Single source of truth**: the ledger, not cached balances
- **Append-only**: no updates, no deletes
- **Durable idempotency**: safe retries without duplication
- **Serialized writes**: no double-spend under concurrency
- **Crash safety**: no partial writes
- **Human-inspectable state**: SQLite file + CSV-friendly data
- **Small surface area**: minimal moving parts

This system is intended as a **stop-gap financial record keeper**, not a production bank.

---

## What This Is (and Isn’t)

### ✅ This is:
- a ledger engine
- single-node
- SQLite-backed
- transactionally correct
- safe under concurrent requests
- resilient to crashes and retries

### ❌ This is not:
- an authentication system
- a cryptographic wallet
- a distributed system
- a fraud detection system
- a feature-complete banking product

Those concerns are deliberately out of scope.

---

## Core Concepts

### Ledger-First Design

Balances are **derived**, not stored.

All state comes from replaying ledger entries:
- deposits
- withdrawals
- transfers

This ensures:
- auditability
- traceability
- deterministic recovery

---

### Durable Idempotency

Every write operation requires an **idempotency key**.

If the same request is replayed:
- it is safely ignored
- no duplicate entries are created
- correctness is preserved

Idempotency is enforced **inside the database transaction**, not in memory.

---

### Serialized Transactions

Writes are serialized at the database level.

This guarantees:
- no double-spend
- no race conditions on balance checks
- deterministic outcomes under concurrency

SQLite is configured in **WAL mode** with immediate write locks to enforce this.

---

## Invariants Enforced

The system is built around explicit invariants, verified both in code and in tests.

### Financial Invariants

- **No negative balances**
- **Transfers are zero-sum**
- **Money is conserved** (except for explicit deposits/withdrawals)
- **Each transfer creates exactly two ledger entries**
- **Idempotency keys apply exactly once**

### System Invariants

- No partial writes
- No duplicate effects
- No silent corruption
- No dependency on in-memory state for correctness

---

## Database Guarantees

SQLite schema enforces additional safety:

- `NOT NULL` on all critical fields
- `CHECK (amount != 0)` on ledger entries
- `UNIQUE (idempotency_key, account_number)` per entry
- Indexed lookups for performance
- WAL mode for concurrent reads
- Busy timeout for contention handling

---

