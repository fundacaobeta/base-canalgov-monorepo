<template>
  <div>
    <DataTable :columns="createColumns(t)" :data="roles" :loading="isLoading" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { createColumns } from '@/features/admin/roles/dataTableColumns.js'
import DataTable from '@/components/datatable/DataTable.vue'
import { useAdminListRefresh } from '@/composables/useAdminListRefresh'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'
import api from '@/api'

const { t } = useI18n()
const roles = ref([])
const isLoading = ref(false)
const { showErrorToast } = useAdminErrorToast()

const getRoles = async () => {
  isLoading.value = true
  try {
    const resp = await api.getRoles()
    roles.value = resp.data.data
  } catch (error) {
    showErrorToast(error)
  } finally {
    isLoading.value = false
  }
}

useAdminListRefresh('team', getRoles)
</script>
