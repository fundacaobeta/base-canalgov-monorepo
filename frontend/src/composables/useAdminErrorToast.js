import { useEmitter } from '@/composables/useEmitter'
import { EMITTER_EVENTS } from '@/constants/emitterEvents.js'
import { handleHTTPError } from '@/utils/http'

/**
 * Returns a helper that fires a destructive toast for HTTP errors.
 */
export function useAdminErrorToast() {
  const emitter = useEmitter()

  const showErrorToast = (error) => {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  }

  const showSuccessToast = (description) => {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, { description })
  }

  return { showErrorToast, showSuccessToast }
}
