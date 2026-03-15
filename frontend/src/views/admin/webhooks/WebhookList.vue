<template>
  <DataTable :columns="createColumns(t)" :data="webhooks" :loading="isLoading" />
</template>

<script setup>
import { ref } from 'vue'
import DataTable from '@/components/datatable/DataTable.vue'
import { createColumns } from '@/features/admin/webhooks/dataTableColumns.js'
import { useAdminListRefresh } from '@/composables/useAdminListRefresh'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'
import api from '@/api'

const webhooks = ref([])
const { t } = useI18n()
const isLoading = ref(false)
const { showErrorToast } = useAdminErrorToast()

const fetchAll = async () => {
  isLoading.value = true
  try {
    const resp = await api.getWebhooks()
    webhooks.value = resp.data.data
  } catch (error) {
    showErrorToast(error)
  } finally {
    isLoading.value = false
  }
}

useAdminListRefresh('webhook', fetchAll)
</script>
