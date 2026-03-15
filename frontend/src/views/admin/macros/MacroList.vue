<template>
  <div>
    <DataTable :columns="createColumns(t)" :data="macros" :loading="isLoading" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import DataTable from '@/components/datatable/DataTable.vue'
import { createColumns } from '@/features/admin/macros/dataTableColumns.js'
import { useAdminListRefresh } from '@/composables/useAdminListRefresh'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'
import api from '@/api'

const { t } = useI18n()
const isLoading = ref(false)
const macros = ref([])
const { showErrorToast } = useAdminErrorToast()

const getMacros = async () => {
  isLoading.value = true
  try {
    const resp = await api.getAllMacros()
    macros.value = resp.data.data
  } catch (error) {
    showErrorToast(error)
  } finally {
    isLoading.value = false
  }
}

useAdminListRefresh('macros', getMacros)
</script>
