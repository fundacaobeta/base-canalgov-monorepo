# Auditoria AppSec - CanalGov

Documento de trabalho da auditoria ofensiva e defensiva do monorepo CanalGov.

## Escopo

- Backend Go 1.25 em `fasthttp`
- Frontend Vue 3 + Vite
- PostgreSQL 17
- Redis 7
- Docker Compose
- Uploads locais em `uploads/`
- Sessao Redis, OIDC/OAuth2, websocket, webhooks, parsing de e-mail e rich content

## Fase 1 - Mapeamento e Modelagem

### Superficie de ataque observada

- Backend HTTP em `9000` e frontend Vite em `8000`.
- Postgres e Redis expostos apenas em loopback no host local.
- Compose principal e compose de dev executam `install` e `upgrade`; o compose de dev ativa `CANALGOV_AUTO_UPGRADE=true`.
- Rotas publicas relevantes:
  - `POST /api/v1/auth/login`
  - `GET /api/v1/config`
  - `GET /api/v1/oidc/{id}/login`
  - `GET /api/v1/oidc/{id}/finish`
  - `GET /uploads/{uuid}` por auth ou signed URL
  - `GET/POST /api/v1/webhooks/whatsapp`
  - `POST /api/v1/webhooks/telegram`
  - `GET/POST /csat/{uuid}`
  - `GET /health`
- Rotas autenticadas por sessao ou API key:
  - `media`, `drafts`, `agents/me`, `tags`, `statuses`, `macros`, `views/shared`, `inboxes`, `ws`
- Rotas protegidas por `perm(...)`:
  - OIDC, conversas, mensagens, contatos, equipes, automacoes, roles, webhooks, relatorios, templates, inboxes, admin

### Ativos de alto valor

- Sessao Redis e cookies de autenticacao
- API keys de agentes
- Tokens de reset de senha
- Segredos OIDC, webhook e OAuth de inbox
- Conversas, anexos, avatares e uploads
- Dados pessoais de contatos e agentes
- Regras de autorizacao, roles, equipes e automacoes
- Logs de atividade

### Trust boundaries

- Browser -> backend: cookie de sessao, CSRF, API key e websocket
- Backend -> Redis: sessao e state OAuth
- Backend -> Postgres: usuarios, mensagens, OIDC, inboxes, webhooks e configuracoes
- Backend -> filesystem/S3: anexos, thumbs e avatares
- Backend -> provedores externos: OIDC discovery, Google/Microsoft OAuth, IMAP/SMTP, webhooks outbound, Telegram e WhatsApp
- E-mail arbitrario -> parser MIME -> banco -> renderizacao HTML no frontend

### Fluxos de maior risco

- Login local com sessao Redis e CSRF double-submit
- Login OIDC com state em sessao e binding final por e-mail
- OAuth de inbox com state em Redis e persistencia de tokens em config
- Upload de arquivo com thumbnail, content-type detectado e serving inline
- Parsing de e-mail com HTML bruto e anexos persistidos
- Websocket autenticado sem validacao de origem

### Abuse cases prioritarios

- Broken Access Control e BOLA em recursos que dependem apenas de `perm(...)`
- Exposicao excessiva de dados em endpoints `auth(...)`
- Stored XSS via HTML de e-mail, TipTap, links arbitrarios e inline images
- SSRF por OIDC discovery, teste de provider e webhooks outbound
- Websocket hijack cross-origin
- Upload de SVG/polyglot/imagem maliciosa e serving inline posterior
- Brute force em login e reset de senha
- Setup/upgrade executado em ambiente inadequado
- Spoofing de webhooks publicos

## Fase 2 - SAST

Status atual: em andamento.

Hotspots de revisao prioritaria:

- `cmd/middlewares.go`
- `internal/auth/auth.go`
- `cmd/auth.go`
- `cmd/oidc.go`
- `cmd/oauth.go`
- `cmd/media.go`
- `internal/media/stores/localfs/fs.go`
- `cmd/websocket.go`
- `internal/inbox/channel/email/imap.go`
- `cmd/macro.go`
- `internal/macro/queries.sql`
- `frontend/src/components/editor/TextEditor.vue`
- `frontend/src/features/conversation/message/MessageBubble.vue`

### Primeiros achados confirmados

- `GET /api/v1/macros` retorna todas as macros para qualquer usuario autenticado; a filtragem por `visibility`, `team_id` e `user_id` ocorre apenas no frontend.
- `reloadAuth()` recarrega a configuracao OIDC sem preservar `SecureCookies`, o que afeta o comportamento do `csrf_token` apos criacao/edicao de provider OIDC.
- O login OIDC vincula contas locais apenas por e-mail e ignora `email_verified`; tambem nao usa `nonce` no fluxo.
- A criacao e atualizacao de providers OIDC fazem discovery server-side em URL controlada pelo usuario com permissao `oidc:manage`, sem protecao SSRF dedicada.
- O websocket aceita qualquer origem (`CheckOrigin == true`).
- Nao ha mecanismo visivel de rate limit, lockout ou throttling nos fluxos de login e reset de senha.

## Achados detalhados da Fase 2

[Upload perigoso por allowlist wildcard, confianca em `Content-Type` do cliente e serving inline]

- Status: Confirmado
- Severidade: Alta
- Categoria: OWASP A05:2021 Security Misconfiguration; OWASP A03:2021 Injection; OWASP API8:2023 Security Misconfiguration; ASVS V12 File and Resources; CWE-434
- Componente afetado: Upload e serving de midia
- Superficie de ataque: `POST /api/v1/media` e `GET /uploads/{uuid}`
- Pre-requisitos: Usuario autenticado com acesso a upload de midia
- Descricao tecnica: O backend aceita extensoes com base em `app.allowed_file_upload_extensions`, cujo valor default em `schema.sql` e `["*"]`. No upload, o backend le `Content-Type` do multipart e a rotina `detectContentType()` confia integralmente no valor enviado pelo cliente sempre que ele nao for um tipo generico. No serving, qualquer arquivo com `Content-Type` iniciado por `image/`, `video/` ou exatamente `application/pdf` e servido com `Content-Disposition: inline`, independentemente do `disposition` originalmente salvo. Isso permite upload de arquivos ativos como SVG com `Content-Type: image/svg+xml` e posterior exposicao inline sob a mesma origem da aplicacao.
- Evidencia: `schema.sql` define `app.allowed_file_upload_extensions` como `["*"]`; `cmd/media.go` valida extensao apenas contra essa lista e passa `fileHeader.Header.Get("Content-Type")` para o pipeline de armazenamento; `internal/media/media.go` retorna o `sourceContentType` diretamente quando ele nao e generico; `cmd/media.go` serve `image/*`, `video/*` e `application/pdf` inline.
- Impacto: Upload irrestrito de tipos perigosos, stored XSS dependente de browser/CSP via SVG/HTML-like payload, phishing interno com PDF inline, exposicao de conteudo ativo hospedado sob a origem da aplicacao e ampliacao de superficie para polyglots e bypass de validacao.
- Exploracao passo a passo:
  1. Autenticar com um usuario comum.
  2. Enviar um arquivo SVG contendo JavaScript ou markup ativo usando multipart e `Content-Type: image/svg+xml`.
  3. Receber o UUID do arquivo retornado pela API.
  4. Acessar `/uploads/{uuid}` com sessao valida ou por signed URL quando aplicavel.
  5. Observar `Content-Type: image/svg+xml` e `Content-Disposition: inline`.
- Requisicao HTTP de teste:
```http
POST /api/v1/media HTTP/1.1
Host: localhost:9000
Cookie: canalgov_session=<sessao>
Content-Type: multipart/form-data; boundary=----cg

------cg
Content-Disposition: form-data; name="inline"

true
------cg
Content-Disposition: form-data; name="files"; filename="poc.svg"
Content-Type: image/svg+xml

<svg xmlns="http://www.w3.org/2000/svg"><script>alert(document.domain)</script></svg>
------cg--
```
- Payloads de teste:
```xml
<svg xmlns="http://www.w3.org/2000/svg"><script>alert(document.domain)</script></svg>
```
```xml
<svg xmlns="http://www.w3.org/2000/svg" onload="fetch('/api/v1/config',{credentials:'include'})"></svg>
```
- Resposta esperada se vulneravel: Upload aceito com `content_type` controlado pelo cliente e `GET /uploads/{uuid}` retornando `Content-Type: image/svg+xml` com `Content-Disposition: inline`.
- Resposta esperada se seguro: Bloqueio do upload para tipos nao permitidos, reclassificacao segura do MIME por magic bytes, e/ou forcing de `Content-Disposition: attachment` para tipos ativos.
- Causa raiz: Politica default-open de extensoes, confianca indevida em `Content-Type` do cliente e decisao de render inline baseada apenas no MIME armazenado.
- Correcao recomendada: Remover wildcard por default; permitir somente lista fechada de tipos seguros; sempre detectar MIME por magic bytes; bloquear SVG/HTML/XML ativos ou servi-los sempre como `attachment`; adicionar `X-Content-Type-Options: nosniff`; revisar signed URLs para uploads sensiveis.
- Estrategia de validacao pos-fix: Repetir upload com SVG, polyglot e PDF malicioso; confirmar rejeicao ou download forçado; validar headers `Content-Type`, `Content-Disposition` e `X-Content-Type-Options`.
- Cobertura: SAST

[Webhooks publicos de WhatsApp e Telegram aceitam spoofing e verificacao incompleta]

- Status: Confirmado
- Severidade: Alta
- Categoria: OWASP A07:2021 Identification and Authentication Failures; OWASP API2:2023 Broken Authentication; OWASP API8:2023 Security Misconfiguration; ASVS V4; CWE-345
- Componente afetado: Webhooks publicos de mensageria
- Superficie de ataque: `GET /api/v1/webhooks/whatsapp`, `POST /api/v1/webhooks/whatsapp`, `POST /api/v1/webhooks/telegram`
- Pre-requisitos: Nenhum
- Descricao tecnica: A verificacao do webhook de WhatsApp aceita qualquer `hub.verify_token` nao vazio quando `hub.mode=subscribe`, e o proprio comentario no handler indica que a comparacao real com um segredo ainda nao existe. Os handlers POST de WhatsApp e Telegram apenas registram o payload e retornam `200 OK`, sem validar assinatura, token, origem ou estrutura. Isso torna o canal publico suscetivel a spoofing, replay e ruido operacional.
- Evidencia: `cmd/webhooks_messaging.go` retorna `200` se `verifyToken != ""`; os handlers POST apenas fazem `app.lo.Debug("received ... webhook", "payload", string(payload))` e devolvem `200 OK`.
- Impacto: Envenenamento de eventos, spoofing de integracao, logs falsos, potencial ativacao indevida de fluxos futuros quando a implementacao for completada sem adicionar autenticacao, e ampliacao de superficie para DAST e fuzzing sem qualquer barreira de autenticidade.
- Exploracao passo a passo:
  1. Chamar o endpoint GET de verificacao com qualquer token nao vazio.
  2. Confirmar que o servidor espelha o `hub.challenge`.
  3. Enviar POST arbitrario em JSON para os endpoints publicos.
  4. Confirmar `200 OK` e ausencia de verificacao de autenticidade.
- Requisicao HTTP de teste:
```http
GET /api/v1/webhooks/whatsapp?hub.mode=subscribe&hub.verify_token=abc&hub.challenge=1337 HTTP/1.1
Host: localhost:9000
```
```http
POST /api/v1/webhooks/telegram HTTP/1.1
Host: localhost:9000
Content-Type: application/json

{"update_id":1,"message":{"text":"spoofed"}}
```
- Payloads de teste:
```json
{"object":"whatsapp_business_account","entry":[{"changes":[{"value":{"messages":[{"text":{"body":"spoof"}}]}}]}]}
```
```json
{"update_id":1,"message":{"text":"spoofed telegram webhook"}}
```
- Resposta esperada se vulneravel: `200 OK` para POST arbitrario e `200 OK` com eco do challenge no GET de verificacao com token qualquer.
- Resposta esperada se seguro: `403` ou `401` sem assinatura/segredo valido, com verificacao estrita do token e do formato assinado do provedor.
- Causa raiz: Endpoints publicos implementados em modo placeholder e sem enforcement criptografico de autenticidade.
- Correcao recomendada: Validar `hub.verify_token` contra segredo configurado; exigir e verificar assinatura HMAC/headers oficiais; rejeitar payloads sem autenticidade; registrar eventos de falha e adicionar replay protection.
- Estrategia de validacao pos-fix: Repetir GET/POST sem segredo e com assinatura invalida; confirmar bloqueio; repetir com segredo correto e assinatura valida; confirmar aceitacao.
- Cobertura: SAST

[Stored XSS provavel via HTML inbound de e-mail renderizado no frontend]

- Status: Provavel
- Severidade: Alta
- Categoria: OWASP A03:2021 Injection; OWASP A07:2021 Identification and Authentication Failures; OWASP API8:2023 Security Misconfiguration; ASVS V5 Validation, Sanitization and Encoding; CWE-79
- Componente afetado: Pipeline de e-mail inbound e renderizacao de mensagens HTML
- Superficie de ataque: Entrada de e-mail IMAP -> mensagens de conversa no frontend
- Pre-requisitos: Atacante capaz de enviar e-mail para uma inbox integrada; vitima autenticada abrindo a conversa no frontend
- Descricao tecnica: O parser IMAP concatena todas as partes `text/html` do e-mail e persiste esse HTML como `Message.Content` sem sanitizacao visivel no backend. No frontend, mensagens HTML sao passadas como `:html="sanitizedContent"` ao componente `Letter`, mas o valor `sanitizedContent` nao sanitiza HTML; ele apenas reescreve `cid:` e caminhos `/uploads/`. Nao ha sanitizacao explicita no codigo local antes da renderizacao. Como o comportamento interno de `vue-letter` nao foi validado nesta fase, o risco fica como provavel e precisa de DAST.
- Evidencia: `internal/inbox/channel/email/imap.go` concatena `htmlParts` e grava `incomingMsg.Message.Content = allHTML.String()`; `frontend/src/features/conversation/message/MessageBubble.vue` envia esse HTML ao componente `Letter`; a computed `sanitizedContent` apenas substitui `cid:` e prefixa `app.root_url`, sem remover tags ou atributos ativos.
- Impacto: Stored XSS contra agentes autenticados, exfiltracao de dados da SPA, acao autenticada via CSRF-like DOM abuse, sequestro de sessao em contextos sem protecoes adicionais e pivot para abuso interno.
- Exploracao passo a passo:
  1. Enviar e-mail HTML para uma inbox integrada contendo payload com evento ativo, URL `javascript:` ou SVG/MathML.
  2. Aguardar ingestao IMAP e persistencia da mensagem.
  3. Abrir a conversa no frontend autenticado.
  4. Observar execucao do payload ou mutacao inesperada do DOM.
- Requisicao HTTP de teste:
```http
POST /api/v1/conversations/<uuid>/messages HTTP/1.1
Host: localhost:9000
Cookie: canalgov_session=<sessao>
Content-Type: application/json

{"message":"<img src=x onerror=alert(1)>","sender_type":"agent","private":false}
```
- Payloads de teste:
```html
<img src=x onerror=alert(document.domain)>
```
```html
<a href="javascript:alert(1)">click</a>
```
```html
<svg xmlns="http://www.w3.org/2000/svg" onload="alert(1)"></svg>
```
- Resposta esperada se vulneravel: Payload persistido e executado/renderizado ativamente ao abrir a mensagem.
- Resposta esperada se seguro: Payload neutralizado por sanitizacao forte, com remocao de `script`, eventos inline, URLs perigosas e markup ativo.
- Causa raiz: HTML inbound e rich content tratados como confiaveis demais, sem sanitizacao forte e centralizada antes da renderizacao.
- Correcao recomendada: Sanitizar HTML no backend ou imediatamente antes da renderizacao com allowlist estrita; remover `script`, atributos `on*`, URLs perigosas, SVG/MathML ativos e CSS perigoso; cobrir tanto mensagens inbound quanto rich text outbound.
- Estrategia de validacao pos-fix: Rodar bateria de payloads em e-mail inbound, editor rich text e notas; confirmar neutralizacao visual e estrutural do HTML.
- Cobertura: SAST

[Exposicao de macros para qualquer usuario autenticado]

- Status: Confirmado
- Severidade: Alta
- Categoria: OWASP A01:2021 Broken Access Control; OWASP API3:2023 Broken Object Property Level Authorization; ASVS V4 Access Control; CWE-200; CWE-284
- Componente afetado: API de macros
- Superficie de ataque: `GET /api/v1/macros`
- Pre-requisitos: Usuario autenticado com qualquer sessao valida
- Descricao tecnica: A rota de listagem de macros usa apenas `auth(...)`, nao `perm(...)`, e o backend retorna todas as macros sem filtrar por `visibility`, `team_id` ou `user_id`. A restricao real esta implementada apenas no frontend.
- Evidencia: `cmd/handlers.go` expõe `GET /api/v1/macros` sob `auth(handleGetMacros)`; `cmd/macro.go` chama `app.macro.GetAll()` e serializa o resultado integral; o frontend filtra client-side por visibilidade.
- Impacto: Vazamento de macros privadas e de equipe, incluindo textos, acoes automatizadas e referencias internas de operacao.
- Exploracao passo a passo:
  1. Autenticar como usuario de baixa permissao.
  2. Chamar `GET /api/v1/macros`.
  3. Comparar os objetos retornados com a filtragem aplicada na UI.
  4. Confirmar acesso a macros fora do escopo do usuario.
- Requisicao HTTP de teste:
```http
GET /api/v1/macros HTTP/1.1
Host: localhost:9000
Cookie: canalgov_session=<sessao>
```
- Payloads de teste:
```text
Sem payload; basta sessao autenticada.
```
- Resposta esperada se vulneravel: `200 OK` com macros de outros usuarios/equipes.
- Resposta esperada se seguro: `200 OK` apenas com macros autorizadas ao usuario ou `403` quando aplicavel.
- Causa raiz: Enforcement de autorizacao deslocado para o frontend.
- Correcao recomendada: Mover a logica de visibilidade para o backend e filtrar por dono, equipe e permissao antes de responder.
- Estrategia de validacao pos-fix: Repetir o endpoint com usuarios de equipes distintas e confirmar escopo correto no backend.
- Cobertura: SAST

[Open redirect via parametro `next` em fluxo OIDC e paginas de autenticacao]

- Status: Confirmado
- Severidade: Media
- Categoria: OWASP A01:2021 Broken Access Control; OWASP A07:2021 Identification and Authentication Failures; ASVS V4; CWE-601
- Componente afetado: Login OIDC e redirecionamentos de pagina de autenticacao
- Superficie de ataque: `GET /api/v1/oidc/{id}/login?next=...`, callback `GET /api/v1/oidc/{id}/finish?...`, e `notAuthPage(...)`
- Pre-requisitos: Usuario autenticado ou vitima induzida a completar login OIDC
- Descricao tecnica: O handler de login OIDC persiste o parametro `next` vindo da query diretamente na sessao e o callback usa esse valor em `RedirectURI` sem validacao de path relativo ou allowlist de host. O middleware `notAuthPage` tambem faz redirect para `next` sem validacao quando o usuario ja esta autenticado. Isso permite open redirect para dominio externo apos login ou ao visitar paginas de auth enquanto autenticado.
- Evidencia: `cmd/auth.go` salva `oidcNextSessKey` como `string(r.RequestCtx.QueryArgs().Peek("next"))`; no callback, `redirectURL` recebe `nextStr` e e usado em `r.RedirectURI(redirectURL, ...)`; `cmd/middlewares.go` em `notAuthPage` le `next` e chama `r.RedirectURI(nextURI, ...)` sem validacao.
- Impacto: Phishing com trusted domain, roubo de fluxo apos login, encadeamento com OIDC/social engineering e desvio de usuarios autenticados para dominios maliciosos.
- Exploracao passo a passo:
  1. Montar URL de login OIDC com `next=https://evil.example/collect`.
  2. Fazer a vitima iniciar o fluxo de SSO por essa URL.
  3. Concluir autenticacao no provider.
  4. Observar redirect final para dominio externo.
- Requisicao HTTP de teste:
```http
GET /api/v1/oidc/1/login?next=https://evil.example/collect HTTP/1.1
Host: localhost:9000
```
- Payloads de teste:
```text
https://evil.example/collect
//evil.example
https:%2f%2fevil.example
```
- Resposta esperada se vulneravel: Redirect final para host externo controlado pelo atacante.
- Resposta esperada se seguro: Redirecionamento apenas para paths relativos internos ou rejeicao do parametro invalido.
- Causa raiz: Uso direto de `next` controlado pelo usuario em redirecionamentos sem normalizacao nem allowlist.
- Correcao recomendada: Aceitar somente paths relativos iniciando com `/`; rejeitar URLs absolutas, `//host`, schemes externos e entradas codificadas equivalentes; centralizar helper seguro de redirect.
- Estrategia de validacao pos-fix: Repetir o fluxo com `https://evil.example`, `//evil.example` e variantes percent-encoded; confirmar fallback para rota interna segura.
- Cobertura: SAST

[Spoofing de IP em trilha de auditoria por confianca irrestrita em headers de proxy]

- Status: Confirmado
- Severidade: Media
- Categoria: OWASP A09:2021 Security Logging and Monitoring Failures; OWASP A05:2021 Security Misconfiguration; ASVS V7 Error Handling and Logging; CWE-290; CWE-117
- Componente afetado: Coleta de IP em login, logout, callback OIDC, alteracoes de usuario e roles
- Superficie de ataque: Requisicoes HTTP com `X-Forwarded-For`, `X-Real-IP`, `CF-Connecting-IP` e headers equivalentes
- Pre-requisitos: Capacidade de enviar requests HTTP diretamente ao backend ou atraves de proxy que repasse headers do cliente
- Descricao tecnica: O projeto usa `github.com/ferluci/fast-realip` para inferir IP do cliente em varios fluxos sensiveis. A biblioteca retorna diretamente o primeiro IP publico encontrado em headers como `X-Forwarded-For`, `X-Real-IP`, `CF-Connecting-IP` e similares, sem qualquer lista de proxies confiaveis. No projeto, esse valor e usado em logs de login, logout, alteracoes de roles e disponibilidade. Isso permite ao cliente falsificar o IP registrado.
- Evidencia: `cmd/login.go`, `cmd/users.go`, `cmd/roles.go` e `cmd/auth.go` usam `realip.FromRequest(r.RequestCtx)`; o codigo da dependencia `fast-realip@v1.0.1` retorna IP de headers de forwarding sem validar trusted proxies.
- Impacto: Corrupcao de trilha de auditoria, dificultacao de investigacao forense, falsos positivos/negativos em monitoramento e bypass imediato de qualquer rate limit ou allowlist futura que reutilize o mesmo helper.
- Exploracao passo a passo:
  1. Enviar request de login ou qualquer fluxo logado com `X-Forwarded-For: 8.8.8.8`.
  2. Concluir a acao.
  3. Inspecionar atividade/log para confirmar IP spoofado.
  4. Repetir com `CF-Connecting-IP` ou `X-Real-IP`.
- Requisicao HTTP de teste:
```http
POST /api/v1/auth/login HTTP/1.1
Host: localhost:9000
Content-Type: application/json
X-Forwarded-For: 8.8.8.8

{"email":"user@example.gov","password":"Secret123!"}
```
- Payloads de teste:
```text
X-Forwarded-For: 8.8.8.8
X-Real-IP: 1.1.1.1
CF-Connecting-IP: 203.0.113.55
```
- Resposta esperada se vulneravel: A acao e registrada com o IP injetado no header, e nao com o peer real.
- Resposta esperada se seguro: O backend ignora headers de forwarding vindos de clientes nao confiaveis e registra apenas o IP do proxy confiavel ou do socket remoto.
- Causa raiz: Resolucao de IP sem modelo de trusted proxies.
- Correcao recomendada: Remover `fast-realip` ou encapsula-lo atras de validacao de proxy confiavel; aceitar headers de forwarding apenas quando `RemoteAddr` pertencer a reverse proxies controlados; logar tambem `remote_addr` bruto.
- Estrategia de validacao pos-fix: Repetir requests com headers spoofados a partir de origem nao confiavel e confirmar que o IP registrado nao muda.
- Cobertura: SAST

[Exposicao de todas as roles e permissoes para qualquer usuario autenticado]

- Status: Confirmado
- Severidade: Media
- Categoria: OWASP A01:2021 Broken Access Control; OWASP API3:2023 Broken Object Property Level Authorization; ASVS V4 Access Control; CWE-200; CWE-284
- Componente afetado: API de roles
- Superficie de ataque: `GET /api/v1/roles`
- Pre-requisitos: Usuario autenticado com qualquer sessao valida
- Descricao tecnica: A listagem de roles esta protegida apenas por `auth(...)`, enquanto as operacoes de leitura individual, criacao, alteracao e exclusao exigem `roles:manage`. O handler retorna todas as roles com o array completo de `permissions`. Isso expõe o modelo de autorizacao interno e o mapeamento de privilegios para qualquer usuario autenticado.
- Evidencia: `cmd/handlers.go` expoe `GET /api/v1/roles` sob `auth(handleGetRoles)`; `cmd/roles.go` chama `app.role.GetAll()`; `internal/role/queries.sql` faz `SELECT id, created_at, updated_at, name, description, permissions FROM roles`; o frontend consome essa rota em telas administrativas e em `AgentForm`, nao como necessidade universal de todos os agentes.
- Impacto: Enumeracao completa do modelo RBAC, facilitacao de escalacao horizontal/vertical guiada por permissoes reais e exposicao desnecessaria de detalhes administrativos a perfis de baixa privilegio.
- Exploracao passo a passo:
  1. Autenticar como usuario comum.
  2. Chamar `GET /api/v1/roles`.
  3. Observar nomes de roles e arrays completos de permissoes.
  4. Usar a informacao para mapear endpoints e combinacoes de privilegio.
- Requisicao HTTP de teste:
```http
GET /api/v1/roles HTTP/1.1
Host: localhost:9000
Cookie: canalgov_session=<sessao>
```
- Payloads de teste:
```text
Sem payload; basta sessao autenticada.
```
- Resposta esperada se vulneravel: `200 OK` com todas as roles e seus arrays de permissoes.
- Resposta esperada se seguro: `403` para usuarios sem `roles:manage`, ou resposta minimizada sem matriz completa de permissoes.
- Causa raiz: Permissao coarse-grained incorreta na rota de listagem.
- Correcao recomendada: Exigir `roles:manage` para listagem completa; se o frontend precisar apenas de nomes selecionaveis, criar endpoint compacto minimizado.
- Estrategia de validacao pos-fix: Repetir `GET /api/v1/roles` com usuario sem `roles:manage` e confirmar `403` ou payload reduzido.
- Cobertura: SAST

[Exposicao de configuracao completa de inboxes para qualquer usuario autenticado]

- Status: Confirmado
- Severidade: Media
- Categoria: OWASP A01:2021 Broken Access Control; OWASP API3:2023 Broken Object Property Level Authorization; ASVS V4 Access Control; CWE-200; CWE-284
- Componente afetado: API de inboxes
- Superficie de ataque: `GET /api/v1/inboxes`
- Pre-requisitos: Usuario autenticado com qualquer sessao valida
- Descricao tecnica: A listagem de inboxes esta protegida apenas por `auth(...)`, enquanto leitura individual e mutacoes exigem `inboxes:manage`. O handler retorna todos os inboxes com `config` quase completo; `ClearPasswords()` apenas mascara senhas, tokens e `client_secret`, mas preserva hosts IMAP/SMTP, usernames, mailbox, `tenant_id`, `client_id`, dominios gerenciados, local parts, flags operacionais e demais metadados de integracao.
- Evidencia: `cmd/handlers.go` expoe `GET /api/v1/inboxes` sob `auth(handleGetInboxes)`; `cmd/inboxes.go` usa `app.inbox.GetAll()` e retorna os objetos apos `ClearPasswords()`; `internal/inbox/models/models.go` mostra que o clear apenas mascara campos sensiveis, preservando a configuracao restante; `frontend/src/stores/inbox.js` usa essa rota globalmente para popular estado da aplicacao, embora a UI administrativa seja quem consome os detalhes completos.
- Impacto: Vazamento de topologia de mensageria/e-mail, enderecos operacionais, usernames de integracao, dominios e detalhes de infraestrutura; facilita phishing interno, ataque a provedores externos, abuse de mail flow e reconhecimento profundo do ambiente.
- Exploracao passo a passo:
  1. Autenticar como usuario comum.
  2. Chamar `GET /api/v1/inboxes`.
  3. Inspecionar o campo `config` de cada inbox.
  4. Confirmar exposicao de hostnames, usernames, mailboxes, dominios e identificadores de OAuth.
- Requisicao HTTP de teste:
```http
GET /api/v1/inboxes HTTP/1.1
Host: localhost:9000
Cookie: canalgov_session=<sessao>
```
- Payloads de teste:
```text
Sem payload; basta sessao autenticada.
```
- Resposta esperada se vulneravel: `200 OK` com todos os inboxes e `config` detalhado, ainda que com segredos mascarados.
- Resposta esperada se seguro: `403` para usuarios sem `inboxes:manage`, ou endpoint compacto retornando apenas `id` e `name` quando necessario para filtros/seletores.
- Causa raiz: Endpoint administrativo completo exposto sob autenticacao generica.
- Correcao recomendada: Exigir `inboxes:manage` para listagem completa; criar endpoint compacto para uso geral contendo apenas identificadores minimos; revisar armazenamento global em Pinia.
- Estrategia de validacao pos-fix: Repetir `GET /api/v1/inboxes` com usuario sem privilegio administrativo e confirmar bloqueio ou resposta reduzida.
- Cobertura: SAST

[Anexos expostos por signed URLs bearer sem binding a usuario ou contexto]

- Status: Confirmado
- Severidade: Media
- Categoria: OWASP A01:2021 Broken Access Control; OWASP API1:2023 Broken Object Level Authorization; ASVS V4 Access Control; CWE-639
- Componente afetado: Serving de uploads e serializacao de anexos em mensagens
- Superficie de ataque: `GET /api/v1/conversations/{uuid}/messages`, `GET /api/v1/conversations/{cuuid}/messages/{uuid}` e `GET /uploads/{uuid}?sig=...&exp=...`
- Pre-requisitos: Acesso inicial a uma mensagem/conversa que contenha anexo, ou vazamento da URL assinada
- Descricao tecnica: O backend injeta URLs assinadas de anexos diretamente nas respostas de mensagens usando `app.media.GetURL(...)`. Para storage local, a URL assinada e calculada apenas sobre `uuid` + `exp` usando uma chave global de aplicacao. Quando `/uploads/{uuid}` e acessado com `sig`/`exp` validos, o middleware marca `auth_method = "signed_url"` e o handler pula completamente as verificacoes de permissao por conversa/modelo. A URL passa a funcionar como bearer token temporario: qualquer portador consegue baixar o anexo ate o vencimento, sem binding ao usuario original, sessao, conversa ou tenant. Mudancas posteriores de permissao nao revogam a URL ja emitida.
- Evidencia: `cmd/messages.go` popula `attachments[j].URL = app.media.GetURL(...)`; `internal/media/stores/localfs/fs.go` assina URLs apenas com `name` e `exp`; `cmd/middlewares.go` aceita signed URL e seta `auth_method = "signed_url"`; `cmd/media.go` verifica `auth_method == "signed_url"` e chama `serveMediaFile(...)` sem `EnforceMediaAccess` nem `EnforceConversationAccess`.
- Impacto: Compartilhamento fora de contexto de anexos privados, bypass de autorizacao por repasse da URL, impossibilidade de revogacao imediata por usuario/conversa e ampliacao de superficie para vazamento acidental em logs, e-mails e clientes terceiros.
- Exploracao passo a passo:
  1. Autenticar em uma conversa com anexo.
  2. Chamar `GET /api/v1/conversations/{uuid}/messages`.
  3. Copiar o campo `attachments[].url`.
  4. Abrir a URL em navegador privado ou enviar a terceiro sem sessao.
  5. Observar que o download continua funcionando ate `exp`.
- Requisicao HTTP de teste:
```http
GET /uploads/550e8400-e29b-41d4-a716-446655440000?sig=<assinatura>&exp=<unix_ts> HTTP/1.1
Host: localhost:9000
```
- Payloads de teste:
```text
URL assinada legitima capturada da resposta de mensagens
```
- Resposta esperada se vulneravel: `200 OK` com o anexo mesmo sem cookie de sessao nem API key.
- Resposta esperada se seguro: Requerimento de autenticacao contextual ou assinatura vinculada ao usuario/sessao/recurso e invalidada quando o contexto deixar de autorizar.
- Causa raiz: Desenho de autorizacao por signed URL como bearer token desacoplado do contexto de negocio.
- Correcao recomendada: Evitar emitir signed URLs diretamente em APIs internas autenticadas; preferir endpoint autenticado que revalide acesso por conversa/modelo; se signed URL for necessaria, reduzir TTL agressivamente, vincular assinatura a claims de usuario/contexto e permitir revogacao por versao/nonce do recurso.
- Estrategia de validacao pos-fix: Capturar uma URL de anexo, testar acesso sem sessao, depois alterar/remover permissao do usuario e confirmar que a URL antiga nao funciona mais.
- Cobertura: SAST

## Fase 3 - DAST

Objetivo desta fase: validar dinamicamente os achados de maior risco e transformar as hipoteses de SAST em evidencias de exploracao controlada.

### Ordem de execucao recomendada

1. Broken Access Control e excess data exposure
2. Open redirect e fluxos de autenticacao
3. Signed URLs e uploads
4. Webhooks publicos
5. Websocket
6. OIDC SSRF e hardening de provider
7. Rate limit, brute force e enumeração
8. XSS inbound/rich content

### Preparacao minima do ambiente

- Subir stack local com backend em `http://127.0.0.1:9000` e frontend em `http://127.0.0.1:8000`
- Criar ao menos:
  - 1 usuario admin
  - 1 usuario agente comum
  - 2 equipes distintas
  - 1 conversa com anexo
  - 1 macro privada e 1 macro de equipe
  - 1 inbox email com config preenchida
- Capturar uma sessao admin e uma sessao low-priv
- Exportar cookies para uso em `curl`

### Casos DAST priorizados

[DAST - BAC - Macros expostas a qualquer autenticado]

- Status: Confirmado
- Severidade: Alta
- Categoria: OWASP A01:2021 Broken Access Control; OWASP API3:2023 Broken Object Property Level Authorization; CWE-200; CWE-284
- Componente afetado: API de macros
- Superficie de ataque: `GET /api/v1/macros`
- Pre-requisitos: Sessao de usuario de baixa permissao
- Descricao tecnica: Validar que um usuario sem `macros:manage` recebe macros de outros escopos.
- Evidencia: Rota `auth(...)` e `GetAll()` sem filtro no backend.
- Impacto: Vazamento de dados operacionais e logica interna.
- Exploracao passo a passo:
  1. Logar como usuario low-priv.
  2. Chamar `GET /api/v1/macros`.
  3. Comparar `team_id`, `user_id` e `visibility` com o perfil do usuario.
- Requisicao HTTP de teste:
```http
GET /api/v1/macros HTTP/1.1
Host: 127.0.0.1:9000
Cookie: canalgov_session=<low_priv_session>
```
- Payloads de teste:
```bash
curl -i -s \
  -H 'Cookie: canalgov_session=<low_priv_session>' \
  http://127.0.0.1:9000/api/v1/macros
```
- Resposta esperada se vulneravel: `200` com macros privadas/de outras equipes.
- Resposta esperada se seguro: Apenas macros autorizadas ao usuario.
- Causa raiz: Filtragem de visibilidade feita no frontend.
- Correcao recomendada: Enforcement no backend.
- Estrategia de validacao pos-fix: Repetir com multiplos usuarios e equipes.
- Cobertura: DAST

[DAST - BAC - Roles completas para qualquer autenticado]

- Status: Confirmado
- Severidade: Media
- Categoria: OWASP A01:2021 Broken Access Control; OWASP API3:2023 Broken Object Property Level Authorization; CWE-200; CWE-284
- Componente afetado: API de roles
- Superficie de ataque: `GET /api/v1/roles`
- Pre-requisitos: Sessao de usuario sem `roles:manage`
- Descricao tecnica: Validar exposicao da matriz de permissao para perfis nao administrativos.
- Evidencia: Listagem sob `auth(...)`.
- Impacto: Enumeracao de RBAC.
- Exploracao passo a passo:
  1. Logar como usuario low-priv.
  2. Chamar `GET /api/v1/roles`.
  3. Confirmar retorno de `permissions`.
- Requisicao HTTP de teste:
```http
GET /api/v1/roles HTTP/1.1
Host: 127.0.0.1:9000
Cookie: canalgov_session=<low_priv_session>
```
- Payloads de teste:
```bash
curl -s \
  -H 'Cookie: canalgov_session=<low_priv_session>' \
  http://127.0.0.1:9000/api/v1/roles | jq .
```
- Resposta esperada se vulneravel: `200` com arrays completos de permissoes.
- Resposta esperada se seguro: `403` ou resposta minimizada.
- Causa raiz: Permissao incorreta na listagem.
- Correcao recomendada: Exigir `roles:manage`.
- Estrategia de validacao pos-fix: Repetir com perfil low-priv.
- Cobertura: DAST

[DAST - BAC - Inboxes completas para qualquer autenticado]

- Status: Confirmado
- Severidade: Media
- Categoria: OWASP A01:2021 Broken Access Control; OWASP API3:2023 Broken Object Property Level Authorization; CWE-200; CWE-284
- Componente afetado: API de inboxes
- Superficie de ataque: `GET /api/v1/inboxes`
- Pre-requisitos: Sessao de usuario sem `inboxes:manage`
- Descricao tecnica: Validar exposicao de `config` detalhado de inboxes para usuario comum.
- Evidencia: Listagem sob `auth(...)`, com apenas mascaramento parcial.
- Impacto: Vazamento de infraestrutura de mail/mensageria.
- Exploracao passo a passo:
  1. Logar como low-priv.
  2. Chamar `GET /api/v1/inboxes`.
  3. Inspecionar `config`.
- Requisicao HTTP de teste:
```http
GET /api/v1/inboxes HTTP/1.1
Host: 127.0.0.1:9000
Cookie: canalgov_session=<low_priv_session>
```
- Payloads de teste:
```bash
curl -s \
  -H 'Cookie: canalgov_session=<low_priv_session>' \
  http://127.0.0.1:9000/api/v1/inboxes | jq .
```
- Resposta esperada se vulneravel: `200` com hostnames, usernames, mailbox, domains, tenant/client ids.
- Resposta esperada se seguro: `403` ou resposta compacta.
- Causa raiz: Endpoint administrativo completo exposto sob autenticacao generica.
- Correcao recomendada: Segregar endpoint compacto e endpoint administrativo.
- Estrategia de validacao pos-fix: Repetir com usuario sem permissao administrativa.
- Cobertura: DAST

[DAST - Open Redirect em `next`]

- Status: Confirmado
- Severidade: Media
- Categoria: OWASP A07:2021 Identification and Authentication Failures; CWE-601
- Componente afetado: Fluxo OIDC e paginas not-auth
- Superficie de ataque: `/api/v1/oidc/{id}/login?next=...`
- Pre-requisitos: Provider OIDC habilitado
- Descricao tecnica: Validar redirect externo apos login usando `next`.
- Evidencia: `next` persistido e reutilizado sem validacao.
- Impacto: Phishing e desvio de fluxo autenticado.
- Exploracao passo a passo:
  1. Acessar URL de login com `next=https://evil.example/collect`.
  2. Completar o fluxo OIDC.
  3. Observar redirect final.
- Requisicao HTTP de teste:
```http
GET /api/v1/oidc/1/login?next=https://evil.example/collect HTTP/1.1
Host: 127.0.0.1:9000
```
- Payloads de teste:
```text
https://evil.example/collect
//evil.example
https:%2f%2fevil.example
```
- Resposta esperada se vulneravel: `302/303` final para host externo.
- Resposta esperada se seguro: Rejeicao ou fallback para rota interna.
- Causa raiz: `RedirectURI` com input controlado pelo usuario.
- Correcao recomendada: Somente paths relativos internos.
- Estrategia de validacao pos-fix: Repetir com varias variantes absolutas e encoded.
- Cobertura: DAST

[DAST - Webhooks publicos spoofaveis]

- Status: Confirmado
- Severidade: Alta
- Categoria: OWASP API2:2023 Broken Authentication; CWE-345
- Componente afetado: Webhooks WhatsApp/Telegram
- Superficie de ataque: `GET/POST /api/v1/webhooks/whatsapp`, `POST /api/v1/webhooks/telegram`
- Pre-requisitos: Nenhum
- Descricao tecnica: Validar que GET aceita qualquer `hub.verify_token` nao vazio e POST aceita payload arbitrario.
- Evidencia: Handlers placeholder sem validacao criptografica.
- Impacto: Spoofing e ruido operacional.
- Exploracao passo a passo:
  1. Chamar GET com `hub.mode=subscribe`.
  2. Enviar POST arbitrario em JSON.
  3. Confirmar `200 OK`.
- Requisicao HTTP de teste:
```http
GET /api/v1/webhooks/whatsapp?hub.mode=subscribe&hub.verify_token=test&hub.challenge=1337 HTTP/1.1
Host: 127.0.0.1:9000
```
```http
POST /api/v1/webhooks/telegram HTTP/1.1
Host: 127.0.0.1:9000
Content-Type: application/json

{"update_id":1,"message":{"text":"spoofed"}}
```
- Payloads de teste:
```bash
curl -i 'http://127.0.0.1:9000/api/v1/webhooks/whatsapp?hub.mode=subscribe&hub.verify_token=test&hub.challenge=1337'
curl -i -X POST \
  -H 'Content-Type: application/json' \
  -d '{"update_id":1,"message":{"text":"spoofed"}}' \
  http://127.0.0.1:9000/api/v1/webhooks/telegram
```
- Resposta esperada se vulneravel: `200 OK` em ambos os casos.
- Resposta esperada se seguro: `403/401` sem segredo ou assinatura valida.
- Causa raiz: Endpoints publicos placeholder sem autenticacao de origem.
- Correcao recomendada: Verificacao estrita por segredo/assinatura.
- Estrategia de validacao pos-fix: Testar sem segredo e com segredo incorreto.
- Cobertura: DAST

[DAST - WebSocket cross-origin]

- Status: Provavel
- Severidade: Alta
- Categoria: OWASP A01:2021 Broken Access Control; CWE-346
- Componente afetado: `/ws`
- Superficie de ataque: Handshake websocket com `Origin` arbitrario
- Pre-requisitos: Sessao autenticada no browser da vitima
- Descricao tecnica: Validar CSWSH com `Origin` externo.
- Evidencia: `CheckOrigin == true`.
- Impacto: Leitura de eventos e abuso do canal autenticado.
- Exploracao passo a passo:
  1. Hospedar HTML malicioso.
  2. Abrir `new WebSocket('ws://127.0.0.1:9000/ws')`.
  3. Observar sucesso do handshake.
- Requisicao HTTP de teste:
```http
GET /ws HTTP/1.1
Host: 127.0.0.1:9000
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: dGVzdGtleXRlc3RrZXk=
Sec-WebSocket-Version: 13
Origin: https://evil.example
Cookie: canalgov_session=<sessao>
```
- Payloads de teste:
```js
new WebSocket('ws://127.0.0.1:9000/ws')
```
- Resposta esperada se vulneravel: `101 Switching Protocols`.
- Resposta esperada se seguro: `403` ou falha de upgrade.
- Causa raiz: Ausencia de validacao de origem.
- Correcao recomendada: Allowlist de `Origin`.
- Estrategia de validacao pos-fix: Repetir com origem valida e invalida.
- Cobertura: DAST

[DAST - Upload SVG inline]

- Status: Confirmado
- Severidade: Alta
- Categoria: OWASP A05:2021 Security Misconfiguration; CWE-434
- Componente afetado: Upload e serving de midia
- Superficie de ataque: `POST /api/v1/media`, `GET /uploads/{uuid}`
- Pre-requisitos: Sessao autenticada
- Descricao tecnica: Validar upload de SVG com `Content-Type: image/svg+xml` e serving inline subsequente.
- Evidencia: Wildcard de extensao, confianca em MIME do cliente e inline por `image/*`.
- Impacto: Hospedagem de conteudo ativo sob a origem da aplicacao.
- Exploracao passo a passo:
  1. Upload de SVG.
  2. Captura do UUID.
  3. GET do upload e inspecao de headers.
- Requisicao HTTP de teste:
```http
POST /api/v1/media HTTP/1.1
Host: 127.0.0.1:9000
Cookie: canalgov_session=<sessao>
Content-Type: multipart/form-data; boundary=----cg

------cg
Content-Disposition: form-data; name="files"; filename="poc.svg"
Content-Type: image/svg+xml

<svg xmlns="http://www.w3.org/2000/svg"><script>alert(1)</script></svg>
------cg--
```
- Payloads de teste:
```xml
<svg xmlns="http://www.w3.org/2000/svg"><script>alert(1)</script></svg>
```
- Resposta esperada se vulneravel: Upload aceito e `GET /uploads/{uuid}` com `Content-Type: image/svg+xml` e `Content-Disposition: inline`.
- Resposta esperada se seguro: Rejeicao, reclassificacao segura ou `attachment`.
- Causa raiz: Politica default-open e MIME trust.
- Correcao recomendada: Allowlist fechada e no-sniff.
- Estrategia de validacao pos-fix: Repetir com SVG e polyglot.
- Cobertura: DAST

[DAST - Signed URL de anexo sem sessao]

- Status: Confirmado
- Severidade: Media
- Categoria: OWASP API1:2023 Broken Object Level Authorization; CWE-639
- Componente afetado: `/uploads/{uuid}?sig=...&exp=...`
- Superficie de ataque: URL assinada retornada em anexos de mensagens
- Pre-requisitos: URL assinada capturada de uma resposta autorizada
- Descricao tecnica: Validar que a URL funciona como bearer token ate expirar.
- Evidencia: `authOrSignedURL` pula verificacoes de contexto.
- Impacto: Bypass de autorizacao por compartilhamento da URL.
- Exploracao passo a passo:
  1. Capturar `attachments[].url` de uma mensagem.
  2. Abrir em navegador sem sessao.
  3. Confirmar download.
- Requisicao HTTP de teste:
```http
GET /uploads/<uuid>?sig=<assinatura>&exp=<exp> HTTP/1.1
Host: 127.0.0.1:9000
```
- Payloads de teste:
```bash
curl -i 'http://127.0.0.1:9000/uploads/<uuid>?sig=<assinatura>&exp=<exp>'
```
- Resposta esperada se vulneravel: `200 OK` sem autenticacao.
- Resposta esperada se seguro: `401/403` ou assinatura vinculada ao contexto do usuario.
- Causa raiz: Bearer URL sem binding a identidade.
- Correcao recomendada: Download autenticado ou signed URL vinculada a contexto.
- Estrategia de validacao pos-fix: Repetir sem sessao e apos mudanca de permissao.
- Cobertura: DAST

[DAST - OIDC SSRF]

- Status: Confirmado
- Severidade: Media
- Categoria: OWASP API7:2023 Server Side Request Forgery; CWE-918
- Componente afetado: Criacao/atualizacao de provider OIDC
- Superficie de ataque: `POST /api/v1/oidc`, `PUT /api/v1/oidc/{id}`
- Pre-requisitos: Sessao admin com `oidc:manage`
- Descricao tecnica: Validar fetch server-side em `provider_url` controlado pelo usuario.
- Evidencia: `oidc.NewProvider(...)` em URL administravel.
- Impacto: SSRF para loopback, RFC1918 e metadata endpoints.
- Exploracao passo a passo:
  1. Subir listener HTTP de teste.
  2. Criar provider apontando para o listener.
  3. Confirmar callback no listener.
- Requisicao HTTP de teste:
```http
POST /api/v1/oidc HTTP/1.1
Host: 127.0.0.1:9000
Cookie: canalgov_session=<admin_session>
Content-Type: application/json

{
  "name":"evil",
  "provider":"Custom",
  "provider_url":"http://127.0.0.1:7777/.well-known/openid-configuration",
  "client_id":"x",
  "client_secret":"y",
  "logo_url":""
}
```
- Payloads de teste:
```text
http://127.0.0.1:7777/.well-known/openid-configuration
http://169.254.169.254/latest/meta-data/
http://<host-interno>/.well-known/openid-configuration
```
- Resposta esperada se vulneravel: Callback visivel no listener ou erro de parse apos conexao.
- Resposta esperada se seguro: Bloqueio local antes de conectar.
- Causa raiz: Discovery server-side sem guard SSRF.
- Correcao recomendada: Allowlist e bloqueio de IPs internos.
- Estrategia de validacao pos-fix: Repetir com loopback, link-local e RFC1918.
- Cobertura: DAST

[DAST - Brute force e reset sem throttling]

- Status: Confirmado
- Severidade: Media
- Categoria: OWASP A07:2021 Identification and Authentication Failures; CWE-307
- Componente afetado: Login e reset de senha
- Superficie de ataque: `POST /api/v1/auth/login`, `POST /api/v1/agents/reset-password`
- Pre-requisitos: Nenhum
- Descricao tecnica: Validar ausencia de rate limit, backoff e lockout.
- Evidencia: Nao ha throttling visivel no codigo.
- Impacto: Password spraying e abuso de reset.
- Exploracao passo a passo:
  1. Enviar 20-50 tentativas seguidas de login.
  2. Repetir reset para o mesmo email.
  3. Observar ausencia de `429` e lockout.
- Requisicao HTTP de teste:
```http
POST /api/v1/auth/login HTTP/1.1
Host: 127.0.0.1:9000
Content-Type: application/json

{"email":"target@example.gov","password":"WrongPass123!"}
```
- Payloads de teste:
```bash
for i in $(seq 1 30); do
  curl -s -o /dev/null -w '%{http_code}\n' \
    -H 'Content-Type: application/json' \
    -d '{"email":"target@example.gov","password":"WrongPass123!"}' \
    http://127.0.0.1:9000/api/v1/auth/login
done
```
- Resposta esperada se vulneravel: Sequencia ilimitada de respostas normais sem `429`.
- Resposta esperada se seguro: `429`, backoff, desafio adicional ou lockout temporario.
- Causa raiz: Ausencia de controle de abuso.
- Correcao recomendada: Rate limiting por IP e conta.
- Estrategia de validacao pos-fix: Repetir burst controlado e confirmar bloqueio.
- Cobertura: DAST

### Automacao sugerida

- `curl` / `bash`: smoke tests e regressao rapida para BAC, open redirect, webhooks e rate limit
- `pytest`: suite HTTP com fixtures de sessao admin e low-priv
- `Cypress`: validacao de websocket cross-origin, open redirect e renderizacao de anexos/HTML
- `nuclei`: templates custom para `/api/v1/macros`, `/api/v1/roles`, `/api/v1/inboxes`, `/api/v1/webhooks/whatsapp`, `/uploads/*?sig=...`

### Exemplos de automacao

```bash
# roles expostas
curl -s -H 'Cookie: canalgov_session=<low_priv_session>' \
  http://127.0.0.1:9000/api/v1/roles | jq '.data[] | {name, permissions}'

# inboxes expostas
curl -s -H 'Cookie: canalgov_session=<low_priv_session>' \
  http://127.0.0.1:9000/api/v1/inboxes | jq '.data[] | {name, from, config}'

# signed URL sem sessao
curl -i 'http://127.0.0.1:9000/uploads/<uuid>?sig=<assinatura>&exp=<exp>'

# webhook whatsapp spoof
curl -i 'http://127.0.0.1:9000/api/v1/webhooks/whatsapp?hub.mode=subscribe&hub.verify_token=test&hub.challenge=1337'
```

### SAST vs DAST — lacunas, complementaridade e estrategia de cobertura

- SAST foi superior para achar falhas de desenho e enforcement:
  - `roles`, `inboxes`, `macros`, signed URLs, `next`, proxy/IP trust, reload de cookies, SSRF em OIDC
- DAST e o melhor caminho para validar impacto operacional:
  - websocket cross-origin, upload SVG inline, brute force, open redirect completo, funcionamento real de signed URLs e stored XSS inbound
- Alguns achados permanecem assimetricos:
  - `fast-realip` e claramente SAST-first; o impacto operacional depende de inspeção de logs
  - stored XSS via inbound HTML permanece dependente de execucao real do frontend e do componente `vue-letter`
- Estrategia recomendada de cobertura:
  1. Validar primeiro todos os confirmados de BAC/exposure e redirect
  2. Em seguida, websocket, upload e signed URLs
  3. Depois, SSRF e brute force
  4. Fechar com XSS inbound/rich content e parser differential tests

## Fase 4 - Correlacao

### Matriz SAST x DAST

| Achado | SAST | DAST | Status de correlacao | Observacao |
|---|---|---|---|---|
| Exposicao de macros | Sim | Sim | Confirmado por ambos | SAST mostrou ausencia de filtro; DAST valida leak real |
| Exposicao de roles | Sim | Sim | Confirmado por ambos | SAST mostrou rota incorreta; DAST valida enumeração |
| Exposicao de inboxes | Sim | Sim | Confirmado por ambos | SAST mostrou payload rico; DAST valida leak operacional |
| Open redirect em `next` | Sim | Sim | Confirmado por ambos | SAST mostrou sink; DAST valida fluxo completo |
| Webhooks spoofaveis | Sim | Sim | Confirmado por ambos | SAST mostrou ausência de verificação; DAST confirma aceitação arbitrária |
| Signed URL bearer | Sim | Sim | Confirmado por ambos | SAST mostrou bypass por design; DAST confirma acesso sem sessão |
| Upload SVG inline | Sim | Sim | Confirmado por ambos | SAST mostrou policy/MIME/inline; DAST confirma comportamento final |
| OIDC SSRF | Sim | Sim | Confirmado por ambos | SAST mostrou fetch server-side; DAST confirma conexão |
| Rate limit ausente | Sim | Sim | Confirmado por ambos | SAST não mostra throttling; DAST confirma ausência de `429` |
| Websocket cross-origin | Sim | Planejado | SAST-only por enquanto | Precisa execução real de handshake cross-origin |
| Stored XSS inbound | Sim | Planejado | SAST-only por enquanto | Cadeia está clara; falta validação runtime do `vue-letter` |
| Spoofing de IP por proxy trust | Sim | Parcial | SAST-only por enquanto | Impacto operacional depende de inspeção de logs |
| Regressao `SecureCookies` em reload OIDC | Sim | Planejado | SAST-only por enquanto | Precisa validar `Set-Cookie` após mutate/reload |
| OIDC linking por e-mail sem `nonce` | Sim | Planejado | SAST-only por enquanto | Precisa harness/IdP controlado para exploração dinâmica |

### Porque alguns bugs aparecem em SAST e ainda nao em DAST

- Bugs de desenho e de enforcement aparecem cedo em SAST porque:
  - a rota e o sink estao explicitamente visiveis;
  - a ausencia de checagem e deterministica;
  - a exploracao nao depende de comportamento de browser ou provider externo.

- Bugs dependentes de ambiente aparecem melhor em DAST porque:
  - exigem comportamento real do browser (`websocket`, XSS, cookies);
  - dependem de IdP controlado ou listener SSRF;
  - precisam de observabilidade externa como logs e callbacks.

- Bugs arquiteturais como signed URLs e proxy trust tendem a:
  - ser detectados com alta confianca em SAST;
  - exigir DAST apenas para demonstrar impacto operacional ao time.

### Classificacao atual

- `Confirmado por ambos`:
  - macros
  - roles
  - inboxes
  - open redirect
  - webhooks spoofaveis
  - signed URLs bearer
  - upload SVG inline
  - OIDC SSRF
  - brute force/reset sem throttling

- `SAST-only`:
  - websocket cross-origin
  - stored XSS inbound HTML
  - proxy/IP spoofing em auditoria
  - regressao `SecureCookies`
  - OIDC account linking por e-mail sem `nonce`

- `DAST-only`:
  - Nenhum ate o momento

## Fase 5 - Relatorio Tecnico

### Resumo executivo

O monorepo CanalGov apresenta uma combinacao de falhas de Broken Access Control, overexposure de catalogo administrativo, autenticacao/SSO fragil, upload perigoso e integracoes publicas inseguras. Os achados de maior risco e melhor explorabilidade real concentram-se em:

- endpoints `auth(...)` que devolvem dados administrativos completos;
- URLs assinadas de anexos tratadas como bearer tokens;
- fluxo OIDC com redirect inseguro e trust excessivo;
- upload default-open com serving inline de tipos ativos;
- webhooks publicos sem verificacao de autenticidade;
- ausencia de controles de abuso em autenticacao.

Nao se trata apenas de bugs isolados de implementacao. Ha um padrao arquitetural recorrente: o backend frequentemente usa autenticacao generica onde deveria usar autorizacao contextual de negocio.

### Top 10 achados por risco

1. Upload perigoso por allowlist wildcard, confianca em `Content-Type` do cliente e serving inline
2. Webhooks publicos de WhatsApp e Telegram aceitam spoofing e verificacao incompleta
3. Exposicao de macros para qualquer usuario autenticado
4. OIDC account linking baseado apenas em e-mail, sem `email_verified` e sem `nonce`
5. SSRF em discovery de provider OIDC administravel
6. WebSocket sem validacao de origem
7. Exposicao de configuracao completa de inboxes para qualquer usuario autenticado
8. Exposicao de todas as roles e permissoes para qualquer usuario autenticado
9. Anexos expostos por signed URLs bearer sem binding a usuario ou contexto
10. Ausencia de rate limit e lockout em login e reset de senha

### Tabela resumida de vulnerabilidades

| Titulo | Status | Severidade | Cobertura | Tipo |
|---|---|---|---|---|
| Upload perigoso / inline ativo | Confirmado | Alta | Ambos | Codigo + arquitetura |
| Webhooks spoofaveis | Confirmado | Alta | Ambos | Codigo |
| Macros expostas | Confirmado | Alta | Ambos | Codigo |
| OIDC linking fraco | Confirmado | Alta | SAST | Arquitetura |
| OIDC SSRF | Confirmado | Media | Ambos | Arquitetura |
| WebSocket sem origin check | Provavel/Confirmado em SAST | Alta | SAST | Codigo |
| Inboxes expostas | Confirmado | Media | Ambos | Codigo |
| Roles expostas | Confirmado | Media | Ambos | Codigo |
| Signed URL bearer | Confirmado | Media | Ambos | Arquitetura |
| Brute force/reset sem throttling | Confirmado | Media | Ambos | Arquitetura |
| Open redirect em `next` | Confirmado | Media | Ambos | Codigo |
| Proxy/IP spoofing em logs | Confirmado | Media | SAST | Arquitetura |
| Regressao `SecureCookies` | Confirmado | Media | SAST | Codigo |
| Stored XSS inbound HTML | Provavel | Alta | SAST | Codigo + arquitetura |

### Backlog de correcoes por prioridade

#### P0 - corrigir imediatamente

- Fechar `GET /api/v1/macros`, `GET /api/v1/roles` e `GET /api/v1/inboxes` com autorizacao correta
- Bloquear upload de SVG e remover wildcard `app.allowed_file_upload_extensions`
- Desabilitar serving inline para tipos ativos nao estritamente seguros
- Corrigir webhooks publicos com verificacao criptografica obrigatoria
- Bloquear `next` externo e centralizar redirect seguro

#### P1 - alta prioridade

- Remover signed URLs bearer de APIs internas autenticadas
- Adicionar rate limit/lockout em login e reset
- Corrigir OIDC account linking: `issuer + sub`, `nonce`, `email_verified`
- Proteger discovery OIDC contra SSRF
- Validar `Origin` no websocket

#### P2 - hardening estrutural

- Sanitizar HTML inbound e outbound com allowlist centralizada
- Introduzir modelo de trusted proxies para IP real
- Corrigir regressao de `SecureCookies` em reload de auth
- Adicionar headers de hardening no serving de upload
- Revisar cache, source maps e artefatos de build/frontend

### Hardening recomendado de codigo, runtime e pipeline

#### Codigo

- Mover autorizacao contextual para o backend em todos os endpoints `auth(...)`
- Adotar helpers centralizados para:
  - redirect seguro
  - validacao de URL externa
  - sanitizacao HTML
  - resolucao de IP real com trusted proxies
- Criar testes unitarios de permissao por endpoint e por objeto

#### Runtime

- Definir allowlists estritas de upload e MIME
- Servir uploads ativos como `attachment`
- Ativar telemetria para brute force, webhook failures e callbacks OIDC
- Reduzir TTL de recursos compartilhaveis e permitir revogacao

#### Pipeline / DevSecOps

- Adicionar testes automatizados de BAC em CI
- Adicionar suites DAST smoke para rotas criticas
- Validar config segura entre dev e prod
- Auditar artefatos `frontend/dist` para source maps e leaks antes de deploy

### Onde procurar primeiro no codigo

```bash
rg -n 'auth\\(' cmd/handlers.go
rg -n 'RedirectURI\\(|Peek\\(\"next\"\\)' cmd
rg -n 'CheckOrigin|/ws' cmd frontend/src
rg -n 'GetURL\\(|SignedURLValidator|authOrSignedURL' cmd internal/media
rg -n 'GetAll\\(' cmd/macro.go cmd/roles.go cmd/inboxes.go
rg -n 'oidc.NewProvider|LoginURL|ExchangeOIDCToken' cmd internal/auth
rg -n 'MultipartForm|Content-Type|application/octet-stream|image/svg' cmd internal/media
rg -n 'ReadEnvelope|text/html|vue-letter' internal frontend/src
rg -n 'fast-realip|FromRequest' cmd
```
