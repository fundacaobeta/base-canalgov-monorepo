<template>
  <div class="space-y-6">
    <div class="box p-6">
      <h2 class="text-xl font-semibold">{{ title }}</h2>
      <p class="mt-2 text-sm text-muted-foreground">
        {{ description }}
      </p>
    </div>

    <div class="grid gap-4 md:grid-cols-2">
      <div class="box p-5">
        <h3 class="font-medium">{{ t('admin.integration.status.title') }}</h3>
        <p class="mt-2 text-sm text-muted-foreground">
          {{ t('admin.integration.status.description', { title: title.toLowerCase() }) }}
        </p>
      </div>
      <div class="box p-5">
        <h3 class="font-medium">{{ t('admin.integration.nextSteps.title') }}</h3>
        <p class="mt-2 text-sm text-muted-foreground">
          {{ t('admin.integration.nextSteps.description') }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
const { t } = useI18n()

defineProps({
  title: {
    type: String,
    required: true
  },
  description: {
    type: String,
    required: true
  }
})
</script>
