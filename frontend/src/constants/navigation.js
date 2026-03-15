export const reportsNavItems = [
  {
    titleKey: 'globals.terms.overview',
    to: { name: 'overview' },
    activeGroup: ['overview'],
    permission: 'reports:manage'
  },
  {
    titleKey: 'globals.terms.customReports',
    to: { name: 'custom-reports' },
    activeGroup: ['custom-reports'],
    permission: 'reports:manage'
  }
]

export const adminNavItems = [
  {
    titleKey: 'globals.terms.setting',
    isTitleKeyPlural: true,
    activeGroup: ['general'],
    children: [
      {
        titleKey: 'globals.terms.general',
        to: { name: 'general' },
        activeGroup: ['general'],
        permission: 'general_settings:manage'
      }
    ]
  },
  {
    titleKey: 'globals.terms.serviceDesk',
    activeGroup: ['status', 'tags', 'shared-view-list', 'new-shared-view', 'edit-shared-view', 'business-hours-list', 'new-business-hours', 'edit-business-hours', 'sla-list', 'new-sla', 'edit-sla'],
    children: [
      {
        titleKey: 'globals.terms.status',
        to: { name: 'status' },
        activeGroup: ['status'],
        permission: 'status:manage',
        isTitleKeyPlural: true
      },
      {
        titleKey: 'globals.terms.tag',
        to: { name: 'tags' },
        activeGroup: ['tags'],
        permission: 'tags:manage',
        isTitleKeyPlural: true
      },
      {
        titleKey: 'globals.terms.sharedView',
        to: { name: 'shared-view-list' },
        activeGroup: ['shared-view-list', 'new-shared-view', 'edit-shared-view'],
        permission: 'shared_views:manage',
        isTitleKeyPlural: true
      },
      {
        titleKey: 'globals.terms.businessHour',
        to: { name: 'business-hours-list' },
        activeGroup: ['business-hours-list', 'new-business-hours', 'edit-business-hours'],
        permission: 'business_hours:manage',
        isTitleKeyPlural: true
      },
      {
        titleKey: 'globals.terms.slaPolicy',
        to: { name: 'sla-list' },
        activeGroup: ['sla-list', 'new-sla', 'edit-sla'],
        permission: 'sla:manage',
        isTitleKeyPlural: true
      }
    ]
  },
  {
    titleKey: 'globals.terms.inbox',
    isTitleKeyPlural: true,
    activeGroup: ['inbox-list', 'new-inbox', 'edit-inbox', 'domains'],
    children: [
      {
        titleKey: 'globals.terms.inbox',
        to: { name: 'inbox-list' },
        activeGroup: ['inbox-list', 'new-inbox', 'edit-inbox'],
        permission: 'inboxes:manage',
        isTitleKeyPlural: true
      },
      {
        titleKey: 'globals.terms.domain',
        to: { name: 'domains' },
        activeGroup: ['domains'],
        permission: 'general_settings:manage',
        isTitleKeyPlural: true
      }
    ]
  },
  {
    titleKey: 'globals.terms.contact',
    isTitleKeyPlural: true,
    activeGroup: ['contact-segments'],
    children: [
      {
        titleKey: 'globals.terms.group',
        to: { name: 'contact-segments' },
        activeGroup: ['contact-segments'],
        permission: 'general_settings:manage',
        isTitleKeyPlural: true
      }
    ]
  },
  {
    titleKey: 'globals.terms.team',
    isTitleKeyPlural: true,
    activeGroup: ['agent-list', 'new-agent', 'edit-agent', 'team-list', 'new-team', 'edit-team', 'role-list', 'new-role', 'edit-role'],
    children: [
      {
        titleKey: 'globals.terms.agent',
        to: { name: 'agent-list' },
        activeGroup: ['agent-list', 'new-agent', 'edit-agent'],
        permission: 'users:manage',
        isTitleKeyPlural: true
      },
      {
        titleKey: 'globals.terms.team',
        to: { name: 'team-list' },
        activeGroup: ['team-list', 'new-team', 'edit-team'],
        permission: 'teams:manage',
        isTitleKeyPlural: true
      },
      {
        titleKey: 'globals.terms.role',
        to: { name: 'role-list' },
        activeGroup: ['role-list', 'new-role', 'edit-role'],
        permission: 'roles:manage',
        isTitleKeyPlural: true
      }
    ]
  },
  {
    titleKey: 'globals.terms.automation',
    isTitleKeyPlural: true,
    activeGroup: ['automations', 'new-automation', 'edit-automation', 'macro-list', 'new-macro', 'edit-macro', 'custom-attributes'],
    children: [
      {
        titleKey: 'globals.terms.automation',
        to: { name: 'automations' },
        activeGroup: ['automations', 'new-automation', 'edit-automation'],
        permission: 'automations:manage',
        isTitleKeyPlural: true
      },
      {
        titleKey: 'globals.terms.macro',
        to: { name: 'macro-list' },
        activeGroup: ['macro-list', 'new-macro', 'edit-macro'],
        permission: 'macros:manage',
        isTitleKeyPlural: true
      },
      {
        titleKey: 'globals.terms.customAttribute',
        to: { name: 'custom-attributes' },
        activeGroup: ['custom-attributes'],
        permission: 'custom_attributes:manage',
        isTitleKeyPlural: true
      }
    ]
  },
  {
    titleKey: 'globals.terms.notification',
    activeGroup: ['notification-email', 'notification-whatsapp', 'notification-telegram', 'notification-sms', 'notification-push', 'notification-letter-notice-summons'],
    children: [
      {
        titleKey: 'globals.terms.email',
        to: { name: 'notification-email' },
        activeGroup: ['notification-email'],
        permission: 'notification_settings:manage'
      },
      {
        titleKey: 'navigation.integrations.whatsapp',
        to: { name: 'notification-whatsapp' },
        activeGroup: ['notification-whatsapp']
      },
      {
        titleKey: 'navigation.integrations.telegram',
        to: { name: 'notification-telegram' },
        activeGroup: ['notification-telegram']
      },
      {
        titleKey: 'navigation.integrations.sms',
        to: { name: 'notification-sms' },
        activeGroup: ['notification-sms']
      },
      {
        titleKey: 'navigation.integrations.push',
        to: { name: 'notification-push' },
        activeGroup: ['notification-push']
      },
      {
        titleKey: 'navigation.notifications.officialCommunications',
        to: { name: 'notification-letter-notice-summons' },
        activeGroup: ['notification-letter-notice-summons']
      }
    ]
  },
  {
    titleKey: 'globals.terms.integration',
    isTitleKeyPlural: true,
    activeGroup: ['integrations-actions', 'webhook-list', 'new-webhook', 'edit-webhook', 'templates', 'new-template', 'edit-template'],
    children: [
      {
        titleKey: 'navigation.integrations.actions',
        to: { name: 'integrations-actions' },
        activeGroup: ['integrations-actions']
      },
      {
        titleKey: 'globals.terms.webhook',
        to: { name: 'webhook-list' },
        activeGroup: ['webhook-list', 'new-webhook', 'edit-webhook'],
        permission: 'webhooks:manage',
        isTitleKeyPlural: true
      },
      {
        titleKey: 'globals.terms.template',
        to: { name: 'templates' },
        activeGroup: ['templates', 'new-template', 'edit-template'],
        permission: 'templates:manage',
        isTitleKeyPlural: true
      }
    ]
  },
  {
    titleKey: 'globals.terms.security',
    activeGroup: ['activity-log', 'sso-list', 'new-sso', 'edit-sso'],
    children: [
      {
        titleKey: 'globals.terms.activityLog',
        to: { name: 'activity-log' },
        activeGroup: ['activity-log'],
        permission: 'activity_logs:manage',
        isTitleKeyPlural: true
      },
      {
        titleKey: 'globals.terms.sso',
        to: { name: 'sso-list' },
        activeGroup: ['sso-list', 'new-sso', 'edit-sso'],
        permission: 'oidc:manage'
      }
    ]
  }
]

export const accountNavItems = [
  {
    titleKey: 'globals.terms.profile',
    to: { name: 'profile' },
    activeGroup: ['profile']
  },
  {
    titleKey: 'account.appPreferences',
    to: { name: 'account-preferences' },
    activeGroup: ['account-preferences']
  }
]

export const contactNavItems = [
  {
    titleKey: 'globals.terms.contact',
    to: { name: 'contacts' },
    activeGroup: ['contacts', 'contact-detail'],
    isTitleKeyPlural: true
  }
]
