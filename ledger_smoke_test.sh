#!/usr/bin/env bash
set -e

BASE_URL="http://localhost:8081"

echo "=== Health check ==="
curl -s $BASE_URL/health | jq
echo

echo "=== Create accounts ==="
curl -s -X POST $BASE_URL/accounts \
  -H "Content-Type: application/json" \
  -d '{"accountNumber":"alice","name":"Alice"}' | jq

curl -s -X POST $BASE_URL/accounts \
  -H "Content-Type: application/json" \
  -d '{"accountNumber":"bob","name":"Bob"}' | jq
echo

echo "=== Deposit into Alice ==="
curl -s -X POST $BASE_URL/accounts/deposit \
  -H "Content-Type: application/json" \
  -d '{"idempotencyKey":"deposit-1","accountNumber":"alice","amount":1000}' | jq
echo

echo "=== Transfer Alice -> Bob ==="
curl -s -X POST $BASE_URL/accounts/transfer \
  -H "Content-Type: application/json" \
  -d '{"idempotencyKey":"tx-1","from":"alice","to":"bob","amount":200}' | jq
echo

echo "=== Replay transfer (idempotency) ==="
curl -s -X POST $BASE_URL/accounts/transfer \
  -H "Content-Type: application/json" \
  -d '{"idempotencyKey":"tx-1","from":"alice","to":"bob","amount":200}' | jq
echo

echo "=== Withdraw from Bob ==="
curl -s -X POST $BASE_URL/accounts/withdraw \
  -H "Content-Type: application/json" \
  -d '{"idempotencyKey":"withdraw-1","accountNumber":"bob","amount":50}' | jq
echo

echo "=== Alice ledger ==="
curl -s $BASE_URL/accounts/alice/ledger | jq
echo

echo "=== Bob ledger ==="
curl -s $BASE_URL/accounts/bob/ledger | jq
echo

echo "=== Overdraw Bob (should fail) ==="
curl -s -X POST $BASE_URL/accounts/withdraw \
  -H "Content-Type: application/json" \
  -d '{"idempotencyKey":"withdraw-2","accountNumber":"bob","amount":1000}' | jq
echo

echo "=== Done ==="
