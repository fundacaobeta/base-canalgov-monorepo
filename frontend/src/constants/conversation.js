export const CONVERSATION_LIST_TYPE = {
  ASSIGNED: 'assigned',
  UNASSIGNED: 'unassigned',
  TEAM_UNASSIGNED: 'team_unassigned',
  VIEW: 'view',
  ALL: 'all',
  MENTIONED: 'mentioned'
}

export const CONVERSATION_DEFAULT_STATUSES = {
  OPEN: 'Open',
  SNOOZED: 'Snoozed',
  RESOLVED: 'Resolved',
  CLOSED: 'Closed',
}

export const CONVERSATION_DEFAULT_STATUSES_LIST = Object.values(CONVERSATION_DEFAULT_STATUSES);

export const MACRO_CONTEXT = {
  REPLY: 'reply',
  NEW_CONVERSATION: 'new-conversation'
}

export const CONVERSATION_SORT_FIELD_MAP = {
  oldest: { model: 'conversations', field: 'last_message_at', order: 'asc' },
  newest: { model: 'conversations', field: 'last_message_at', order: 'desc' },
  started_first: { model: 'conversations', field: 'created_at', order: 'asc' },
  started_last: { model: 'conversations', field: 'created_at', order: 'desc' },
  waiting_longest: { model: 'conversations', field: 'waiting_since', order: 'asc' },
  next_sla_target: { model: 'conversations', field: 'next_sla_deadline_at', order: 'asc' },
  priority_first: { model: 'conversations', field: 'priority_id', order: 'desc' }
}

export const REFRESH_MODEL = {
  VIEW: 'view'
}