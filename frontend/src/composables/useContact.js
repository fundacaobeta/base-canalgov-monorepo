import { ref } from 'vue'
import api from '@/api'

export function useContact() {
  const loading = ref(false)
  const error = ref(null)

  const getContact = async (id) => {
    loading.value = true
    try {
      const resp = await api.getContact(id)
      return resp.data.data
    } catch (err) {
      error.value = err
      throw err
    } finally {
      loading.value = false
    }
  }

  const getConversations = async (id) => {
    try {
      const resp = await api.getContactConversations(id)
      return resp.data.data
    } catch (err) {
      console.error('Error fetching contact conversations', err)
      return []
    }
  }

  const getStats = async (id) => {
    try {
      const resp = await api.getContactStats(id)
      return resp.data.data
    } catch (err) {
      console.error('Error fetching contact stats', err)
      return { total_conversations: 0, resolved_conversations: 0, open_conversations: 0 }
    }
  }

  const getNotes = async (id) => {
    try {
      const resp = await api.getContactNotes(id)
      return resp.data.data
    } catch (err) {
      console.error('Error fetching contact notes', err)
      return []
    }
  }

  const createNote = async (id, note) => {
    try {
      const resp = await api.createContactNote(id, { note })
      return resp.data.data
    } catch (err) {
      console.error('Error creating contact note', err)
      throw err
    }
  }

  const deleteNote = async (contactId, noteId) => {
    try {
      await api.deleteContactNote(contactId, noteId)
    } catch (err) {
      console.error('Error deleting contact note', err)
      throw err
    }
  }

  const updateContact = async (id, data) => {
    loading.value = true
    try {
      const resp = await api.updateContact(id, data)
      return resp.data.data
    } catch (err) {
      error.value = err
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    error,
    getContact,
    getConversations,
    getStats,
    getNotes,
    createNote,
    deleteNote,
    updateContact
  }
}
