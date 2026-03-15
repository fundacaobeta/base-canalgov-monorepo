# AI Collaboration Plan & Module Scopes

Este documento define a divisão de escopos entre a **Gemini CLI** e a **Outra IA** (ex: Claude, GitHub Copilot, etc.) para a finalização e polimento do CanalGov (versão final).

A estratégia visa maximizar a eficiência, alocando tarefas baseadas nas forças de cada agente: a Gemini CLI cuidará da infraestrutura, arquitetura backend em Go e integrações sistêmicas, enquanto a Outra IA ficará responsável pela experiência do usuário, design de interface em Vue e estado do frontend.

---

## 🤖 Gemini CLI (Backend, Arquitetura e Infraestrutura)

**Foco Principal:** Go, PostgreSQL, Redis, Docker, CI/CD, e integrações de sistema.

### Escopo por Módulos:
1. **Core Backend & APIs (`internal/`, `cmd/`):**
   - Roteamento e middlewares (`cmd/handlers.go`, `cmd/middlewares.go`).
   - Gerenciamento de Banco de Dados (`internal/dbutil`, migrations).
   - Websockets (`internal/ws`).
   - Segurança e Autorização (`internal/auth`, `internal/authz`, `internal/oidc`).

2. **Motores de Regras e Background Workers:**
   - Automações e Autoatribuição (`internal/automation`, `internal/autoassigner`).
   - SLAs e Horário Comercial (`internal/sla`, `internal/business_hours`).
   - Relatórios e Logs (`internal/report`, `internal/activity_log`).

3. **Integrações e Canais:**
   - Provedores de Notificação (`internal/notification/providers/*` - WhatsApp, Telegram, SMS, Push).
   - Gerenciamento de Mídia e S3 (`internal/media`).
   - IA e Assistentes (`internal/ai`).

4. **Infraestrutura e DevOps:**
   - Docker (`Dockerfile`, `docker-compose.yml`).
   - CI/CD Pipelines (`.github/workflows/`).
   - Scripts de automação e Seed (`Makefile`, `scripts/`).

---

## 🧠 Outra IA (Frontend, UI/UX e Acessibilidade)

**Foco Principal:** Vue 3, Tailwind CSS, Shadcn Vue, Pinia, Vite, e Acessibilidade (a11y).

### Escopo por Módulos:
1. **Views e Layouts (`frontend/src/views/`, `frontend/src/layouts/`):**
   - Polimento do layout principal (`InboxLayout`, `AdminLayout`, `AccountLayout`).
   - Responsividade em dispositivos móveis.
   - Refinamento visual da lista de contatos e detalhes (`ContactsView`, `ContactDetailView`).

2. **Componentes UI e Design System (`frontend/src/components/`, `frontend/src/features/`):**
   - Integração e ajuste de componentes Shadcn Vue (`ui/*`).
   - Experiência de chat (`features/conversation/message/MessageBubble.vue`, `ReplyBox.vue`).
   - Painéis de Administração (`features/admin/*`).

3. **Estado e Lógica de Interface (`frontend/src/stores/`, `frontend/src/composables/`):**
   - Gerenciamento de estado complexo via Pinia (`users.js`, `conversation.js`).
   - Formatação e internacionalização (`vue-i18n`, `composables/`).

4. **Testes e Ferramentas Frontend:**
   - Testes E2E (Cypress) e unitários (Vitest).
   - Configurações do Vite e plugins (`vite.config.js`).

---

## 🤝 Fluxo de Trabalho e Sincronização

- **Comunicação:** Atualizar o `todo.md` assim que uma tarefa for concluída.
- **Tipos compartilhados:** Quaisquer mudanças em contratos de API ou payloads devem ser comunicadas imediatamente, pois afetam ambos os escopos. A Gemini CLI atualiza as rotas Go e a Outra IA atualiza o `api/index.js` do frontend.