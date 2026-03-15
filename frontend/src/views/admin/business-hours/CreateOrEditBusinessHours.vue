<template>
  <div class="mb-5">
    <CustomBreadcrumb :links="breadcrumbLinks" />
  </div>
  <Spinner v-if="isLoading"></Spinner>
  <BusinessHoursForm
    :initial-values="businessHours"
    :submitForm="submitForm"
    :isNewForm="isNewForm"
    :class="{ 'opacity-50 transition-opacity duration-300': isLoading }"
    :isLoading="formLoading"
  />
</template>

<script setup>
import { onMounted, ref, computed } from 'vue'
import api from '@/api'
import BusinessHoursForm from '@/features/admin/business-hours/BusinessHoursForm.vue'
import { useRouter } from 'vue-router'
import { Spinner } from '@/components/ui/spinner'
import { CustomBreadcrumb } from '@/components/ui/breadcrumb'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const businessHours = ref({})
const { showErrorToast, showSuccessToast } = useAdminErrorToast()
const isLoading = ref(false)
const formLoading = ref(false)
const router = useRouter()
const props = defineProps({
  id: {
    type: String,
    required: false
  }
})

const submitForm = async (values) => {
  try {
    formLoading.value = true
    if (props.id) {
      await api.updateBusinessHours(props.id, values)
      showSuccessToast(t('globals.messages.updatedSuccessfully', {
        name: t('globals.terms.businessHour', 2)
      }))
    } else {
      await api.createBusinessHours(values)
      showSuccessToast(t('globals.messages.createdSuccessfully', {
        name: t('globals.terms.businessHour', 2)
      }))
      router.push({ name: 'business-hours-list' })
    }
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
  { path: 'business-hours-list', label: t('globals.terms.businessHour', 2) },
  { path: '', label: breadCrumLabel() }
]

onMounted(async () => {
  if (props.id) {
    try {
      isLoading.value = true
      const resp = await api.getBusinessHours(props.id)
      businessHours.value = resp.data.data
    } catch (error) {
      showErrorToast(error)
    } finally {
      isLoading.value = false
    }
  }
})
</script>
