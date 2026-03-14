---
name: run-tests
description: Run Go and/or frontend tests for the project and report results
argument-hint: [go|frontend|all]
disable-model-invocation: true
allowed-tools: Bash
---

# Run Tests — CanalGov

Execute os testes conforme o argumento `$ARGUMENTS` (padrão: `all`).

## Go Tests
```bash
cd /Users/jobs/Dev/a-publico/base-canalgov-monorepo
go test -count=1 -race ./...
```
- `-count=1` desabilita cache
- `-race` detecta race conditions

## Frontend Unit Tests (Vitest)
```bash
cd frontend
pnpm test:unit --run
```

## Frontend E2E (Cypress)
```bash
cd frontend
pnpm test:e2e
```
Requer servidor rodando — avise o usuário se necessário.

## Processo
1. Se `$ARGUMENTS` for `go` → rodar apenas Go tests
2. Se `$ARGUMENTS` for `frontend` → rodar apenas Vitest
3. Se `$ARGUMENTS` for `all` ou vazio → rodar Go + Vitest
4. E2E só roda se explicitamente pedido com `e2e`

## Reportar
- Quantos testes passaram/falharam
- Output completo de falhas com arquivo:linha
- Sugerir causa provável de falhas comuns (ex: banco não conectado, variável de env faltando)
