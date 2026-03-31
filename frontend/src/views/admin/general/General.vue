<template>
  <div>
    <AdminPageHeader
      :title="$t('admin.general.title', 'Configurações Gerais')"
      :description="$t('admin.general.description', 'Personalize a identidade e o comportamento do seu workspace.')"
      :breadcrumbs="[{ label: $t('globals.terms.administration'), to: '/admin' }, { label: $t('admin.general.title', 'Configurações Gerais') }]"
    />

    <AdminPageWithHelp>
      <template #content>
        <div :class="{ 'opacity-50 pointer-events-none transition-opacity': isLoading }">
          <GeneralSettingForm :submitForm="submitForm" :initial-values="initialValues" />
        </div>
      </template>
      <template #help>
        <p>{{ $t('admin.general.help') }}</p>
        <p class="text-xs text-muted-foreground/70">
          {{ $t('admin.general.helpTimezone', 'O fuso horário afeta o cálculo de SLA e exibição de datas.') }}
        </p>
      </template>
    </AdminPageWithHelp>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import GeneralSettingForm from '@/features/admin/general/GeneralSettingForm.vue'
import AdminPageWithHelp from '@/layouts/admin/AdminPageWithHelp.vue'
import AdminPageHeader from '@/components/layout/AdminPageHeader.vue'
import { useAppSettingsStore } from '@/stores/appSettings'
import api from '@/api'

const initialValues = ref({})
const isLoading = ref(false)
const settingsStore = useAppSettingsStore()

onMounted(async () => {
  isLoading.value = true
  await settingsStore.fetchSettings('general')
  const data = settingsStore.settings
  isLoading.value = false
  initialValues.value = Object.keys(data).reduce((acc, key) => {
    // Remove 'app.' prefix
    const newKey = key.replace(/^app\./, '')
    acc[newKey] = data[key]
    return acc
  }, {})
})

const submitForm = async (values) => {
  // Prepend keys with `app.`
  const updatedValues = Object.fromEntries(
    Object.entries(values).map(([key, value]) => [`app.${key}`, value])
  )
  await api.updateSettings(updatedValues)
}
</script>
