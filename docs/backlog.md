# Backlog

Este e o backlog canônico do CanalGov. O arquivo [todo.md](/Users/jobs/Dev/a-publico/base-canalgov-monorepo/todo.md) fica apenas como índice curto e ponte para este documento.

## 1. Bloqueadores Ativos

### Merge review: live-chat, help-center e AI
- Corrigir os blockers documentados em [docs/merge-review.md](/Users/jobs/Dev/a-publico/base-canalgov-monorepo/docs/merge-review.md) antes de aceitar novas integrações nessa frente.
- Restaurar build dos pacotes novos de `AI` e `conversation`.
- Fechar wiring em `cmd/handlers.go`, `cmd/init.go` e `cmd/main.go` para qualquer módulo realmente ativo.

### Qualidade de validação
- Restaurar um `go test ./...` limpo ou, no mínimo, separar claramente dívida legada de regressão nova.
- Adicionar validações mínimas de integração para módulos públicos novos antes de seguir com mais merges.

## 2. Backlog de Código

### Backend
- Otimizar consultas ao banco e revisar índices ainda faltantes.
- Automatizar carregamento dinâmico de permissões no fluxo de roles em vez de manter lista hardcoded no frontend.
- Consolidar inicialização de módulos opcionais para reduzir código parcialmente ligado.

### Frontend
- Melhorar acessibilidade dos fluxos de conversa e administração.
- Ajustar responsividade de sidebar, inbox e painéis laterais.
- Limpar dependências legadas e revisar impacto no bundle final.

## 3. Segurança e Cobertura

- Validar CSRF e políticas CORS em frontend e backend.
- Expandir testes E2E do fluxo principal.
- Manter a auditoria AppSec e os smoke tests em [docs/security-audit-appsec.md](/Users/jobs/Dev/a-publico/base-canalgov-monorepo/docs/security-audit-appsec.md) e [scripts/security/README.md](/Users/jobs/Dev/a-publico/base-canalgov-monorepo/scripts/security/README.md).

## 4. Migração Seletiva de Branches

### Encerrado
- Revisão concluída das branches remotas fora de `live-chat` e `help-articles`.
- `origin/fix/reply-draft-persistence` reavaliada e fechada sem ação restante.
- Branches 100% absorvidas ou sem delta útil já foram removidas do remoto.

### Delegado para outra IA
- `origin/feat/live-chat-channel`
- `origin/help-articles-and-ai-responses`

Referências:
- [docs/branch-migration-roadmap.md](/Users/jobs/Dev/a-publico/base-canalgov-monorepo/docs/branch-migration-roadmap.md)
- [docs/branch-pr-recommendations.md](/Users/jobs/Dev/a-publico/base-canalgov-monorepo/docs/branch-pr-recommendations.md)

## 5. Documentação de Produto e Planejamento

- Manter direção de produto resumida em [docs/roadmap.md](/Users/jobs/Dev/a-publico/base-canalgov-monorepo/docs/roadmap.md).
- Manter backlog técnico e operacional neste arquivo.
- Usar `todo.md` apenas como ponte rápida para quem entrar no repositório.
