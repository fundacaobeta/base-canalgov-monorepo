# Handoff Prompt — Continuação do Redesign Admin CanalGov

## Contexto do Projeto

CanalGov é uma plataforma de atendimento ao cidadão (helpdesk governamental) com:
- **Backend**: Go + FastHTTP + PostgreSQL + Redis
- **Frontend**: Vue 3 + Vite + Shadcn UI + TailwindCSS + Pinia
- **Localização**: `/Users/jobs/Dev/a-publico/base-canalgov-monorepo`

## O que já foi feito

### Design System (documentado em `docs/design/`)
- `00-design-system.md` — tokens completos (cores, tipografia, espaçamento, sombras, acessibilidade)
- `01-admin-pages-specs.md` — specs de UX para todas as 24 páginas admin
- `02-components-spec.md` — specs de componentes reutilizáveis

### Skills criadas (`.claude/skills/`)
- `design-system` — aplica o design system ao revisar/criar componentes
- `review-go`, `review-vue`, `add-migration`, `i18n-check`, `run-tests`, `gov-system-design`

### Mudanças já implementadas

1. **`frontend/src/assets/styles/main.scss`**
   - Cor primária mudada para azul Gov.BR: `--primary: 220 79% 39%`
   - Adicionados tokens semânticos: `--success`, `--warning`, `--info` (com variantes `-bg`, `-foreground`)
   - Novas classes CSS: `.page-header`, `.page-title`, `.page-description`, `.form-section`, `.form-section-title`, `.badge-success`, `.badge-warning`, `.badge-info`, `.badge-neutral`, `.help-panel`, `.admin-content`

2. **`frontend/src/layouts/admin/AdminLayout.vue`**
   - Fundo `bg-muted/30`, padding melhorado, max-width `max-w-6xl` centralizado

3. **`frontend/src/layouts/admin/AdminPageWithHelp.vue`**
   - Grid responsivo `lg:grid-cols-[1fr_280px]` com painel de ajuda sticky
   - Painel de ajuda com estilo `help-panel` e título "Ajuda"

4. **`frontend/src/components/layout/AdminPageHeader.vue`** ← NOVO COMPONENTE
   - Props: `title`, `description`, `breadcrumbs`
   - Slot `#actions` para botões
   - Breadcrumb com chevron, H1, descrição

5. **`frontend/src/components/ui/badge/index.js`**
   - Badges arredondados (`rounded-full`)
   - Novas variantes: `success`, `warning`, `info`, `neutral`

6. **`frontend/src/views/admin/tags/TagsView.vue`**
   - Usa `AdminPageHeader` com breadcrumb, título, descrição
   - Botão "Nova Tag" com ícone `Plus` no header
   - Dialog com botão Cancelar ao lado do Salvar
   - Removido `Spinner` sobreposto (loading via tabela)

7. **`frontend/src/views/admin/general/General.vue`**
   - Usa `AdminPageHeader` com breadcrumb e descrição
   - Removido `Spinner` sobreposto

8. **`frontend/src/features/admin/general/GeneralSettingForm.vue`**
   - Formulário agrupado em 4 seções com títulos: Identidade, Localização, Atendimento, Arquivos
   - Placeholders descritivos nos inputs
   - Input de MB com unidade ao lado
   - Botão salvar separado por divisor

## O que fazer a seguir

### Prioridade Alta — Páginas mais usadas

Para cada uma das páginas abaixo, aplicar o mesmo padrão de `TagsView.vue`:
- Adicionar `AdminPageHeader` com título, descrição, breadcrumb e ação primária no header
- Remover Spinner sobreposto (usar loading state nativo da tabela)
- Adicionar botão Cancelar nos dialogs/forms
- Usar ícone `Plus` no botão de criar

**Páginas a fazer:**

1. **`/views/admin/teams/Teams.vue`** e `features/admin/teams/TeamForm.vue`
   - Header: "Equipes" + "Organize agentes em equipes para roteamento de conversas."
   - Form em seções: Dados Básicos | Membros | Configurações

2. **`/views/admin/agents/Agents.vue`** e `features/admin/agents/AgentForm.vue`
   - Header: "Agentes" + "Gerencie os atendentes do seu workspace."
   - Botão secundário "Importar CSV" antes do botão principal

3. **`/views/admin/status/StatusView.vue`**
   - Header: "Status de Conversa" + descrição
   - Adicionar badge de preview de cor no formulário

4. **`/views/admin/inbox/InboxView.vue`** e `NewInbox.vue`
   - Header: "Caixas de Entrada" + descrição
   - Wizard com step indicator (componente `WizardSteps`)

5. **`/views/admin/teams/Roles.vue`**
   - Header: "Papéis e Permissões"
   - Permissões agrupadas por recurso (fieldsets)

### Prioridade Média

6. **`/views/admin/automations/Automation.vue`**
   - Header com tabs como ações secundárias

7. **`/views/admin/templates/Templates.vue`**
   - Header por tab (Respostas / Layouts / Notificações)

8. **`/views/admin/webhooks/Webhooks.vue`**
   - Header padrão

### Componentes a criar

```
frontend/src/components/layout/
  AdminPageHeader.vue  ← já existe
  WizardSteps.vue      ← criar (ver spec em docs/design/02-components-spec.md)
  HelpItem.vue         ← criar (item de ajuda com ícone + título + texto)
  EmptyState.vue       ← criar (ícone + título + descrição + CTA)
```

### AuthLayout — melhorar login

**`/src/layouts/auth/AuthLayout.vue`** — adicionar identidade gov:
- Logo do workspace centralizado
- Fundo com padrão sutil (gradiente ou pattern)
- Rodapé com "Powered by CanalGov" + link de privacidade

## Padrão de implementação

Para cada página, seguir este template:

```vue
<template>
  <div>
    <!-- 1. Header com breadcrumb, título, descrição e ação -->
    <AdminPageHeader
      title="Título da Página"
      description="Descrição breve do que esta página faz."
      :breadcrumbs="[{ label: 'Admin', to: '/admin' }, { label: 'Título' }]"
    >
      <template #actions>
        <Button @click="openDialog">
          <Plus class="h-4 w-4 mr-1.5" />
          Criar Novo
        </Button>
      </template>
    </AdminPageHeader>

    <!-- 2. Conteúdo com painel de ajuda -->
    <AdminPageWithHelp>
      <template #content>
        <!-- tabela ou form -->
      </template>
      <template #help>
        <p>Texto de ajuda contextual.</p>
      </template>
    </AdminPageWithHelp>
  </div>
</template>
```

## Convenções

- Importar `AdminPageHeader` de `@/components/layout/AdminPageHeader.vue`
- Importar `Plus` de `lucide-vue-next`
- Botão cancelar sempre ao lado do submit: `variant="outline"`
- Loading: nunca `<Spinner>` sobreposto — usar `opacity-50 pointer-events-none` no container
- Formulários com campos agrupados em `<div class="form-section">` + `<h2 class="form-section-title">`
- Badges de status: usar `variant="success"`, `variant="warning"`, `variant="info"`, `variant="neutral"`

## Arquivos de referência importantes

- Design tokens: `docs/design/00-design-system.md`
- Specs de páginas: `docs/design/01-admin-pages-specs.md`
- Specs de componentes: `docs/design/02-components-spec.md`
- Skill de design: `.claude/skills/design-system/SKILL.md`
- Exemplo implementado: `frontend/src/views/admin/tags/TagsView.vue`
- CSS global: `frontend/src/assets/styles/main.scss`
