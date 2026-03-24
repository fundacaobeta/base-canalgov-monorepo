# Security Smoke Tests

Artefatos leves para repetir parte da auditoria AppSec sem depender de harness complexo.

## Conteudo

- `smoke.sh`
  - valida exposicao de `macros`, `roles` e `inboxes`
  - valida webhook publico de WhatsApp
  - valida brute force simples em login
  - valida signed URL sem sessao quando `SIGNED_URL` for fornecida
- `nuclei/`
  - templates simples para BAC/data exposure e webhook spoofing

## Uso rapido

```bash
export BASE_URL='http://127.0.0.1:9000'
export LOW_PRIV_COOKIE='canalgov_session=<sessao_low_priv>'
export ADMIN_COOKIE='canalgov_session=<sessao_admin>'
export TEST_EMAIL='target@example.gov'
export TEST_PASSWORD='WrongPass123!'
export SIGNED_URL='http://127.0.0.1:9000/uploads/<uuid>?sig=<sig>&exp=<exp>'

bash scripts/security/smoke.sh
```

## Observacoes

- Os testes marcam como `VULNERABLE` quando o comportamento observado bate com os achados atuais.
- `SIGNED_URL` e opcional. Sem ele, o teste de upload assinado e pulado.
- Os templates `nuclei` assumem que o cookie pode ser passado via variable ou header custom.

## Nuclei

Exemplos:

```bash
nuclei -u http://127.0.0.1:9000 \
  -t scripts/security/nuclei/webhooks-whatsapp-spoof.yaml

nuclei -u http://127.0.0.1:9000 \
  -var cookie='canalgov_session=<sessao_low_priv>' \
  -t scripts/security/nuclei/roles-exposure.yaml
```
