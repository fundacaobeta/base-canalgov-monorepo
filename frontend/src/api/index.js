import axios from 'axios'

const http = axios.create({
  timeout: 10000,
  withCredentials: true
})

http.interceptors.request.use((config) => {
  const csrfCookie = document.cookie.split('; ').find((row) => row.startsWith('csrf_token='))
  const csrfToken = csrfCookie ? csrfCookie.substring('csrf_token='.length) : null

  if (csrfToken) {
    config.headers['X-CSRFTOKEN'] = csrfToken
  }

  return config
})

// Reports & Overview
const getOverviewCounts = () => http.get('/api/v1/reports/overview/counts')
const getOverviewCharts = (params) => http.get('/api/v1/reports/overview/charts', { params })
const getOverviewSLA = (params) => http.get('/api/v1/reports/overview/sla', { params })
const getOverviewCSAT = (params) => http.get('/api/v1/reports/overview/csat', { params })
const getOverviewMessageVolume = (params) => http.get('/api/v1/reports/overview/messages', { params })
const getOverviewTagDistribution = (params) => http.get('/api/v1/reports/overview/tags', { params })

// Custom Reports
const getCustomReports = () => http.get('/api/v1/reports/custom')
const getCustomReport = (id) => http.get(`/api/v1/reports/custom/${id}`)
const createCustomReport = (data) => http.post('/api/v1/reports/custom', data)
const updateCustomReport = (id, data) => http.put(`/api/v1/reports/custom/${id}`, data)
const deleteCustomReport = (id) => http.delete(`/api/v1/reports/custom/${id}`)
const executeCustomReport = (id) => http.get(`/api/v1/reports/custom/${id}/execute`)

// Auth & Language
const login = (data) => http.post('/api/v1/auth/login', data, { headers: { 'Content-Type': 'application/json' } })
const logout = () => http.get('/logout')
const getLanguage = (lang) => http.get(`/api/v1/lang/${lang}`)
const getConfig = () => http.get('/api/v1/config')

// Conversations
const getConversations = (params) => http.get('/api/v1/conversations/all', { params })
const getAllConversations = getConversations
const getUnassignedConversations = (params) => http.get('/api/v1/conversations/unassigned', { params })
const getAssignedConversations = (params) => http.get('/api/v1/conversations/assigned', { params })
const getMentionedConversations = (params) => http.get('/api/v1/conversations/mentioned', { params })
const getTeamUnassignedConversations = (id, params) => http.get(`/api/v1/teams/${id}/conversations/unassigned`, { params })
const getViewConversations = (id, params) => http.get(`/api/v1/views/${id}/conversations`, { params })
const getConversation = (uuid) => http.get(`/api/v1/conversations/${uuid}`)
const getConversationParticipants = (uuid) => http.get(`/api/v1/conversations/${uuid}/participants`)
const updateUserAssignee = (uuid, data) => http.put(`/api/v1/conversations/${uuid}/assignee/user`, data)
const updateTeamAssignee = (uuid, data) => http.put(`/api/v1/conversations/${uuid}/assignee/team`, data)
const updateAssignee = (uuid, type, data) => {
  if (type === 'team') {
    return updateTeamAssignee(uuid, data)
  }

  return updateUserAssignee(uuid, data)
}
const removeUserAssignee = (uuid) => http.put(`/api/v1/conversations/${uuid}/assignee/user/remove`)
const removeTeamAssignee = (uuid) => http.put(`/api/v1/conversations/${uuid}/assignee/team/remove`)
const removeAssignee = removeUserAssignee
const updatePriority = (uuid, data) => http.put(`/api/v1/conversations/${uuid}/priority`, data)
const updateConversationPriority = updatePriority
const updateStatus = (uuid, data) =>
  http.put(`/api/v1/conversations/${uuid}/status`, data, {
    headers: { 'Content-Type': 'application/json' }
  })
const updateConversationStatus = updateStatus
const updateLastSeen = (uuid) => http.put(`/api/v1/conversations/${uuid}/last-seen`)
const updateAssigneeLastSeen = updateLastSeen
const markUnread = (uuid) => http.put(`/api/v1/conversations/${uuid}/mark-unread`)
const markConversationAsUnread = markUnread
const updateTags = (uuid, data) => http.post(`/api/v1/conversations/${uuid}/tags`, data)
const upsertTags = updateTags
const getConversationMessages = (uuid, params) => http.get(`/api/v1/conversations/${uuid}/messages`, { params })
const getConversationMessage = (cuuid, uuid) => http.get(`/api/v1/conversations/${cuuid}/messages/${uuid}`)
const sendMessage = (uuid, data) => http.post(`/api/v1/conversations/${uuid}/messages`, data)
const retryMessage = (cuuid, uuid) => http.put(`/api/v1/conversations/${cuuid}/messages/${uuid}/retry`)
const createConversation = (data) =>
  http.post('/api/v1/conversations', JSON.parse(JSON.stringify(data)), {
    headers: { 'Content-Type': 'application/json' }
  })
const updateConversationCustomAttributes = (uuid, data) => http.put(`/api/v1/conversations/${uuid}/custom-attributes`, data)
const updateConversationCustomAttribute = updateConversationCustomAttributes
const updateContactCustomAttributes = (uuid, data) => http.put(`/api/v1/conversations/${uuid}/contacts/custom-attributes`, data)
const searchConversations = (params) => http.get('/api/v1/conversations/search', { params })
const searchMessages = (params) => http.get('/api/v1/messages/search', { params })
const searchContacts = (params) => http.get('/api/v1/contacts/search', { params })

// Settings
const getSettings = () => http.get('/api/v1/settings/general')
const updateSettings = (data) => http.put('/api/v1/settings/general', data)
const getEmailNotificationSettings = () => http.get('/api/v1/settings/notifications/email')
const updateEmailNotificationSettings = (data) => http.put('/api/v1/settings/notifications/email', data)
const getWhatsAppNotificationSettings = () => http.get('/api/v1/settings/notifications/whatsapp')
const updateWhatsAppNotificationSettings = (data) => http.put('/api/v1/settings/notifications/whatsapp', data)
const getTelegramNotificationSettings = () => http.get('/api/v1/settings/notifications/telegram')
const updateTelegramNotificationSettings = (data) => http.put('/api/v1/settings/notifications/telegram', data)
const getSMSNotificationSettings = () => http.get('/api/v1/settings/notifications/sms')
const updateSMSNotificationSettings = (data) => http.put('/api/v1/settings/notifications/sms', data)
const getPushNotificationSettings = () => http.get('/api/v1/settings/notifications/push')
const updatePushNotificationSettings = (data) => http.put('/api/v1/settings/notifications/push', data)
const getOfficialCommunicationsNotificationSettings = () => http.get('/api/v1/settings/notifications/official-communications')
const updateOfficialCommunicationsNotificationSettings = (data) => http.put('/api/v1/settings/notifications/official-communications', data)
const getMailDomainsSettings = () => http.get('/api/v1/settings/mail/domains')
const updateMailDomainsSettings = (data) => http.put('/api/v1/settings/mail/domains', data)

// Agents / Users
const getAgents = () => http.get('/api/v1/agents')
const getAgent = (id) => http.get(`/api/v1/agents/${id}`)
const getUser = getAgent
const createAgent = (data) => http.post('/api/v1/agents', data)
const createUser = createAgent
const updateAgent = (id, data) => http.put(`/api/v1/agents/${id}`, data)
const updateUser = updateAgent
const deleteAgent = (id) => http.delete(`/api/v1/agents/${id}`)
const deleteUser = deleteAgent
const importAgents = (data) => http.post('/api/v1/agents/import', data, { headers: { 'Content-Type': 'multipart/form-data' } })
const getAgentImportStatus = () => http.get('/api/v1/agents/import/status')
const getCurrentAgent = () => http.get('/api/v1/agents/me')
const getCurrentUser = getCurrentAgent
const updateCurrentAgent = (data) => http.put('/api/v1/agents/me', data, { headers: { 'Content-Type': 'multipart/form-data' } })
const updateCurrentUser = updateCurrentAgent
const getCurrentAgentTeams = () => http.get('/api/v1/agents/me/teams')
const updateAgentAvailability = (data) => http.put('/api/v1/agents/me/availability', data)
const updateCurrentUserAvailability = updateAgentAvailability
const deleteCurrentAgentAvatar = () => http.delete('/api/v1/agents/me/avatar')
const deleteUserAvatar = deleteCurrentAgentAvatar
const getUsersCompact = () => http.get('/api/v1/agents/compact')
const resetPassword = (data) => http.post('/api/v1/agents/reset-password', data)
const setPassword = (data) => http.post('/api/v1/agents/set-password', data)
const generateAPIKey = (id) => http.post(`/api/v1/agents/${id}/api-key`)
const revokeAPIKey = (id) => http.delete(`/api/v1/agents/${id}/api-key`)

// Media
const uploadMedia = (data) => {
  const formData = new FormData()
  formData.append('files', data.files)
  formData.append('inline', String(Boolean(data.inline)))
  if (data.linked_model) {
    formData.append('linked_model', data.linked_model)
  }

  return http.post('/api/v1/media', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

// Contacts
const getContacts = (params) => http.get('/api/v1/contacts', { params })
const getContact = (id) => http.get(`/api/v1/contacts/${id}`)
const getContactConversations = (id) => http.get(`/api/v1/contacts/${id}/conversations`)
const getContactStats = (id) => http.get(`/api/v1/contacts/${id}/stats`)
const updateContact = (id, data) => http.put(`/api/v1/contacts/${id}`, data, { headers: { 'Content-Type': 'multipart/form-data' } })
const blockContact = (id, data) => http.put(`/api/v1/contacts/${id}/block`, data)
const getContactNotes = (id) => http.get(`/api/v1/contacts/${id}/notes`)
const createContactNote = (id, data) => http.post(`/api/v1/contacts/${id}/notes`, data)
const deleteContactNote = (contactId, noteId) => http.delete(`/api/v1/contacts/${contactId}/notes/${noteId}`)
const getContactSegments = () => http.get('/api/v1/admin/contact-segments')
const createContactSegment = (data) => http.post('/api/v1/admin/contact-segments', data)
const getContactSegment = (id) => http.get(`/api/v1/admin/contact-segments/${id}`)
const updateContactSegment = (id, data) => http.put(`/api/v1/admin/contact-segments/${id}`, data)
const deleteContactSegment = (id) => http.delete(`/api/v1/admin/contact-segments/${id}`)
const updateContactCustomAttribute = (id, data) => http.put(`/api/v1/contacts/${id}/custom-attributes`, data)

// Teams
const getTeams = () => http.get('/api/v1/teams')
const getTeam = (id) => http.get(`/api/v1/teams/${id}`)
const createTeam = (data) => http.post('/api/v1/teams', data)
const updateTeam = (id, data) => http.put(`/api/v1/teams/${id}`, data)
const deleteTeam = (id) => http.delete(`/api/v1/teams/${id}`)
const getTeamsCompact = () => http.get('/api/v1/teams/compact')

// Inboxes & Webhooks
const getInboxes = () => http.get('/api/v1/inboxes')
const getInbox = (id) => http.get(`/api/v1/inboxes/${id}`)
const createInbox = (data) => http.post('/api/v1/inboxes', data)
const updateInbox = (id, data) => http.put(`/api/v1/inboxes/${id}`, data)
const toggleInbox = (id) => http.put(`/api/v1/inboxes/${id}/toggle`)
const deleteInbox = (id) => http.delete(`/api/v1/inboxes/${id}`)
const initiateOAuthFlow = (provider) => http.post(`/api/v1/inboxes/oauth/${provider}/authorize`)
const getWebhooks = () => http.get('/api/v1/webhooks')
const getWebhook = (id) => http.get(`/api/v1/webhooks/${id}`)
const createWebhook = (data) => http.post('/api/v1/webhooks', data)
const updateWebhook = (id, data) => http.put(`/api/v1/webhooks/${id}`, data)
const deleteWebhook = (id) => http.delete(`/api/v1/webhooks/${id}`)
const toggleWebhook = (id) => http.put(`/api/v1/webhooks/${id}/toggle`)
const testWebhook = (id) => http.post(`/api/v1/webhooks/${id}/test`)

// Roles & OIDC
const getRoles = () => http.get('/api/v1/roles')
const getRole = (id) => http.get(`/api/v1/roles/${id}`)
const createRole = (data) => http.post('/api/v1/roles', data)
const updateRole = (id, data) => http.put(`/api/v1/roles/${id}`, data)
const deleteRole = (id) => http.delete(`/api/v1/roles/${id}`)
const getAllOIDC = () => http.get('/api/v1/oidc')
const getOIDC = (id) => http.get(`/api/v1/oidc/${id}`)
const createOIDC = (data) => http.post('/api/v1/oidc', data)
const updateOIDC = (id, data) => http.put(`/api/v1/oidc/${id}`, data)
const deleteOIDC = (id) => http.delete(`/api/v1/oidc/${id}`)

// Templates & Categories
const getTemplateCategories = () => http.get('/api/v1/admin/template-categories')
const createTemplateCategory = (data) => http.post('/api/v1/admin/template-categories', data)
const updateTemplateCategory = (id, data) => http.put(`/api/v1/admin/template-categories/${id}`, data)
const deleteTemplateCategory = (id) => http.delete(`/api/v1/admin/template-categories/${id}`)
const getTemplates = (type, params = {}) => http.get('/api/v1/templates', { params: { type, ...params } })
const getTemplate = (id) => http.get(`/api/v1/templates/${id}`)
const createTemplate = (data) => http.post('/api/v1/templates', data)
const updateTemplate = (id, data) => http.put(`/api/v1/templates/${id}`, data)
const deleteTemplate = (id) => http.delete(`/api/v1/templates/${id}`)

// Business Hours & SLA
const getBusinessHours = () => http.get('/api/v1/business-hours')
const getAllBusinessHours = getBusinessHours
const getBusinessHour = (id) => http.get(`/api/v1/business-hours/${id}`)
const createBusinessHours = (data) => http.post('/api/v1/business-hours', data)
const updateBusinessHours = (id, data) => http.put(`/api/v1/business-hours/${id}`, data)
const deleteBusinessHour = (id) => http.delete(`/api/v1/business-hours/${id}`)
const deleteBusinessHours = deleteBusinessHour
const getSLAs = () => http.get('/api/v1/sla')
const getAllSLAs = getSLAs
const getSLA = (id) => http.get(`/api/v1/sla/${id}`)
const createSLA = (data) => http.post('/api/v1/sla', data)
const updateSLA = (id, data) => http.put(`/api/v1/sla/${id}`, data)
const deleteSLA = (id) => http.delete(`/api/v1/sla/${id}`)

// Automation & AI
const getAutomationRules = () => http.get('/api/v1/automations/rules')
const getAutomationRule = (id) => http.get(`/api/v1/automations/rules/${id}`)
const createAutomationRule = (data) => http.post('/api/v1/automations/rules', data)
const updateAutomationRule = (id, data) => http.put(`/api/v1/automations/rules/${id}`, data)
const toggleAutomationRule = (id) => http.put(`/api/v1/automations/rules/${id}/toggle`)
const updateAutomationRuleWeights = (data) => http.put('/api/v1/automations/rules/weights', data)
const updateAutomationRuleExecutionMode = (data) => http.put('/api/v1/automations/rules/execution-mode', data)
const updateAutomationRulesExecutionMode = updateAutomationRuleExecutionMode
const deleteAutomationRule = (id) => http.delete(`/api/v1/automations/rules/${id}`)
const getAiPrompts = () => http.get('/api/v1/ai/prompts')
const aiCompletion = (data) => http.post('/api/v1/ai/completion', data)
const updateAiProvider = (data) => http.put('/api/v1/ai/provider', data)
const updateAIProvider = updateAiProvider

// Views
const getCurrentUserViews = () => http.get('/api/v1/views/me')
const createView = (data) => http.post('/api/v1/views/me', data)
const updateView = (id, data) => http.put(`/api/v1/views/me/${id}`, data)
const deleteView = (id) => http.delete(`/api/v1/views/me/${id}`)
const getSharedViews = () => http.get('/api/v1/views/shared')
const getAllSharedViews = () => http.get('/api/v1/shared-views')
const getSharedView = (id) => http.get(`/api/v1/shared-views/${id}`)
const createSharedView = (data) => http.post('/api/v1/shared-views', data)
const updateSharedView = (id, data) => http.put(`/api/v1/shared-views/${id}`, data)
const deleteSharedView = (id) => http.delete(`/api/v1/shared-views/${id}`)

// Tags, Statuses, Priorities
const getTags = () => http.get('/api/v1/tags')
const createTag = (data) => http.post('/api/v1/tags', data)
const updateTag = (id, data) => http.put(`/api/v1/tags/${id}`, data)
const deleteTag = (id) => http.delete(`/api/v1/tags/${id}`)
const getStatuses = () => http.get('/api/v1/statuses')
const createStatus = (data) => http.post('/api/v1/statuses', data)
const updateStatusEntity = (id, data) => http.put(`/api/v1/statuses/${id}`, data)
const deleteStatus = (id) => http.delete(`/api/v1/statuses/${id}`)
const getPriorities = () => http.get('/api/v1/priorities')

// Custom Attributes
const getCustomAttributes = () => http.get('/api/v1/custom-attributes')
const getCustomAttribute = (id) => http.get(`/api/v1/custom-attributes/${id}`)
const createCustomAttribute = (data) => http.post('/api/v1/custom-attributes', data)
const updateCustomAttribute = (id, data) => http.put(`/api/v1/custom-attributes/${id}`, data)
const deleteCustomAttribute = (id) => http.delete(`/api/v1/custom-attributes/${id}`)

// Notifications & Logs
const getActivityLogs = (params) => http.get('/api/v1/activity-logs', { params })
const getNotifications = () => http.get('/api/v1/notifications')
const getNotificationStats = () => http.get('/api/v1/notifications/stats')
const markNotificationAsRead = (id) => http.put(`/api/v1/notifications/${id}/read`)
const markAllNotificationsAsRead = () => http.put('/api/v1/notifications/read-all')
const deleteNotification = (id) => http.delete(`/api/v1/notifications/${id}`)
const deleteAllNotifications = () => http.delete('/api/v1/notifications')

// Macros
const getMacros = () => http.get('/api/v1/macros')
const getAllMacros = getMacros
const getMacro = (id) => http.get(`/api/v1/macros/${id}`)
const createMacro = (data) => http.post('/api/v1/macros', data)
const updateMacro = (id, data) => http.put(`/api/v1/macros/${id}`, data)
const deleteMacro = (id) => http.delete(`/api/v1/macros/${id}`)
const applyMacro = (conversationUuid, id) => http.post(`/api/v1/conversations/${conversationUuid}/macros/${id}/apply`)

// Drafts
const getAllDrafts = () => http.get('/api/v1/drafts')
const saveDraft = (uuid, data) => http.post(`/api/v1/conversations/${uuid}/draft`, data)
const deleteDraft = (uuid) => http.delete(`/api/v1/conversations/${uuid}/draft`)

// Page visits (livechat)
const getContactPageVisits = (uuid) => http.get(`/api/v1/conversations/${uuid}/page-visits`)

// Available languages
const getAvailableLanguages = () => http.get('/api/v1/lang')

// OIDC enabled
const getAllEnabledOIDC = () => http.get('/api/v1/oidc/enabled')

// Help Center
const getHelpCenters = () => http.get('/api/v1/help-centers')
const getHelpCenter = (id) => http.get(`/api/v1/help-centers/${id}`)
const createHelpCenter = (data) => http.post('/api/v1/help-centers', data)
const updateHelpCenter = (id, data) => http.put(`/api/v1/help-centers/${id}`, data)
const deleteHelpCenter = (id) => http.delete(`/api/v1/help-centers/${id}`)
const getHelpCenterTree = (id, params) => http.get(`/api/v1/help-centers/${id}/tree`, { params })

// Collections
const getCollections = (helpCenterId, params) => http.get(`/api/v1/help-centers/${helpCenterId}/collections`, { params })
const getCollection = (id) => http.get(`/api/v1/help-centers/*/collections/${id}`)
const createCollection = (helpCenterId, data) => http.post(`/api/v1/help-centers/${helpCenterId}/collections`, data)
const updateCollection = (helpCenterId, id, data) => http.put(`/api/v1/help-centers/${helpCenterId}/collections/${id}`, data)
const deleteCollection = (helpCenterId, id) => http.delete(`/api/v1/help-centers/${helpCenterId}/collections/${id}`)
const toggleCollection = (id) => http.put(`/api/v1/collections/${id}/toggle`)

// Articles
const getArticles = (collectionId, params) => http.get(`/api/v1/collections/${collectionId}/articles`, { params })
const getArticle = (id) => http.get(`/api/v1/collections/*/articles/${id}`)
const createArticle = (collectionId, data) => http.post(`/api/v1/collections/${collectionId}/articles`, data)
const updateArticle = (collectionId, id, data) => http.put(`/api/v1/collections/${collectionId}/articles/${id}`, data)
const updateArticleByID = (id, data) => http.put(`/api/v1/articles/${id}`, data)
const deleteArticle = (collectionId, id) => http.delete(`/api/v1/collections/${collectionId}/articles/${id}`)
const updateArticleStatus = (id, data) => http.put(`/api/v1/articles/${id}/status`, data)

// AI Assistants
const getAIAssistants = () => http.get('/api/v1/ai-assistants')
const getAIAssistant = (id) => http.get(`/api/v1/ai-assistants/${id}`)
const createAIAssistant = (data) => http.post('/api/v1/ai-assistants', data)
const updateAIAssistant = (id, data) => http.put(`/api/v1/ai-assistants/${id}`, data)
const deleteAIAssistant = (id) => http.delete(`/api/v1/ai-assistants/${id}`)

// AI Snippets
const getAISnippets = () => http.get('/api/v1/ai-snippets')
const getAISnippet = (id) => http.get(`/api/v1/ai-snippets/${id}`)
const createAISnippet = (data) => http.post('/api/v1/ai-snippets', data)
const updateAISnippet = (id, data) => http.put(`/api/v1/ai-snippets/${id}`, data)
const deleteAISnippet = (id) => http.delete(`/api/v1/ai-snippets/${id}`)

export default {
  getOverviewCounts, getOverviewCharts, getOverviewSLA, getOverviewCSAT, getOverviewMessageVolume, getOverviewTagDistribution,
  getCustomReports, getCustomReport, createCustomReport, updateCustomReport, deleteCustomReport, executeCustomReport,
  login, logout, getLanguage, getConfig,
  getConversations, getAllConversations, getUnassignedConversations, getAssignedConversations, getMentionedConversations,
  getTeamUnassignedConversations, getViewConversations, getConversation, getConversationParticipants,
  updateUserAssignee, updateAssignee, updateTeamAssignee, removeUserAssignee, removeTeamAssignee, removeAssignee,
  updatePriority, updateConversationPriority, updateStatus, updateConversationStatus, updateLastSeen, updateAssigneeLastSeen,
  markUnread, markConversationAsUnread, updateTags, upsertTags, getConversationMessages, getConversationMessage,
  sendMessage, retryMessage, createConversation, updateConversationCustomAttributes, updateConversationCustomAttribute,
  updateContactCustomAttributes, searchConversations, searchMessages, searchContacts,
  getSettings, updateSettings, getEmailNotificationSettings, updateEmailNotificationSettings,
  getWhatsAppNotificationSettings, updateWhatsAppNotificationSettings, getTelegramNotificationSettings, updateTelegramNotificationSettings,
  getSMSNotificationSettings, updateSMSNotificationSettings, getPushNotificationSettings, updatePushNotificationSettings,
  getOfficialCommunicationsNotificationSettings, updateOfficialCommunicationsNotificationSettings, getMailDomainsSettings, updateMailDomainsSettings,
  getAgents, getAgent, getUser, createAgent, createUser, updateAgent, updateUser, deleteAgent, deleteUser,
  importAgents, getAgentImportStatus, getCurrentAgent, getCurrentUser, updateCurrentAgent, updateCurrentUser,
  getCurrentAgentTeams, updateAgentAvailability, updateCurrentUserAvailability, deleteCurrentAgentAvatar, deleteUserAvatar,
  getUsersCompact, resetPassword, setPassword, generateAPIKey, revokeAPIKey, uploadMedia,
  getContacts, getContact, getContactConversations, getContactStats, updateContact, blockContact,
  getContactNotes, createContactNote, deleteContactNote, getContactSegments, createContactSegment, getContactSegment, updateContactSegment, deleteContactSegment, updateContactCustomAttribute,
  getTeams, getTeam, createTeam, updateTeam, deleteTeam, getTeamsCompact,
  getInboxes, getInbox, createInbox, updateInbox, toggleInbox, deleteInbox, initiateOAuthFlow,
  getWebhooks, getWebhook, createWebhook, updateWebhook, deleteWebhook, toggleWebhook, testWebhook,
  getRoles, getRole, createRole, updateRole, deleteRole, getAllOIDC, getOIDC, createOIDC, updateOIDC, deleteOIDC,
  getTemplateCategories, createTemplateCategory, updateTemplateCategory, deleteTemplateCategory,
  getTemplates, getTemplate, createTemplate, updateTemplate, deleteTemplate,
  getBusinessHours, getAllBusinessHours, getBusinessHour, createBusinessHours, updateBusinessHours, deleteBusinessHour, deleteBusinessHours,
  getSLAs, getAllSLAs, getSLA, createSLA, updateSLA, deleteSLA,
  getAutomationRules, getAutomationRule, createAutomationRule, updateAutomationRule, toggleAutomationRule,
  updateAutomationRuleWeights, updateAutomationRuleExecutionMode, updateAutomationRulesExecutionMode, deleteAutomationRule,
  getAiPrompts, aiCompletion, updateAiProvider, updateAIProvider,
  getCurrentUserViews, createView, updateView, deleteView, getSharedViews, getAllSharedViews, getSharedView, createSharedView, updateSharedView, deleteSharedView,
  getTags, createTag, updateTag, deleteTag, getStatuses, createStatus, deleteStatus, getPriorities,
  getCustomAttributes, getCustomAttribute, createCustomAttribute, updateCustomAttribute, deleteCustomAttribute,
  getActivityLogs, getNotifications, getNotificationStats, markNotificationAsRead, markAllNotificationsAsRead, deleteNotification, deleteAllNotifications,
  getMacros, getAllMacros, getMacro, createMacro, updateMacro, deleteMacro, applyMacro,
  getAllDrafts, saveDraft, deleteDraft,
  getContactPageVisits,
  getAvailableLanguages,
  getAllEnabledOIDC,
  getHelpCenters, getHelpCenter, createHelpCenter, updateHelpCenter, deleteHelpCenter, getHelpCenterTree,
  getCollections, getCollection, createCollection, updateCollection, deleteCollection, toggleCollection,
  getArticles, getArticle, createArticle, updateArticle, updateArticleByID, deleteArticle, updateArticleStatus,
  getAIAssistants, getAIAssistant, createAIAssistant, updateAIAssistant, deleteAIAssistant,
  getAISnippets, getAISnippet, createAISnippet, updateAISnippet, deleteAISnippet
}
