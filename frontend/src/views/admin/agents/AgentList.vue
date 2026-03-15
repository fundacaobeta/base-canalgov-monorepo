<template>
  <div>
    <div class="flex justify-end mb-5 gap-2">
      <Importer
        entity-key="globals.terms.agent"
        :upload-fn="api.importAgents"
        :get-status-fn="api.getAgentImportStatus"
        @import-complete="getData"
      />
    </div>
    <DataTable :columns="createColumns(t)" :data="data" :loading="isLoading" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { createColumns } from '@/features/admin/agents/dataTableColumns.js'
import DataTable from '@/components/datatable/DataTable.vue'
import { useAgentsStore } from '@/stores/agents'
import { useI18n } from 'vue-i18n'
import Importer from '@/components/importer/Importer.vue'
import { useAdminListRefresh } from '@/composables/useAdminListRefresh'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import api from '@/api'

const isLoading = ref(false)
const usersStore = useAgentsStore()
const { t } = useI18n()
const data = ref([])
const { showErrorToast } = useAdminErrorToast()

const getData = async () => {
  isLoading.value = true
  try {
    await usersStore.fetchAgents(true)
    data.value = usersStore.agents
  } catch (error) {
    showErrorToast(error)
  } finally {
    isLoading.value = false
  }
}

useAdminListRefresh('agent', getData)
</script>
