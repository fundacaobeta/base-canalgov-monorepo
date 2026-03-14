---
name: review-vue
description: Review Vue 3 components following project conventions (Shadcn UI, Pinia, Zod, VeeValidate, TipTap)
argument-hint: [component-path]
allowed-tools: Read, Grep
---

# Vue 3 Component Review — CanalGov

Review `$ARGUMENTS` (ou arquivos Vue modificados) contra os padrões do projeto.

## Estrutura do Componente
- Usar `<script setup>` (Composition API)
- Props tipadas com TypeScript/JSDoc ou objeto de validação Vue
- Emits declarados explicitamente com `defineEmits`
- Sem lógica de negócio direta — usar composables em `src/composables/`

## Estado e Stores (Pinia)
- Estado global em stores (`src/stores/`) — não em `provide/inject` ad-hoc
- Stores acessadas via composable gerado por `useXxxStore()`
- Evitar mutações diretas fora de actions

## Formulários (VeeValidate + Zod)
- Schema de validação em arquivo separado `formSchema.js` ou `formSchema.ts`
- Usar `useForm` do VeeValidate com `zodResolver`
- Erros exibidos com componente `FormMessage` do Shadcn

## UI (Shadcn UI + TailwindCSS)
- Usar componentes de `src/components/ui/` antes de criar novos
- Classes Tailwind — sem `style` inline exceto valores dinâmicos
- Responsividade: mobile-first (`sm:`, `md:`, `lg:`)
- Ícones do `lucide-vue-next` ou `@radix-icons/vue`

## API Calls
- Chamadas via funções em `src/api/index.js` — nunca axios direto no componente
- Loading state com `ref(false)` + try/finally
- Erros exibidos via toast/notificação, não `console.error`

## i18n
- Textos visíveis ao usuário via `$t('chave')` ou `useI18n().t('chave')`
- Chave nova deve ser adicionada em `i18n/en.json` (e `pt-BR.json` se existir)

## O que verificar
1. Componente usa `<script setup>`?
2. Textos estão traduzidos via i18n?
3. Formulários usam VeeValidate + Zod?
4. Componentes Shadcn reutilizados onde possível?
5. API calls passam por `src/api/`?
6. Não há lógica de negócio inline — está em composable/store?

Reporte problemas com arquivo:linha e sugestão de correção.
