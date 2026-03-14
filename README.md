<a href="https://zerodha.tech"><img src="https://zerodha.tech/static/images/github-badge.svg" align="right" alt="Zerodha Tech Badge" /></a>


# CanalGov

Central de atendimento moderna, open source e self-hosted. Aplicacao em binario unico.

![image](https://libredesk.io/hero.png?v=1)


Consulte a [documentacao local](./docs/README.md) para operacao e desenvolvimento.

## Recursos

- **Caixas compartilhadas**
  Gerencie conversas em multiplas filas e equipes.
- **Permissoes granulares**
  Crie funcoes com controle fino por equipe e por agente.
- **Automacoes**
  Automatize marcacao, atribuicao e roteamento de conversas.
- **Pesquisas CSAT**
  Meça satisfacao com pesquisas automaticas.
- **Macros**
  Padronize respostas e acoes frequentes em um clique.
- **Organizacao**
  Use tags, status personalizados e adiamento de conversas.
- **Autoatribuicao**
  Distribua carga conforme regras e capacidade.
- **SLA**
  Defina metas e receba alertas antes de violacoes.
- **Atributos personalizados**
  Amplie contatos e conversas com campos de negocio.
- **Assistencia com IA**
  Reescreva respostas para ajustar tom e clareza.
- **Logs de atividade**
  Audite alteracoes e eventos relevantes.
- **Webhooks**
  Integre eventos de conversa e mensagem em tempo real.
- **Barra de comando**
  Execute acoes rapidas com atalhos.


## Instalacao

### Docker

The latest image is available on DockerHub at [`canalgov/canalgov:latest`](https://hub.docker.com/)

```shell
# Download the compose file and sample config file in the current directory.
curl -LO https://github.com/abhinavxd/libredesk/raw/main/docker-compose.yml
curl -LO https://github.com/abhinavxd/libredesk/raw/main/config.sample.toml

# Copy the config.sample.toml to config.toml and edit it as needed.
cp config.sample.toml config.toml

# Run the services in the background.
docker compose up -d

# Setting System user password.
docker exec -it canalgov_app ./canalgov --set-system-user-password
```

Go to `http://localhost:9000` and login with username `System` and the password you set using the `--set-system-user-password` command.

Veja a documentacao local em [docs](./docs/README.md)

__________________

### Binary
- Download the latest release and extract the `canalgov` binary.
- Copy config.sample.toml to config.toml and edit as needed.
- `./canalgov --install` to setup the Postgres DB (or `--upgrade` to upgrade an existing DB. Upgrades are idempotent and running them multiple times have no side effects).
- Run `./canalgov --set-system-user-password` to set the password for the System user.
- Run `./canalgov` and visit `http://localhost:9000` and login with username `System` and the password you set using the --set-system-user-password command.

Veja a documentacao local em [docs](./docs/README.md)
__________________

## Desenvolvimento

- Para contribuir, leia [docs/contributing.md](./docs/contributing.md).
- Para usuarios e acessos de desenvolvimento, veja [docs/usuarios-dev.md](./docs/usuarios-dev.md).
- Para o frontend, veja [docs/frontend.md](./docs/frontend.md).
- Para direcao do projeto, veja [docs/roadmap.md](./docs/roadmap.md).

The backend is written in Go and the frontend is Vue.js 3 with Shadcn UI.



## Documentacao

A documentacao operacional e de desenvolvimento fica em [docs/README.md](./docs/README.md).
