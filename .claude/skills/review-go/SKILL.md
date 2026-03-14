---
name: review-go
description: Review Go code following project conventions (FastHTTP, fastglue, sqlx, Casbin, Koanf)
argument-hint: [file-or-package]
allowed-tools: Read, Grep, Bash
---

# Go Code Review — CanalGov

Review `$ARGUMENTS` (or current changed files if no argument) against the project's patterns.

## HTTP Handlers (cmd/)
- Handlers recebem `*fastglue.Request` e retornam `error`
- Usar `r.SendEnvelope()` para respostas de sucesso, `r.SendErrorEnvelope()` para erros
- Autenticação via middleware — nunca validar manualmente dentro do handler
- Permissões verificadas via Casbin (`app.authz.Enforce(...)`)

## Business Logic (internal/)
- Cada módulo tem: `manager.go` (ou arquivo principal), `models/models.go`, `queries.sql`
- Queries SQL separadas em arquivo `.sql`, carregadas via `goyesql`
- Usar `sqlx` para scans — nunca iterar rows manualmente quando `sqlx.Select` resolve
- Transações explícitas para operações que afetam múltiplas tabelas

## Padrões de Erro
- Wrapping com `fmt.Errorf("context: %w", err)`
- Erros de negócio com tipos personalizados quando necessário
- Nunca expor detalhes internos em mensagens de erro para o cliente

## Configuração
- Ler config via `app.ko.MustString(...)`, `app.ko.MustInt(...)`, etc. (Koanf)
- Nunca hardcodar valores que deveriam vir de config

## Concorrência
- Background workers usam `context.Context` para graceful shutdown
- Channels bufferizados para filas (ver padrões em `internal/inbox/`)
- Goroutines devem ter recover para evitar crashes silenciosos

## O que verificar
1. Handler segue o padrão de resposta do projeto?
2. SQL usa prepared statements / parâmetros nomeados?
3. Erros estão sendo propagados corretamente?
4. Há algum vazamento de goroutine?
5. Permissões Casbin estão sendo verificadas?
6. Config está vindo de Koanf, não hardcoded?

Reporte problemas com arquivo:linha e sugestão de correção.
