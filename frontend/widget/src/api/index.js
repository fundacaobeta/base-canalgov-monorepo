import axios from 'axios'
import { parseJWT } from '@shared-ui/utils/string'

const VISITOR_JWT_KEY = 'canalgov_visitor_jwt'

function getInboxIDFromQuery () {
    const params = new URLSearchParams(window.location.search)
    const inboxId = params.get('inbox_id')
    return inboxId ? parseInt(inboxId, 10) : null
}

export function setVisitorJWT (jwt) {
    localStorage.setItem(VISITOR_JWT_KEY, jwt)
}

export function clearVisitorJWT () {
    localStorage.removeItem(VISITOR_JWT_KEY)
}

export function getVisitorJWT () {
    return localStorage.getItem(VISITOR_JWT_KEY)
}

// Returns visitor JWT if current user is authenticated (for merge).
function getVisitorJWTForMerge (sessionToken) {
    const visitorJWT = getVisitorJWT()
    if (!visitorJWT || !sessionToken) {
        return null
    }
    try {
        const claims = parseJWT(sessionToken)
        if (claims && !claims.is_visitor && claims.external_user_id) {
            return visitorJWT
        }
    } catch {
        // Ignore JWT parse errors
    }
    return null
}

const http = axios.create({
    timeout: 10000,
    responseType: 'json'
})

// Set content type and authentication headers
http.interceptors.request.use((request) => {
    if ((request.method === 'post' || request.method === 'put') && !request.headers['Content-Type']) {
        request.headers['Content-Type'] = 'application/json'
    }

    // Add authentication headers for widget API endpoints
    if (request.url && request.url.includes('/api/v1/widget/')) {
        const session = localStorage.getItem('canalgov_session')
        const inboxId = getInboxIDFromQuery()

        // Add JWT to Authorization header
        if (session) {
            request.headers['Authorization'] = `Bearer ${session}`
        }

        // Add inbox ID to custom header
        if (inboxId) {
            request.headers['X-Inbox-ID'] = inboxId.toString()
        }

        const visitorJWTForMerge = getVisitorJWTForMerge(session)
        if (visitorJWTForMerge) {
            request.headers['X-Visitor-JWT'] = visitorJWTForMerge
        }
    }

    return request
})

http.interceptors.response.use((response) => {
    if (response.headers['x-clear-visitor']) {
        clearVisitorJWT()
    }
    return response
})

const getWidgetSettings = (inboxID) => http.get('/api/v1/widget/chat/settings', {
    params: { inbox_id: inboxID }
})
const getLanguage = (lang) => http.get(`/api/v1/lang/${lang}`)
const getAvailableLanguages = () => http.get('/api/v1/lang')
const initChatConversation = (data) => http.post('/api/v1/widget/chat/conversations/init', data)
const getChatConversations = () => http.get('/api/v1/widget/chat/conversations')
const getChatConversation = (uuid) => http.get(`/api/v1/widget/chat/conversations/${uuid}`)
const sendChatMessage = (uuid, data) => http.post(`/api/v1/widget/chat/conversations/${uuid}/message`, data)
const closeChatConversation = (uuid) => http.post(`/api/v1/widget/chat/conversations/${uuid}/close`)
const uploadMedia = (conversationUUID, files) => {
    const formData = new FormData()

    formData.append('conversation_uuid', conversationUUID)

    for (let i = 0; i < files.length; i++) {
        formData.append('files', files[i])
    }

    const session = localStorage.getItem('canalgov_session')
    const inboxId = getInboxIDFromQuery()

    const headers = {
        'Content-Type': 'multipart/form-data'
    }

    if (session) {
        headers['Authorization'] = `Bearer ${session}`
    }
    if (inboxId) {
        headers['X-Inbox-ID'] = inboxId.toString()
    }

    const visitorJWTForMerge = getVisitorJWTForMerge(session)
    if (visitorJWTForMerge) {
        headers['X-Visitor-JWT'] = visitorJWTForMerge
    }

    return axios.post('/api/v1/widget/media/upload', formData, {
        headers,
        timeout: 30000
    }).then((response) => {
        if (response.headers['x-clear-visitor']) {
            clearVisitorJWT()
        }
        return response
    })
}
const updateConversationLastSeen = (uuid) => http.post(`/api/v1/widget/chat/conversations/${uuid}/update-last-seen`)
const submitCSATResponse = (csatUuid, rating, feedback) =>
    http.post(`/api/v1/csat/${csatUuid}/response`, {
        rating,
        feedback,
    })

export default {
    getWidgetSettings,
    getLanguage,
    getAvailableLanguages,
    initChatConversation,
    getChatConversations,
    getChatConversation,
    sendChatMessage,
    closeChatConversation,
    uploadMedia,
    updateConversationLastSeen,
    submitCSATResponse
}
