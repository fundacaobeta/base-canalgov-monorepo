import { onMounted } from 'vue'
import { useEmitter } from '@/composables/useEmitter'
import { EMITTER_EVENTS } from '@/constants/emitterEvents.js'

/**
 * Registers a REFRESH_LIST listener for a specific model.
 * Calls fetchFn immediately on mount and whenever the matching event fires.
 *
 * @param {string} model - The model discriminator (e.g. 'tags', 'status', 'team')
 * @param {Function} fetchFn - Async function to call when refresh is triggered
 */
export function useAdminListRefresh(model, fetchFn) {
  const emitter = useEmitter()

  onMounted(() => {
    fetchFn()
    emitter.on(EMITTER_EVENTS.REFRESH_LIST, (data) => {
      if (data?.model === model) fetchFn()
    })
  })

  const emitRefresh = () => {
    emitter.emit(EMITTER_EVENTS.REFRESH_LIST, { model })
  }

  return { emitRefresh }
}
