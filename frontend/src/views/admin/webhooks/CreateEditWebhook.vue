<template>
  <div class="mb-5">
    <CustomBreadcrumb :links="breadcrumbLinks" />
  </div>
  <div :class="{ 'opacity-50 transition-opacity duration-300': isLoading }">
    <WebhookForm @submit.prevent="onSubmit" :form="form" :isNewForm="isNewForm">
      <template #footer>
        <div class="flex space-x-3">
          <Button type="submit" :isLoading="formLoading">
            {{ isNewForm ? t('globals.messages.create') : t('globals.messages.update') }}
          </Button>
          <Button
            v-if="!isNewForm"
            type="button"
            variant="outline"
            :isLoading="testLoading"
            @click="handleTestWebhook"
          >
            {{ $t('globals.messages.send', { name: t('globals.terms.test').toLowerCase() }) }}
          </Button>
        </div>
      </template>
    </WebhookForm>
  </div>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue'
import api from '@/api'
import WebhookForm from '@/features/admin/webhooks/WebhookForm.vue'
import { CustomBreadcrumb } from '@/components/ui/breadcrumb'
import { Button } from '@/components/ui/button'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { createFormSchema } from '@/features/admin/webhooks/formSchema.js'

const router = useRouter()
const { t } = useI18n()
const isLoading = ref(false)
const formLoading = ref(false)
const testLoading = ref(false)
const { showErrorToast, showSuccessToast } = useAdminErrorToast()

const props = defineProps({ id: { type: String, required: false } })

const form = useForm({
  validationSchema: toTypedSchema(createFormSchema(t)),
  initialValues: { name: '', url: '', events: [], secret: '', is_active: true, headers: '{}' }
})

const isNewForm = computed(() => !props.id)

const breadcrumbLinks = [
  { path: 'webhook-list', label: t('globals.terms.webhook') },
  { path: '', label: props.id ? t('globals.messages.edit') : t('globals.messages.new') }
]

const onSubmit = form.handleSubmit(async (values) => {
  formLoading.value = true
  try {
    if (props.id) {
      if (values.secret && values.secret.includes('•')) values.secret = ''
      await api.updateWebhook(props.id, values)
      showSuccessToast(t('globals.messages.updatedSuccessfully', { name: t('globals.terms.webhook') }))
    } else {
      await api.createWebhook(values)
      showSuccessToast(t('globals.messages.createdSuccessfully', { name: t('globals.terms.webhook') }))
      router.push({ name: 'webhook-list' })
    }
  } catch (error) {
    showErrorToast(error)
  } finally {
    formLoading.value = false
  }
})

const handleTestWebhook = async () => {
  if (!props.id) return
  testLoading.value = true
  try {
    await api.testWebhook(props.id)
    showSuccessToast(t('globals.messages.sentSuccessfully', { name: t('globals.terms.webhook') }))
  } catch (error) {
    showErrorToast(error)
  } finally {
    testLoading.value = false
  }
}

onMounted(async () => {
  if (props.id) {
    isLoading.value = true
    try {
      const resp = await api.getWebhook(props.id)
      form.setValues(resp.data.data)
    } catch (error) {
      showErrorToast(error)
    } finally {
      isLoading.value = false
    }
  }
})
</script>
