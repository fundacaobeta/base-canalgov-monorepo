# CanalGov - TODO List (Versão Final)

Lista de tarefas pendentes para o lançamento da versão final, dividida de acordo com os escopos definidos no `agents.md`.

## 🤖 Tarefas Backend (Gemini CLI)

### Débitos Técnicos (Baseados no código-fonte)
- [x] **`internal/media/models/models.go`:** Refatorado para usar constantes dos pacotes de modelo.
- [x] **`internal/authz/authz.go`:** Refatorado para usar constantes de modelo em vez de strings estáticas.
- [x] **`internal/user/user.go`:** Adicionada a constraint `UNIQUE` para `reset_password_token` via nova migration (`v1.0.9`).
- [x] **`internal/conversation/conversation.go`:** Removido o comentário `TODO` sobre `squirrel`.

### Polimento Final Backend
- [ ] Otimizar consultas ao banco de dados e garantir que todos os índices necessários estão criados.
- [x] **`internal/ws`:** Logs de erro em conexões de Websocket foram aprimorados para usar o logger estruturado.
- [x] **`internal/automation`:** Testes unitários expandidos para cobrir mais operadores (`not equal`, `less than`), ações de manipulação de tags e triggers baseados em eventos.
- [x] **`internal/notification/providers`:** Adicionado um timeout de 10 segundos para as requisições HTTP nos provedores de WhatsApp e Telegram.

---

## 🧠 Tarefas Frontend (Outra IA)

### Débitos Técnicos (Baseados no código-fonte)
- [x] **`frontend/src/stores/users.js`:** Renomeada para `agents.js` — export `useAgentsStore`, store ID `agents`, todos os imports atualizados.
- [x] **`frontend/src/stores/conversation.js` & `App.vue`:** `CONVERSATION_SORT_FIELD_MAP` e `REFRESH_MODEL` movidos para `constants/conversation.js`.
- [x] **`frontend/src/features/admin/activity-log/ActivityLog.vue`:** Paginação extraída para `composables/usePagination.js` + componente `AdminPagination.vue` reutilizável.
- [ ] **`frontend/src/features/admin/roles/RoleForm.vue`:** Automatizar o carregamento de permissões buscando-as dinamicamente em vez de hardcode (`// TODO: Prepare this by fetching all perms from the file`).

### Polimento Final Frontend
- [ ] **Acessibilidade (a11y):** Garantir suporte completo a leitores de tela e navegação por teclado nos componentes de chat e menus de administração.
- [ ] **Design Responsivo:** Testar e ajustar quebras de layout do Sidebar e caixas de entrada em telas menores.
- [ ] **UX Improvements:** Adicionar animações suaves (transitions) em ações de concluir conversa e navegação de menus.
- [ ] **Limpeza de Código:** Remover dependências legadas e otimizar bundle final no `vite.config.js`.

---

## 🔄 Tarefas Compartilhadas / Integração
- [ ] **Auditoria de Segurança:** Validar CSRF tokens e políticas CORS no frontend e backend.
- [x] **Auditoria de Segurança:** Documentação AppSec consolidada em `docs/security-audit-appsec.md` e smoke tests iniciais adicionados em `scripts/security/`.
- [ ] **Revisão de i18n:** Garantir que 100% das strings visíveis e de erro (incluindo as mensagens de backend) estejam mapeadas no `vue-i18n` ou devolvidas corretamente traduzidas.
- [ ] **Testes E2E:** Criar casos de teste cobrindo o fluxo principal (login -> receber mensagem -> responder -> resolver).
