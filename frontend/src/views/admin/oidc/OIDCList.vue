<template>
  <div>
    <DataTable :columns="createColumns(t)" :data="oidc" :loading="isLoading" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import DataTable from '@/components/datatable/DataTable.vue'
import { createColumns } from '@/features/admin/oidc/dataTableColumns.js'
import { useAdminListRefresh } from '@/composables/useAdminListRefresh'
import { useI18n } from 'vue-i18n'
import api from '@/api'

const oidc = ref([])
const { t } = useI18n()
const isLoading = ref(false)

const fetchAll = async () => {
  isLoading.value = true
  try {
    const resp = await api.getAllOIDC()
    oidc.value = resp.data.data
  } finally {
    isLoading.value = false
  }
}

useAdminListRefresh('oidc', fetchAll)
</script>
