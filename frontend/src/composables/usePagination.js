import { ref, computed, watch } from 'vue'
import { getVisiblePages } from '@/utils/pagination'

/**
 * Reusable pagination composable.
 * @param {Function} fetchFn - async function called when page/perPage changes
 * @param {Object} [options]
 * @param {number} [options.defaultPerPage=15]
 */
export function usePagination(fetchFn, { defaultPerPage = 15 } = {}) {
  const page = ref(1)
  const perPage = ref(defaultPerPage)
  const totalCount = ref(0)
  const totalPages = ref(0)

  const visiblePages = computed(() => getVisiblePages(page.value, totalPages.value))

  function goToPage(p) {
    if (p >= 1 && p <= totalPages.value && p !== page.value) {
      page.value = p
    }
  }

  function handlePerPageChange() {
    page.value = 1
    fetchFn()
  }

  function setPaginationMeta({ count, total_pages }) {
    totalCount.value = count
    totalPages.value = total_pages
  }

  watch([page], fetchFn)

  return {
    page,
    perPage,
    totalCount,
    totalPages,
    visiblePages,
    goToPage,
    handlePerPageChange,
    setPaginationMeta
  }
}
