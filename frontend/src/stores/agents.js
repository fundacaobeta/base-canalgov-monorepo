import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { handleHTTPError } from '@/utils/http'
import { useEmitter } from '@/composables/useEmitter'
import { EMITTER_EVENTS } from '@/constants/emitterEvents'
import api from '@/api'

export const useAgentsStore = defineStore('agents', () => {
    const agents = ref([])
    const emitter = useEmitter()
    const options = computed(() => agents.value.map(agent => ({
        label: agent.first_name + ' ' + agent.last_name,
        value: String(agent.id),
        avatar_url: agent.avatar_url,
    })))
    const fetchAgents = async (force = false) => {
        if (!force && agents.value.length) return
        try {
            const response = await api.getUsersCompact()
            agents.value = response?.data?.data || []
        } catch (error) {
            emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
                variant: 'destructive',
                description: handleHTTPError(error).message
            })
        }
    }
    return {
        agents,
        options,
        fetchAgents,
    }
})
