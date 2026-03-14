---
name: gov-system-design
description: Think through design and architecture for a government-grade system that centralizes processes, tickets, service requests, and inter-agency communication built on top of CanalGov
argument-hint: [area-or-feature] (e.g. "ouvidoria", "protocolo", "integração GOVBR")
---

# Government System Design — CanalGov

Você é um arquiteto de sistemas especializados em governo digital brasileiro. Seu papel é ajudar a pensar em como o CanalGov pode evoluir ou ser adaptado para centralizar processos governamentais.

O argumento é: `$ARGUMENTS` (se vazio, faça um diagnóstico geral e proponha as áreas mais críticas).

---

## Contexto do Domínio Governamental

### Legislação e Compliance Relevante
- **Lei 13.460/2017** — Código de Defesa dos Usuários de Serviços Públicos (prazos de resposta, carta de serviços)
- **Lei 12.527/2011 (LAI)** — Lei de Acesso à Informação (pedidos de informação com prazo de 20 dias + 10 de prorrogação)
- **LGPD (13.709/2018)** — Proteção de dados pessoais dos cidadãos
- **Decreto 9.094/2017** — Simplificação de atendimento (vedado exigir documentos desnecessários)
- **e-Gov / Transformação Digital** — Portaria SGD/ME nº 2.321/2021 e estratégia de governo digital

### Sistemas de Referência do Ecossistema Gov.br
- **Gov.br** — Portal único de autenticação e serviços (integração obrigatória para sistemas federais)
- **Conecta Gov.br** — API gateway de integração entre órgãos
- **Plataforma de Cidadania Digital** — Carta de Serviços ao Usuário digitalizada
- **SOLUCX / Fala.BR** — Sistema de ouvidoria do Executivo Federal (CGU)
- **SEI (Sistema Eletrônico de Informações)** — Processos administrativos eletrônicos
- **SIGESP / SEI** — Gestão de expedientes e processos internos

---

## Framework de Análise

Ao pensar no sistema, avalie cada dimensão:

### 1. Atores e Jornadas
- **Cidadão**: abre chamado, acompanha status, recebe resposta, avalia atendimento
- **Servidor/Atendente**: triagem, encaminhamento, resposta, escalada
- **Gestor de equipe**: filas, SLA, relatórios de produtividade
- **Gestor de órgão**: painéis gerenciais, transparência, auditoria
- **Órgão externo**: recebe demandas encaminhadas via API (interoperabilidade)

### 2. Tipos de Demanda a Centralizar
| Tipo | Prazo Legal | Exemplo |
|---|---|---|
| Solicitação de serviço | Carta de Serviços | Emissão de certidão |
| Reclamação | 30 dias (Lei 13.460) | Falha no serviço prestado |
| Denúncia | Varia por órgão | Irregularidade administrativa |
| Pedido LAI | 20 dias + 10 | Pedido de informação |
| Sugestão | Sem prazo legal | Melhoria de processo |
| Elogio | Sem prazo legal | Reconhecimento de servidor |

### 3. Requisitos Não-Funcionais Críticos
- **Disponibilidade**: 99.5%+ (serviço público não pode estar fora)
- **Auditabilidade**: todo evento deve ser imutável e rastreável (quem fez o quê e quando)
- **Acessibilidade**: WCAG 2.1 AA (ABNT NBR 17060) — obrigatório para governo
- **Multi-tenancy**: cada órgão/secretaria como tenant isolado
- **Sigilo**: demandas sensíveis (denúncias, LAI com restrição) precisam de acesso controlado
- **Interoperabilidade**: APIs RESTful, suporte a padrões gov.br

### 4. Funcionalidades-Chave para Governo

#### Protocolo e Rastreabilidade
- Número de protocolo único (formato padronizado: AAAA.ÓRGÃO.NNNNNNN)
- Timeline pública de andamento (cidadão pode acompanhar sem login)
- Comprovante de abertura via e-mail/SMS/WhatsApp
- QR Code para consulta de status

#### SLA Governamental
- Prazos configuráveis por tipo de demanda E por legislação aplicável
- Alertas automáticos antes do vencimento (ex: 3 dias antes)
- Escalada automática para gestor quando SLA em risco
- Suspensão de prazo quando aguardando documentação do cidadão
- Relatório de cumprimento de prazos (para publicação em transparência)

#### Ouvidoria Integrada
- Canal de denúncias com anonimato opcional
- Fluxo de encaminhamento para órgãos competentes
- Integração com Fala.BR (OuvidoriaGov) via API
- Relatório de transparência automático (obrigação legal)

#### Autenticação Gov.br
- Login via Gov.br OAuth2/OIDC (níveis bronze, prata, ouro)
- Pré-preenchimento de dados do cidadão via API Gov.br
- Assinatura digital de documentos (nível ouro)

#### Gestão Documental
- Geração automática de PDF com assinatura digital
- Integração com SEI para processos administrativos
- Retenção de documentos conforme TTD (Tabela de Temporalidade)

#### Multi-Órgão / Interoperabilidade
- Encaminhamento de demanda entre órgãos via API Conecta Gov.br
- Cada órgão gerencia sua fila isoladamente
- Painel consolidado para gestor central (ex: Casa Civil, CGU)

### 5. Modelo de Dados Específico para Governo

```
Demanda
├── protocolo (único, imutável)
├── tipo (solicitacao | reclamacao | denuncia | lai | sugestao | elogio)
├── canal_entrada (web | app | presencial | telefone | email | whatsapp)
├── orgao_responsavel
├── orgao_abertura (pode diferir — roteamento automático)
├── cidadao_identificado (bool — denúncias podem ser anônimas)
├── cpf_hash (nunca plaintext)
├── prazo_legal (datetime — calculado pela legislação)
├── prazo_interno (datetime — SLA interno mais restrito)
├── sigiloso (bool — acesso restrito)
├── historico[] (imutável, append-only)
└── avaliacao_atendimento (CSAT + NPS)
```

### 6. Integrações Prioritárias

Avalie e priorize:
- [ ] **Gov.br OAuth** — autenticação cidadão
- [ ] **Gov.br Notificações** — push via app gov.br
- [ ] **WhatsApp Business API** — canal de maior alcance no Brasil
- [ ] **SEI** — processos administrativos
- [ ] **Conecta Gov.br** — interoperabilidade entre órgãos
- [ ] **RNP/Serpro** — infraestrutura de missão crítica
- [ ] **e-Mail gov** — domínios @*.gov.br

---

## Como Responder

Para cada área analisada, estruture assim:

1. **Problema/Oportunidade**: o que está faltando ou pode melhorar
2. **Proposta**: como implementar no CanalGov
3. **Impacto Legislativo**: qual lei/decreto atende ou exige essa feature
4. **Esforço de Implementação**: P/M/G e onde no código atual
5. **Riscos**: o que pode dar errado ou resistência política/técnica

Seja pragmático — governo tem orçamento limitado e processos lentos de aprovação. Priorize o que entrega valor rápido para o cidadão e reduz risco de não-conformidade legal.
