<template>
  <div :class="wrapperClass" :title="absoluteLabel || undefined">
    <span v-if="showAbsolute && absoluteLabel">{{ absoluteLabel }}</span>
    <span v-if="showRelative && relativeLabel" :class="relativeClass">{{ relativeLabel }}</span>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { format, formatDistanceToNow } from 'date-fns'
import { ptBR, enUS } from 'date-fns/locale'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  value: {
    type: [String, Date],
    default: null
  },
  formatStr: {
    type: String,
    default: 'PPpp'
  },
  showAbsolute: {
    type: Boolean,
    default: true
  },
  showRelative: {
    type: Boolean,
    default: true
  },
  inline: {
    type: Boolean,
    default: false
  },
  centered: {
    type: Boolean,
    default: false
  },
  compact: {
    type: Boolean,
    default: false
  }
})

const { locale } = useI18n()

const dateValue = computed(() => {
  if (!props.value) return null
  const parsed = props.value instanceof Date ? props.value : new Date(props.value)
  return Number.isNaN(parsed.getTime()) ? null : parsed
})

const dateLocale = computed(() => (locale.value === 'pt-BR' ? ptBR : enUS))

const absoluteLabel = computed(() => {
  if (!dateValue.value) return ''
  return format(dateValue.value, props.formatStr, { locale: dateLocale.value })
})

const relativeLabel = computed(() => {
  if (!dateValue.value) return ''
  return formatDistanceToNow(dateValue.value, {
    addSuffix: true,
    locale: dateLocale.value
  })
})

const wrapperClass = computed(() => ({
  'flex items-center gap-2': props.inline,
  'flex flex-col': !props.inline,
  'text-center': props.centered,
  'text-xs text-muted-foreground': props.compact
}))

const relativeClass = computed(() => ({
  'text-xs text-muted-foreground': true
}))
</script>
