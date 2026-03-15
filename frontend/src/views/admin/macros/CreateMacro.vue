<template>
  <div class="mb-5">
    <CustomBreadcrumb :links="breadcrumbLinks" />
  </div>
  <MacroForm :submitForm="onSubmit" :isLoading="formLoading" />
</template>

<script setup>
import { ref } from 'vue'
import MacroForm from '@/features/admin/macros/MacroForm.vue'
import { CustomBreadcrumb } from '@/components/ui/breadcrumb'
import { useRouter } from 'vue-router'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'
import { useMacroStore } from '@/stores/macro'
import api from '@/api'

const router = useRouter()
const { showErrorToast, showSuccessToast } = useAdminErrorToast()
const { t } = useI18n()
const macroStore = useMacroStore()
const formLoading = ref(false)
const breadcrumbLinks = [
  { path: 'macro-list', label: t('globals.terms.macro', 2) },
  {
    path: '',
    label: t('globals.messages.new', {
      name: t('globals.terms.macro')
    })
  }
]

const onSubmit = (values) => {
  createMacro(values)
}

const createMacro = async (values) => {
  try {
    formLoading.value = true
    await api.createMacro(values)

    await macroStore.loadMacros(true)

    showSuccessToast(t('globals.messages.createdSuccessfully', {
      name: t('globals.terms.macro')
    }))
    router.push({ name: 'macro-list' })
  } catch (error) {
    showErrorToast(error)
  } finally {
    formLoading.value = false
  }
}
</script>
