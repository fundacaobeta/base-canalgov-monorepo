<script setup>
import { DotsHorizontalIcon } from '@radix-icons/vue'
import { cn } from '@/lib/utils'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  class: { type: null, required: false }
})

const { t } = useI18n();
</script>

<template>
  <span
    role="presentation"
    aria-hidden="true"
    :class="cn('flex h-9 w-9 items-center justify-center', props.class)"
  >
    <slot>
      <DotsHorizontalIcon class="h-4 w-4" />
    </slot>
    <span class="sr-only">{{ t('globals.ui.more') }}</span>
  </span>
</template>
