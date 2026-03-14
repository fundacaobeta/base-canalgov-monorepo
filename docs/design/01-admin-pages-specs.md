# CanalGov — Specs de UX das Páginas Admin

> Especificações de interface para todas as 24+ páginas do painel administrativo.
> Cada spec define: objetivo, anatomia, estados, comportamentos e violações a corrigir.

---

## Convenções deste documento

```
[P] = Prioridade Alta (impacto em uso diário)
[M] = Prioridade Média
[B] = Prioridade Baixa (melhorias estéticas)

⚠ = Violação de heurística/acessibilidade a corrigir
✦ = Feature nova a implementar
```

---

## SEÇÃO 1 — WORKSPACE

---

### 1.1 Configurações Gerais `/admin/general` [P]

**Objetivo:** Configurar identidade e comportamento global do workspace (nome, logo, cores, fuso horário).

**Anatomia da Página:**
```
PageHeader
  Breadcrumb: Admin > Configurações Gerais
  H1: "Configurações Gerais"
  Desc: "Personalize a identidade e o comportamento do seu workspace."

Grid [8fr / 4fr]
  ┌─ Formulário Principal ──────────────────────────────────────┐
  │  Fieldset: Identidade Visual                                │
  │    [Logo do workspace]         ← upload com preview         │
  │    [Nome do workspace *]       ← text input                 │
  │    [Cor principal *]           ← color picker + hex input   │
  │    [Favicon]                   ← upload (32×32)             │
  │                                                             │
  │  Fieldset: Localização                                      │
  │    [Fuso horário *]            ← searchable select          │
  │    [Idioma padrão *]           ← select                     │
  │    [Formato de data]           ← select (dd/mm/aaaa etc.)   │
  │                                                             │
  │  Fieldset: Comportamento de Conversas                       │
  │    [Assinatura padrão de e-mail] ← textarea + rich text     │
  │    [Dias para auto-fechar]     ← number input               │
  │    [Habilitar CSAT]            ← toggle                     │
  │                                                             │
  │  [Salvar Configurações]                                     │
  └─────────────────────────────────────────────────────────────┘
  ┌─ Painel de Ajuda ─────────────────────────────────────────┐
  │  ℹ A cor principal é aplicada em todos os elementos        │
  │    da interface.                                           │
  │  ℹ O fuso horário afeta exibição de datas e cálculos SLA. │
  │  🔗 Ver documentação                                       │
  └────────────────────────────────────────────────────────────┘
```

**Estados:**
- Loading: skeleton no form
- Saving: botão com spinner + disabled
- Success: toast "Configurações salvas com sucesso."
- Error: toast de erro + campos com erro em destaque

**Violações a corrigir:**
- ⚠ Color picker atual pode ter problemas de contraste — validar ratio em tempo real
- ⚠ Preview da cor principal deve refletir imediatamente no sidebar (CSS var update)
- ✦ Adicionar preview ao vivo da identidade (logo + cor + nome na sidebar)

---

### 1.2 Horários de Atendimento `/admin/business-hours` [M]

**Objetivo:** Definir quando a equipe está disponível; usado pelo motor de SLA.

**Anatomia — Lista:**
```
PageHeader
  H1: "Horários de Atendimento"
  Desc: "Configure os períodos de operação para cálculo de SLA."
  [+ Novo Horário]

DataTable
  Colunas: Nome | Fuso Horário | Dias ativos | Padrão | Criado em | Ações
  Ações por linha: [Editar] [···]→[Definir como padrão | Excluir]
  Empty state: "Nenhum horário cadastrado. O SLA usa 24×7 como padrão."
```

**Anatomia — Formulário (modal lg ou página):**
```
H2: "Criar Horário de Atendimento"

[Nome *]                   ← ex: "Atendimento Padrão"
[Fuso horário *]           ← searchable select
[Definir como padrão]      ← toggle

Fieldset: Dias e Horários
  ┌────────────────────────────────────────────────────────┐
  │ ☑ Segunda-feira   [08:00] até [18:00]                 │
  │ ☑ Terça-feira     [08:00] até [18:00]                 │
  │ ☑ Quarta-feira    [08:00] até [18:00]                 │
  │ ☑ Quinta-feira    [08:00] até [18:00]                 │
  │ ☑ Sexta-feira     [08:00] até [18:00]                 │
  │ ☐ Sábado          [---] até [---]                     │
  │ ☐ Domingo         [---] até [---]                     │
  └────────────────────────────────────────────────────────┘

Fieldset: Feriados
  [+ Adicionar feriado]
  Lista de feriados adicionados (data + descrição)

[Criar Horário]  [Cancelar]
```

**Violações a corrigir:**
- ⚠ Time inputs devem ter validação: hora início < hora fim
- ⚠ Ao desmarcar um dia, os campos de hora devem ser desabilitados e esmaecidos (não apenas escondidos)

---

### 1.3 SLA `/admin/sla` [P]

**Objetivo:** Definir acordos de nível de serviço com alertas e escalonamentos.

**Anatomia — Lista:**
```
PageHeader
  H1: "Políticas de SLA"
  [+ Nova Política]

DataTable
  Colunas: Nome | 1ª Resposta | Resolução | Próx. Resposta | Equipes | Ações
  Badge de tempo: ex "4h" pill brand-100
  Empty: "Sem políticas de SLA. Crie a primeira para começar a monitorar prazos."
```

**Anatomia — Formulário:**
```
H2: "Criar Política de SLA"

[Nome da política *]

Fieldset: Metas de Tempo
  [Primeira resposta *]    ← number + select [minutos/horas/dias]
  [Resolução *]            ← number + select
  [Próxima resposta]       ← number + select (opcional)

Fieldset: Alertas
  [Alerta de vencimento]   ← número de horas antes
  [Escalar para]           ← select equipe/supervisor

Fieldset: Aplicação
  [Equipes]                ← multi-select (ou "Todas as equipes")
  [Aplicar apenas em horário comercial] ← toggle

[Criar Política]  [Cancelar]
```

**✦ Features novas:**
- Gráfico sparkline de cumprimento (%) na listagem
- Indicador visual de SLA mais agressivo (< 1h) em vermelho

---

## SEÇÃO 2 — CONVERSAS

---

### 2.1 Tags `/admin/conversations/tags` [P]

**Anatomia — Página única (sem subpáginas):**
```
PageHeader
  H1: "Tags"
  Desc: "Categorize conversas para filtrar e reportar."
  [+ Nova Tag]

DataTable
  Colunas: Cor | Nome | Conversas marcadas | Criada em | Ações
  Ações: [Editar] [Excluir]

  Empty: "Nenhuma tag criada. Tags ajudam a categorizar e filtrar conversas."
```

**Formulário (Dialog sm):**
```
Título: "Criar Tag"

[Nome *]          ← input com preview de badge ao lado
[Cor *]           ← color swatches pré-definidas (12 cores) + input hex
                     Preview: ● [nome da tag]

[Criar Tag]  [Cancelar]
```

**Violações a corrigir:**
- ⚠ Color picker atual sem amostra de como a tag ficará — adicionar preview inline
- ⚠ Sem indicação de quantas conversas usam a tag (impede deleção com contexto)

---

### 2.2 Status de Conversa `/admin/conversations/statuses` [P]

**Anatomia:**
```
PageHeader
  H1: "Status de Conversa"
  Desc: "Defina os estados possíveis de uma conversa no seu fluxo."
  [+ Novo Status]

DataTable
  Colunas: Cor | Nome | Tipo | Conversas | Padrão | Ações
  Tipos: ● Aberto  ● Em andamento  ● Resolvido  ● Personalizado
  Badge "Padrão do sistema" em status não editáveis

  Empty: "Use os status padrão ou crie personalizados para seu fluxo."

⚠ IMPORTANTE: status padrão do sistema (Aberto, Resolvido) não podem ser excluídos
   → mostrar badge "Sistema" + desabilitar botão excluir com tooltip explicativo
```

**Formulário (Dialog md):**
```
Título: "Criar Status"

[Nome *]             ← input
[Cor *]              ← color swatches + hex
[Ícone]              ← picker de ícones (lucide)
[Tipo *]             ← select: Aberto / Em Andamento / Resolvido
                        Tooltip: "O tipo determina como o SLA trata conversas neste status"
[Ordem de exibição]  ← number (drag-to-reorder na lista)

Preview:
  ● [Nome do status]  ← badge com cor escolhida

[Criar Status]  [Cancelar]
```

---

### 2.3 Macros `/admin/conversations/macros` [M]

**Objetivo:** Ações em sequência executadas com um clique na conversa.

**Anatomia — Lista:**
```
PageHeader
  H1: "Macros"
  Desc: "Automatize ações repetitivas com um clique durante o atendimento."
  [+ Nova Macro]

DataTable
  Colunas: Nome | Ações | Visibilidade | Criada por | Ações
  Badge de visibilidade: Privado / Equipe / Global

  Empty: "Sem macros. Crie macros para acelerar o atendimento."
```

**Anatomia — Editor de Macro (página dedicada):**
```
PageHeader com botão voltar
  ← Voltar para Macros
  H1: "Criar Macro"

Grid [7fr / 5fr]
  Esquerda:
    [Nome da macro *]
    [Visibilidade *]   ← Só eu / Minha equipe / Todos

    Fieldset: Sequência de Ações
      ┌────────────────────────────────────────────────────────┐
      │  1. [Tipo de ação ▼]  [Parâmetro]              [✕]   │
      │     ex: Adicionar Tag ▼  [select: tag]                │
      │                                                        │
      │  2. [Tipo de ação ▼]  [Parâmetro]              [✕]   │
      │     ex: Atribuir para ▼  [select: agente]             │
      │                                                        │
      │  [+ Adicionar ação]                                    │
      └────────────────────────────────────────────────────────┘

    [Criar Macro]  [Cancelar]

  Direita:
    Painel: Ações disponíveis
      • Atribuir para agente/equipe
      • Alterar status
      • Adicionar/remover tag
      • Adicionar nota interna
      • Enviar mensagem
      • Definir prioridade
      • Disparar webhook
```

---

### 2.4 Visões Compartilhadas `/admin/conversations/shared-views` [M]

**Anatomia — Lista:**
```
PageHeader
  H1: "Visões Compartilhadas"
  Desc: "Filtros de conversa salvos e visíveis para toda a equipe."
  [+ Nova Visão]

DataTable
  Colunas: Nome | Filtros aplicados | Visibilidade | Equipe | Ações
```

**Formulário (Dialog lg):**
```
Título: "Criar Visão"

[Nome *]
[Visibilidade] ← Todos / Equipe específica
[Equipe]       ← select (aparece se visibilidade = equipe)

Fieldset: Filtros
  Cada filtro: [Campo ▼] [Operador ▼] [Valor]  [✕]
  ex: [Status ▼] [é ▼] [Aberto]
  ex: [Tag ▼] [contém ▼] [urgente]
  [+ Adicionar filtro]
  Operador global: [Todos os filtros (E)] / [Qualquer filtro (OU)]

Preview: "Mostrará conversas abertas com tag urgente"

[Criar Visão]  [Cancelar]
```

---

## SEÇÃO 3 — CANAIS / INBOXES

---

### 3.1 Caixas de Entrada `/admin/inboxes` [P]

**Anatomia — Lista:**
```
PageHeader
  H1: "Caixas de Entrada"
  Desc: "Gerencie os canais de comunicação do seu workspace."
  [+ Nova Caixa]

DataTable
  Colunas: Canal | Nome | Endereço | Status | Criada em | Ações
  Canal badge: ✉ E-mail / 🌐 Web / 💬 WhatsApp
  Status: ● Ativa / ○ Inativa

  Ações: [Configurar] [···]→[Ativar/Desativar | Excluir]
```

**Fluxo de Criação — Wizard (3 passos):**
```
Step 1: Tipo de Canal
  ┌─────────────────────────────────────────────────────────────┐
  │  Selecione o tipo de canal                                  │
  │                                                             │
  │  [✉ E-mail]      [🌐 Web Widget]   [💬 WhatsApp]           │
  │  Conecte uma     Chat no seu site  Via API Meta/Z-API       │
  │  caixa de e-mail                                            │
  │                                                             │
  │  [📱 Telegram]   [📟 SMS]          [+ Em breve]            │
  └─────────────────────────────────────────────────────────────┘

Step 2: Configuração (varia por tipo)
  Para E-mail:
    [Nome da caixa *]
    [Tipo de conexão] ← OAuth (Google/Microsoft) | IMAP/SMTP manual
    ↳ OAuth: botão "Conectar com Google" / "Conectar com Microsoft"
    ↳ Manual:
        Fieldset: Recebimento (IMAP)
          [Servidor] [Porta] [SSL]
          [Usuário] [Senha]
          [Pasta] (padrão: INBOX)
        Fieldset: Envio (SMTP)
          [Servidor] [Porta] [TLS/STARTTLS]
          [Usuário] [Senha]
          [Endereço de envio]
        [🧪 Testar Conexão] → retorna badge ✓ ou ✗ com detalhe

Step 3: Configurações adicionais
  [Equipe padrão]    ← select
  [Agente padrão]    ← select
  [Assinatura]       ← textarea
  [Auto-resposta]    ← toggle + textarea se ativo

Indicador de progresso: ① Tipo → ② Configuração → ③ Finalizar
```

**Violações a corrigir:**
- ⚠ Botão "Testar Conexão" precisa de feedback visual claro (loading → ✓/✗)
- ⚠ Mensagens de erro OAuth devem ser em português claro
- ⚠ Wizard sem step indicator visível

---

### 3.2 Domínios `/admin/domains` [M]

**Anatomia:**
```
PageHeader
  H1: "Domínios de E-mail"
  Desc: "Registre os domínios usados para envio e recebimento de e-mails."
  [+ Adicionar Domínio]

Grid de cards (1-3 colunas)
  Card por domínio:
    [Nome] — [dominio.gov.br]
    Badge provedor: AWS SES / Mailgun / Custom
    Badge status: ● Ativo / ○ Inativo / ★ Padrão
    Badge estratégia: Caixa gerenciada / IMAP / Webhook
    [Notas de configuração] (expansível)
    [Editar] [Excluir]
```

**Formulário (Dialog lg):**
```
Título: "Adicionar Domínio"

[Nome para identificação *]   ← ex: "Secretaria de Saúde"
[Domínio (FQDN) *]           ← ex: saude.prefeitura.gov.br
[Provedor *]                  ← select: AWS SES / Mailgun / SendGrid / Personalizado
[Estratégia de entrada *]     ← select: Caixa gerenciada / IMAP / Webhook
[Ativo]                       ← toggle
[Padrão]                      ← toggle (apenas um por vez)
[Notas de configuração DNS]   ← textarea (MX, SPF, DKIM)

[Adicionar Domínio]  [Cancelar]
```

---

## SEÇÃO 4 — EQUIPES & AGENTES

---

### 4.1 Agentes `/admin/teams/agents` [P]

**Anatomia — Lista:**
```
PageHeader
  H1: "Agentes"
  Desc: "Gerencie os atendentes do seu workspace."
  [Importar CSV]  [+ Convidar Agente]

DataTable
  Colunas: Agente | E-mail | Função | Equipes | Status | Criado em | Ações
  Agente: [Avatar] Nome
  Status: ● Ativo / ○ Inativo / ◷ Pendente (aguardando aceite)

  Filtros rápidos: [Todos ▼] [Função ▼] [Equipe ▼] [Status ▼]
  Busca: [🔍 Buscar por nome ou e-mail...]

  Ações por linha: [Editar] [···]→[Redefinir senha | Ativar/Desativar | Excluir]
```

**Formulário de Criação/Edição (Dialog md ou página dedicada):**
```
Título: "Convidar Agente"

[Foto do agente]   ← upload avatar circular

Fieldset: Dados Pessoais
  [Nome completo *]
  [E-mail *]        ← receberá convite
  [Telefone]

Fieldset: Acesso
  [Função *]        ← select: Administrador / Agente / [Funções customizadas]
  [Equipes]         ← multi-select com busca

[Enviar Convite]  [Cancelar]

Para edição — substituir "Enviar Convite" por "Salvar Alterações"
+ mostrar data de cadastro e último acesso como metadata
```

**Violações a corrigir:**
- ⚠ Status "Pendente" não claramente diferenciado de "Ativo" na tabela atual
- ⚠ Sem busca na listagem
- ✦ Bulk actions: selecionar múltiplos → Ativar em massa / Atribuir equipe / Excluir

---

### 4.2 Equipes `/admin/teams/teams` [P]

**Anatomia — Lista:**
```
PageHeader
  H1: "Equipes"
  Desc: "Organize agentes em equipes para roteamento de conversas."
  [+ Nova Equipe]

DataTable
  Colunas: Nome | Agentes | Responsável | Horário | SLA | Ações
  Agentes: [avatar] [avatar] [avatar] +N (hover = lista)
```

**Formulário:**
```
Título: "Criar Equipe"

[Nome da equipe *]
[Descrição]
[Responsável]        ← select agente (1)

Fieldset: Membros
  [Buscar agentes...]
  Lista de agentes selecionados (com badge de função)
  Cada agente: [Avatar] Nome  [Remover ✕]

Fieldset: Configurações de Atendimento
  [Horário de atendimento] ← select (lista do /admin/business-hours)
  [Política de SLA]        ← select
  [Modo de atribuição]     ← radio: Manual / Round-robin / Menor fila
  [Capacidade máxima]      ← number (conversas simultâneas por agente)

[Criar Equipe]  [Cancelar]
```

---

### 4.3 Papéis & Permissões `/admin/teams/roles` [P]

**Anatomia — Lista:**
```
PageHeader
  H1: "Papéis e Permissões"
  Desc: "Controle de acesso baseado em funções (RBAC)."
  [+ Novo Papel]

DataTable
  Colunas: Nome | Permissões (resumo) | Agentes | Sistema | Ações
  Badge "Sistema" em papéis não editáveis (Administrador, Agente padrão)
```

**Editor de Papel (página dedicada — não modal):**
```
PageHeader
  ← Voltar para Papéis
  H1: "Criar Papel"

Grid [6fr / 6fr]
  Esquerda:
    [Nome do papel *]
    [Descrição]

    Fieldset: Permissões
      ┌─ Conversas ────────────────────────────────────────────┐
      │  ☑ Ver conversas atribuídas                            │
      │  ☑ Ver todas as conversas                              │
      │  ☐ Criar conversas                                     │
      │  ☑ Responder conversas                                 │
      │  ☐ Excluir conversas                                   │
      └────────────────────────────────────────────────────────┘
      ┌─ Contatos ─────────────────────────────────────────────┐
      │  ☑ Ver contatos       ☐ Criar      ☐ Editar  ☐ Excluir│
      └────────────────────────────────────────────────────────┘
      ┌─ Relatórios ───────────────────────────────────────────┐
      │  ☐ Ver relatórios     ☐ Exportar                       │
      └────────────────────────────────────────────────────────┘
      ┌─ Administração ────────────────────────────────────────┐
      │  ☐ Gerenciar equipes  ☐ Gerenciar agentes              │
      │  ☐ Gerenciar inboxes  ☐ Configurações do workspace     │
      └────────────────────────────────────────────────────────┘

  Direita:
    Resumo de permissões ativas (contador)
    Agentes com este papel (lista)

[Criar Papel]  [Cancelar]
```

---

### 4.4 Log de Atividades `/admin/teams/activity-log` [M]

```
PageHeader
  H1: "Log de Atividades"
  Desc: "Auditoria de ações realizadas por agentes e administradores."
  [Exportar CSV]

Filtros:
  [🔍 Buscar...]  [Agente ▼]  [Tipo de ação ▼]  [📅 Período]  [Limpar]

DataTable
  Colunas: Data/Hora | Agente | Ação | Recurso | Detalhes | IP
  Linha expandível: mostra diff (antes/depois) para edições

  Empty: "Nenhuma atividade registrada no período selecionado."
```

---

## SEÇÃO 5 — AUTOMAÇÕES

---

### 5.1 Automações `/admin/automations` [P]

**Anatomia — Lista com Tabs:**
```
PageHeader
  H1: "Automações"
  Desc: "Regras que executam ações automaticamente baseadas em eventos."
  [+ Nova Automação]

Tabs: [Nova Conversa] [Conversa Atualizada] [Resposta Recebida] [Todas]

DataTable por tab:
  Colunas: Nome | Condições (resumo) | Ações (resumo) | Ativo | Ord. | Editar
  Toggle "Ativo" direto na linha (com otimistic update)
  Drag handle (⋮⋮) para reordenar prioridade
```

**Editor de Automação (página dedicada — xl):**
```
PageHeader
  ← Voltar para Automações
  H1: "Criar Automação"

[Nome da automação *]
[Evento gatilho *]    ← select: Nova conversa / Atualizada / Resposta / ...
[Ativo]               ← toggle

Grid [6fr / 6fr]
  Esquerda — Condições (SE):
    Operador: [Todas as condições (E) ▼] / [Qualquer condição (OU) ▼]

    Cada condição:
    [Atributo ▼]     [Operador ▼]    [Valor]           [✕]
    ex: Prioridade    é igual a       Alta
    ex: Tag           contém          urgente

    [+ Adicionar Condição]

  Direita — Ações (ENTÃO):
    Cada ação:
    [Tipo de Ação ▼]    [Parâmetro]                    [✕]
    ex: Atribuir equipe  [select: Suporte Nível 1]
    ex: Adicionar tag    [select: escalado]
    ex: Definir status   [select: Em andamento]
    ex: Enviar e-mail    [para: agente | template: ...]

    [+ Adicionar Ação]

[Criar Automação]   [Cancelar]
```

---

### 5.2 Templates `/admin/templates` [M]

**Anatomia — Lista com Tabs:**
```
PageHeader
  H1: "Templates"
  [+ Novo Template]

Tabs:
  [Respostas]  [Layouts de E-mail]  [Notificações de E-mail]

DataTable:
  Colunas: Nome | Preview | Escopo | Equipe | Atualizado | Ações
  Preview: truncated text (150 chars) ou thumbnail para layouts HTML
```

**Editor de Template (página dedicada — wide):**
```
PageHeader
  ← Voltar para Templates
  H1: "Criar Template de Resposta"

Grid [7fr / 5fr]
  Esquerda:
    [Nome *]
    [Escopo]  ← Global / Equipe específica
    [Equipe]  ← select (se escopo = equipe)

    [Conteúdo *]
    ↳ Para "Resposta": editor TipTap rico (bold, italic, listas, links)
    ↳ Para "Layout HTML": editor de código com syntax highlight + preview
    ↳ Para "Notificação": campos pré-definidos (subject, body)

    [Criar Template]  [Cancelar]

  Direita:
    Preview ao vivo

    Variáveis disponíveis:
    ┌──────────────────────────────────┐
    │ {{contact.name}}                 │
    │ {{conversation.id}}              │
    │ {{agent.name}}                   │
    │ {{inbox.name}}                   │
    │  ... (clique para inserir)       │
    └──────────────────────────────────┘
```

---

## SEÇÃO 6 — SEGURANÇA

---

### 6.1 SSO / OIDC `/admin/sso` [M]

**Anatomia:**
```
PageHeader
  H1: "Single Sign-On (SSO)"
  Desc: "Configure provedores OIDC para login centralizado."
  [+ Adicionar Provedor]

DataTable
  Colunas: Nome | Discovery URL | Status | Última sincronização | Ações
  Ações: [Editar] [Testar] [Excluir]
```

**Formulário (Dialog lg):**
```
Título: "Configurar Provedor OIDC"

[Nome do provedor *]        ← ex: "Gov.br", "Azure AD Prefeitura"
[Discovery URL *]           ← ex: https://sso.gov.br/.well-known/openid-configuration
[Client ID *]
[Client Secret *]           ← input type="password" + toggle mostrar
[Escopos]                   ← tags input, padrão: openid email profile
[Habilitado]                ← toggle

Fieldset: Mapeamento de Atributos
  [Campo e-mail]  ← input, padrão: email
  [Campo nome]    ← input, padrão: name

[🧪 Testar Conexão]

[Salvar Provedor]  [Cancelar]
```

---

## SEÇÃO 7 — NOTIFICAÇÕES

---

### 7.1 Shell padrão para canais de notificação

Todos os canais de notificação seguem o mesmo shell:

```
PageHeader
  H1: "Notificações via [Canal]"
  Desc: "[O que este canal faz]"

Grid [8fr / 4fr]
  Esquerda: Formulário de configuração do canal
  Direita:
    Painel de Ajuda

    Eventos cobertos:
    ✓ Nova mensagem recebida
    ✓ Conversa atribuída
    ✓ SLA em risco
    ✓ Menção em nota interna
    ✓ (etc.)
```

### 7.2 E-mail `/admin/notification` [P]

```
Fieldset: Servidor SMTP
  [Habilitar notificações por e-mail]  ← toggle principal
  [Host *]          ← ex: smtp.gov.br
  [Porta *]         ← 465 / 587 / 25
  [Segurança]       ← select: SSL/TLS / STARTTLS / Nenhuma
  [Usuário *]
  [Senha *]
  [Endereço de envio *]  ← ex: noreply@prefeitura.gov.br
  [Nome do remetente]    ← ex: "CanalGov Prefeitura"

[🧪 Enviar E-mail de Teste]   → input email de destino + enviar
[Salvar Configuração]
```

### 7.3 Comunicações Oficiais `/admin/notification/official-communications` [P]

```
Fieldset: Configuração Geral
  [Habilitar módulo]   ← toggle
  [Caixa de destino *] ← select inbox
  [Prioridade padrão]  ← select: Baixa / Normal / Alta / Urgente
  [Status inicial]     ← select
  [Prefixo do assunto] ← ex: "[OFÍCIO]"
  [Meta de prazo]      ← number + select h/dias

Fieldset: Tipos de Documento
  ☑ Ofício          ☑ Carta         ☑ Notificação
  ☑ Intimação       ☐ Memorando     ☐ Circular

Fieldset: Regras de Roteamento
  Cada regra: [Tipo de doc ▼] → [Equipe ▼]  [✕]
  [+ Adicionar regra]

[Salvar Configuração]
```

---

## SEÇÃO 8 — INTEGRAÇÕES

---

### 8.1 Webhooks `/admin/webhooks` [M]

**Anatomia:**
```
PageHeader
  H1: "Webhooks"
  Desc: "Receba notificações em tempo real quando eventos ocorrerem."
  [+ Novo Webhook]

DataTable
  Colunas: Nome | URL | Eventos | Status | Último disparo | Ações
  Ações: [Editar] [Ver logs] [Excluir]
```

**Formulário:**
```
[Nome *]
[URL de destino *]   ← https://...
[Ativo]              ← toggle

Fieldset: Eventos
  ☑ Conversa criada        ☑ Conversa atualizada
  ☑ Mensagem recebida      ☑ Mensagem enviada
  ☑ Conversa atribuída     ☑ Status alterado
  ☐ CSAT respondido        ☐ Agente criado

Fieldset: Segurança
  [Segredo HMAC]   ← auto-gerado + toggle mostrar + botão copiar
  [Headers customizados]  ← key-value editor

[Salvar Webhook]  [Cancelar]
```

**Página de Logs (subpágina):**
```
DataTable de entregas:
  Colunas: Data | Evento | Status HTTP | Tempo de resposta | Ações
  Status: ✓ 200 / ✗ 500 / ⟳ Tentando novamente
  Ações: [Ver payload] [Reenviar]
```

### 8.2 Ações de Integração `/admin/integrations/actions` [M]

```
PageHeader
  H1: "Ações de Integração"
  Desc: "Defina chamadas HTTP acionadas por eventos nas conversas."
  [+ Nova Ação]

Agrupado por integração:
  ▼ WhatsApp
    Card por ação: Nome | Método | URL | Gatilhos | [Editar] [Excluir]
  ▼ CRM Externo
    ...
```

---

## PADRÕES TRANSVERSAIS

### Estado Vazio (empty state)

Toda tabela/lista DEVE ter um empty state que inclui:
1. **Ícone** relevante (lucide, 48px, neutral-300)
2. **Título** descritivo ("Nenhuma equipe cadastrada ainda")
3. **Subtexto** explicando o valor ("Equipes organizam agentes e roteiam conversas automaticamente")
4. **CTA** quando aplicável ("[+ Criar primeira equipe]")

### Confirmação de Exclusão

Dialog sm padrão:
```
Título: "Excluir [Nome do recurso]?"
Body:   "Esta ação não pode ser desfeita. [Consequências específicas, ex:
         '143 conversas perderão esta tag.']"
Footer: [Excluir]        ← variant="danger"
        [Cancelar]       ← variant="secondary"
```

### Feedback de Ações

| Ação | Feedback |
|---|---|
| Criar | Toast: "[Recurso] criado com sucesso." |
| Editar | Toast: "[Recurso] atualizado com sucesso." |
| Excluir | Toast: "[Recurso] excluído." (sem undo — confirmar antes) |
| Erro API | Toast danger: "Não foi possível [ação]. [Detalhe]. Tente novamente." |
| Toggle | Mudança instantânea (optimistic) + toast silencioso |

### Loading States

| Contexto | Loading pattern |
|---|---|
| Tabela inicial | Skeleton rows (5 linhas) |
| Busca/filtro | Spinner no campo + skeleton |
| Botão de submit | Spinner no botão + disabled |
| Página inteira | Skeleton do layout completo |
| Operação em background | Toast "Processando..." → atualiza para resultado |
