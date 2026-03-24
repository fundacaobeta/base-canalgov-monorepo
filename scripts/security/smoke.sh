#!/usr/bin/env bash
set -u

BASE_URL="${BASE_URL:-http://127.0.0.1:9000}"
LOW_PRIV_COOKIE="${LOW_PRIV_COOKIE:-}"
ADMIN_COOKIE="${ADMIN_COOKIE:-}"
TEST_EMAIL="${TEST_EMAIL:-target@example.gov}"
TEST_PASSWORD="${TEST_PASSWORD:-WrongPass123!}"
SIGNED_URL="${SIGNED_URL:-}"
BRUTE_FORCE_ATTEMPTS="${BRUTE_FORCE_ATTEMPTS:-10}"

pass() {
  printf '[SAFE] %s\n' "$1"
}

fail() {
  printf '[VULNERABLE] %s\n' "$1"
}

warn() {
  printf '[SKIP] %s\n' "$1"
}

code_for() {
  curl -s -o /tmp/canalgov-security-body.$$ -w '%{http_code}' "$@"
}

body_contains() {
  pattern="$1"
  grep -q "$pattern" /tmp/canalgov-security-body.$$
}

test_roles_exposure() {
  if [ -z "$LOW_PRIV_COOKIE" ]; then
    warn 'roles exposure requires LOW_PRIV_COOKIE'
    return
  fi

  code=$(code_for -H "Cookie: $LOW_PRIV_COOKIE" "$BASE_URL/api/v1/roles")
  if [ "$code" = "200" ] && body_contains '"permissions"'; then
    fail 'GET /api/v1/roles exposed permissions to low-priv user'
  else
    pass 'GET /api/v1/roles blocked or minimized'
  fi
}

test_inboxes_exposure() {
  if [ -z "$LOW_PRIV_COOKIE" ]; then
    warn 'inboxes exposure requires LOW_PRIV_COOKIE'
    return
  fi

  code=$(code_for -H "Cookie: $LOW_PRIV_COOKIE" "$BASE_URL/api/v1/inboxes")
  if [ "$code" = "200" ] && body_contains '"config"'; then
    fail 'GET /api/v1/inboxes exposed inbox config to low-priv user'
  else
    pass 'GET /api/v1/inboxes blocked or minimized'
  fi
}

test_macros_exposure() {
  if [ -z "$LOW_PRIV_COOKIE" ]; then
    warn 'macros exposure requires LOW_PRIV_COOKIE'
    return
  fi

  code=$(code_for -H "Cookie: $LOW_PRIV_COOKIE" "$BASE_URL/api/v1/macros")
  if [ "$code" = "200" ] && body_contains '"visibility"'; then
    fail 'GET /api/v1/macros returned macro catalog to low-priv user'
  else
    pass 'GET /api/v1/macros filtered or blocked'
  fi
}

test_webhook_spoof() {
  code=$(code_for "$BASE_URL/api/v1/webhooks/whatsapp?hub.mode=subscribe&hub.verify_token=test&hub.challenge=1337")
  if [ "$code" = "200" ] && body_contains '1337'; then
    fail 'WhatsApp webhook verification accepted arbitrary verify_token'
  else
    pass 'WhatsApp webhook verification rejected arbitrary token'
  fi
}

test_signed_url() {
  if [ -z "$SIGNED_URL" ]; then
    warn 'signed URL test skipped; set SIGNED_URL to a captured attachment URL'
    return
  fi

  code=$(code_for "$SIGNED_URL")
  if [ "$code" = "200" ]; then
    fail 'signed attachment URL worked without session'
  else
    pass 'signed attachment URL rejected without session/context'
  fi
}

test_bruteforce_login() {
  tmp_codes="/tmp/canalgov-security-codes.$$"
  : > "$tmp_codes"

  i=1
  while [ "$i" -le "$BRUTE_FORCE_ATTEMPTS" ]; do
    curl -s -o /dev/null -w '%{http_code}\n' \
      -H 'Content-Type: application/json' \
      -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\"}" \
      "$BASE_URL/api/v1/auth/login" >> "$tmp_codes"
    i=$((i + 1))
  done

  if grep -q '^429$' "$tmp_codes"; then
    pass 'login flow returned 429 during brute-force probe'
  else
    fail 'login flow showed no 429 during brute-force probe'
  fi
}

cleanup() {
  rm -f /tmp/canalgov-security-body.$$ /tmp/canalgov-security-codes.$$
}

trap cleanup EXIT

printf 'BASE_URL=%s\n' "$BASE_URL"
test_macros_exposure
test_roles_exposure
test_inboxes_exposure
test_webhook_spoof
test_signed_url
test_bruteforce_login
