import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { getLocalePaths } from '@/router/paths'

// Provides locale-aware inbox type URL values and helpers.
// Components always work with the internal values (assigned/unassigned/all/mentioned)
// and this composable translates them to/from the localized URL params.
export function useInboxTypes() {
  const { locale } = useI18n()

  const types = computed(() => getLocalePaths(locale.value).inboxTypes)

  // Maps internal value → localized URL param
  const toParam = (internalType) => types.value[internalType] ?? internalType

  return { types, toParam }
}
