<template>
  <div class="mb-5">
    <CustomBreadcrumb :links="breadcrumbLinks" />
  </div>
  <Spinner v-if="isLoading" />
  <OIDCForm
    :initial-values="oidc"
    :submitForm="submitForm"
    :isNewForm="isNewForm"
    :class="{ 'opacity-50 transition-opacity duration-300': isLoading }"
    :isLoading="formLoading"
  />
</template>

<script setup>
import { onMounted, ref, computed } from 'vue'
import api from '@/api'
import OIDCForm from '@/features/admin/oidc/OIDCForm.vue'
import { Spinner } from '@/components/ui/spinner'
import { CustomBreadcrumb } from '@/components/ui/breadcrumb'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const router = useRouter()
const { t } = useI18n()
const oidc = ref({
  provider: 'Google'
})
const { showErrorToast, showSuccessToast } = useAdminErrorToast()
const isLoading = ref(false)
const formLoading = ref(false)
const props = defineProps({
  id: {
    type: String,
    required: false
  }
})

const submitForm = async (values) => {
  try {
    let toastDescription = ''
    if (props.id) {
      if (values.client_secret.includes('•')) {
        values.client_secret = ''
      }
      await api.updateOIDC(props.id, values)
      toastDescription = t('globals.messages.updatedSuccessfully', {
        name: t('globals.terms.provider')
      })
    } else {
      await api.createOIDC(values)
      router.push({ name: 'sso-list' })
      toastDescription = t('globals.messages.createdSuccessfully', {
        name: t('globals.terms.provider')
      })
    }
    showSuccessToast(toastDescription)
  } catch (error) {
    showErrorToast(error)
  } finally {
    formLoading.value = false
  }
}

const breadCrumLabel = () => {
  return props.id ? t('globals.messages.edit') : t('globals.messages.new')
}

const isNewForm = computed(() => {
  return props.id ? false : true
})

const breadcrumbLinks = [
  { path: 'sso-list', label: t('globals.terms.sso') },
  { path: '', label: breadCrumLabel() }
]

onMounted(async () => {
  if (props.id) {
    try {
      isLoading.value = true
      const resp = await api.getOIDC(props.id)
      oidc.value = resp.data.data
    } catch (error) {
      showErrorToast(error)
    } finally {
      isLoading.value = false
    }
  }
})
</script>
