# Branch Migration Roadmap

Objetivo: transformar o backlog de migracao seletiva de branches em um plano executavel para o CanalGov, separando o que ja foi absorvido, o que cabe em PR pequeno e o que depende de decisao de produto ou refactor estrutural.

Base de comparacao usada: `origin/main`

## Matriz de Todas as Branches Revisadas

Classificacao:
- `Mergeado/Sem delta util`: branch sem delta pratico contra `origin/main` ou com conteudo ja absorvido.
- `PR pequeno`: branch nao deve ser mergeada inteira, mas tem commits ou trechos que cabem em PR isolado.
- `Roadmap`: branch com valor parcial, mas dependente de escopo, arquitetura ou produto.
- `Descartar`: sem valor de produto/engenharia para a base principal.

| Branch | Status | Observacao |
| --- | --- | --- |
| `origin/copilot/sub-pr-179` | Descartar | Conteudo de planejamento, sem valor de runtime. |
| `origin/feat-fs-urls` | Mergeado/Sem delta util | Sem delta util restante contra `origin/main`. |
| `origin/feat-mentions` | Mergeado/Sem delta util | Sem delta util restante contra `origin/main`. |
| `origin/feat-notifications` | Mergeado/Sem delta util | Sem delta util restante contra `origin/main`. |
| `origin/feat/activity/audit-log` | Mergeado/Sem delta util | Sem delta util restante contra `origin/main`. |
| `origin/feat/allow-macro-in-new-conversations` | Mergeado/Sem delta util | Sem delta util restante contra `origin/main`. |
| `origin/feat/api-user` | Mergeado/Sem delta util | Sem delta util restante contra `origin/main`. |
| `origin/feat/dark-mode` | Mergeado/Sem delta util | Sem delta util restante contra `origin/main`. |
| `origin/feat/live-chat-channel` | Roadmap | Branch muito grande; extracao apenas por recorte pequeno ou iniciativa de produto. |
| `origin/fix-mimetype-detection` | Mergeado/Sem delta util | Sem delta util restante contra `origin/main`. |
| `origin/fix/empty-message-id` | Mergeado/Sem delta util | Sem delta util restante contra `origin/main`. |
| `origin/fix/imap-ignore-inbox-emails` | Ja migrado | Correcao incorporada em `internal/inbox/channel/email/imap.go`. |
| `origin/fix/max-header-size` | Mergeado/Sem delta util | Sem delta util restante contra `origin/main`. |
| `origin/fix/post-put-handlers-return-objects` | Mergeado/Sem delta util | Sem delta util restante contra `origin/main`. |
| `origin/fix/reply-draft-persistence` | Fechada sem acao restante | O que tinha valor ja esta na base atual; restante e ruido de merge. |
| `origin/help-articles-and-ai-responses` | Roadmap | Fixes de AI aproveitaveis; branch inteira nao cabe no produto atual. |
| `origin/main` | Base | Branch principal. |
| `origin/mvp` | Mergeado/Sem delta util | Sem delta util restante contra `origin/main`. |
| `origin/oauth-inbox-and-encryption` | Mergeado/Sem delta util | Sem delta util restante contra `origin/main`. |
| `origin/patch-2` | Mergeado/Sem delta util | Sem delta util restante contra `origin/main`. |
| `origin/patch-api-1` | Mergeado/Sem delta util | Sem delta util restante contra `origin/main`. |
| `origin/patch-api-2` | Mergeado/Sem delta util | Sem delta util restante contra `origin/main`. |
| `origin/refactor-apis` | Mergeado/Sem delta util | Sem delta util restante contra `origin/main`. |
| `origin/update-plus-addressing` | Mergeado/Sem delta util | Sem delta util restante contra `origin/main`. |

## Ja Migrado

Estes itens ja entraram na base principal por migracao seletiva:

- ignorar e-mails da propria inbox IMAP;
- preservar query string de signed URLs em thumbs;
- preview de rascunho na lista de conversas;
- envio otimista com mensagem pendente e rollback;
- avatar/nome do agente clicaveis para quem pode gerenciar usuarios;
- bloqueio de reply por e-mail sem destinatario `To`.

## Branches Fora de `live-chat` e `help-articles`

Conclusao final da revisao:

- fora `origin/feat/live-chat-channel` e `origin/help-articles-and-ai-responses`, nao restou branch com backlog tecnico relevante;
- `origin/fix/reply-draft-persistence` foi reavaliada e fechada:
  - `conversation_last_seen` e `UpdateUserLastSeen(...)` ja existem na base atual;
  - o `MessageBubble.vue` unificado ja esta presente;
  - o preview HTML com `Letter` em `CommandBox.vue` ja esta presente;
  - o restante do diff e ruido de merge, layout ou ajustes ja absorvidos por outros caminhos;
- as demais branches remotas fora desses dois grupos ja estavam mergeadas, sem delta util ou sem valor de produto.

Implicacao pratica:

- nao ha mais trabalho pendente de migracao seletiva fora dos blocos `live-chat` e `help-articles`;
- qualquer nova rodada fora desses blocos deve começar por bug report ou necessidade concreta, nao por branch archaeology.

## Itens de Roadmap e Mudancas Necessarias

### 1. Refactor mais amplo do fluxo inbound/livechat

Referencia principal:
- `origin/feat/live-chat-channel` / commit `67823902`

Escopo:
- upgrade de visitante para contato em cenarios de reply por e-mail ou autenticacao posterior;
- rate limit especifico do ciclo de criacao de conversa livechat;
- revisao do fluxo inbound para reduzir duplicidade e estados quebrados.

Porque nao cabe em cherry-pick direto:
- o commit mexe em `cmd/chat.go`, `cmd/conversation.go`, `internal/conversation/conversation.go`, `internal/conversation/message.go`, `internal/conversation/queries.sql` e `internal/envelope/envelope.go`;
- o comportamento depende de semantica do widget/livechat e de dados de usuario autenticado vs visitante;
- ha mistura de regra de negocio, query changes e mudanca de entrada HTTP.

Mudancas necessarias antes de implementar:
- definir se CanalGov quer tratar visitante e contato como continuidade da mesma identidade no livechat;
- definir se `external_user_id` sera a chave canônica de reconciliação no widget;
- validar impacto de schema e consultas em conversas abertas do livechat;
- revisar o fluxo atual de `cmd/chat.go` para nao quebrar compatibilidade com JWTs ja emitidos;
- documentar trust boundary do widget: inbox secret, JWT, visitor JWT e merge de identidade.

Subtarefas recomendadas:
- mapear todos os endpoints publicos do widget e seu fluxo de autenticacao em um diagrama curto;
- isolar o caso de uso "visitor -> authenticated user" e escrever teste de integracao para merge;
- definir uma policy de rate limit por `inbox_id`, `remote IP` e `contact identity`;
- revisar `GetContactChatConversations(...)` e queries relacionadas antes de mexer em reuso de conversa;
- separar qualquer mudanca de schema em migration independente, sem acoplar com UX do widget.

Critérios de aceite:
- visitante autenticado depois do inicio da conversa nao gera duplicidade de contato/conversa indevida;
- criacao de conversa livechat nao permite burst abusivo por inbox;
- compatibilidade mantida com widget atual e com JWTs existentes;
- testes cobrindo criacao, reconexao e upgrade de identidade.

Risco:
- alto

Prioridade:
- somente quando livechat autenticado e continuidade de identidade entrarem no roadmap ativo.

### 2. Revisao do rate limit atual do widget

Referencias:
- `origin/feat/live-chat-channel` / commits `1962abdc` e `67823902`
- implementacao atual em [cmd/widget_middleware.go](/Users/jobs/Dev/a-publico/base-canalgov-monorepo/cmd/widget_middleware.go)

Estado atual observado:
- ja existe `rateLimitWidget(...)` aplicando regra nomeada `widget`;
- falta documentacao clara de cobertura endpoint a endpoint;
- ainda nao esta explicitado se ha limite distinto para `init`, `message send`, `last_seen`, `settings` e websocket handshake.

Mudancas necessarias:
- inventariar quais handlers publicos do widget estao atras do limitador atual;
- separar endpoints de leitura barata de endpoints que criam estado;
- definir chaves de rate limit adequadas:
  - por IP para trafego anonimo;
  - por `inbox_id`;
  - por identidade resolvida quando houver JWT valido;
- adicionar observabilidade minima para rate limit hit por inbox.

Subtarefas recomendadas:
- revisar `cmd/handlers.go`, `cmd/chat.go` e `cmd/widget_ws.go` para listar cobertura real do middleware;
- adicionar teste de integracao para burst em `conversations/init`;
- adicionar teste de burst em envio de mensagem/widget API;
- decidir se websocket handshake entra no mesmo bucket ou em bucket separado;
- adicionar counters/logs estruturados para `widget rate limit exceeded`.

Critérios de aceite:
- `init` e `send message` protegidos contra burst realista;
- endpoints de config/public bootstrap nao sofrem throttling desnecessario;
- limite observavel por logs/metricas;
- testes automatizados cobrindo estouro de limite.

Risco:
- medio

Prioridade:
- alta apenas se houver abuso operacional do widget; caso contrario, tratar como hardening.

### 3. Continuidade por e-mail, batch de unread e recursos avancados de widget

Referencias:
- `origin/feat/live-chat-channel` / commits `16ca6b6d`, `54e61442`, `d7067bce`, `9c43b885`

Escopo:
- e-mail fallback/continuity inbox;
- envio batelado de unread messages;
- ocultacao/ajuste visual de mensagens de continuidade;
- refinamentos do ciclo de desconexao do widget.

Porque fica em roadmap:
- mistura produto, notificacao, fluxo operacional e expectativa de atendimento;
- depende de decisao do CanalGov sobre uso de livechat como canal primario ou complementar;
- exige alinhamento com notificacoes, SLA, templates e politicas de contato.

Mudancas necessarias:
- definir politicas de continuidade:
  - quando enviar e-mail;
  - para quem enviar;
  - em que janela de inatividade;
  - se mensagens devem ser bateladas;
- alinhar com regras de privacidade e registro formal do CanalGov;
- definir templates e localizacao i18n dessas notificacoes;
- revisar impacto em `internal/notification/providers/email` e no pipeline de conversation updates.

Subtarefas recomendadas:
- escrever ADR curta de "conversation continuity";
- validar necessidade com produto/operação antes de mexer em código;
- levantar requisitos de template, opt-out e auditoria;
- definir se isso pertence ao inbox de livechat ou ao modelo geral de conversa.

Critérios de aceite:
- regra de negocio aprovada por produto;
- templates e traducoes definidos;
- sem duplicidade de notificacao nem ruido operacional;
- metricas minimas de disparo/entrega.

Risco:
- alto

Prioridade:
- baixa enquanto o livechat nao for iniciativa central do produto.

### 4. Fixes adicionais de AI fora do escopo ativo

Referencia principal:
- `origin/help-articles-and-ai-responses`

Observacao:
- esta branch mistura help center, widget, AI e novas features de produto;
- nao deve ser reaberta inteira;
- apenas fixes pontuais de AI devem ser extraidos quando o modulo estiver de fato no escopo ativo.

Mudancas necessarias para revisitar:
- decidir se o modulo de assistentes/help center sera evoluido agora ou fica congelado;
- separar fixes de AI em commits pequenos e autocontenidos;
- evitar carregar dependencias de help center/widget por efeito colateral.

Subtarefas recomendadas:
- extrair primeiro os fixes de menor risco:
  - `6f62a777` destinatarios automaticos para AI/automações/CSAT;
  - `d0df6f93` rate limit de completions;
  - `af137327` correcoes de enqueue/primeira mensagem/contexto;
- revisar `8bf0255b` se a ordenacao do contexto de completion ainda divergir da base atual;
- deixar `30902310` e afins fora do escopo ate haver decisao explicita sobre markdown/html/help center.

Critérios de aceite:
- cada fix de AI entra em PR separado;
- sem novas features de help center/widget acopladas;
- testes cobrindo enqueue, reply automatico e limite de completions.

Risco:
- medio

Prioridade:
- media se AI continuar ativa no CanalGov; baixa se o modulo permanecer secundario.

## Ordem Recomendada de Execucao Futura

1. `SendAutoReply` para AI/automações/CSAT.
2. Rate limit de completions de AI.
3. Validacao de trusted domain do widget via `Referrer`.
4. Revisao de cobertura do rate limit do widget.
5. CSAT visivel no bubble do agente.
6. `external_user_id` no sidebar do contato.
7. Refactor amplo de identidade/continuity do livechat apenas com roadmap aprovado.

## O Que Nao Fazer

- nao abrir PR da branch `origin/feat/live-chat-channel` inteira;
- nao abrir PR da branch `origin/help-articles-and-ai-responses` inteira;
- nao misturar migration de schema, UX de widget e autenticacao em um mesmo PR;
- nao adicionar i18n novo parcial: qualquer texto novo precisa entrar em todos os idiomas suportados ou reaproveitar chave existente.

## Observacao Final

Se a meta for reduzir risco e manter velocidade, a estrategia correta continua sendo:
- cherry-pick apenas de recortes pequenos;
- documentacao clara do que depende de produto;
- testes automatizados a cada migracao;
- nenhuma adocao de branch inteira como fonte de merge.
