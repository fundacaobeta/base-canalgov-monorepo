<script setup lang="ts">
import { TriangleAlert } from 'lucide-vue-next'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { useI18n } from 'vue-i18n'

defineProps({
  errorMessage: String,
  border: Boolean
})

const { t } = useI18n()
</script>

<template>
  <Alert variant="destructive" v-if="errorMessage" :class="{ 'border-0': !border }">
    <TriangleAlert class="w-4 h-4" />
    <AlertTitle>{{ t('globals.ui.error') }}</AlertTitle>
    <AlertDescription>
      {{ errorMessage }}
    </AlertDescription>
  </Alert>
</template>
