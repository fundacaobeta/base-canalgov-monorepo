import { createRouter, createWebHistory } from 'vue-router'
import App from '@/App.vue'
import OuterApp from '@/OuterApp.vue'
import InboxLayout from '@/layouts/inbox/InboxLayout.vue'
import AccountLayout from '@/layouts/account/AccountLayout.vue'
import AdminLayout from '@/layouts/admin/AdminLayout.vue'
import { useAppSettingsStore } from '@/stores/appSettings'

const routes = [
  {
    path: '/',
    component: OuterApp,
    children: [
      {
        path: '',
        name: 'login',
        component: () => import('@/views/auth/UserLoginView.vue'),
        meta: { title: 'Login' }
      },
      {
        path: 'reset-password',
        name: 'reset-password',
        component: () => import('@/views/auth/ResetPasswordView.vue'),
        meta: { title: 'Redefinir senha' }
      },
      {
        path: 'set-password',
        name: 'set-password',
        component: () => import('@/views/auth/SetPasswordView.vue'),
        meta: { title: 'Definir senha' }
      }
    ]
  },
  {
    path: '/',
    component: App,
    children: [
      {
        path: 'contacts',
        name: 'contacts',
        component: () => import('@/views/contact/ContactsView.vue'),
        meta: { title: 'Todos os contatos' }
      },
      {
        path: 'contacts/:id',
        name: 'contact-detail',
        component: () => import('@/views/contact/ContactDetailView.vue'),
        meta: { title: 'Contatos' }
      },
      {
        path: '/reports',
        name: 'reports',
        redirect: '/reports/overview',
        children: [
          {
            path: 'overview',
            name: 'overview',
            component: () => import('@/views/reports/OverviewView.vue'),
            meta: { title: 'Visão geral' }
          }
        ]
      },
      {
        path: '/inboxes/teams/:teamID',
        name: 'teams',
        props: true,
        component: InboxLayout,
        meta: { title: 'Caixa da equipe', hidePageHeader: true },
        children: [
          {
            path: '',
            name: 'team-inbox',
            component: () => import('@/views/inbox/InboxView.vue'),
            meta: { title: 'Caixa da equipe' },
            children: [
              {
                path: 'conversation/:uuid',
                name: 'team-inbox-conversation',
                component: () => import('@/views/conversation/ConversationDetailView.vue'),
                props: true,
                meta: { title: 'Caixa da equipe', hidePageHeader: true }
              }
            ]
          }
        ]
      },
      {
        path: '/inboxes/views/:viewID',
        name: 'views',
        props: true,
        component: InboxLayout,
        meta: { title: 'Caixa da visão', hidePageHeader: true },
        children: [
          {
            path: '',
            name: 'view-inbox',
            component: () => import('@/views/inbox/InboxView.vue'),
            meta: { title: 'Caixa da visão' },
            children: [
              {
                path: 'conversation/:uuid',
                name: 'view-inbox-conversation',
                component: () => import('@/views/conversation/ConversationDetailView.vue'),
                props: true,
                meta: { title: 'Caixa da visão', hidePageHeader: true }
              }
            ]
          }
        ]
      },
      {
        path: 'inboxes/search',
        name: 'search',
        component: () => import('@/views/search/SearchView.vue'),
        meta: { title: 'Busca', hidePageHeader: true }
      },
      {
        path: '/inboxes/:type(assigned|unassigned|all|mentioned)?',
        name: 'inboxes',
        redirect: '/inboxes/assigned',
        component: InboxLayout,
        props: true,
        meta: { title: 'Caixa de entrada', hidePageHeader: true },
        children: [
          {
            path: '',
            name: 'inbox',
            component: () => import('@/views/inbox/InboxView.vue'),
            meta: {
              title: 'Caixa de entrada',
              type: (route) => {
                if (route.params.type === 'assigned') return 'Minha caixa'
                if (route.params.type === 'mentioned') return 'Menções'
                return route.params.type
              }
            },
            children: [
              {
                path: 'conversation/:uuid',
                name: 'inbox-conversation',
                component: () => import('@/views/conversation/ConversationDetailView.vue'),
                props: true,
                meta: {
                  title: 'Caixa de entrada',
                  type: (route) => {
                    if (route.params.type === 'assigned') return 'Minha caixa'
                    if (route.params.type === 'mentioned') return 'Menções'
                    return route.params.type
                  },
                  hidePageHeader: true
                }
              }
            ]
          }
        ]
      },
      {
        path: '/account/:page?',
        name: 'account',
        redirect: '/account/profile',
        component: AccountLayout,
        props: true,
        meta: { title: 'Conta' },
        children: [
          {
            path: 'profile',
            name: 'profile',
            component: () => import('@/views/account/profile/ProfileEditView.vue'),
            meta: { title: 'Editar perfil' }
          },
          {
            path: 'preferences',
            name: 'account-preferences',
            component: () => import('@/views/account/preferences/AccountPreferencesView.vue'),
            meta: { title: 'Preferências do app' }
          }
        ]
      },
      {
        path: '/admin',
        name: 'admin',
        component: AdminLayout,
        meta: { title: 'Administração' },
        children: [
          {
            path: 'custom-attributes',
            name: 'custom-attributes',
            component: () => import('@/views/admin/custom-attributes/CustomAttributes.vue'),
            meta: { title: 'Atributos personalizados' }
          },
          {
            path: 'general',
            name: 'general',
            component: () => import('@/views/admin/general/General.vue'),
            meta: { title: 'Geral' }
          },
          {
            path: 'business-hours',
            component: () => import('@/views/admin/business-hours/BusinessHours.vue'),
            meta: { title: 'Horários comerciais' },
            children: [
              {
                path: '',
                name: 'business-hours-list',
                component: () => import('@/views/admin/business-hours/BusinessHoursList.vue')
              },
              {
                path: 'new',
                name: 'new-business-hours',
                component: () =>
                  import('@/views/admin/business-hours/CreateOrEditBusinessHours.vue'),
                meta: { title: 'Novo horário comercial' }
              },
              {
                path: ':id/edit',
                name: 'edit-business-hours',
                props: true,
                component: () =>
                  import('@/views/admin/business-hours/CreateOrEditBusinessHours.vue'),
                meta: { title: 'Editar horário comercial' }
              }
            ]
          },
          {
            path: 'sla',
            component: () => import('@/views/admin/sla/SLA.vue'),
            meta: { title: 'SLA' },
            children: [
              {
                path: '',
                name: 'sla-list',
                component: () => import('@/views/admin/sla/SLAList.vue')
              },
              {
                path: 'new',
                name: 'new-sla',
                component: () => import('@/views/admin/sla/CreateEditSLA.vue'),
                meta: { title: 'Novo SLA' }
              },
              {
                path: ':id/edit',
                props: true,
                name: 'edit-sla',
                component: () => import('@/views/admin/sla/CreateEditSLA.vue'),
                meta: { title: 'Editar SLA' }
              }
            ]
          },
          {
            path: 'inboxes',
            component: () => import('@/views/admin/inbox/InboxView.vue'),
            meta: { title: 'Caixas de entrada' },
            children: [
              {
                path: '',
                name: 'inbox-list',
                component: () => import('@/views/admin/inbox/InboxList.vue')
              },
              {
                path: 'new',
                name: 'new-inbox',
                component: () => import('@/views/admin/inbox/NewInbox.vue'),
                meta: { title: 'Nova caixa de entrada' }
              },
              {
                path: ':id/edit',
                props: true,
                name: 'edit-inbox',
                component: () => import('@/views/admin/inbox/EditInbox.vue'),
                meta: { title: 'Editar caixa de entrada' }
              }
            ]
          },
          {
            path: 'domains',
            name: 'domains',
            component: () => import('@/views/admin/domains/DomainsView.vue'),
            meta: { title: 'Domínios' }
          },
          {
            path: 'notification',
            component: () => import('@/features/admin/notification/NotificationSetting.vue'),
            meta: { title: 'E-mail' }
          },
          {
            path: 'notification/whatsapp',
            name: 'notification-whatsapp',
            component: () => import('@/views/admin/notification/NotificationChannelConfigView.vue'),
            props: { channel: 'whatsapp' },
            meta: { title: 'WhatsApp' }
          },
          {
            path: 'notification/telegram',
            name: 'notification-telegram',
            component: () => import('@/views/admin/notification/NotificationChannelConfigView.vue'),
            props: { channel: 'telegram' },
            meta: { title: 'Telegram' }
          },
          {
            path: 'notification/sms',
            name: 'notification-sms',
            component: () => import('@/views/admin/notification/NotificationChannelConfigView.vue'),
            props: { channel: 'sms' },
            meta: { title: 'SMS' }
          },
          {
            path: 'notification/push',
            name: 'notification-push',
            component: () => import('@/views/admin/notification/NotificationChannelConfigView.vue'),
            props: { channel: 'push' },
            meta: { title: 'Notificação push' }
          },
          {
            path: 'notification/letter-notice-summons',
            redirect: '/admin/notification/official-communications'
          },
          {
            path: 'notification/official-communications',
            name: 'notification-letter-notice-summons',
            component: () =>
              import('@/views/admin/notification/NotificationOfficialCommunicationsView.vue'),
            meta: { title: 'Comunicações oficiais' }
          },
          {
            path: 'teams',
            meta: { title: 'Equipes' },
            children: [
              {
                path: 'agents',
                component: () => import('@/views/admin/agents/Agents.vue'),
                meta: { title: 'Agentes' },
                children: [
                  {
                    path: '',
                    name: 'agent-list',
                    component: () => import('@/views/admin/agents/AgentList.vue')
                  },
                  {
                    path: 'new',
                    name: 'new-agent',
                    component: () => import('@/views/admin/agents/CreateAgent.vue'),
                    meta: { title: 'Criar agente' }
                  },
                  {
                    path: ':id/edit',
                    props: true,
                    name: 'edit-agent',
                    component: () => import('@/views/admin/agents/EditAgent.vue'),
                    meta: { title: 'Editar agente' }
                  }
                ]
              },
              {
                path: 'teams',
                component: () => import('@/views/admin/teams/Teams.vue'),
                meta: { title: 'Equipes' },
                children: [
                  {
                    path: '',
                    name: 'team-list',
                    component: () => import('@/views/admin/teams/TeamList.vue')
                  },
                  {
                    path: 'new',
                    name: 'new-team',
                    component: () => import('@/views/admin/teams/CreateTeamForm.vue'),
                    meta: { title: 'Criar equipe' }
                  },
                  {
                    path: ':id/edit',
                    props: true,
                    name: 'edit-team',
                    component: () => import('@/views/admin/teams/EditTeamForm.vue'),
                    meta: { title: 'Editar equipe' }
                  }
                ]
              },
              {
                path: 'roles',
                component: () => import('@/views/admin/roles/Roles.vue'),
                meta: { title: 'Funções' },
                children: [
                  {
                    path: '',
                    name: 'role-list',
                    component: () => import('@/views/admin/roles/RoleList.vue')
                  },
                  {
                    path: 'new',
                    name: 'new-role',
                    component: () => import('@/views/admin/roles/NewRole.vue'),
                    meta: { title: 'Criar função' }
                  },
                  {
                    path: ':id/edit',
                    props: true,
                    name: 'edit-role',
                    component: () => import('@/views/admin/roles/EditRole.vue'),
                    meta: { title: 'Editar função' }
                  }
                ]
              },
              {
                path: 'activity-log',
                name: 'activity-log',
                component: () => import('@/views/admin/activity-log/ActivityLog.vue'),
                meta: { title: 'Log de atividades' }
              }
            ]
          },
          {
            path: 'automations',
            component: () => import('@/views/admin/automations/Automation.vue'),
            name: 'automations',
            meta: { title: 'Automações' },
            children: [
              {
                path: 'new',
                props: true,
                name: 'new-automation',
                component: () => import('@/views/admin/automations/CreateOrEditRule.vue'),
                meta: { title: 'Criar automação' }
              },
              {
                path: ':id/edit',
                props: true,
                name: 'edit-automation',
                component: () => import('@/views/admin/automations/CreateOrEditRule.vue'),
                meta: { title: 'Editar automação' }
              }
            ]
          },
          {
            path: 'templates',
            component: () => import('@/views/admin/templates/Templates.vue'),
            name: 'templates',
            meta: { title: 'Modelos' },
            children: [
              {
                path: ':id/edit',
                name: 'edit-template',
                props: true,
                component: () => import('@/views/admin/templates/CreateEditTemplate.vue'),
                meta: { title: 'Editar modelo' }
              },
              {
                path: 'new',
                name: 'new-template',
                props: true,
                component: () => import('@/views/admin/templates/CreateEditTemplate.vue'),
                meta: { title: 'Novo modelo' }
              }
            ]
          },
          {
            path: 'sso',
            component: () => import('@/views/admin/oidc/OIDC.vue'),
            name: 'sso',
            meta: { title: 'SSO' },
            children: [
              {
                path: '',
                name: 'sso-list',
                component: () => import('@/views/admin/oidc/OIDCList.vue')
              },
              {
                path: ':id/edit',
                props: true,
                name: 'edit-sso',
                component: () => import('@/views/admin/oidc/CreateEditOIDC.vue'),
                meta: { title: 'Editar SSO' }
              },
              {
                path: 'new',
                name: 'new-sso',
                component: () => import('@/views/admin/oidc/CreateEditOIDC.vue'),
                meta: { title: 'Novo SSO' }
              }
            ]
          },
          {
            path: 'webhooks',
            component: () => import('@/views/admin/webhooks/Webhooks.vue'),
            name: 'webhooks',
            meta: { title: 'Webhooks' },
            children: [
              {
                path: '',
                name: 'webhook-list',
                component: () => import('@/views/admin/webhooks/WebhookList.vue')
              },
              {
                path: ':id/edit',
                props: true,
                name: 'edit-webhook',
                component: () => import('@/views/admin/webhooks/CreateEditWebhook.vue'),
                meta: { title: 'Editar webhook' }
              },
              {
                path: 'new',
                name: 'new-webhook',
                component: () => import('@/views/admin/webhooks/CreateEditWebhook.vue'),
                meta: { title: 'Novo webhook' }
              }
            ]
          },
          {
            path: 'integrations/actions',
            name: 'integrations-actions',
            component: () => import('@/views/admin/integrations/IntegrationActionsView.vue'),
            meta: { title: 'Ações' }
          },
          {
            path: 'conversations',
            meta: { title: 'Conversas' },
            children: [
              {
                path: 'tags',
                component: () => import('@/views/admin/tags/TagsView.vue'),
                meta: { title: 'Tags' }
              },
              {
                path: 'statuses',
                component: () => import('@/views/admin/status/StatusView.vue'),
                meta: { title: 'Status' }
              },
              {
                path: 'macros',
                component: () => import('@/views/admin/macros/Macros.vue'),
                meta: { title: 'Macros' },
                children: [
                  {
                    path: '',
                    name: 'macro-list',
                    component: () => import('@/views/admin/macros/MacroList.vue')
                  },
                  {
                    path: 'new',
                    name: 'new-macro',
                    component: () => import('@/views/admin/macros/CreateMacro.vue'),
                    meta: { title: 'Criar macro' }
                  },
                  {
                    path: ':id/edit',
                    props: true,
                    name: 'edit-macro',
                    component: () => import('@/views/admin/macros/EditMacro.vue'),
                    meta: { title: 'Editar macro' }
                  }
                ]
              },
              {
                path: 'shared-views',
                component: () => import('@/views/admin/shared-views/SharedViews.vue'),
                meta: { title: 'Visões compartilhadas' },
                children: [
                  {
                    path: '',
                    name: 'shared-view-list',
                    component: () => import('@/views/admin/shared-views/SharedViewList.vue')
                  },
                  {
                    path: 'new',
                    name: 'new-shared-view',
                    component: () => import('@/views/admin/shared-views/CreateSharedView.vue'),
                    meta: { title: 'Criar visão compartilhada' }
                  },
                  {
                    path: ':id/edit',
                    props: true,
                    name: 'edit-shared-view',
                    component: () => import('@/views/admin/shared-views/EditSharedView.vue'),
                    meta: { title: 'Editar visão compartilhada' }
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
    redirect: () => {
      return '/inboxes/assigned'
    }
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routes
})

router.beforeEach((to, from, next) => {
  // Make page title with the route name and site name
  const appSettingsStore = useAppSettingsStore()
  const siteName = appSettingsStore.settings?.['app.site_name'] || 'CanalGov'
  const pageTitle = to.meta?.title || ''
  document.title = `${pageTitle} - ${siteName}`
  next()
})

export default router
