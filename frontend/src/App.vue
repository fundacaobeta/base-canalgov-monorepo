<template>
  <div
    class="flex w-full h-screen text-foreground bg-canvas p-1.5"
    style="background-image: radial-gradient(circle at top left, hsl(var(--primary) / 0.18), transparent 24%), radial-gradient(circle at bottom right, hsl(var(--accent) / 0.16), transparent 22%), linear-gradient(180deg, hsl(var(--canvas)), hsl(var(--background)));"
  >
    <!-- Icon sidebar always visible -->
    <SidebarProvider style="--sidebar-width: 3rem" class="w-auto z-50">
      <ShadcnSidebar collapsible="none" class="border rounded-lg overflow-hidden">
        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupContent>
              <SidebarMenu>
                <SidebarMenuItem>
                  <Tooltip>
                    <TooltipTrigger as-child>
                      <SidebarMenuButton asChild :isActive="route.matched.some(r => r.meta?.area === 'inboxes')">
                        <router-link :to="{ name: 'inboxes' }">
                          <Inbox />
                        </router-link>
                      </SidebarMenuButton>
                    </TooltipTrigger>
                    <TooltipContent side="right">
                      <p>{{ t('globals.terms.inbox', 2) }}</p>
                    </TooltipContent>
                  </Tooltip>
                </SidebarMenuItem>
                <SidebarMenuItem v-if="userStore.can('contacts:read_all')">
                  <Tooltip>
                    <TooltipTrigger as-child>
                      <SidebarMenuButton asChild :isActive="route.matched.some(r => r.meta?.area === 'contacts')">
                        <router-link :to="{ name: 'contacts' }">
                          <BookUser />
                        </router-link>
                      </SidebarMenuButton>
                    </TooltipTrigger>
                    <TooltipContent side="right">
                      <p>{{ t('globals.terms.contact', 2) }}</p>
                    </TooltipContent>
                  </Tooltip>
                </SidebarMenuItem>
                <SidebarMenuItem v-if="userStore.hasReportTabPermissions">
                  <Tooltip>
                    <TooltipTrigger as-child>
                      <SidebarMenuButton asChild :isActive="route.matched.some(r => r.meta?.area === 'reports')">
                        <router-link :to="{ name: 'reports' }">
                          <FileLineChart />
                        </router-link>
                      </SidebarMenuButton>
                    </TooltipTrigger>
                    <TooltipContent side="right">
                      <p>{{ t('globals.terms.report', 2) }}</p>
                    </TooltipContent>
                  </Tooltip>
                </SidebarMenuItem>
                <SidebarMenuItem v-if="userStore.hasAdminTabPermissions">
                  <Tooltip>
                    <TooltipTrigger as-child>
                      <SidebarMenuButton asChild :isActive="route.matched.some(r => r.meta?.area === 'admin')">
                        <router-link
                          :to="{
                            name: userStore.can('general_settings:manage') ? 'general' : 'admin'
                          }"
                        >
                          <Shield />
                        </router-link>
                      </SidebarMenuButton>
                    </TooltipTrigger>
                    <TooltipContent side="right">
                      <p>{{ t('globals.terms.admin') }}</p>
                    </TooltipContent>
                  </Tooltip>
                </SidebarMenuItem>
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
        <SidebarFooter>
          <SidebarMenu>
            <SidebarMenuItem>
              <Tooltip>
                <TooltipTrigger as-child>
                  <NotificationBell />
                </TooltipTrigger>
                <TooltipContent side="right">
                  <p>{{ t('globals.terms.notification', 2) }}</p>
                </TooltipContent>
              </Tooltip>
            </SidebarMenuItem>
            <SidebarMenuItem>
              <SidebarNavUser />
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarFooter>
      </ShadcnSidebar>
    </SidebarProvider>

    <!-- Main sidebar that collapses -->
    <div class="flex-1">
      <Sidebar
        :userTeams="userStore.teams"
        :userViews="userViews"
        :sharedViews="sharedViewStore.sharedViewList"
        @create-view="openCreateViewForm = true"
        @edit-view="editView"
        @delete-view="deleteView"
        @create-conversation="() => (openCreateConversationDialog = true)"
      >
        <div class="flex flex-col h-full rounded-lg overflow-hidden bg-background">
          <!-- Show admin banner only in admin routes -->
          <AdminBanner v-if="route.matched.some(r => r.meta?.area === 'admin')" />

          <!-- Common header for all pages -->
          <PageHeader />

          <!-- Main content -->
          <RouterView class="flex-grow" />
        </div>
        <ViewForm v-model:openDialog="openCreateViewForm" v-model:view="view" />
      </Sidebar>
    </div>
  </div>

  <!-- Command box -->
  <Command />

  <!-- Create conversation dialog -->
  <CreateConversation v-model="openCreateConversationDialog" v-if="openCreateConversationDialog" />
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { RouterView } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useAppSettingsStore } from '@/stores/appSettings'
import { initWS } from '@/websocket.js'
import { EMITTER_EVENTS } from '@/constants/emitterEvents.js'
import { REFRESH_MODEL } from '@/constants/conversation'
import { useEmitter } from '@/composables/useEmitter'
import { handleHTTPError } from '@/utils/http'
import { useConversationStore } from './stores/conversation'
import { useInboxStore } from '@/stores/inbox'
import { useAgentsStore } from '@/stores/agents'
import { useTeamStore } from '@/stores/team'
import { useSlaStore } from '@/stores/sla'
import { useMacroStore } from '@/stores/macro'
import { useSharedViewStore } from '@/stores/sharedView'
import { useTagStore } from '@/stores/tag'
import { useCustomAttributeStore } from '@/stores/customAttributes'
import { useIdleDetection } from '@/composables/useIdleDetection'
import PageHeader from './components/layout/PageHeader.vue'
import ViewForm from '@/features/view/ViewForm.vue'
import AdminBanner from '@/components/banner/AdminBanner.vue'
import api from '@/api'
import { toast as sooner } from 'vue-sonner'
import Sidebar from '@/components/sidebar/Sidebar.vue'
import Command from '@/features/command/CommandBox.vue'
import CreateConversation from '@/features/conversation/CreateConversation.vue'
import { Inbox, Shield, FileLineChart, BookUser } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import {
  Sidebar as ShadcnSidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarMenu,
  SidebarGroupContent,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider
} from '@/components/ui/sidebar'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import SidebarNavUser from '@/components/sidebar/SidebarNavUser.vue'
import NotificationBell from '@/components/sidebar/NotificationBell.vue'

const route = useRoute()
const emitter = useEmitter()
const userStore = useUserStore()
const appSettingsStore = useAppSettingsStore()
const conversationStore = useConversationStore()
const usersStore = useAgentsStore()
const teamStore = useTeamStore()
const inboxStore = useInboxStore()
const slaStore = useSlaStore()
const macroStore = useMacroStore()
const sharedViewStore = useSharedViewStore()
const tagStore = useTagStore()
const customAttributeStore = useCustomAttributeStore()
const userViews = ref([])
const view = ref({})
const openCreateViewForm = ref(false)
const openCreateConversationDialog = ref(false)
const { t } = useI18n()

initWS()
useIdleDetection()

onMounted(() => {
  initToaster()
  listenViewRefresh()
  initStores()
})

// Initialize data stores
const initStores = async () => {
  if (!userStore.userID) {
    await userStore.getCurrentUser()
  }
  await Promise.allSettled([
    getUserViews(),
    appSettingsStore.fetchPublicConfig(),
    sharedViewStore.loadSharedViews(),
    conversationStore.fetchStatuses(),
    conversationStore.fetchPriorities(),
    conversationStore.fetchAllDrafts(),
    usersStore.fetchAgents(),
    teamStore.fetchTeams(),
    inboxStore.fetchInboxes(),
    slaStore.fetchSlas(),
    macroStore.loadMacros(),
    tagStore.fetchTags(),
    customAttributeStore.fetchCustomAttributes()
  ])
}

const editView = (v) => {
  view.value = { ...v }
  openCreateViewForm.value = true
}

const deleteView = async (view) => {
  try {
    await api.deleteView(view.id)
    emitter.emit(EMITTER_EVENTS.REFRESH_LIST, { model: 'view' })
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      description: t('globals.messages.deletedSuccessfully', {
        name: t('globals.terms.view')
      })
    })
  } catch (err) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(err).message
    })
  }
}

const getUserViews = async () => {
  try {
    const response = await api.getCurrentUserViews()
    userViews.value = response.data.data
  } catch (err) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(err).message
    })
  }
}

const initToaster = () => {
  emitter.on(EMITTER_EVENTS.SHOW_TOAST, (message) => {
    if (message.variant === 'destructive') {
      sooner.error(message.description)
    } else if (message.variant === 'warning') {
      sooner.warning(message.description)
    } else {
      sooner.success(message.description)
    }
  })
}

const listenViewRefresh = () => {
  emitter.on(EMITTER_EVENTS.REFRESH_LIST, refreshViews)
}

const refreshViews = (data) => {
  openCreateViewForm.value = false
  if (data?.model === REFRESH_MODEL.VIEW) {
    getUserViews()
  }
}
</script>

<style scoped>
:deep(.group\/sidebar-wrapper) {
  min-height: auto !important;
  height: 100%;
}
</style>
