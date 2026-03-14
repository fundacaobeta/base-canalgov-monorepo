---
name: design-system
description: Apply the CanalGov government design system — tokens, components, layout, accessibility, Nielsen heuristics — when creating or refactoring Vue components and admin pages
argument-hint: [component-or-page]
allowed-tools: Read, Grep, Glob, Edit, Write
---

# CanalGov Design System Skill

Você é um designer/engenheiro front-end especializado no design system governamental CanalGov.
Quando aplicar este skill em `$ARGUMENTS`, siga **todos** os princípios abaixo.

---

## 1. Princípios Fundamentais

### Heurísticas de Nielsen (obrigatórias)
| # | Heurística | Como aplicar no CanalGov |
|---|---|---|
| 1 | Visibilidade do status | Loading states em todo fetch; badges de SLA; progress steps em wizards |
| 2 | Correspondência com o mundo real | Labels em pt-BR; terminologia governamental (Protocolo, Demanda, Expediente) |
| 3 | Controle e liberdade | Botão cancelar em todo modal/form; confirmação antes de deletar |
| 4 | Consistência e padrões | Mesmo padrão de tabela/form/modal em todas as páginas |
| 5 | Prevenção de erros | Validação inline; desabilitar submit em form inválido; confirm em ações destrutivas |
| 6 | Reconhecimento em vez de lembrança | Labels sempre visíveis; tooltips em ícones; breadcrumb em todo header |
| 7 | Flexibilidade e eficiência | Atalhos de teclado; ações em bulk; filtros persistidos na URL |
| 8 | Design estético e minimalista | Sem informação redundante; hierarquia visual clara; whitespace generoso |
| 9 | Ajuda no reconhecimento de erros | Mensagens de erro em linguagem simples + ação sugerida |
| 10 | Ajuda e documentação | Painel de ajuda contextual à direita; links para docs |

### Leis de UX aplicadas
- **Lei de Fitts**: Botões primários grandes (min 44×44px); ações destrutivas pequenas/afastadas
- **Lei de Hick**: Máx 7 itens por nível de navegação; agrupar opções avançadas em accordion
- **Lei de Miller**: Chunks de informação (ex: protocolo em formato `AAAA.ÓRGÃO.NNNNN`)
- **Efeito de posição serial**: Ações mais importantes no topo e no fim da lista
- **Lei de proximidade (Gestalt)**: Campos relacionados agrupados em fieldsets com label
- **Lei de Jakob**: Seguir padrões do DS GOV.BR que servidores já conhecem

---

## 2. Design Tokens

### Cores (CSS Custom Properties em `frontend/src/style.css` ou Tailwind config)

```css
/* Paleta Primária — White-label (customizável via settings) */
--color-brand-50:  #EEF4FF;
--color-brand-100: #D9E8FF;
--color-brand-200: #ADCFFF;
--color-brand-300: #6BA5F5;
--color-brand-400: #3B7DE8;
--color-brand-500: #1351B4;  /* Azul Gov.BR padrão */
--color-brand-600: #0D3F8F;
--color-brand-700: #0A2D6B;
--color-brand-800: #071D47;
--color-brand-900: #040E24;

/* Semânticas */
--color-success:   #168821;
--color-warning:   #FFCD07;
--color-danger:    #E52207;
--color-info:      #155BCB;

/* Neutros */
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
```

### Tipografia
```css
--font-family-base: 'Inter', system-ui, -apple-system, sans-serif;
--font-family-mono: 'JetBrains Mono', 'Fira Code', monospace;

/* Escala Modular (base 16px, ratio 1.25) */
--text-xs:   0.75rem;   /* 12px — labels, badges */
--text-sm:   0.875rem;  /* 14px — corpo, tabelas */
--text-base: 1rem;      /* 16px — default */
--text-lg:   1.125rem;  /* 18px — subtítulos */
--text-xl:   1.25rem;   /* 20px — títulos de seção */
--text-2xl:  1.5rem;    /* 24px — títulos de página */
--text-3xl:  1.875rem;  /* 30px — títulos de destaque */

/* Pesos */
--font-normal:   400;
--font-medium:   500;
--font-semibold: 600;
--font-bold:     700;
```

### Espaçamento (grid de 4px)
```css
--space-1:  0.25rem;  /* 4px */
--space-2:  0.5rem;   /* 8px */
--space-3:  0.75rem;  /* 12px */
--space-4:  1rem;     /* 16px */
--space-5:  1.25rem;  /* 20px */
--space-6:  1.5rem;   /* 24px */
--space-8:  2rem;     /* 32px */
--space-10: 2.5rem;   /* 40px */
--space-12: 3rem;     /* 48px */
--space-16: 4rem;     /* 64px */
```

### Elevação / Sombras
```css
--shadow-sm:  0 1px 2px rgba(0,0,0,0.05);
--shadow-md:  0 4px 6px rgba(0,0,0,0.07), 0 1px 3px rgba(0,0,0,0.06);
--shadow-lg:  0 10px 15px rgba(0,0,0,0.07), 0 4px 6px rgba(0,0,0,0.05);
--shadow-xl:  0 20px 25px rgba(0,0,0,0.07), 0 10px 10px rgba(0,0,0,0.04);
```

### Border Radius
```css
--radius-sm:   0.25rem;  /* 4px */
--radius-md:   0.375rem; /* 6px — inputs, cards */
--radius-lg:   0.5rem;   /* 8px — modals */
--radius-xl:   0.75rem;  /* 12px — painéis */
--radius-full:  9999px;  /* pills/badges */
```

---

## 3. Componentes — Padrões

### Page Header (todo admin page)
```vue
<PageHeader>
  <Breadcrumb items="[{ label: 'Admin', to: '/admin' }, { label: 'Equipes' }]" />
  <h1>Título da Página</h1>
  <p class="text-neutral-500 text-sm">Descrição breve do que esta página faz.</p>
  <template #actions>
    <Button variant="default">+ Nova Equipe</Button>
  </template>
</PageHeader>
```

### Data Table padrão
- Cabeçalho com título + botão de ação primária (canto direito)
- Linha de pesquisa/filtros abaixo do cabeçalho (colapsável)
- Colunas: max 5-6 visíveis; resto em "..." dropdown
- Ações por linha: ícone edit + dropdown (3 pontos) para mais
- Estado vazio: ilustração + texto explicativo + CTA
- Paginação: "Mostrando X-Y de Z resultados" + controles
- Loading: skeleton rows (não spinner sobre a tabela)

### Formulários
- Label sempre acima do input (não floating)
- Placeholder apenas para formato/exemplo (ex: "ex: atendimento@prefeitura.gov.br")
- Mensagem de erro inline abaixo do input (vermelho, ícone ⚠)
- Campos obrigatórios: asterisco (*) no label, não no placeholder
- Fieldsets com `<legend>` para grupos de campos relacionados
- Botões: primário à esquerda (submit), secundário ao lado (cancelar)
- Largura máxima de formulário standalone: 640px

### Modais / Dialogs
- Título claro (verbo + substantivo: "Criar Tag", "Editar Equipe")
- Sem informação desnecessária no body
- Ações no footer: destrutivas (vermelho) à esquerda, confirmação à direita
- Esc para fechar; click fora fecha (exceto forms com dados)
- Tamanho: sm (400px) para confirmações, md (560px) padrão, lg (720px) para forms complexos

### Badges / Status
```
Status da Conversa:
  ● Aberto      → bg-blue-50  text-blue-700  border-blue-200
  ● Em Andamento → bg-yellow-50 text-yellow-700 border-yellow-200
  ● Resolvido   → bg-green-50 text-green-700  border-green-200
  ● Arquivado   → bg-gray-50  text-gray-600  border-gray-200

SLA Status:
  ▲ Vencido     → bg-red-50   text-red-700   border-red-200
  ◆ Em risco    → bg-orange-50 text-orange-700 border-orange-200
  ✓ No prazo    → bg-green-50 text-green-700  border-green-200
```

### Sidebar de Navegação
- Largura: 240px expandida, 64px colapsada (ícone + tooltip)
- Grupos com label de seção em uppercase xs
- Item ativo: fundo brand-50, borda esquerda 3px brand-500, texto brand-700
- Hover: fundo neutral-100
- Ícones: lucide-vue (24px, stroke-width 1.5)

---

## 4. Padrões de Layout Admin

### Layout de Página com Ajuda (padrão atual, manter)
```
┌─────────────────────────────────────┬────────────────┐
│  PageHeader                          │                │
│  ─────────────────────────────────  │   Painel de    │
│  Conteúdo principal                  │   Ajuda        │
│  (tabela / form / lista)             │   Contextual   │
│                                     │                │
└─────────────────────────────────────┴────────────────┘
```

### Layout de Formulário Dedicado (criar/editar)
```
┌────────────────────────────────────────────────────────┐
│  ← Voltar para [Lista]                                 │
│  ─────────────────────────────────────────────────     │
│  Título: Criar/Editar X                                │
│  ┌────────────────────────────────┐  ┌──────────────┐ │
│  │  Formulário principal          │  │  Resumo /    │ │
│  │  (max-w-2xl)                   │  │  Preview     │ │
│  │                                │  │  (opcional)  │ │
│  │  [Salvar]  [Cancelar]          │  └──────────────┘ │
│  └────────────────────────────────┘                    │
└────────────────────────────────────────────────────────┘
```

---

## 5. Acessibilidade (WCAG 2.1 AA — obrigatório para governo)

- **Contraste**: mínimo 4.5:1 para texto normal, 3:1 para texto grande/ícones
- **Foco visível**: `focus-visible:ring-2 ring-brand-500 ring-offset-2` em todo elemento interativo
- **Landmarks ARIA**: `<main>`, `<nav>`, `<aside>` corretos
- **Labels associados**: todo `<input>` com `<label for="">` ou `aria-label`
- **Mensagens de erro**: `role="alert"` ou `aria-live="polite"` para anúncio de screen reader
- **Tabelas**: `<th scope="col">` para cabeçalhos; `<caption>` descritiva
- **Modais**: foco aprisionado dentro (`focus-trap`); `aria-modal="true"`; `aria-labelledby` no título
- **Ícones decorativos**: `aria-hidden="true"`; ícones funcionais com `aria-label`

---

## 6. Responsividade

| Breakpoint | Largura | Comportamento |
|---|---|---|
| mobile | < 768px | Sidebar colapsada; tabelas scrolláveis; forms full-width |
| tablet | 768-1024px | Sidebar colapsada por padrão; grid 2 colunas |
| desktop | > 1024px | Sidebar expandida; layout full |
| wide | > 1440px | Max-width container; não esticar além de 1280px de conteúdo |

---

## 7. Microcopy Governamental

- Títulos de página: imperativos ou substantivos ("Gerenciar Equipes", "Equipes")
- Botões: verbo + objeto ("Criar Equipe", não "Criar" ou "Salvar")
- Estados vazios: título + explicação + CTA ("Nenhuma equipe cadastrada. Crie a primeira equipe para organizar seus atendentes.")
- Confirmações de exclusão: "Tem certeza que deseja excluir **[Nome]**? Esta ação não pode ser desfeita."
- Erros: "Não foi possível salvar. [Motivo específico]. [Ação sugerida]."
- Sucesso: "[Recurso] [ação com sucesso]." ex: "Equipe criada com sucesso."

---

## Processo de Aplicação

1. Ler o componente alvo
2. Identificar violações de heurísticas e padrões acima
3. Aplicar tokens, corrigir estrutura HTML, melhorar acessibilidade
4. Verificar se i18n está aplicado
5. Checar se estados (loading, empty, error, success) estão implementados
6. Reportar o que foi mudado e por quê (qual heurística/lei justifica)
