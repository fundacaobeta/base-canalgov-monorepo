// Path segments per locale.
// These are used by the router factory to build locale-aware URL paths.
// Route NAMES never change — only these path strings do.

export const localePaths = {
  'pt-BR': {
    // Auth
    resetPassword: 'recuperar-senha',
    setPassword: 'definir-senha',
    // Top-level sections
    contacts: 'contatos',
    reports: 'relatorios',
    inboxes: 'caixas',
    account: 'conta',
    admin: 'admin',
    // Reports
    reportOverview: 'visao-geral',
    // Inbox
    inboxTeams: 'equipes',
    inboxViews: 'visoes',
    inboxSearch: 'busca',
    inboxConversation: 'conversa',
    // Account
    accountProfile: 'perfil',
    accountPreferences: 'preferencias',
    // Admin
    adminGeneral: 'geral',
    adminCustomAttributes: 'atributos-personalizados',
    adminBusinessHours: 'horarios-comerciais',
    adminSLA: 'sla',
    adminInboxes: 'caixas-de-entrada',
    adminDomains: 'dominios',
    adminCustomReports: 'relatorios/personalizados',
    adminContactSegments: 'contatos/grupos',
    adminNotifications: 'notificacoes',
    adminNotifOfficialComms: 'comunicacoes-oficiais',
    adminPeople: 'times',
    adminAgents: 'agentes',
    adminTeams: 'equipes',
    adminRoles: 'funcoes',
    adminActivityLog: 'log-de-atividades',
    adminAutomations: 'automacoes',
    adminTemplates: 'modelos',
    adminSSO: 'sso',
    adminWebhooks: 'webhooks',
    adminIntegrationActions: 'integracoes/acoes',
    adminConversations: 'conversas',
    adminTags: 'tags',
    adminStatus: 'status',
    adminMacros: 'macros',
    adminSharedViews: 'visoes-compartilhadas',
    // Inbox type param values (used in URL)
    inboxTypes: {
      assigned: 'atribuidas',
      unassigned: 'nao-atribuidas',
      all: 'todas',
      mentioned: 'mencoes',
    },
    // CRUD suffixes
    new: 'novo',
    newFem: 'nova',
    edit: 'editar',
  },
  en: {
    // Auth
    resetPassword: 'reset-password',
    setPassword: 'set-password',
    // Top-level sections
    contacts: 'contacts',
    reports: 'reports',
    inboxes: 'inboxes',
    account: 'account',
    admin: 'admin',
    // Reports
    reportOverview: 'overview',
    // Inbox
    inboxTeams: 'teams',
    inboxViews: 'views',
    inboxSearch: 'search',
    inboxConversation: 'conversation',
    // Account
    accountProfile: 'profile',
    accountPreferences: 'preferences',
    // Admin
    adminGeneral: 'general',
    adminCustomAttributes: 'custom-attributes',
    adminBusinessHours: 'business-hours',
    adminSLA: 'sla',
    adminInboxes: 'inboxes',
    adminDomains: 'domains',
    adminCustomReports: 'reports/custom',
    adminContactSegments: 'contacts/groups',
    adminNotifications: 'notifications',
    adminNotifOfficialComms: 'official-communications',
    adminPeople: 'people',
    adminAgents: 'agents',
    adminTeams: 'teams',
    adminRoles: 'roles',
    adminActivityLog: 'activity-log',
    adminAutomations: 'automations',
    adminTemplates: 'templates',
    adminSSO: 'sso',
    adminWebhooks: 'webhooks',
    adminIntegrationActions: 'integrations/actions',
    adminConversations: 'conversations',
    adminTags: 'tags',
    adminStatus: 'statuses',
    adminMacros: 'macros',
    adminSharedViews: 'shared-views',
    // Inbox type param values (used in URL)
    inboxTypes: {
      assigned: 'assigned',
      unassigned: 'unassigned',
      all: 'all',
      mentioned: 'mentioned',
    },
    // CRUD suffixes
    new: 'new',
    newFem: 'new',
    edit: 'edit',
  },
}

export function getLocalePaths(locale) {
  return localePaths[locale] ?? localePaths['pt-BR']
}
