# TODO

Arquivo curto de navegaÃ§Ã£o. O backlog canÃ´nico foi reorganizado e movido para:

- [docs/backlog.md](/Users/jobs/Dev/a-publico/base-canalgov-monorepo/docs/backlog.md)

ReferÃªncias operacionais:

- DireÃ§Ã£o de produto: [docs/roadmap.md](/Users/jobs/Dev/a-publico/base-canalgov-monorepo/docs/roadmap.md)
- RevisÃ£o dos merges em andamento: [docs/merge-review.md](/Users/jobs/Dev/a-publico/base-canalgov-monorepo/docs/merge-review.md)
- MigraÃ§Ã£o seletiva de branches: [docs/branch-migration-roadmap.md](/Users/jobs/Dev/a-publico/base-canalgov-monorepo/docs/branch-migration-roadmap.md)
- RecomendaÃ§Ãµes de PR por branch: [docs/branch-pr-recommendations.md](/Users/jobs/Dev/a-publico/base-canalgov-monorepo/docs/branch-pr-recommendations.md)
- Auditoria AppSec: [docs/security-audit-appsec.md](/Users/jobs/Dev/a-publico/base-canalgov-monorepo/docs/security-audit-appsec.md)

Prioridade imediata:

1. Resolver os blockers descritos em `docs/merge-review.md`.
2. Manter `live-chat` e `help-articles` isolados atÃ© o build e o wiring estarem completos.
3. Usar `docs/backlog.md` para atualizaÃ§Ã£o de tarefas daqui em diante.

Backlog adicional preservado do branch local:

## ðŸ§  Tarefas Frontend (Outra IA)

### DÃ©bitos TÃ©cnicos (Baseados no cÃ³digo-fonte)
- [x] **`frontend/src/stores/users.js`:** Renomeada para `agents.js` â€” export `useAgentsStore`, store ID `agents`, todos os imports atualizados.
- [x] **`frontend/src/stores/conversation.js` & `App.vue`:** `CONVERSATION_SORT_FIELD_MAP` e `REFRESH_MODEL` movidos para `constants/conversation.js`.
- [x] **`frontend/src/features/admin/activity-log/ActivityLog.vue`:** PaginaÃ§Ã£o extraÃ­da para `composables/usePagination.js` + componente `AdminPagination.vue` reutilizÃ¡vel.
- [ ] **`frontend/src/features/admin/roles/RoleForm.vue`:** Automatizar o carregamento de permissÃµes buscando-as dinamicamente em vez de hardcode (`// TODO: Prepare this by fetching all perms from the file`).

### Polimento Final Frontend
- [ ] **Acessibilidade (a11y):** Garantir suporte completo a leitores de tela e navegaÃ§Ã£o por teclado nos componentes de chat e menus de administraÃ§Ã£o.
- [ ] **Design Responsivo:** Testar e ajustar quebras de layout do Sidebar e caixas de entrada em telas menores.
- [ ] **UX Improvements:** Adicionar animaÃ§Ãµes suaves (transitions) em aÃ§Ãµes de concluir conversa e navegaÃ§Ã£o de menus.
- [ ] **Limpeza de CÃ³digo:** Remover dependÃªncias legadas e otimizar bundle final no `vite.config.js`.

---

## ðŸ”„ Tarefas Compartilhadas / IntegraÃ§Ã£o
- [ ] **Auditoria de SeguranÃ§a:** Validar CSRF tokens e polÃ­ticas CORS no frontend e backend.
- [x] **Auditoria de SeguranÃ§a:** DocumentaÃ§Ã£o AppSec consolidada em `docs/security-audit-appsec.md` e smoke tests iniciais adicionados em `scripts/security/`.
- [ ] **RevisÃ£o de i18n:** Garantir que 100% das strings visÃ­veis e de erro (incluindo as mensagens de backend) estejam mapeadas no `vue-i18n` ou devolvidas corretamente traduzidas.
- [ ] **Testes E2E:** Criar casos de teste cobrindo o fluxo principal (login -> receber mensagem -> responder -> resolver).

- [x] **`frontend/src/views/admin/inbox/NewInbox.vue` + `cmd/inboxes.go`:** Corrigido o fluxo de criação de caixa de entrada para aceitar apenas `email` e validar o canal antes do `INSERT`, evitando `GeneralException` para canais não suportados.
