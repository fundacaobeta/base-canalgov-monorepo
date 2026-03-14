---
name: i18n-check
description: Check for missing or unused i18n translation keys across Vue components and translation files
allowed-tools: Read, Grep, Glob
---

# i18n Consistency Check — CanalGov

Verifique a consistência das chaves de tradução no projeto.

## O que verificar

### 1. Chaves usadas no frontend mas ausentes nos arquivos de tradução
- Busque por `$t('...')`, `t('...')`, `useI18n()` nos arquivos `.vue` e `.js`
- Compare com as chaves em `i18n/en.json`
- Se `i18n/pt-BR.json` existir, verificar também

### 2. Chaves nos arquivos de tradução mas não usadas no código
- Listar chaves em `i18n/en.json` que não aparecem em nenhum arquivo Vue/JS
- Isso ajuda a limpar chaves obsoletas

### 3. Inconsistência entre arquivos de tradução
- Chaves presentes em `en.json` mas ausentes em `pt-BR.json` (ou vice-versa)
- Valores em branco/vazios

### 4. Textos hardcoded em português/inglês
- Buscar strings literais em português dentro de templates Vue que deveriam ser i18n
- Padrão suspeito: texto entre aspas dentro de `{{ }}` ou atributos `:placeholder`, `:label`, etc.

## Processo
1. Extrair todas as chaves usadas no frontend
2. Comparar com `i18n/en.json`
3. Listar: chaves faltando | chaves não utilizadas | inconsistências entre idiomas
4. Sugerir as entradas que precisam ser adicionadas

Reporte em formato de lista organizada por severidade.
