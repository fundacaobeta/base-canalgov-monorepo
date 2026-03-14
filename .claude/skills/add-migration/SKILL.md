---
name: add-migration
description: Scaffold a new database migration file following the project's versioned migration pattern
argument-hint: [version] (e.g. v1.1.0)
disable-model-invocation: true
allowed-tools: Read, Glob, Write, Bash
---

# Add Database Migration — CanalGov

Crie uma nova migration para a versão `$ARGUMENTS`.

## Processo

1. **Leia as migrations existentes** em `internal/migrations/` para entender o padrão:
   - Cada arquivo é `v{version}.go` (ex: `v1.0.4.go`)
   - Cada arquivo tem uma função `V{version}(db *sqlx.DB) error`
   - Migrations são idempotentes (seguro rodar múltiplas vezes)
   - Usar `IF NOT EXISTS`, `IF EXISTS`, ou verificar antes de alterar

2. **Verifique onde as migrations são registradas** (provavelmente `cmd/install.go` ou `cmd/upgrade.go`) e adicione a nova versão.

3. **Crie o arquivo** `internal/migrations/v{version}.go` com:
   ```go
   package migrations

   import "github.com/jmoirons/sqlx"

   func V{Version}(db *sqlx.DB) error {
       // TODO: adicionar SQL da migration aqui
       return nil
   }
   ```

4. **Padrões SQL a seguir**:
   - `ALTER TABLE ... ADD COLUMN IF NOT EXISTS ...`
   - `CREATE INDEX IF NOT EXISTS ...`
   - `CREATE TABLE IF NOT EXISTS ...`
   - Nunca DROP sem verificar existência
   - Adicionar comentário explicando o propósito da migration

5. **Atualize `schema.sql`** se a migration adiciona tabelas/colunas permanentes.

Peça ao usuário a descrição do que a migration deve fazer se não for fornecida.
