# Merge Review

Objetivo: registrar a revisão técnica do trabalho de merge em andamento nas frentes de `live-chat`, `widget`, `help-center` e `AI`, destacando blockers reais antes de seguir com integração adicional.

## Achados

### 1. Build quebrado por dependências ausentes no pipeline de AI

Severidade: alta

Arquivos:
- [internal/ai/conversation_completions.go](/Users/jobs/Dev/a-publico/base-canalgov-monorepo/internal/ai/conversation_completions.go#L257)
- [internal/ai/conversation_completions.go](/Users/jobs/Dev/a-publico/base-canalgov-monorepo/internal/ai/conversation_completions.go#L558)
- [internal/ai/conversation_completions.go](/Users/jobs/Dev/a-publico/base-canalgov-monorepo/internal/ai/conversation_completions.go#L640)

Problema:
- o código novo chama `stringutil.CleanJSONResponse` e `stringutil.MarkdownToHTML`;
- essas funções não existem na base atual;
- `go test ./...` falha por isso.

Impacto:
- o pacote `internal/ai` não compila;
- qualquer tentativa de integrar completions assíncronas está bloqueada no estado atual.

Direção recomendada:
- não seguir adicionando chamadas ao serviço de completions enquanto as funções utilitárias e testes não existirem na base;
- primeiro fechar o contrato mínimo de `stringutil`.

### 2. Continuidade de conversa entrou incompleta e não compila

Severidade: alta

Arquivo:
- [internal/conversation/continuity.go](/Users/jobs/Dev/a-publico/base-canalgov-monorepo/internal/conversation/continuity.go#L20)

Problema:
- o arquivo referencia `m.continuityConfig`, queries como `GetOfflineLiveChatConversations`, `GetUnreadMessages`, `DeleteMessage`, `UpdateContinuityEmailTracking`, além de métodos ausentes como `fetchMessageAttachments`;
- esses campos e queries não existem no `Manager` atual;
- o pacote `internal/conversation` quebra no build.

Impacto:
- a feature está parcialmente copiada, mas sem a infraestrutura necessária;
- no estado atual, isso não é código pronto para merge, é branch fragmentada.

Direção recomendada:
- isolar `continuity.go` atrás de um branch de feature próprio até completar schema, queries, config e chamadas de inicialização;
- não misturar continuidade de conversa com merges parciais de AI/widget.

### 3. Código novo de widget/help-center foi adicionado, mas ainda não está exposto por rotas

Severidade: média

Arquivos:
- [cmd/chat.go](/Users/jobs/Dev/a-publico/base-canalgov-monorepo/cmd/chat.go)
- [cmd/widget_ws.go](/Users/jobs/Dev/a-publico/base-canalgov-monorepo/cmd/widget_ws.go)
- [cmd/handlers.go](/Users/jobs/Dev/a-publico/base-canalgov-monorepo/cmd/handlers.go#L286)

Problema:
- existem handlers para widget/live-chat, mas `cmd/handlers.go` ainda só expõe `/ws` do app principal;
- não aparecem rotas públicas do widget nem websocket dedicado nessa tabela;
- isso indica integração parcial ou código morto no estado atual.

Impacto:
- falsa percepção de entrega: existe código, mas ele não está acessível pela aplicação;
- aumenta o custo de manutenção porque parece mergeado sem realmente estar funcional.

Direção recomendada:
- qualquer entrega de widget deve ser validada a partir de `handlers` e testes de integração HTTP, não só pela presença dos arquivos.

### 4. Inicialização do AI não acompanha o novo fluxo de completions

Severidade: média

Arquivos:
- [cmd/init.go](/Users/jobs/Dev/a-publico/base-canalgov-monorepo/cmd/init.go#L921)
- [cmd/main.go](/Users/jobs/Dev/a-publico/base-canalgov-monorepo/cmd/main.go#L189)

Problema:
- `initAI(...)` ainda instancia o manager apenas com DB/i18n/logger/chave;
- não há ligação explícita com `helpcenter`, `conversation store` ou start de workers de completions;
- o código novo de completions existe, mas a inicialização principal não o sustenta.

Impacto:
- parte do código de AI está solta no repositório;
- mesmo corrigindo o build, o fluxo continuará incompleto sem wiring explícito.

Direção recomendada:
- definir o fluxo canônico de inicialização antes de aceitar novos arquivos de completions;
- só integrar quando `main/init` e os contratos de store estiverem fechados.

### 5. O review encontra dívida anterior fora do merge, mas ela continua quebrando a validação global

Severidade: baixa

Arquivos:
- `internal/automation/evaluator_extended_test.go`
- `internal/inbox/channel/whatsapp/whatsapp.go`

Problema:
- há testes quebrados e imports não usados fora do escopo da outra IA;
- isso contamina `go test ./...` e dificulta separar regressão nova de dívida preexistente.

Impacto:
- sinal de build menos confiável;
- custo maior para validar merges grandes.

Direção recomendada:
- manter registro de que esses erros existem fora do escopo do merge;
- ao validar integração da outra IA, distinguir sempre `novos blockers` de `dívida já existente`.

## Verificação executada

Comando tentado:

```bash
go test ./...
```

Resultado:
- falhou por restrição de cache do sandbox e, além disso, revelou erros reais de compilação nos pacotes novos de `AI` e `conversation continuity`.

## Conclusão

O trabalho da outra IA ainda não está em estado de merge seguro. O problema principal não é qualidade superficial de código, e sim integração incompleta:

- código novo foi adicionado sem completar dependências internas;
- parte do fluxo não está exposta por rotas;
- a inicialização principal não acompanha os novos módulos;
- o build global já quebra antes de qualquer validação funcional.

Critério mínimo antes de seguir:

1. restaurar `go test`/`go build` do conjunto novo;
2. fechar wiring de `handlers` e `main/init`;
3. separar o que é live-chat/help-center/AI experimental do que realmente está pronto para produção.
