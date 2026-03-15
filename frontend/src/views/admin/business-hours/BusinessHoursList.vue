<template>
  <div>
    <DataTable :columns="createColumns(t)" :data="businessHours" :loading="isLoading" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import DataTable from '@/components/datatable/DataTable.vue'
import { useI18n } from 'vue-i18n'
import { createColumns } from '@/features/admin/business-hours/dataTableColumns.js'
import { useAdminListRefresh } from '@/composables/useAdminListRefresh'
import api from '@/api'

const { t } = useI18n()
const businessHours = ref([])
const isLoading = ref(false)

const fetchAll = async () => {
  isLoading.value = true
  try {
    const resp = await api.getAllBusinessHours()
    businessHours.value = resp.data.data
  } finally {
    isLoading.value = false
  }
}

useAdminListRefresh('business_hours', fetchAll)
</script>
