import { createRouter, createWebHistory } from 'vue-router'
import App from '@/App.vue'
import OuterApp from '@/OuterApp.vue'
import InboxLayout from '@/layouts/inbox/InboxLayout.vue'
import AccountLayout from '@/layouts/account/AccountLayout.vue'
import AdminLayout from '@/layouts/admin/AdminLayout.vue'
import { useAppSettingsStore } from '@/stores/appSettings'
import i18n from '@/plugins/i18n'
import { getLocalePaths } from './paths'

export function createAppRouter(locale = 'pt-BR') {
  const p = getLocalePaths(locale)

  // Maps localized URL type back to internal API value
  const internalType = Object.fromEntries(
    Object.entries(p.inboxTypes).map(([internal, localized]) => [localized, internal])
  )

  const inboxTypeKey = (route) => {
    const type = internalType[route.params.type] ?? route.params.type
    if (type === 'assigned') return 'route.myInbox'
    if (type === 'mentioned') return 'route.mentions'
    if (type === 'unassigned') return 'globals.terms.unassigned'
    return 'route.allConversations'
  }

  // Props transformer: component always receives the internal English type value
  const inboxProps = (route) => ({
    type: internalType[route.params.type] ?? route.params.type
  })

  const routes = [
    {
      path: '/',
      component: OuterApp,
      children: [
        {
          path: '',
          name: 'login',
          component: () => import('@/views/auth/UserLoginView.vue'),
          meta: { title: 'route.login' }
        },
        {
          path: p.resetPassword,
          name: 'reset-password',
          component: () => import('@/views/auth/ResetPasswordView.vue'),
          meta: { title: 'route.resetPassword' }
        },
        {
          path: p.setPassword,
          name: 'set-password',
          component: () => import('@/views/auth/SetPasswordView.vue'),
          meta: { title: 'route.setPassword' }
        }
      ]
    },
    {
      path: '/',
      component: App,
      children: [
        {
          path: p.contacts,
          name: 'contacts',
          component: () => import('@/views/contact/ContactsView.vue'),
          meta: { title: 'route.contacts', area: 'contacts' }
        },
        {
          path: `${p.contacts}/:id`,
          name: 'contact-detail',
          component: () => import('@/views/contact/ContactDetailView.vue'),
          meta: { title: 'route.contact', area: 'contacts' }
        },
        {
          path: `/${p.reports}`,
          name: 'reports',
          redirect: { name: 'overview' },
          children: [
            {
              path: p.reportOverview,
              name: 'overview',
              component: () => import('@/views/reports/OverviewView.vue'),
              meta: { title: 'route.overview', area: 'reports' }
            }
          ]
        },
        {
          path: `/${p.inboxes}/${p.inboxTeams}/:teamID`,
          name: 'teams',
          props: true,
          component: InboxLayout,
          meta: { title: 'route.teamInbox', hidePageHeader: true, area: 'inboxes' },
          children: [
            {
              path: '',
              name: 'team-inbox',
              component: () => import('@/views/inbox/InboxView.vue'),
              meta: { title: 'route.teamInbox', area: 'inboxes' },
              children: [
                {
                  path: `${p.inboxConversation}/:uuid`,
                  name: 'team-inbox-conversation',
                  component: () => import('@/views/conversation/ConversationDetailView.vue'),
                  props: true,
                  meta: { title: 'route.teamInbox', hidePageHeader: true, area: 'inboxes' }
                }
              ]
            }
          ]
        },
        {
          path: `/${p.inboxes}/${p.inboxViews}/:viewID`,
          name: 'views',
          props: true,
          component: InboxLayout,
          meta: { title: 'route.viewInbox', hidePageHeader: true, area: 'inboxes' },
          children: [
            {
              path: '',
              name: 'view-inbox',
              component: () => import('@/views/inbox/InboxView.vue'),
              meta: { title: 'route.viewInbox', area: 'inboxes' },
              children: [
                {
                  path: `${p.inboxConversation}/:uuid`,
                  name: 'view-inbox-conversation',
                  component: () => import('@/views/conversation/ConversationDetailView.vue'),
                  props: true,
                  meta: { title: 'route.viewInbox', hidePageHeader: true, area: 'inboxes' }
                }
              ]
            }
          ]
        },
        {
          path: `${p.inboxes}/${p.inboxSearch}`,
          name: 'search',
          component: () => import('@/views/search/SearchView.vue'),
          meta: { title: 'route.search', hidePageHeader: true, area: 'inboxes' }
        },
        {
          path: `/${p.inboxes}`,
          name: 'inboxes',
          redirect: { name: 'inbox', params: { type: p.inboxTypes.assigned } }
        },
        {
          path: `/${p.inboxes}/:type(${p.inboxTypes.assigned}|${p.inboxTypes.unassigned}|${p.inboxTypes.all}|${p.inboxTypes.mentioned})`,
          name: 'inbox',
          component: InboxLayout,
          props: inboxProps,
          meta: { title: 'route.inbox', hidePageHeader: true, area: 'inboxes' },
          children: [
            {
              path: '',
              name: 'inbox-view',
              component: () => import('@/views/inbox/InboxView.vue'),
              meta: { title: 'route.inbox', type: inboxTypeKey, area: 'inboxes' },
              children: [
                {
                  path: `${p.inboxConversation}/:uuid`,
                  name: 'inbox-conversation',
                  component: () => import('@/views/conversation/ConversationDetailView.vue'),
                  props: true,
                  meta: {
                    title: 'route.inbox',
                    type: inboxTypeKey,
                    hidePageHeader: true,
                    area: 'inboxes'
                  }
                }
              ]
            }
          ]
        },
        {
          path: `/${p.account}/:page?`,
          name: 'account',
          redirect: { name: 'profile' },
          component: AccountLayout,
          props: true,
          meta: { title: 'route.account', area: 'account' },
          children: [
            {
              path: p.accountProfile,
              name: 'profile',
              component: () => import('@/views/account/profile/ProfileEditView.vue'),
              meta: { title: 'route.editProfile', area: 'account' }
            },
            {
              path: p.accountPreferences,
              name: 'account-preferences',
              component: () => import('@/views/account/preferences/AccountPreferencesView.vue'),
              meta: { title: 'route.appPreferences', area: 'account' }
            }
          ]
        },
        {
          path: `/${p.admin}`,
          name: 'admin',
          redirect: { name: 'general' },
          component: AdminLayout,
          meta: { title: 'route.admin', area: 'admin' },
          children: [
            {
              path: p.adminCustomAttributes,
              name: 'custom-attributes',
              component: () => import('@/views/admin/custom-attributes/CustomAttributes.vue'),
              meta: { title: 'route.customAttributes', area: 'admin' }
            },
            {
              path: p.adminGeneral,
              name: 'general',
              component: () => import('@/views/admin/general/General.vue'),
              meta: { title: 'route.general', area: 'admin' }
            },
            {
              path: p.adminBusinessHours,
              component: () => import('@/views/admin/business-hours/BusinessHours.vue'),
              meta: { title: 'route.businessHours', area: 'admin' },
              children: [
                {
                  path: '',
                  name: 'business-hours-list',
                  component: () => import('@/views/admin/business-hours/BusinessHoursList.vue')
                },
                {
                  path: p.newFem,
                  name: 'new-business-hours',
                  component: () =>
                    import('@/views/admin/business-hours/CreateOrEditBusinessHours.vue'),
                  meta: { title: 'route.newBusinessHours', area: 'admin' }
                },
                {
                  path: `:id/${p.edit}`,
                  name: 'edit-business-hours',
                  props: true,
                  component: () =>
                    import('@/views/admin/business-hours/CreateOrEditBusinessHours.vue'),
                  meta: { title: 'route.editBusinessHours', area: 'admin' }
                }
              ]
            },
            {
              path: p.adminSLA,
              component: () => import('@/views/admin/sla/SLA.vue'),
              meta: { title: 'route.sla', area: 'admin' },
              children: [
                {
                  path: '',
                  name: 'sla-list',
                  component: () => import('@/views/admin/sla/SLAList.vue')
                },
                {
                  path: p.new,
                  name: 'new-sla',
                  component: () => import('@/views/admin/sla/CreateEditSLA.vue'),
                  meta: { title: 'route.newSLA', area: 'admin' }
                },
                {
                  path: `:id/${p.edit}`,
                  props: true,
                  name: 'edit-sla',
                  component: () => import('@/views/admin/sla/CreateEditSLA.vue'),
                  meta: { title: 'route.editSLA', area: 'admin' }
                }
              ]
            },
            {
              path: p.adminInboxes,
              component: () => import('@/views/admin/inbox/InboxView.vue'),
              meta: { title: 'route.inboxes', area: 'admin' },
              children: [
                {
                  path: '',
                  name: 'inbox-list',
                  component: () => import('@/views/admin/inbox/InboxList.vue')
                },
                {
                  path: p.newFem,
                  name: 'new-inbox',
                  component: () => import('@/views/admin/inbox/NewInbox.vue'),
                  meta: { title: 'route.newInbox', area: 'admin' }
                },
                {
                  path: `:id/${p.edit}`,
                  props: true,
                  name: 'edit-inbox',
                  component: () => import('@/views/admin/inbox/EditInbox.vue'),
                  meta: { title: 'route.editInbox', area: 'admin' }
                }
              ]
            },
            {
              path: p.adminDomains,
              name: 'domains',
              component: () => import('@/views/admin/domains/DomainsView.vue'),
              meta: { title: 'route.domains', area: 'admin' }
            },
            {
              path: p.adminCustomReports,
              name: 'custom-reports',
              component: () => import('@/views/admin/reports/CustomReportsView.vue'),
              meta: { title: 'route.customReports', area: 'admin' }
            },
            {
              path: p.adminContactSegments,
              name: 'contact-segments',
              component: () => import('@/views/admin/contacts/ContactSegmentsView.vue'),
              meta: { title: 'route.contactSegments', area: 'admin' }
            },
            {
              path: p.adminNotifications,
              name: 'notification-email',
              component: () => import('@/features/admin/notification/NotificationSetting.vue'),
              meta: { title: 'route.notificationEmail', area: 'admin' }
            },
            {
              path: `${p.adminNotifications}/whatsapp`,
              name: 'notification-whatsapp',
              component: () =>
                import('@/views/admin/notification/NotificationChannelConfigView.vue'),
              props: { channel: 'whatsapp' },
              meta: { title: 'globals.terms.whatsapp', area: 'admin' }
            },
            {
              path: `${p.adminNotifications}/telegram`,
              name: 'notification-telegram',
              component: () =>
                import('@/views/admin/notification/NotificationChannelConfigView.vue'),
              props: { channel: 'telegram' },
              meta: { title: 'globals.terms.telegram', area: 'admin' }
            },
            {
              path: `${p.adminNotifications}/sms`,
              name: 'notification-sms',
              component: () =>
                import('@/views/admin/notification/NotificationChannelConfigView.vue'),
              props: { channel: 'sms' },
              meta: { title: 'globals.terms.sms', area: 'admin' }
            },
            {
              path: `${p.adminNotifications}/push`,
              name: 'notification-push',
              component: () =>
                import('@/views/admin/notification/NotificationChannelConfigView.vue'),
              props: { channel: 'push' },
              meta: { title: 'globals.terms.pushNotification', area: 'admin' }
            },
            {
              path: `${p.adminNotifications}/${p.adminNotifOfficialComms}`,
              name: 'notification-letter-notice-summons',
              component: () =>
                import(
                  '@/views/admin/notification/NotificationOfficialCommunicationsView.vue'
                ),
              meta: { title: 'route.notificationOfficialCommunications', area: 'admin' }
            },
            {
              path: p.adminPeople,
              meta: { title: 'route.teams', area: 'admin', section: 'people' },
              children: [
                {
                  path: p.adminAgents,
                  component: () => import('@/views/admin/agents/Agents.vue'),
                  meta: { title: 'route.agents', area: 'admin', section: 'people' },
                  children: [
                    {
                      path: '',
                      name: 'agent-list',
                      component: () => import('@/views/admin/agents/AgentList.vue')
                    },
                    {
                      path: p.new,
                      name: 'new-agent',
                      component: () => import('@/views/admin/agents/CreateAgent.vue'),
                      meta: { title: 'route.newAgent', area: 'admin', section: 'people' }
                    },
                    {
                      path: `:id/${p.edit}`,
                      props: true,
                      name: 'edit-agent',
                      component: () => import('@/views/admin/agents/EditAgent.vue'),
                      meta: { title: 'route.editAgent', area: 'admin', section: 'people' }
                    }
                  ]
                },
                {
                  path: p.adminTeams,
                  component: () => import('@/views/admin/teams/Teams.vue'),
                  meta: { title: 'route.teams', area: 'admin', section: 'people' },
                  children: [
                    {
                      path: '',
                      name: 'team-list',
                      component: () => import('@/views/admin/teams/TeamList.vue')
                    },
                    {
                      path: p.newFem,
                      name: 'new-team',
                      component: () => import('@/views/admin/teams/CreateTeamForm.vue'),
                      meta: { title: 'route.newTeam', area: 'admin', section: 'people' }
                    },
                    {
                      path: `:id/${p.edit}`,
                      props: true,
                      name: 'edit-team',
                      component: () => import('@/views/admin/teams/EditTeamForm.vue'),
                      meta: { title: 'route.editTeam', area: 'admin', section: 'people' }
                    }
                  ]
                },
                {
                  path: p.adminRoles,
                  component: () => import('@/views/admin/roles/Roles.vue'),
                  meta: { title: 'route.roles', area: 'admin', section: 'people' },
                  children: [
                    {
                      path: '',
                      name: 'role-list',
                      component: () => import('@/views/admin/roles/RoleList.vue')
                    },
                    {
                      path: p.newFem,
                      name: 'new-role',
                      component: () => import('@/views/admin/roles/NewRole.vue'),
                      meta: { title: 'route.newRole', area: 'admin', section: 'people' }
                    },
                    {
                      path: `:id/${p.edit}`,
                      props: true,
                      name: 'edit-role',
                      component: () => import('@/views/admin/roles/EditRole.vue'),
                      meta: { title: 'route.editRole', area: 'admin', section: 'people' }
                    }
                  ]
                },
                {
                  path: p.adminActivityLog,
                  name: 'activity-log',
                  component: () => import('@/views/admin/activity-log/ActivityLog.vue'),
                  meta: { title: 'route.activityLog', area: 'admin', section: 'people' }
                }
              ]
            },
            {
              path: p.adminAutomations,
              component: () => import('@/views/admin/automations/Automation.vue'),
              name: 'automations',
              meta: { title: 'route.automations', area: 'admin' },
              children: [
                {
                  path: p.newFem,
                  props: true,
                  name: 'new-automation',
                  component: () => import('@/views/admin/automations/CreateOrEditRule.vue'),
                  meta: { title: 'route.newAutomation', area: 'admin' }
                },
                {
                  path: `:id/${p.edit}`,
                  props: true,
                  name: 'edit-automation',
                  component: () => import('@/views/admin/automations/CreateOrEditRule.vue'),
                  meta: { title: 'route.editAutomation', area: 'admin' }
                }
              ]
            },
            {
              path: p.adminTemplates,
              component: () => import('@/views/admin/templates/Templates.vue'),
              name: 'templates',
              meta: { title: 'route.templates', area: 'admin' },
              children: [
                {
                  path: `:id/${p.edit}`,
                  name: 'edit-template',
                  props: true,
                  component: () => import('@/views/admin/templates/CreateEditTemplate.vue'),
                  meta: { title: 'route.editTemplate', area: 'admin' }
                },
                {
                  path: p.new,
                  name: 'new-template',
                  props: true,
                  component: () => import('@/views/admin/templates/CreateEditTemplate.vue'),
                  meta: { title: 'route.newTemplate', area: 'admin' }
                }
              ]
            },
            {
              path: p.adminSSO,
              component: () => import('@/views/admin/oidc/OIDC.vue'),
              name: 'sso',
              meta: { title: 'route.sso', area: 'admin' },
              children: [
                {
                  path: '',
                  name: 'sso-list',
                  component: () => import('@/views/admin/oidc/OIDCList.vue')
                },
                {
                  path: `:id/${p.edit}`,
                  props: true,
                  name: 'edit-sso',
                  component: () => import('@/views/admin/oidc/CreateEditOIDC.vue'),
                  meta: { title: 'route.editSSO', area: 'admin' }
                },
                {
                  path: p.new,
                  name: 'new-sso',
                  component: () => import('@/views/admin/oidc/CreateEditOIDC.vue'),
                  meta: { title: 'route.newSSO', area: 'admin' }
                }
              ]
            },
            {
              path: p.adminWebhooks,
              component: () => import('@/views/admin/webhooks/Webhooks.vue'),
              name: 'webhooks',
              meta: { title: 'route.webhooks', area: 'admin' },
              children: [
                {
                  path: '',
                  name: 'webhook-list',
                  component: () => import('@/views/admin/webhooks/WebhookList.vue')
                },
                {
                  path: `:id/${p.edit}`,
                  props: true,
                  name: 'edit-webhook',
                  component: () => import('@/views/admin/webhooks/CreateEditWebhook.vue'),
                  meta: { title: 'route.editWebhook', area: 'admin' }
                },
                {
                  path: p.new,
                  name: 'new-webhook',
                  component: () => import('@/views/admin/webhooks/CreateEditWebhook.vue'),
                  meta: { title: 'route.newWebhook', area: 'admin' }
                }
              ]
            },
            {
              path: p.adminIntegrationActions,
              name: 'integrations-actions',
              component: () =>
                import('@/views/admin/integrations/IntegrationActionsView.vue'),
              meta: { title: 'route.integrationActions', area: 'admin' }
            },
            {
              path: p.adminConversations,
              meta: { title: 'route.conversations', area: 'admin' },
              children: [
                {
                  path: p.adminTags,
                  name: 'tags',
                  component: () => import('@/views/admin/tags/TagsView.vue'),
                  meta: { title: 'route.tags', area: 'admin' }
                },
                {
                  path: p.adminStatus,
                  name: 'status',
                  component: () => import('@/views/admin/status/StatusView.vue'),
                  meta: { title: 'route.status', area: 'admin' }
                },
                {
                  path: p.adminMacros,
                  component: () => import('@/views/admin/macros/Macros.vue'),
                  meta: { title: 'route.macros', area: 'admin' },
                  children: [
                    {
                      path: '',
                      name: 'macro-list',
                      component: () => import('@/views/admin/macros/MacroList.vue')
                    },
                    {
                      path: p.newFem,
                      name: 'new-macro',
                      component: () => import('@/views/admin/macros/CreateMacro.vue'),
                      meta: { title: 'route.newMacro', area: 'admin' }
                    },
                    {
                      path: `:id/${p.edit}`,
                      props: true,
                      name: 'edit-macro',
                      component: () => import('@/views/admin/macros/EditMacro.vue'),
                      meta: { title: 'route.editMacro', area: 'admin' }
                    }
                  ]
                },
                {
                  path: p.adminSharedViews,
                  component: () => import('@/views/admin/shared-views/SharedViews.vue'),
                  meta: { title: 'route.sharedViews', area: 'admin' },
                  children: [
                    {
                      path: '',
                      name: 'shared-view-list',
                      component: () => import('@/views/admin/shared-views/SharedViewList.vue')
                    },
                    {
                      path: p.newFem,
                      name: 'new-shared-view',
                      component: () =>
                        import('@/views/admin/shared-views/CreateSharedView.vue'),
                      meta: { title: 'route.newSharedView', area: 'admin' }
                    },
                    {
                      path: `:id/${p.edit}`,
                      props: true,
                      name: 'edit-shared-view',
                      component: () =>
                        import('@/views/admin/shared-views/EditSharedView.vue'),
                      meta: { title: 'route.editSharedView', area: 'admin' }
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: { name: 'inbox', params: { type: p.inboxTypes.assigned } }
    }
  ]

  const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes
  })

  router.beforeEach((to, from, next) => {
    const appSettingsStore = useAppSettingsStore()
    const siteName = appSettingsStore.settings?.['app.site_name'] || 'CanalGov'
    const titleKey = to.meta?.title || ''
    const pageTitle = titleKey ? i18n.global.t(titleKey) : ''
    document.title = `${pageTitle} - ${siteName}`
    next()
  })

  return router
}
