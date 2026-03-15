<template>
  <div class="page-header">
    <nav
      v-if="breadcrumbs?.length"
      class="flex items-center gap-1.5 mb-2 text-xs text-muted-foreground"
      :aria-label="$t('globals.terms.navigation')"
    >
      <template v-for="(crumb, i) in breadcrumbs" :key="i">
        <RouterLink
          v-if="crumb.to"
          :to="crumb.to"
          class="hover:text-foreground transition-colors"
        >
          {{ crumb.label }}
        </RouterLink>
        <span
          v-else
          :class="i === breadcrumbs.length - 1 ? 'text-foreground font-medium' : ''"
        >
          {{ crumb.label }}
        </span>
        <ChevronRight
          v-if="i < breadcrumbs.length - 1"
          class="h-3 w-3 shrink-0"
          aria-hidden="true"
        />
      </template>
    </nav>
    <div class="page-header-inner">
      <div class="min-w-0 flex-1">
        <h1 class="page-title">{{ title }}</h1>
        <p v-if="description" class="page-description">{{ description }}</p>
      </div>
      <div v-if="$slots.actions" class="flex shrink-0 items-center gap-2 mt-0.5">
        <slot name="actions" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ChevronRight } from 'lucide-vue-next'

defineProps({
  title: { type: String, required: true },
  description: { type: String, default: null },
  breadcrumbs: { type: Array, default: null }
})
</script>
