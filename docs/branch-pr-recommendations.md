# Branch Review and PR Recommendations for CanalGov

Objetivo: revisar as branches remotas existentes e separar o que vale virar PR para o CanalGov, o que deve ser cherry-pickado seletivamente e o que nao deve entrar.

Base usada para comparacao: `origin/main`

## Resumo executivo

Resultado da triagem:

- `Abrir PR agora`:
  - `origin/fix/imap-ignore-inbox-emails`
- `Nao abrir PR direto; apenas extrair commits/trechos uteis`:
  - `origin/feat/live-chat-channel`
  - `origin/help-articles-and-ai-responses`
  - `origin/fix/reply-draft-persistence`
- `Nao vale PR`:
  - `origin/copilot/sub-pr-179`
- `Ja estao mergeadas ou sem delta util`:
  - `origin/feat-fs-urls`
  - `origin/feat-mentions`
  - `origin/feat-notifications`
  - `origin/feat/activity/audit-log`
  - `origin/feat/allow-macro-in-new-conversations`
  - `origin/feat/api-user`
  - `origin/feat/dark-mode`
  - `origin/fix-mimetype-detection`
  - `origin/fix/empty-message-id`
  - `origin/fix/max-header-size`
  - `origin/fix/post-put-handlers-return-objects`
  - `origin/mvp`
  - `origin/oauth-inbox-and-encryption`
  - `origin/patch-2`
  - `origin/patch-api-1`
  - `origin/patch-api-2`
  - `origin/refactor-apis`
  - `origin/update-plus-addressing`

## Branch-by-branch

### `origin/fix/imap-ignore-inbox-emails`

- Delta real: 2 commits, 1 arquivo, mudanca pequena e localizada
- Arquivo principal:
  - `internal/inbox/channel/email/imap.go`
- O que faz:
  - ignora e-mails recebidos cujo remetente e o proprio endereco da inbox
  - melhora logging no fluxo IMAP
- Valor para CanalGov:
  - alto
  - reduz loops de ingestao de e-mail
  - reduz ruido operacional e risco de auto-processamento
- Risco de integracao:
  - baixo
- Recomendacao:
  - abrir PR direto

#### PR recomendado

- Titulo:
  - `fix(email): ignore incoming emails sent from the inbox address`
- Escopo:
  - comparar `From` com o endereco da propria inbox e descartar a mensagem
  - melhorar logs do pipeline IMAP
- Justificativa:
  - ganho operacional claro, baixo risco, pouco acoplamento

#### Descricao pronta do PR

```md
## Summary

This PR hardens the IMAP intake path by ignoring emails sent from the inbox's own address.

## Why

Without this guard, the inbox can process messages originating from itself, which increases the risk of loops, duplicate processing, and noisy operational behavior.

## Changes

- extract the inbox email address from the configured `From` field
- compare incoming sender email with the inbox email
- ignore self-originated messages
- improve log messages in the IMAP processing path

## Risk

Low. The change is isolated to the IMAP email intake flow.

## Validation

- send a message from the inbox address to itself
- confirm the message is ignored
- send a message from a different address
- confirm normal intake still works
```

### `origin/feat/live-chat-channel`

- Delta real: muito grande, dezenas de commits, mudancas backend e frontend
- Areas impactadas:
  - novo canal de live chat / widget
  - endpoints publicos de chat
  - websocket/widget middleware
  - rate limiting especifico para widget
  - reestruturacao forte do frontend (`frontend/apps/main`)
- Valor para CanalGov:
  - misto
  - existem ideas uteis, mas o pacote completo diverge bastante do produto atual
- Partes potencialmente aproveitaveis:
  - guardas de reply / UX de mensagem enviada imediatamente
  - fallback de avatar
  - melhorias de draft preview
  - fixes relacionados a thumb assinada e fluxo de mensagem
- Risco de integracao:
  - muito alto se entrar como PR unico
  - mistura feature de produto upstream com refatoracao estrutural
- Recomendacao:
  - nao abrir PR direto
  - decompor em PRs pequenos apenas se houver demanda de produto

#### Itens que poderiam virar PRs separados

1. `feat(ui): show pending sent messages immediately in conversation thread`
2. `fix(media): ensure signed thumbnails resolve correctly`
3. `feat(conversation): add reply guards to prevent invalid sends`
4. `feat(ui): improve avatar fallback behavior`

### `origin/help-articles-and-ai-responses`

- Delta real: muito grande, backend + frontend + help center + AI assistants + widget
- Areas impactadas:
  - help center
  - AI assistants
  - markdown/html conversion
  - rate limit de completions
  - widget/public endpoints
- Valor para CanalGov:
  - parcial
  - ha correcoes pontuais possivelmente uteis no modulo AI
- Partes potencialmente aproveitaveis:
  - `fix(ai): compute email recipients for AI and automated replies`
  - rate limiting para requests de completion
  - correcoes na ordem/seleção de mensagens para contexto do AI
- Risco de integracao:
  - muito alto como PR unico
  - branch traz features novas demais e superficie publica adicional
- Recomendacao:
  - nao abrir PR direto
  - extrair somente fixes de AI se o CanalGov estiver usando esse fluxo

#### PRs seletivos possiveis

1. `fix(ai): compute recipients correctly for AI and automated email replies`
2. `feat(ai): rate limit completion requests`
3. `fix(ai): include the first message when building assistant context`

### `origin/fix/reply-draft-persistence`

- Delta observado:
  - branch contaminada por merge com historico paralelo
  - diff inclui documentacao, i18n, media, config, UI e conversas
- Valor para CanalGov:
  - baixo como PR unico
  - alto risco de carregar ruido e regressao
- Pontos possivelmente uteis dentro da branch:
  - per-user last seen em conversa
  - pequenas correcoes de media/url
  - ajustes de renderizacao HTML com `vue-letter`
- Estado para CanalGov:
  - parte desse conteudo ja aparece na base atual
  - o restante precisa ser reavaliado commit a commit
- Recomendacao:
  - nao abrir PR direto
  - usar apenas como referencia tecnica para cherry-picks cirurgicos

### `origin/copilot/sub-pr-179`

- Delta real: 1 commit
- Conteudo:
  - `Initial plan`
- Valor para CanalGov:
  - nenhum como codigo de producao
- Recomendacao:
  - descartar

## Priorizacao de PRs

### PR 1 - recomendado abrir ja

- Branch fonte:
  - `origin/fix/imap-ignore-inbox-emails`
- Motivo:
  - pequena
  - baixo risco
  - ganho operacional imediato

### PR 2 - so se houver demanda clara de produto

- Fonte:
  - recorte manual de `origin/help-articles-and-ai-responses`
- Tema:
  - fixes de AI isolados

### PR 3 - so se houver roadmap para live chat/widget

- Fonte:
  - recorte manual de `origin/feat/live-chat-channel`
- Tema:
  - UX/fixes pequenos, nunca a branch inteira

## Comandos uteis

```bash
git log --oneline origin/main..origin/fix/imap-ignore-inbox-emails
git diff --stat origin/main...origin/fix/imap-ignore-inbox-emails

git log --oneline origin/main..origin/help-articles-and-ai-responses
git diff --stat origin/main...origin/help-articles-and-ai-responses

git log --oneline origin/main..origin/feat/live-chat-channel
git diff --stat origin/main...origin/feat/live-chat-channel
```

## Recomendacao final

Se a meta for abrir PRs que realmente agreguem ao CanalGov sem puxar divergencia de produto do upstream, a estrategia correta e:

1. abrir apenas o PR de `fix/imap-ignore-inbox-emails`;
2. tratar `help-articles-and-ai-responses` e `feat/live-chat-channel` como minas de cherry-pick, nao como branches mergeaveis;
3. ignorar `fix/reply-draft-persistence` como PR e usar apenas para consulta manual.
