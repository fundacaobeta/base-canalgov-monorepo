# Usuarios de Desenvolvimento

Este ambiente local do CanalGov foi preparado com os usuarios abaixo para teste funcional da central de chamados.

## URL local

- Aplicacao: `http://localhost:9000`
- Swagger/OpenAPI: `http://localhost:9000/static/public/swagger.html`

## Senha padrao

Todos os usuarios abaixo usam a mesma senha padrao:

`CanalGov@123`

## Usuarios disponiveis

### 1. Admin principal

- Nome: Ana Silva
- Perfil: `Admin`
- Equipe: `Atendimento`
- Login: `ana.silva@canalgov.local`

### 2. Atendente de ouvidoria

- Nome: Bruno Costa
- Perfil: `Agent`
- Equipe: `Ouvidoria`
- Login: `bruno.costa@canalgov.local`

### 3. Atendente de protocolos

- Nome: Carla Souza
- Perfil: `Agent`
- Equipe: `Protocolos`
- Login: `carla.souza@canalgov.local`

### 4. Usuario tecnico do sistema

- Nome: System
- Perfil: usuario interno do sistema
- Login: `System`

## Observacoes

- Os usuarios `Ana`, `Bruno` e `Carla` sao mantidos pela seed de desenvolvimento.
- Para recriar os dados mock, use:

```bash
make dev-seed
```

- Para alterar a senha padrao dos usuarios seedados:

```bash
make dev-seed DEV_SEED_PASSWORD='SuaSenha@123'
```
