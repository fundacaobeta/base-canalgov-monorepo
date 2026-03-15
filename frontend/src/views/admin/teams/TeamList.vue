<template>
  <div>
    <DataTable :columns="columns" :data="data" :loading="isLoading" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { columns } from '@/features/admin/teams/TeamsDataTableColumns.js'
import DataTable from '@/components/datatable/DataTable.vue'
import { useAdminListRefresh } from '@/composables/useAdminListRefresh'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import api from '@/api'

const data = ref([])
const isLoading = ref(false)
const { showErrorToast } = useAdminErrorToast()

const getData = async () => {
  isLoading.value = true
  try {
    const response = await api.getTeams()
    data.value = response.data.data
  } catch (error) {
    showErrorToast(error)
  } finally {
    isLoading.value = false
  }
}

useAdminListRefresh('team', getData)
</script>
