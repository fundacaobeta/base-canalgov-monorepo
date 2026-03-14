# CanalGov — Design System

> Sistema de design governamental white-label para plataformas de atendimento cidadão.
> Baseado no DS GOV.BR, Heurísticas de Nielsen, Leis de UX e WCAG 2.1 AA.

---

## 1. Fundamentos

### 1.1 Princípios de Design

| Princípio | Definição | Aplicação |
|---|---|---|
| **Confiança** | Servidores e cidadãos precisam confiar no sistema | Visual limpo, oficial, sem adornos desnecessários |
| **Clareza** | Cada tela tem um objetivo primário | Hierarquia visual forte; sem informação redundante |
| **Acessibilidade** | Qualquer pessoa deve conseguir usar | WCAG 2.1 AA obrigatório; ARIA correto |
| **Consistência** | Mesmos padrões em todas as páginas | Design tokens; componentes reutilizáveis |
| **Eficiência** | Servidores usam o sistema dezenas de vezes ao dia | Atalhos; ações em bulk; filtros persistidos |

---

### 1.2 Heurísticas de Nielsen — Checklist de Implementação

- [ ] **H1 – Visibilidade do status**: Loading skeleton em tabelas, spinner em botões, badges de prazo SLA, step indicators em wizards
- [ ] **H2 – Correspondência com o mundo real**: Terminologia pt-BR (Protocolo, Demanda, Expediente, Prazo), ícones semânticos, datas em dd/mm/aaaa
- [ ] **H3 – Controle e liberdade**: Botão Cancelar em todo modal, confirmação antes de delete, undo em ações reversíveis, breadcrumb em toda página
- [ ] **H4 – Consistência e padrões**: Mesma posição para botão primário, mesmo padrão de tabela, mesmo layout de formulário
- [ ] **H5 – Prevenção de erros**: Validação inline ao blur, submit desabilitado se inválido, confirm dialog para ações destrutivas, preview antes de salvar
- [ ] **H6 – Reconhecimento em vez de lembrança**: Labels sempre visíveis, help text contextual, breadcrumb, tooltips em ícones sem label
- [ ] **H7 – Flexibilidade e eficiência**: Ações em bulk na tabela, filtros rápidos, busca global, atalhos de teclado documentados
- [ ] **H8 – Design estético e minimalista**: Máx 3 ações primárias visíveis, whitespace generoso, paleta limitada
- [ ] **H9 – Ajuda no reconhecimento de erros**: Erro inline com texto claro + ação sugerida, highlight do campo problemático
- [ ] **H10 – Ajuda e documentação**: Painel de ajuda à direita, links externos, tooltips de contexto

---

## 2. Design Tokens

### 2.1 Cores

O sistema usa **CSS Custom Properties** para suportar white-label. O arquivo `frontend/src/style.css` define os tokens; o painel Admin > Geral pode sobrescrever `--color-brand-*` via settings.

```css
/* ─── Marca (white-label) ──────────────────────────────── */
:root {
  --color-brand-50:   #EEF4FF;
  --color-brand-100:  #D9E8FF;
  --color-brand-200:  #ADCFFF;
  --color-brand-300:  #6BA5F5;
  --color-brand-400:  #3B7DE8;
  --color-brand-500:  #1351B4;  /* padrão gov.br */
  --color-brand-600:  #0D3F8F;
  --color-brand-700:  #0A2D6B;
  --color-brand-800:  #071D47;
  --color-brand-900:  #040E24;

  /* ─── Semânticas ────────────────────────────────────── */
  --color-success-bg:   #EAFAF0;
  --color-success:      #168821;
  --color-success-dark: #0F5F17;

  --color-warning-bg:   #FFF9E6;
  --color-warning:      #B47800;
  --color-warning-dark: #7A5200;

  --color-danger-bg:    #FFF0EE;
  --color-danger:       #E52207;
  --color-danger-dark:  #A31505;

  --color-info-bg:      #EEF4FF;
  --color-info:         #155BCB;
  --color-info-dark:    #0D3F8F;

  /* ─── Neutros ───────────────────────────────────────── */
  --color-neutral-0:   #FFFFFF;
  --color-neutral-50:  #F8F9FA;
  --color-neutral-100: #EDEDED;
  --color-neutral-200: #D8D8D8;
  --color-neutral-300: #B0B0B0;
  --color-neutral-400: #888888;
  --color-neutral-500: #666666;
  --color-neutral-600: #444444;
  --color-neutral-700: #333333;
  --color-neutral-800: #1B1B1B;
  --color-neutral-900: #000000;

  /* ─── Superfícies ────────────────────────────────────── */
  --color-surface-body:    var(--color-neutral-50);
  --color-surface-card:    var(--color-neutral-0);
  --color-surface-sidebar: var(--color-neutral-0);
  --color-surface-header:  var(--color-neutral-0);
  --color-surface-overlay: rgba(0,0,0,0.4);

  /* ─── Bordas ─────────────────────────────────────────── */
  --color-border-subtle:  var(--color-neutral-100);
  --color-border-default: var(--color-neutral-200);
  --color-border-strong:  var(--color-neutral-300);
}
```

#### Contraste mínimo (WCAG 2.1 AA)
| Uso | Cor texto | Cor fundo | Contraste |
|---|---|---|---|
| Corpo | `neutral-800` | `neutral-0` | 14:1 ✓ |
| Label | `neutral-700` | `neutral-0` | 10:1 ✓ |
| Placeholder | `neutral-400` | `neutral-0` | 4.7:1 ✓ |
| Brand em branco | `neutral-0` | `brand-500` | 5.2:1 ✓ |
| Danger texto | `danger-dark` | `danger-bg` | 4.8:1 ✓ |

---

### 2.2 Tipografia

```css
:root {
  --font-sans: 'Inter', system-ui, -apple-system, 'Segoe UI', sans-serif;
  --font-mono: 'JetBrains Mono', 'Fira Code', 'Courier New', monospace;

  /* Escala (base 16px, ratio 1.25 — Major Third) */
  --text-2xs:  0.625rem;  /*  10px — legal, footnotes */
  --text-xs:   0.75rem;   /*  12px — badges, timestamps */
  --text-sm:   0.875rem;  /*  14px — tabelas, help text */
  --text-base: 1rem;      /*  16px — corpo padrão */
  --text-lg:   1.125rem;  /*  18px — lead text */
  --text-xl:   1.25rem;   /*  20px — subtítulos de seção */
  --text-2xl:  1.5rem;    /*  24px — título de página */
  --text-3xl:  1.875rem;  /*  30px — destaque (dashboard) */

  /* Pesos */
  --font-normal:   400;
  --font-medium:   500;
  --font-semibold: 600;
  --font-bold:     700;

  /* Line heights */
  --leading-tight:  1.25;
  --leading-snug:   1.375;
  --leading-normal: 1.5;
  --leading-relaxed:1.625;
}
```

**Hierarquia de texto nas páginas:**
- `text-2xl font-bold` → Título da página (H1)
- `text-sm text-neutral-500` → Descrição da página
- `text-xl font-semibold` → Título de seção (H2)
- `text-base font-medium` → Título de card / grupo (H3)
- `text-sm` → Corpo, tabelas, inputs
- `text-xs text-neutral-400` → Metadados, timestamps, badges

---

### 2.3 Espaçamento

Grid de **4px**. Tailwind classes mapeadas:

| Token | Valor | Tailwind | Uso |
|---|---|---|---|
| 1 | 4px | `p-1 gap-1` | Intra-elemento |
| 2 | 8px | `p-2 gap-2` | Ícone + label |
| 3 | 12px | `p-3 gap-3` | Badge padding |
| 4 | 16px | `p-4 gap-4` | Padding de card |
| 5 | 20px | `p-5 gap-5` | Espaço entre campos de form |
| 6 | 24px | `p-6 gap-6` | Padding de seção |
| 8 | 32px | `p-8 gap-8` | Entre seções |
| 10 | 40px | `p-10` | Padding de página |
| 12 | 48px | `p-12` | Espaço extra-large |
| 16 | 64px | `p-16` | Seções de destaque |

---

### 2.4 Elevação e Sombras

```css
:root {
  --shadow-xs: 0 1px 2px rgba(0,0,0,0.04);                          /* inputs */
  --shadow-sm: 0 1px 3px rgba(0,0,0,0.06), 0 1px 2px rgba(0,0,0,0.04); /* cards */
  --shadow-md: 0 4px 8px rgba(0,0,0,0.06), 0 2px 4px rgba(0,0,0,0.04); /* dropdowns */
  --shadow-lg: 0 10px 20px rgba(0,0,0,0.07),0 4px 8px rgba(0,0,0,0.04); /* modais */
  --shadow-xl: 0 24px 48px rgba(0,0,0,0.08),0 8px 16px rgba(0,0,0,0.04);/* side panels */
  --shadow-focus: 0 0 0 3px rgba(19,81,180,0.25);                    /* focus ring */
}
```

---

### 2.5 Border Radius

```css
:root {
  --radius-xs:   2px;    /* badges pequenos */
  --radius-sm:   4px;    /* tags, chips */
  --radius-md:   6px;    /* inputs, botões, cards (padrão) */
  --radius-lg:   8px;    /* modais, painéis */
  --radius-xl:   12px;   /* cards de destaque */
  --radius-2xl:  16px;   /* banners */
  --radius-full: 9999px; /* pills, avatares */
}
```

---

## 3. Componentes

### 3.1 Botões

```
Variantes:
  default   → brand-500 bg, branco text        (ação primária)
  secondary → branco bg, neutral-300 border    (ação secundária)
  ghost     → transparente bg, brand-500 text  (ação terciária/link)
  danger    → danger bg, branco text           (ações destrutivas)
  outline   → transparente bg, danger border   (destructive secundária)

Tamanhos:
  sm  → h-8  px-3 text-xs   (ações em tabela)
  md  → h-9  px-4 text-sm   (padrão)
  lg  → h-10 px-5 text-base (CTA principal)
  icon → h-9 w-9            (ícone sem label)

Estados:
  :disabled → opacity-50 cursor-not-allowed
  :loading  → spinner icon substituindo label (não desabilitar o botão)
  :focus-visible → shadow-focus ring-2
```

### 3.2 Inputs / Campos de Formulário

**Anatomia:**
```
[Label*]
[Placeholder text                               ] ← input
[Texto de ajuda / ⚠ Mensagem de erro            ]
```

**Regras:**
- Label: `text-sm font-medium text-neutral-700`
- Input: `h-9 px-3 border border-neutral-200 rounded-md text-sm`
- Focus: `border-brand-500 shadow-focus outline-none`
- Erro: `border-danger` + mensagem `text-xs text-danger mt-1`
- Disabled: `bg-neutral-50 text-neutral-400 cursor-not-allowed`
- Obrigatório: asterisco vermelho após label (`text-danger`)

### 3.3 Data Table

**Estrutura:**
```
┌─ Toolbar ──────────────────────────────────────────────────┐
│  [🔍 Buscar...]   [Filtros ▼]              [+ Nova ação]  │
└────────────────────────────────────────────────────────────┘
┌─ Table ────────────────────────────────────────────────────┐
│  ☐ │ Coluna A ↕ │ Coluna B ↕ │ Coluna C │  Ações        │
│────┼────────────┼────────────┼──────────┼───────────────│
│  ☐ │ valor      │ valor      │ Badge    │  ✏ ···        │
│  ☐ │ valor      │ valor      │ Badge    │  ✏ ···        │
└────────────────────────────────────────────────────────────┘
┌─ Footer ───────────────────────────────────────────────────┐
│  Mostrando 1–20 de 143 resultados       ← 1 2 3 ... 8 →  │
└────────────────────────────────────────────────────────────┘
```

**Estados obrigatórios:**
- **Loading**: skeleton rows (3-5 linhas) — NÃO spinner sobre a tabela
- **Empty (sem dados)**: ilustração centralizada + texto + CTA
- **Empty (sem resultados de busca)**: ícone busca + "Nenhum resultado para '[termo]'" + botão "Limpar filtros"
- **Error**: banner de erro + botão "Tentar novamente"
- **Seleção múltipla**: barra de ação aparece no topo quando há ☐ marcados

### 3.4 Modal / Dialog

```
Tamanhos:
  sm (400px)  → Confirmação de exclusão
  md (560px)  → Formulário simples (padrão)
  lg (720px)  → Formulário complexo, previews
  xl (900px)  → Builders (automações, macros)
  full        → Editores (TipTap templates)

Estrutura obrigatória:
  Header: título (verb + noun) + botão X fechar
  Body:   conteúdo scrollável se > 70vh
  Footer: ações alinhadas à direita
            [Cancelar]  [Ação Principal]
          Para destrutivas:
            [Excluir]   [Cancelar]
          (destrutiva à esquerda para evitar clique acidental)
```

### 3.5 Badges / Pills de Status

```css
/* Base */
.badge {
  display: inline-flex; align-items: center; gap: 4px;
  padding: 2px 8px; border-radius: 9999px; border: 1px solid;
  font-size: 0.75rem; font-weight: 500; line-height: 1.5;
  white-space: nowrap;
}

/* Variantes semânticas */
.badge-success { background: var(--color-success-bg); color: var(--color-success-dark); border-color: #A3D9A5; }
.badge-warning { background: var(--color-warning-bg); color: var(--color-warning-dark); border-color: #F5D98A; }
.badge-danger  { background: var(--color-danger-bg);  color: var(--color-danger-dark);  border-color: #F4A89A; }
.badge-info    { background: var(--color-info-bg);    color: var(--color-info-dark);    border-color: #ADCFFF; }
.badge-neutral { background: var(--color-neutral-100); color: var(--color-neutral-600); border-color: var(--color-neutral-200); }
.badge-brand   { background: var(--color-brand-50);   color: var(--color-brand-700);   border-color: var(--color-brand-200); }
```

---

## 4. Layout do Sistema Admin

### 4.1 Shell Principal

```
┌─────────────────────────────────────────────────────────────┐
│  HEADER                                                      │
│  [Logo]  CanalGov — Prefeitura X    [Notif] [Avatar] [?]   │
├────────────────────┬────────────────────────────────────────┤
│  SIDEBAR (240px)   │  MAIN CONTENT AREA                     │
│                    │                                        │
│  ▼ Workspace       │  PageHeader                            │
│    General         │  ─────────────────────────────────     │
│    Horários        │  Conteúdo                              │
│    SLA             │                                        │
│  ▼ Conversas       │                                        │
│    Tags            │                                        │
│    Status          │                                        │
│    ...             │                                        │
│                    │                                        │
│  ─────────────── │                                        │
│  [Avatar]          │                                        │
│  Nome Servidor     │                                        │
│  ● Online ▼        │                                        │
└────────────────────┴────────────────────────────────────────┘
```

### 4.2 Page Header (padrão obrigatório)

```vue
<div class="flex items-start justify-between mb-6">
  <div>
    <Breadcrumb :items="breadcrumbs" class="mb-1" />
    <h1 class="text-2xl font-bold text-neutral-800">Título</h1>
    <p class="text-sm text-neutral-500 mt-1">Descrição breve.</p>
  </div>
  <div class="flex gap-2">
    <!-- Ações secundárias antes, primária por último -->
    <Button variant="secondary">Importar</Button>
    <Button variant="default">+ Criar Novo</Button>
  </div>
</div>
```

### 4.3 Layout com Painel de Ajuda

```vue
<div class="grid grid-cols-1 lg:grid-cols-[1fr_320px] gap-8">
  <div><!-- Conteúdo principal --></div>
  <aside>
    <HelpPanel>
      <HelpItem icon="info" title="O que são equipes?">
        Equipes agrupam atendentes e organizam filas de atendimento.
      </HelpItem>
      <HelpItem icon="link" title="Documentação" :href="docsUrl" />
    </HelpPanel>
  </aside>
</div>
```

---

## 5. Sidebar de Navegação

### 5.1 Estrutura de Grupos

```
Workspace
  ├ General (⚙)
  ├ Horários de Atendimento (🕐)
  └ SLA (📊)

Conversas
  ├ Tags (🏷)
  ├ Status (●)
  ├ Macros (⚡)
  └ Visões Compartilhadas (👁)

Canais
  ├ Caixas de Entrada (📥)
  └ Domínios (🌐)

Equipes & Agentes
  ├ Agentes (👤)
  ├ Equipes (👥)
  └ Papéis (🔑)

Automação
  ├ Automações (🤖)
  └ Templates (📄)

Segurança
  ├ SSO / OIDC (🔒)
  ├ Log de Atividades (📋)
  └ Papéis & Permissões (🛡)

Notificações
  ├ E-mail (✉)
  ├ WhatsApp (💬)
  ├ Telegram (📱)
  ├ SMS (📟)
  ├ Push (🔔)
  └ Comunicações Oficiais (🏛)

Integrações
  ├ Webhooks (🔗)
  └ Ações (⚡)
```

### 5.2 CSS do Item de Navegação

```css
.nav-item {
  display: flex; align-items: center; gap: 10px;
  padding: 8px 12px; border-radius: 6px;
  font-size: 0.875rem; font-weight: 500; color: var(--color-neutral-600);
  transition: background 150ms, color 150ms;
  cursor: pointer; text-decoration: none;
  border-left: 3px solid transparent;
}
.nav-item:hover {
  background: var(--color-neutral-50);
  color: var(--color-neutral-800);
}
.nav-item.active {
  background: var(--color-brand-50);
  color: var(--color-brand-700);
  border-left-color: var(--color-brand-500);
  font-weight: 600;
}
.nav-group-label {
  font-size: 0.6875rem; font-weight: 700; letter-spacing: 0.06em;
  text-transform: uppercase; color: var(--color-neutral-400);
  padding: 16px 12px 4px;
}
```

---

## 6. Tokens Tailwind (tailwind.config.js)

```js
module.exports = {
  theme: {
    extend: {
      colors: {
        brand: {
          50:  'var(--color-brand-50)',
          100: 'var(--color-brand-100)',
          500: 'var(--color-brand-500)',
          600: 'var(--color-brand-600)',
          700: 'var(--color-brand-700)',
        },
        neutral: {
          50:  'var(--color-neutral-50)',
          100: 'var(--color-neutral-100)',
          200: 'var(--color-neutral-200)',
          400: 'var(--color-neutral-400)',
          500: 'var(--color-neutral-500)',
          600: 'var(--color-neutral-600)',
          700: 'var(--color-neutral-700)',
          800: 'var(--color-neutral-800)',
        },
        success: 'var(--color-success)',
        warning: 'var(--color-warning)',
        danger:  'var(--color-danger)',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
      boxShadow: {
        focus: 'var(--shadow-focus)',
      },
    },
  },
}
```

---

## 7. Microcopy — Guia de Voz

| Contexto | ❌ Evitar | ✓ Usar |
|---|---|---|
| Botão criar | "Salvar", "Criar", "OK" | "Criar Equipe", "Adicionar Tag" |
| Botão cancelar | "Voltar", "Não" | "Cancelar" (sempre) |
| Confirmação delete | "Deseja excluir?" | "Excluir **Nome**? Esta ação não pode ser desfeita." |
| Estado vazio | "Sem dados", "Nenhum" | "Nenhuma equipe cadastrada ainda. [Criar primeira equipe]" |
| Erro de campo | "Campo inválido" | "Informe um e-mail válido (ex: nome@gov.br)" |
| Erro de API | "Erro ao salvar" | "Não foi possível salvar. Verifique sua conexão e tente novamente." |
| Loading | "Carregando..." | Skeleton (sem texto), ou "Carregando equipes..." em contexto específico |
| Sucesso | "OK!", "Pronto!" | "Equipe criada com sucesso." |
| Tooltip botão | "Clique aqui" | "Editar equipe", "Excluir tag" |

---

## 8. Acessibilidade — Checklist

### HTML Semântico
- [ ] `<main>` envolve o conteúdo principal
- [ ] `<nav>` com `aria-label="Navegação administrativa"` na sidebar
- [ ] `<h1>` único por página (título da página)
- [ ] Hierarquia h1→h2→h3 respeitada (sem pular níveis)
- [ ] Tabelas com `<thead>`, `<th scope="col">`, `<caption>` opcional

### Interatividade
- [ ] Todo elemento interativo acessível por teclado (Tab, Enter, Space, Esc)
- [ ] `focus-visible` visível em todos os elementos
- [ ] Modais com `focus-trap` e `aria-modal="true"`
- [ ] Dropdowns com `aria-expanded` e `aria-haspopup`
- [ ] Mensagens de erro com `role="alert"`
- [ ] Notificações toast com `aria-live="polite"`

### Imagens e Ícones
- [ ] Ícones decorativos: `aria-hidden="true"`
- [ ] Ícones funcionais: `aria-label="Editar equipe"`
- [ ] Avatares: `alt="Nome do usuário"`
- [ ] Logos: `alt="CanalGov — [Nome do órgão]"`

### Contraste e Visual
- [ ] Modo de alto contraste suportado (prefers-contrast)
- [ ] Não depender apenas de cor para transmitir estado (usar ícone + cor + texto)
- [ ] Tamanho mínimo de toque: 44×44px em mobile

---

## 9. Dark Mode (futuro)

Estrutura preparada via CSS Custom Properties — basta sobrescrever em `[data-theme="dark"]`:
```css
[data-theme="dark"] {
  --color-surface-body:    #0F1117;
  --color-surface-card:    #1A1D27;
  --color-surface-sidebar: #13151F;
  --color-neutral-0:       #1A1D27;
  --color-neutral-800:     #F0F2F5;
  /* ... */
}
```
