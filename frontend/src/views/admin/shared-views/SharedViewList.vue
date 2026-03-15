<template>
  <div>
    <DataTable :columns="createColumns(t)" :data="sharedViews" :loading="isLoading" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import DataTable from '@/components/datatable/DataTable.vue'
import { createColumns } from '@/features/admin/shared-views/dataTableColumns.js'
import { useAdminListRefresh } from '@/composables/useAdminListRefresh'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'
import api from '@/api'

const { t } = useI18n()
const isLoading = ref(false)
const sharedViews = ref([])
const { showErrorToast } = useAdminErrorToast()

const getSharedViews = async () => {
  isLoading.value = true
  try {
    const resp = await api.getAllSharedViews()
    sharedViews.value = resp.data.data
  } catch (error) {
    showErrorToast(error)
  } finally {
    isLoading.value = false
  }
}

useAdminListRefresh('shared-views', getSharedViews)
</script>
