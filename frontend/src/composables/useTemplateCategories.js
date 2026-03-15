import { ref } from 'vue'
import api from '@/api'

export function useTemplateCategories() {
  const loading = ref(false)
  const categories = ref([])

  const fetchCategories = async () => {
    loading.value = true
    try {
      const resp = await api.getTemplateCategories()
      categories.value = resp.data.data
      return categories.value
    } catch (err) {
      console.error('Error fetching template categories', err)
      return []
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    categories,
    fetchCategories
  }
}
