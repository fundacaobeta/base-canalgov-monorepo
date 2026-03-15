<template>
  <div class="mb-5">
    <CustomBreadcrumb :links="breadcrumbLinks" />
  </div>
  <TemplateForm
    :initial-values="template"
    :submitForm="submitForm"
    :isLoading="formLoading"
  />
</template>

<script setup>
import { onMounted, ref } from 'vue'
import api from '@/api'
import TemplateForm from '@/features/admin/templates/TemplateForm.vue'
import { useRouter, useRoute } from 'vue-router'
import { CustomBreadcrumb } from '@/components/ui/breadcrumb'
import { useAdminListRefresh } from '@/composables/useAdminListRefresh'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useEmitter } from '@/composables/useEmitter'
import { EMITTER_EVENTS } from '@/constants/emitterEvents.js'
import { useI18n } from 'vue-i18n'

const template = ref({})
const { t } = useI18n()
const formLoading = ref(false)
const emitter = useEmitter()
const router = useRouter()
const route = useRoute()
const { emitRefresh } = useAdminListRefresh('templates', () => {})
const { showErrorToast, showSuccessToast } = useAdminErrorToast()

const props = defineProps({
  id: { type: String, required: false, default: null }
})

const submitForm = async (values) => {
  formLoading.value = true
  try {
    if (props.id) {
      await api.updateTemplate(props.id, values)
      showSuccessToast(t('globals.messages.updatedSuccessfully', { name: t('globals.terms.template') }))
    } else {
      await api.createTemplate(values)
      showSuccessToast(t('globals.messages.createdSuccessfully', { name: t('globals.terms.template') }))
      emitRefresh()
      router.push({ name: 'templates' })
    }
  } catch (error) {
    showErrorToast(error)
  } finally {
    formLoading.value = false
  }
}

const breadCrumLabel = () => props.id ? t('globals.messages.edit') : t('globals.messages.new')

const breadcrumbLinks = [
  { path: 'templates', label: t('globals.terms.template') },
  { path: '', label: breadCrumLabel() }
]

onMounted(async () => {
  if (props.id) {
    formLoading.value = true
    try {
      const resp = await api.getTemplate(props.id)
      template.value = resp.data.data
    } catch (error) {
      showErrorToast(error)
    } finally {
      formLoading.value = false
    }
  } else {
    template.value = { type: route.query.type || 'response', team_id: null, is_default: false }
  }
})
</script>
