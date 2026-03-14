# CanalGov — Spec de Componentes

> Especificações técnicas e de comportamento dos componentes reutilizáveis.
> Para cada componente: props, slots, estados, acessibilidade e exemplo de uso.

---

## 1. PageHeader

Componente obrigatório no topo de toda página admin.

**Props:**
```ts
interface PageHeaderProps {
  title: string           // H1 da página
  description?: string    // Subtexto descritivo
  breadcrumbs?: Array<{ label: string; to?: string }>
}
// Slot: #actions — botões e ações do lado direito
```

**Anatomia:**
```
┌────────────────────────────────────────────────────────────────┐
│  Admin > Equipes > Criar                      ← breadcrumb     │
│  Criar Equipe                                 ← h1             │
│  Configure uma nova equipe de atendimento.    ← description    │
│                                    [Cancelar] [Criar Equipe]   │
└────────────────────────────────────────────────────────────────┘
```

**Implementação Vue:**
```vue
<!-- src/components/layout/PageHeader.vue -->
<template>
  <div class="mb-6">
    <Breadcrumb v-if="breadcrumbs?.length" :items="breadcrumbs" class="mb-2" />
    <div class="flex items-start justify-between gap-4">
      <div class="min-w-0">
        <h1 class="text-2xl font-bold text-neutral-800 truncate">{{ title }}</h1>
        <p v-if="description" class="mt-1 text-sm text-neutral-500">{{ description }}</p>
      </div>
      <div v-if="$slots.actions" class="flex shrink-0 items-center gap-2">
        <slot name="actions" />
      </div>
    </div>
  </div>
</template>
```

---

## 2. DataTable (versão gov)

**Props:**
```ts
interface DataTableProps {
  columns: ColumnDef[]
  data: unknown[]
  loading?: boolean       // mostra skeleton
  empty?: {
    icon?: string         // lucide icon name
    title: string
    description?: string
    action?: { label: string; onClick: () => void }
  }
  searchable?: boolean    // mostra campo de busca
  searchPlaceholder?: string
  selectable?: boolean    // checkbox por linha
}
// Slot: #toolbar-extra — filtros adicionais
// Slot: #bulk-actions — ações quando há seleção
```

**Estado de Loading (Skeleton):**
```vue
<template v-if="loading">
  <TableRow v-for="i in 5" :key="i">
    <TableCell v-for="col in columns" :key="col.id">
      <Skeleton class="h-4 w-full rounded" />
    </TableCell>
  </TableRow>
</template>
```

**Estado Empty:**
```vue
<template v-if="!loading && data.length === 0">
  <TableRow>
    <TableCell :colspan="columns.length">
      <div class="flex flex-col items-center py-16 text-center">
        <component :is="icon" class="h-12 w-12 text-neutral-300 mb-4" />
        <p class="text-base font-medium text-neutral-700">{{ empty.title }}</p>
        <p v-if="empty.description" class="mt-1 text-sm text-neutral-500 max-w-sm">
          {{ empty.description }}
        </p>
        <Button v-if="empty.action" class="mt-4" @click="empty.action.onClick">
          {{ empty.action.label }}
        </Button>
      </div>
    </TableCell>
  </TableRow>
</template>
```

---

## 3. HelpPanel

Painel lateral de ajuda contextual.

**Props:**
```ts
interface HelpPanelProps {
  title?: string  // padrão: "Ajuda"
}
// Slot: default — HelpItem components
```

**HelpItem props:**
```ts
interface HelpItemProps {
  icon?: string   // lucide icon
  title: string
  href?: string   // se externo, abre nova aba
}
// Slot: default — descrição do item
```

**Implementação:**
```vue
<!-- Exemplo de uso -->
<HelpPanel>
  <HelpItem icon="Users" title="O que são equipes?">
    Equipes agrupam agentes e definem filas de atendimento automático.
  </HelpItem>
  <HelpItem icon="ArrowRight" title="Ver documentação" href="https://docs.canalgov.com/equipes" />
</HelpPanel>
```

**Estilo:**
```css
.help-panel {
  background: var(--color-neutral-50);
  border: 1px solid var(--color-border-subtle);
  border-radius: var(--radius-lg);
  padding: 1.25rem;
}
.help-item-icon {
  width: 1.25rem; height: 1.25rem;
  color: var(--color-brand-500);
  flex-shrink: 0;
}
```

---

## 4. StatusBadge

Badge semântico para status de entidades.

**Props:**
```ts
type StatusVariant =
  | 'success'   // verde — ativo, resolvido, online
  | 'warning'   // amarelo — em andamento, atenção, snoozed
  | 'danger'    // vermelho — erro, vencido, inativo
  | 'info'      // azul — informação, processando
  | 'neutral'   // cinza — arquivado, rascunho
  | 'brand'     // marca — destacado, especial

interface StatusBadgeProps {
  variant: StatusVariant
  dot?: boolean     // mostrar bolinha colorida antes do texto
  size?: 'sm' | 'md'
}
```

**Implementação:**
```vue
<template>
  <span :class="['badge', `badge-${variant}`, size === 'sm' ? 'badge-sm' : '']">
    <span v-if="dot" :class="['badge-dot', `dot-${variant}`]" aria-hidden="true" />
    <slot />
  </span>
</template>
```

**Exemplos:**
```vue
<StatusBadge variant="success" dot>Ativo</StatusBadge>
<StatusBadge variant="warning" dot>Em Andamento</StatusBadge>
<StatusBadge variant="danger">SLA Vencido</StatusBadge>
<StatusBadge variant="neutral">Arquivado</StatusBadge>
```

---

## 5. ConfirmDialog

Dialog de confirmação padrão para ações destrutivas.

**Props:**
```ts
interface ConfirmDialogProps {
  open: boolean
  title: string
  description: string
  confirmLabel?: string     // padrão: "Confirmar"
  cancelLabel?: string      // padrão: "Cancelar"
  variant?: 'danger' | 'default'  // padrão: 'danger' para exclusões
  loading?: boolean
}
// Emits: confirm, cancel, update:open
```

**Implementação:**
```vue
<template>
  <AlertDialog :open="open" @update:open="$emit('update:open', $event)">
    <AlertDialogContent class="max-w-md">
      <AlertDialogHeader>
        <AlertDialogTitle>{{ title }}</AlertDialogTitle>
        <AlertDialogDescription>{{ description }}</AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <!-- Destrutiva à ESQUERDA para evitar clique acidental (Heurística 5) -->
        <Button
          v-if="variant === 'danger'"
          variant="danger"
          :loading="loading"
          @click="$emit('confirm')"
        >
          {{ confirmLabel ?? 'Excluir' }}
        </Button>
        <AlertDialogCancel @click="$emit('cancel')">
          {{ cancelLabel ?? 'Cancelar' }}
        </AlertDialogCancel>
        <Button
          v-if="variant !== 'danger'"
          :loading="loading"
          @click="$emit('confirm')"
        >
          {{ confirmLabel ?? 'Confirmar' }}
        </Button>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
```

---

## 6. FormField (wrapper)

Wrapper que combina label + input + mensagem de erro + help text.

**Props:**
```ts
interface FormFieldProps {
  label: string
  name: string           // para associar label ao input via for/id
  required?: boolean
  helpText?: string      // texto de ajuda abaixo do input
  error?: string         // mensagem de erro
}
// Slot: default — o input/select/textarea
```

**Implementação:**
```vue
<template>
  <div class="form-field">
    <label :for="name" class="form-label">
      {{ label }}
      <span v-if="required" class="text-danger ml-0.5" aria-hidden="true">*</span>
      <span v-if="required" class="sr-only">(obrigatório)</span>
    </label>
    <slot />
    <p v-if="helpText && !error" class="form-help">{{ helpText }}</p>
    <p v-if="error" class="form-error" role="alert">
      <AlertCircle class="inline h-3.5 w-3.5 mr-1" aria-hidden="true" />
      {{ error }}
    </p>
  </div>
</template>

<style scoped>
.form-field { display: flex; flex-direction: column; gap: 0.25rem; }
.form-label { font-size: 0.875rem; font-weight: 500; color: var(--color-neutral-700); }
.form-help  { font-size: 0.75rem; color: var(--color-neutral-500); }
.form-error { font-size: 0.75rem; color: var(--color-danger); display: flex; align-items: center; }
</style>
```

---

## 7. WizardSteps

Indicador de progresso para fluxos multi-passo (ex: criação de inbox).

**Props:**
```ts
interface WizardStepsProps {
  steps: Array<{ label: string; description?: string }>
  current: number  // 0-indexed
}
```

**Implementação:**
```vue
<template>
  <nav aria-label="Progresso" class="wizard-steps">
    <ol class="flex items-center">
      <li v-for="(step, i) in steps" :key="i" class="wizard-step">
        <div :class="stepClass(i)" :aria-current="i === current ? 'step' : undefined">
          <!-- Completo -->
          <Check v-if="i < current" class="h-4 w-4" aria-hidden="true" />
          <!-- Atual -->
          <span v-else-if="i === current" aria-hidden="true">{{ i + 1 }}</span>
          <!-- Futuro -->
          <span v-else aria-hidden="true">{{ i + 1 }}</span>
        </div>
        <span class="wizard-label">{{ step.label }}</span>
        <!-- Linha conectora -->
        <div v-if="i < steps.length - 1" :class="['wizard-line', i < current ? 'completed' : '']" />
      </li>
    </ol>
  </nav>
</template>
```

**Estilos:**
```css
.wizard-step {
  display: flex; align-items: center; gap: 8px; flex: 1;
}
.wizard-step-indicator {
  width: 2rem; height: 2rem; border-radius: 9999px;
  display: flex; align-items: center; justify-content: center;
  font-size: 0.875rem; font-weight: 600;
  border: 2px solid;
  transition: all 200ms;
}
.wizard-step-indicator.completed {
  background: var(--color-brand-500); border-color: var(--color-brand-500); color: white;
}
.wizard-step-indicator.current {
  background: white; border-color: var(--color-brand-500); color: var(--color-brand-500);
}
.wizard-step-indicator.upcoming {
  background: white; border-color: var(--color-neutral-200); color: var(--color-neutral-400);
}
.wizard-line {
  flex: 1; height: 2px; background: var(--color-neutral-200);
}
.wizard-line.completed {
  background: var(--color-brand-500);
}
```

---

## 8. SLATimer

Exibe o tempo restante de SLA com cor semântica.

**Props:**
```ts
interface SLATimerProps {
  deadline: Date | string
  unit?: 'auto' | 'minutes' | 'hours' | 'days'
  showIcon?: boolean
}
```

**Lógica:**
```ts
const timeLeft = computed(() => {
  const diff = new Date(deadline).getTime() - Date.now()
  if (diff < 0) return { label: 'Vencido', variant: 'danger', icon: 'AlertTriangle' }

  const hours = diff / 3_600_000
  if (hours < 1)  return { label: `${Math.ceil(diff/60000)}min`, variant: 'danger', icon: 'Clock' }
  if (hours < 4)  return { label: `${Math.ceil(hours)}h`, variant: 'warning', icon: 'Clock' }
  if (hours < 24) return { label: `${Math.ceil(hours)}h`, variant: 'info', icon: 'Clock' }
  return { label: `${Math.ceil(hours/24)}d`, variant: 'neutral', icon: 'Calendar' }
})
```

**Renderização:**
```vue
<StatusBadge :variant="timeLeft.variant">
  <component :is="timeLeft.icon" v-if="showIcon" class="h-3 w-3 mr-1" />
  {{ timeLeft.label }}
</StatusBadge>
```

---

## 9. BulkActionsBar

Barra de ações que aparece quando há itens selecionados na tabela.

**Props:**
```ts
interface BulkActionsBarProps {
  selectedCount: number
  actions: Array<{
    label: string
    icon?: string
    variant?: 'default' | 'danger'
    onClick: () => void
  }>
}
```

**Implementação:**
```vue
<Transition name="slide-up">
  <div v-if="selectedCount > 0" class="bulk-bar" role="toolbar" :aria-label="`${selectedCount} itens selecionados`">
    <span class="text-sm font-medium text-neutral-700">
      {{ selectedCount }} {{ selectedCount === 1 ? 'item selecionado' : 'itens selecionados' }}
    </span>
    <div class="flex gap-2">
      <Button
        v-for="action in actions"
        :key="action.label"
        :variant="action.variant ?? 'secondary'"
        size="sm"
        @click="action.onClick"
      >
        <component :is="action.icon" v-if="action.icon" class="h-3.5 w-3.5 mr-1.5" />
        {{ action.label }}
      </Button>
    </div>
  </div>
</Transition>

<style scoped>
.bulk-bar {
  position: sticky; bottom: 0;
  display: flex; align-items: center; justify-content: space-between;
  padding: 0.75rem 1rem;
  background: var(--color-neutral-800);
  color: white;
  border-radius: var(--radius-lg) var(--radius-lg) 0 0;
  box-shadow: var(--shadow-xl);
}
.slide-up-enter-active, .slide-up-leave-active { transition: transform 200ms, opacity 200ms; }
.slide-up-enter-from, .slide-up-leave-to { transform: translateY(100%); opacity: 0; }
</style>
```

---

## 10. ConnectionTestButton

Botão para testar conexões (SMTP, IMAP, OIDC, etc).

**Props:**
```ts
interface ConnectionTestButtonProps {
  testFn: () => Promise<{ ok: boolean; message?: string }>
  label?: string  // padrão: "Testar Conexão"
}
```

**Estados:**
```
Inicial:     [🧪 Testar Conexão]          ← secondary button
Testando:    [⟳ Testando...]              ← loading spinner, disabled
Sucesso:     [✓ Conexão bem-sucedida]     ← success variant, auto-reset 3s
Erro:        [✗ Falha na conexão]         ← danger variant + tooltip com detalhe
```

**Implementação:**
```vue
<script setup>
const status = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const errorMessage = ref('')

async function test() {
  status.value = 'loading'
  try {
    const result = await props.testFn()
    status.value = result.ok ? 'success' : 'error'
    errorMessage.value = result.message ?? ''
  } catch (e) {
    status.value = 'error'
    errorMessage.value = 'Não foi possível conectar.'
  }
  if (status.value === 'success') {
    setTimeout(() => { status.value = 'idle' }, 3000)
  }
}
</script>

<template>
  <TooltipProvider>
    <Tooltip :disabled="status !== 'error'">
      <TooltipTrigger as-child>
        <Button
          :variant="status === 'success' ? 'success' : status === 'error' ? 'danger' : 'secondary'"
          :loading="status === 'loading'"
          :disabled="status === 'loading'"
          @click="test"
        >
          <Check v-if="status === 'success'" class="h-4 w-4 mr-1.5" />
          <X v-else-if="status === 'error'" class="h-4 w-4 mr-1.5" />
          <FlaskConical v-else class="h-4 w-4 mr-1.5" />
          {{ status === 'success' ? 'Conexão OK' : status === 'error' ? 'Falha na conexão' : (label ?? 'Testar Conexão') }}
        </Button>
      </TooltipTrigger>
      <TooltipContent v-if="errorMessage">{{ errorMessage }}</TooltipContent>
    </Tooltip>
  </TooltipProvider>
</template>
```

---

## 11. Padrões de Animação

```css
/* Transições padrão */
:root {
  --transition-fast:   150ms ease;
  --transition-base:   200ms ease;
  --transition-slow:   300ms ease;
  --transition-spring: 300ms cubic-bezier(0.34, 1.56, 0.64, 1); /* slight overshoot */
}

/* Só animar se o usuário não preferir redução de movimento */
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }
}

/* Page transitions */
.page-enter-active { transition: opacity var(--transition-base); }
.page-enter-from   { opacity: 0; }

/* Modal */
.dialog-enter-active { transition: all var(--transition-base); }
.dialog-enter-from   { opacity: 0; transform: scale(0.95) translateY(-8px); }

/* Sidebar item hover */
.nav-item { transition: background var(--transition-fast), color var(--transition-fast); }

/* Toasts */
.toast-enter-active { transition: all var(--transition-spring); }
.toast-enter-from   { opacity: 0; transform: translateX(calc(100% + 1rem)); }
.toast-leave-active { transition: all var(--transition-base); }
.toast-leave-to     { opacity: 0; transform: translateX(calc(100% + 1rem)); }
```

---

## 12. Checklist de Qualidade (pre-merge)

Antes de fazer PR de qualquer tela admin:

### Visual
- [ ] PageHeader presente com breadcrumb, H1 e description
- [ ] Botão primário segue padrão (verbo + substantivo, variant="default")
- [ ] Cores usando design tokens (não hex hardcoded)
- [ ] Ícones são do lucide-vue-next (não misturar libs)
- [ ] Espaçamentos em múltiplos de 4px (gap-2, p-4, etc.)

### Comportamento
- [ ] Estado loading implementado (skeleton ou botão spinner)
- [ ] Estado empty implementado com CTA
- [ ] Estado de erro com mensagem clara e ação
- [ ] Confirmação antes de DELETE
- [ ] Toast de feedback pós-ação (sucesso e erro)
- [ ] Formulário com validação inline (não apenas no submit)
- [ ] Submit desabilitado enquanto form inválido

### Acessibilidade
- [ ] Todo input tem label associado (for + id)
- [ ] Ícones decorativos com aria-hidden="true"
- [ ] Ícones funcionais com aria-label
- [ ] Erros com role="alert"
- [ ] Focus visível em todos os elementos interativos
- [ ] Tab order lógico (sem tabindex > 0)
- [ ] Contraste mínimo 4.5:1 para texto

### i18n
- [ ] Todos os textos via $t() ou t()
- [ ] Chave adicionada em i18n/en.json
- [ ] Chave adicionada em i18n/pt-BR.json
- [ ] Sem string literal em português no template
