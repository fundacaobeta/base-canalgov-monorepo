# Especificação UX/UI — Sistema Administrativo CanalGov

> Versão 1.0 · Março 2026
> Contexto: plataforma de atendimento ao cidadão para órgãos e entidades públicas brasileiras.

---

## 1. Princípios de Design

### 1.1 Mentalidade do usuário-alvo
O usuário administrador é um **servidor público** (gerente de atendimento, analista de TI ou gestor de ouvidoria). Ele:
- Prefere clareza a flexibilidade técnica
- Precisa de ajuda contextual junto ao formulário — não em documentação separada
- Opera em contexto institucional (LGPD, acessibilidade, conformidade)
- Usa o sistema principalmente em desktop durante horário de expediente

### 1.2 Padrões de layout obrigatórios

#### Tela de Listagem
```
[ AdminPageHeader: título + descrição + breadcrumbs + botão "Novo" ]
─────────────────────────────────────────────────────────────────────
[ conteúdo principal (DataTable / cards) ]  │  [ painel de ajuda  ]
                                            │  • O que é este recurso
                                            │  • Quando usar
                                            │  • Dicas de configuração
                                            │  • Links externos
```

#### Tela de Criar / Editar
```
[ AdminPageHeader: título contextual + breadcrumbs + link "Voltar" ]
─────────────────────────────────────────────────────────────────────
[ formulário com seções                    ]  │  [ ajuda contextual ]
  Section 1: campos de identificação       ]  │  • Explicação de cada
  Section 2: configurações avançadas       ]  │    campo ou seção
  Section 3: toggles/status               ]  │  • Exemplos de valores
  [Botão Salvar]  [Botão Cancelar]         ]  │  • Advertências
```

#### Padrão de help panel
- Título: **"Ajuda"** (localizado via i18n)
- Largura: 280px fixo, sticky no scroll
- Conteúdo: parágrafos curtos + badges de atenção quando necessário
- Sempre presente em telas de config/form, opcional em listas simples

### 1.3 Regras de i18n
- **Zero strings hardcoded** em templates Vue — toda string visível usa `$t()`
- Exceção: labels de terceiros (nomes de provedores: "AWS SES", "Meta Cloud API", etc.)
- Idioma padrão: pt-BR. Inglês como fallback obrigatório em `en.json`
- Mensagens de erro, sucesso e confirmação: sempre i18n

---

## 2. Inventário de Telas e Estado Atual

### 2.1 Telas em estado BOAS (sem ação necessária)
| Tela | Rota | Notas |
|------|------|-------|
| Geral | `/admin/general` | ✅ Completa com i18n, help, loading |
| Etiquetas | `/admin/conversations/tags` | ✅ Padrão de referência |
| Status | `/admin/conversations/statuses` | ✅ Padrão de referência |
| Segmentos de Contatos | `/admin/contacts/segments` | ✅ Recém-melhorada |
| Domínios | `/admin/domains` | ✅ Recém-melhorada |
| Ações e Integrações | `/admin/integrations/actions` | ✅ Recém-melhorada |
| Log de Atividades | `/admin/teams/activity-log` | ✅ Com AdminPageHeader |
| SSO/OIDC | `/admin/sso` | ✅ Com AdminPageHeader |

### 2.2 Telas com MELHORIAS NECESSÁRIAS
| Tela | Rota | Problemas | Prioridade |
|------|------|-----------|------------|
| Layout base (AdminPageWithHelp) | — | "Ajuda" hardcoded em pt | 🔴 Alta |
| NotificationConfigShell | — | "Ajuda integrada", "Ocultar/Abrir ajuda" hardcoded | 🔴 Alta |
| Webhooks | `/admin/webhooks` | Title/description/help hardcoded em pt | 🔴 Alta |
| Caixas de Entrada | `/admin/inboxes` | Title/description/help hardcoded em pt | 🔴 Alta |
| Atributos Personalizados | `/admin/custom-attributes` | Sem AdminPageHeader; help em inglês | 🔴 Alta |
| Templates | `/admin/templates` | Help text hardcoded em pt | 🟡 Média |
| Notif. Canal (WhatsApp etc.) | `/admin/notification/*` | 90% strings hardcoded | 🟡 Média |
| Agentes — criar/editar | `/admin/teams/agents/new` | Help não muda ao criar | 🟡 Média |
| Equipes — criar/editar | `/admin/teams/teams/new` | Help não muda ao criar | 🟡 Média |
| Papéis — criar/editar | `/admin/teams/roles/new` | Help não muda ao criar | 🟡 Média |
| SLA — criar/editar | `/admin/sla/new` | Help não muda ao criar | 🟡 Média |
| Horário — criar/editar | `/admin/business-hours/new` | Help não muda ao criar | 🟡 Média |

---

## 3. Especificação por Tela

### 3.1 AdminPageWithHelp (layout global)

**Problema:** Título "Ajuda" hardcoded em português.

**Solução:**
```vue
<!-- antes -->
<p class="help-panel-title">Ajuda</p>

<!-- depois -->
<p class="help-panel-title">{{ $t('globals.terms.help') }}</p>
```

**i18n a adicionar:**
- `globals.terms.help` = "Help" / "Ajuda"

---

### 3.2 NotificationConfigShell

**Problema:** Strings hardcoded: "Ajuda integrada", "Ocultar ajuda", "Abrir ajuda", "Referências rápidas..."

**Solução:** Usar i18n para todas essas strings.

**i18n a adicionar:**
- `admin.notification.integratedHelp` = "Integrated help" / "Ajuda integrada"
- `admin.notification.helpDescription` = "Quick reference..." / "Referências rápidas..."
- `admin.notification.hideHelp` = "Hide help" / "Ocultar ajuda"
- `admin.notification.showHelp` = "Show help" / "Abrir ajuda"

---

### 3.3 Webhooks

**Problema:** Title, description e help hardcoded.

**Solução:** Usar i18n existente + novos keys:
- `admin.webhook.title` = "Webhooks" / "Webhooks"
- `admin.webhook.description` = "Receive real-time notifications..." / "Receba notificações em tempo real..."
- `admin.webhook.help` = "Configure webhooks..." / "Configure webhooks para..."
- `admin.webhook.help2` = "Webhooks allow..." / "Os webhooks permitem integrar..."

**Melhorias no WebhookList:**
- Mover botão "Novo webhook" para o `#actions` slot do AdminPageHeader do pai (Webhooks.vue)
- Detectar rota atual para mostrar o botão apenas na listagem

---

### 3.4 Caixas de Entrada (InboxView)

**Problema:** Title, description e help hardcoded em português.

**Solução:**
- `admin.inbox.title` = "Inboxes" / "Caixas de Entrada"
- `admin.inbox.description` = "Manage the communication channels..." / "Gerencie os canais de comunicação..."
- `admin.inbox.help` = "Connect your email..." / "Conecte sua conta de e-mail..."
- `admin.inbox.help2` = "Each added inbox..." / "Cada caixa de entrada adicionada..."

---

### 3.5 Atributos Personalizados

**Problemas:**
1. Missing `AdminPageHeader` (é a única tela de lista sem header próprio)
2. Help text em inglês

**Solução:**
- Adicionar `AdminPageHeader` com i18n
- Mover botão "Novo atributo" do content para o `#actions` slot
- Corrigir help text → i18n

**i18n a adicionar:**
- `admin.customAttributes.title` = "Custom Attributes" / "Atributos Personalizados"
- `admin.customAttributes.description` = "Define extra fields..." / "Defina campos extras para contatos e conversas."
- `admin.customAttributes.help` = "Custom attributes let you capture..." / "Atributos personalizados permitem registrar informações adicionais sobre contatos e conversas."
- `admin.customAttributes.help2` = "Use them in automations..." / "Use-os em regras de automação, filtros e visões compartilhadas."

---

### 3.6 Templates

**Problema:** Help text hardcoded em português.

**Solução:** Mover para i18n.
- `admin.template.help` = "Response templates are the default base for service." / "Modelos de resposta passam a ser a base padrão do atendimento."
- `admin.template.help2` = "Use global scope..." / "Use escopo global quando o texto servir para toda a operação..."
- `admin.template.help3` = "Email layouts and notifications..." / "Layouts e notificações de e-mail continuam separados..."

---

### 3.7 Help Contextual em Criar/Editar (padrão de rota)

**Problema:** Quando o usuário navega para `/admin/teams/teams/new`, o help panel do pai (Teams.vue) ainda mostra o help genérico sobre equipes, não ajuda sobre o formulário.

**Solução:** Cada parent view detecta a rota atual e mostra help diferente:

```vue
<!-- Teams.vue -->
<template #help>
  <!-- Lista -->
  <template v-if="isListRoute">
    <p>{{ $t('admin.team.help') }}</p>
    <p>{{ $t('admin.team.help2') }}</p>
  </template>
  <!-- Criar / Editar -->
  <template v-else>
    <p>{{ $t('admin.team.form.help') }}</p>
    <p>{{ $t('admin.team.form.help2') }}</p>
  </template>
</template>
```

**Telas afetadas:**
- Teams.vue (rotas: teams, new-team, edit-team)
- Agents.vue (rotas: agents, new-agent, edit-agent)
- Roles.vue (rotas: roles, new-role, edit-role)
- SLA.vue (rotas: sla, new-sla, edit-sla)
- BusinessHours.vue (rotas: business-hours, new-business-hours, edit-business-hours)
- Macros.vue (rotas: macros, new-macro, edit-macro)
- SharedViews.vue
- OIDC.vue
- Webhooks.vue (rotas: webhooks, new-webhook, edit-webhook)
- Templates.vue (rotas: templates, new-template, edit-template)

---

## 4. Design System Reference

### Cores e tokens
```css
--primary: 220 79% 39%;        /* Azul Gov.BR */
--muted-foreground: ...;        /* Texto secundário */
--border: ...;                  /* Bordas */
```

### Classes CSS utilitárias
- `.box` — card com borda e fundo sutil
- `.page-header` — container do AdminPageHeader
- `.page-title` — `text-2xl font-semibold tracking-tight`
- `.page-description` — `text-sm text-muted-foreground mt-1`
- `.help-panel` — aside sticky no desktop
- `.help-panel-title` — label "Ajuda" do help panel
- `.link-style` — links com sublinhado e cor primária

### Componentes de referência (padrão de ouro)
- **TagsView.vue** — padrão de listagem com dialog
- **DomainsView.vue** — padrão de listagem com modal + AlertDialog
- **ContactSegmentsView.vue** — padrão com Sheet + FilterBuilder

---

## 5. Checklist de Qualidade por Tela

Para cada tela, verificar:
- [ ] AdminPageHeader com title, description, breadcrumbs (i18n)
- [ ] AdminPageWithHelp com split layout
- [ ] Help panel com conteúdo real e útil (não genérico)
- [ ] Help muda conforme contexto (lista vs criar vs editar)
- [ ] Zero strings hardcoded visíveis
- [ ] Loading state presente
- [ ] Empty state presente quando lista pode estar vazia
- [ ] Confirmação via AlertDialog antes de destruir dados
- [ ] Erros tratados via useAdminErrorToast
- [ ] Formulários com validação via vee-validate + zod

---

## 6. i18n Keys a Adicionar

```json
{
  "globals.terms.help": "Help",
  "admin.webhook.title": "Webhooks",
  "admin.webhook.description": "Receive real-time HTTP notifications when events occur in the system.",
  "admin.webhook.help": "Configure webhooks to integrate CanalGov with external services.",
  "admin.webhook.help2": "Webhooks send HTTP POST requests when specific events occur.",
  "admin.inbox.title": "Inboxes",
  "admin.inbox.description": "Manage the email channels for your workspace.",
  "admin.inbox.help": "Connect your email account or configure IMAP/SMTP manually.",
  "admin.inbox.help2": "Each added inbox creates a new email channel for receiving messages.",
  "admin.customAttributes.title": "Custom Attributes",
  "admin.customAttributes.description": "Define extra fields to capture on contacts and conversations.",
  "admin.customAttributes.help": "Custom attributes let you record additional information beyond the default fields.",
  "admin.customAttributes.help2": "Use them in automation rules, filters, and shared views.",
  "admin.template.help": "Response templates are the default base for service.",
  "admin.template.help2": "Use global scope when the text serves the entire operation, or assign to a team for automatic prioritization.",
  "admin.template.help3": "Email layouts and notification templates are kept separate to avoid mixing reply content with technical structure.",
  "admin.notification.integratedHelp": "Integrated help",
  "admin.notification.helpDescription": "Quick references to configure this channel safely and consistently.",
  "admin.notification.hideHelp": "Hide help",
  "admin.notification.showHelp": "Show help",
  "admin.team.form.help": "Set a name and optional description to identify this team in assignments and reports.",
  "admin.team.form.help2": "After creating a team, assign agents and configure business hours and SLA policy.",
  "admin.agent.form.help": "Fill in the agent's name, email, and role. An invitation will be sent to the email provided.",
  "admin.agent.form.help2": "The role defines which screens and actions the agent can access.",
  "admin.role.form.help": "Create custom roles to restrict or expand agent access beyond the system defaults.",
  "admin.role.form.help2": "Permissions are granular — review each group before saving.",
  "admin.sla.form.help": "Define time targets for first response, next response, and resolution.",
  "admin.sla.form.help2": "At least one time target is required. Alerts are optional but recommended.",
  "admin.businessHours.form.help": "Set the days and hours your team is available to respond.",
  "admin.businessHours.form.help2": "Teams that don't have their own hours inherit the workspace schedule.",
  "admin.macro.form.help": "Define a sequence of actions to execute with a single click during a conversation.",
  "admin.macro.form.help2": "Macros can send messages, assign agents, add tags, and change status.",
  "admin.sharedViews.form.help": "A shared view is a saved filter visible to all agents or a specific team.",
  "admin.webhook.form.help": "Configure the endpoint URL, events, and optional secret for signature verification.",
  "admin.webhook.form.help2": "Test your webhook after saving to confirm delivery is working.",
  "admin.template.form.help": "Write the template content using plain text or HTML depending on the type.",
  "admin.template.form.help2": "Templates support variables like {contact_name} and {agent_name}."
}
```
