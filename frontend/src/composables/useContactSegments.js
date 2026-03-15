import { ref } from 'vue'
import api from '@/api'

export function useContactSegments() {
  const loading = ref(false)
  const segments = ref([])

  const fetchSegments = async () => {
    loading.value = true
    try {
      const resp = await api.getContactSegments()
      segments.value = resp.data.data
      return segments.value
    } catch (err) {
      console.error('Error fetching contact segments', err)
      return []
    } finally {
      loading.value = false
    }
  }

  const createSegment = async (data) => {
    try {
      const resp = await api.createContactSegment(data)
      return resp.data.data
    } catch (err) {
      console.error('Error creating contact segment', err)
      throw err
    }
  }

  const updateSegment = async (id, data) => {
    try {
      const resp = await api.updateContactSegment(id, data)
      return resp.data.data
    } catch (err) {
      console.error('Error updating contact segment', id, err)
      throw err
    }
  }

  const deleteSegment = async (id) => {
    try {
      await api.deleteContactSegment(id)
    } catch (err) {
      console.error('Error deleting contact segment', id, err)
      throw err
    }
  }

  return {
    loading,
    segments,
    fetchSegments,
    createSegment,
    updateSegment,
    deleteSegment
  }
}
