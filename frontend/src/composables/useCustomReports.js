import { ref } from 'vue'
import api from '@/api'

export function useCustomReports() {
  const loading = ref(false)
  const reports = ref([])

  const fetchReports = async () => {
    loading.value = true
    try {
      const resp = await api.getCustomReports()
      reports.value = resp.data.data
      return reports.value
    } catch (err) {
      console.error('Error fetching custom reports', err)
      return []
    } finally {
      loading.value = false
    }
  }

  const createReport = async (data) => {
    try {
      const resp = await api.createCustomReport(data)
      return resp.data.data
    } catch (err) {
      console.error('Error creating custom report', err)
      throw err
    }
  }

  const updateReport = async (id, data) => {
    try {
      const resp = await api.updateCustomReport(id, data)
      return resp.data.data
    } catch (err) {
      console.error('Error updating custom report', id, err)
      throw err
    }
  }

  const deleteReport = async (id) => {
    try {
      await api.deleteCustomReport(id)
    } catch (err) {
      console.error('Error deleting custom report', id, err)
      throw err
    }
  }

  const executeReport = async (id) => {
    try {
      const resp = await api.executeCustomReport(id)
      return resp.data.data
    } catch (err) {
      console.error('Error executing custom report', id, err)
      return []
    }
  }

  return {
    loading,
    reports,
    fetchReports,
    createReport,
    updateReport,
    deleteReport,
    executeReport
  }
}
