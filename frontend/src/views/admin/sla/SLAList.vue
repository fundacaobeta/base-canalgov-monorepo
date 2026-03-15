<template>
  <div>
    <DataTable :columns="createColumns(t)" :data="slas" :loading="isLoading" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import DataTable from '@/components/datatable/DataTable.vue'
import { createColumns } from '@/features/admin/sla/dataTableColumns.js'
import { useAdminListRefresh } from '@/composables/useAdminListRefresh'
import { useI18n } from 'vue-i18n'
import api from '@/api'

const { t } = useI18n()
const slas = ref([])
const isLoading = ref(false)

const fetchAll = async () => {
  isLoading.value = true
  try {
    const resp = await api.getAllSLAs()
    slas.value = resp.data.data
  } finally {
    isLoading.value = false
  }
}

useAdminListRefresh('sla', fetchAll)
</script>
