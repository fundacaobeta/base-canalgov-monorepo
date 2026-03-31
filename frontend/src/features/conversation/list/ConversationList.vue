<template>
  <div class="h-screen flex flex-col">
    <div class="flex items-center space-x-4 px-2 h-12 border-b shrink-0">
      <SidebarTrigger class="cursor-pointer" />
      <span class="text-xl font-semibold">{{ title }}</span>
    </div>

    <div class="p-2 flex flex-wrap justify-between items-start gap-2">
      <div v-if="!route.params.viewID" class="flex flex-wrap items-center gap-2">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" class="w-auto min-w-30">
              <div class="flex items-center gap-2">
                <span class="text-xs text-muted-foreground">Caixa:</span>
                <Badge variant="outline">
                  {{ currentBoxLabel }}
                </Badge>
              </div>
              <ChevronDown class="w-4 h-4 ml-2 opacity-50" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem @click="handleBoxChange(CONVERSATION_LIST_TYPE.ASSIGNED)">
              Minha caixa
            </DropdownMenuItem>
            <DropdownMenuItem @click="handleBoxChange(CONVERSATION_LIST_TYPE.UNASSIGNED)">
              Nao atribuidos
            </DropdownMenuItem>
            <DropdownMenuItem @click="handleBoxChange(CONVERSATION_LIST_TYPE.MENTIONED)">
              Menções
            </DropdownMenuItem>
            <DropdownMenuItem @click="handleBoxChange(CONVERSATION_LIST_TYPE.ALL)">
              Todos os chamados
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>

        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" class="w-auto min-w-30">
              <div class="flex items-center gap-2">
                <span class="text-xs text-muted-foreground">Status:</span>
                <Badge :variant="conversationStore.getListStatus ? getConversationStatusBadgeVariant(conversationStore.getListStatus) : 'outline'">
                  {{ currentStatusLabel }}
                </Badge>
                <span class="text-xs text-muted-foreground">({{ conversationStore.conversations.total }})</span>
              </div>
              <ChevronDown class="w-4 h-4 ml-2 opacity-50" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem @click="handleStatusChange(null)">
              Todos os status
            </DropdownMenuItem>
            <DropdownMenuItem
              v-for="status in conversationStore.statusOptions"
              :key="status.value"
              @click="handleStatusChange(status)"
            >
              {{ status.label }}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      <div v-else>
        <Button variant="ghost" class="w-30">
          <span>{{ conversationStore.conversations.total }}</span>
        </Button>
      </div>

      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" class="w-full sm:w-auto min-w-[13rem] justify-between ml-auto">
            <span class="truncate">{{ $t(conversationStore.getListSortField) }}</span>
            <ChevronDown class="w-4 h-4 ml-2 opacity-50" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuItem @click="handleSortChange('oldest')">
            {{ $t('conversation.sort.oldestActivity') }}
          </DropdownMenuItem>
          <DropdownMenuItem @click="handleSortChange('newest')">
            {{ $t('conversation.sort.newestActivity') }}
          </DropdownMenuItem>
          <DropdownMenuItem @click="handleSortChange('started_first')">
            {{ $t('conversation.sort.startedFirst') }}
          </DropdownMenuItem>
          <DropdownMenuItem @click="handleSortChange('started_last')">
            {{ $t('conversation.sort.startedLast') }}
          </DropdownMenuItem>
          <DropdownMenuItem @click="handleSortChange('waiting_longest')">
            {{ $t('conversation.sort.waitingLongest') }}
          </DropdownMenuItem>
          <DropdownMenuItem @click="handleSortChange('next_sla_target')">
            {{ $t('conversation.sort.nextSLATarget') }}
          </DropdownMenuItem>
          <DropdownMenuItem @click="handleSortChange('priority_first')">
            {{ $t('conversation.sort.priorityFirst') }}
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>

    <div class="flex-grow overflow-y-auto">
      <EmptyList
        v-if="!hasConversations && !hasErrored && !isLoading"
        key="empty"
        class="px-4 py-8"
        :title="t('conversation.noConversationsFound')"
        :message="t('conversation.tryAdjustingFilters')"
        :icon="MessageCircleQuestion"
      />

      <EmptyList
        v-if="conversationStore.conversations.errorMessage"
        key="error"
        class="px-4 py-8"
        :title="t('conversation.couldNotFetch')"
        :message="conversationStore.conversations.errorMessage"
        :icon="MessageCircleWarning"
      />

      <TransitionGroup
        enter-active-class="transition-all duration-300 ease-in-out"
        enter-from-class="opacity-0 transform translate-y-4"
        enter-to-class="opacity-100 transform translate-y-0"
        leave-active-class="transition-all duration-300 ease-in-out"
        leave-from-class="opacity-100 transform translate-y-0"
        leave-to-class="opacity-0 transform translate-y-4"
      >
        <div
          v-if="!conversationStore.conversations.errorMessage"
          key="list"
          class="divide-y divide-gray-200 dark:divide-gray-700"
          :class="{ 'border-b dark:border-gray-700': hasConversations }"
        >
          <ConversationListItem
            v-for="conversation in conversationStore.conversationsList"
            :key="conversation.uuid"
            :conversation="conversation"
            :currentConversation="conversationStore.current"
            :contactFullName="conversationStore.getContactFullName(conversation.uuid)"
            class="transition-colors duration-200 hover:bg-gray-50 dark:hover:bg-gray-600"
          />
        </div>

        <div v-if="isLoading" key="loading" class="space-y-4">
          <ConversationListItemSkeleton v-for="index in 5" :key="index" />
        </div>
      </TransitionGroup>

      <div
        v-if="!hasErrored && (conversationStore.conversations.hasMore || hasConversations)"
        class="flex justify-center items-center p-5"
      >
        <Button
          v-if="conversationStore.conversations.hasMore"
          variant="outline"
          @click="loadNextPage"
          :disabled="isLoading"
          class="transition-all duration-200 ease-in-out transform hover:scale-105"
        >
          <Loader2 v-if="isLoading" class="mr-2 h-4 w-4 animate-spin" />
          {{ isLoading ? t('globals.terms.loading') : t('globals.terms.loadMore') }}
        </Button>
        <p
          class="text-sm text-gray-500"
          v-else-if="conversationStore.conversationsList.length > 10"
        >
          {{ $t('conversation.allLoaded') }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, watch } from 'vue'
import { useConversationStore } from '@/stores/conversation'
import { MessageCircleQuestion, MessageCircleWarning, ChevronDown, Loader2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
import { SidebarTrigger } from '@/components/ui/sidebar'
import EmptyList from '@/features/conversation/list/ConversationEmptyList.vue'
import ConversationListItem from '@/features/conversation/list/ConversationListItem.vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import ConversationListItemSkeleton from '@/features/conversation/list/ConversationListItemSkeleton.vue'
import { getConversationStatusBadgeVariant, translateConversationStatus } from '@/utils/conversationStatus'
import { CONVERSATION_LIST_TYPE, CONVERSATION_DEFAULT_STATUSES } from '@/constants/conversation'
import { getLocalePaths } from '@/router/paths'
import { useInboxTypes } from '@/composables/useInboxTypes'

const conversationStore = useConversationStore()
const route = useRoute()
const router = useRouter()
const { t, locale } = useI18n()
const { toParam: inboxTypeParam } = useInboxTypes()

const listType = computed(() => {
  const paramsType = route.params.type
  const p = getLocalePaths(locale.value)
  const internalType = Object.fromEntries(
    Object.entries(p.inboxTypes).map(([internal, localized]) => [localized, internal])
  )
  return internalType[paramsType] ?? paramsType
})

const fetchCurrentList = () => {
  if (route.params.viewID) {
    conversationStore.setListStatus('', false)
    conversationStore.fetchConversationsList(
      true,
      CONVERSATION_LIST_TYPE.VIEW,
      0,
      [],
      route.params.viewID
    )
    return
  }

  if (route.params.teamID) {
    if (!conversationStore.getListStatus) {
      conversationStore.setListStatus(CONVERSATION_DEFAULT_STATUSES.OPEN, false)
    }
    conversationStore.fetchConversationsList(
      true,
      CONVERSATION_LIST_TYPE.TEAM_UNASSIGNED,
      route.params.teamID
    )
    return
  }

  if (listType.value) {
    if (!conversationStore.getListStatus && listType.value !== CONVERSATION_LIST_TYPE.ALL) {
      conversationStore.setListStatus(CONVERSATION_DEFAULT_STATUSES.OPEN, false)
    }
    conversationStore.fetchConversationsList(true, listType.value)
  }
}

const title = computed(() => {
  const key = route.meta?.type?.(route) || route.meta?.title || ''
  return key ? t(key) : ''
})

const currentBoxLabel = computed(() => {
  const map = {
    [CONVERSATION_LIST_TYPE.ASSIGNED]: 'Minha caixa',
    [CONVERSATION_LIST_TYPE.UNASSIGNED]: 'Nao atribuidos',
    [CONVERSATION_LIST_TYPE.MENTIONED]: 'Menções',
    [CONVERSATION_LIST_TYPE.ALL]: 'Todos os chamados'
  }
  return map[listType.value] || 'Minha caixa'
})

const currentStatusLabel = computed(() => {
  if (!conversationStore.getListStatus) return 'Todos os status'
  return translateConversationStatus(conversationStore.getListStatus, t)
})

const handleStatusChange = (status) => {
  const nextStatus = status?.name || ''
  conversationStore.setListStatus(nextStatus, false)

  if (!nextStatus && listType.value !== CONVERSATION_LIST_TYPE.ALL && !route.params.viewID && !route.params.teamID) {
    router.push({
      name: 'inbox',
      params: { type: inboxTypeParam(CONVERSATION_LIST_TYPE.ALL) }
    })
    return
  }

  conversationStore.fetchFirstPageConversations()
}

const handleBoxChange = (type) => {
  if (type === CONVERSATION_LIST_TYPE.ALL) {
    conversationStore.setListStatus('', false)
  }
  router.push({
    name: 'inbox',
    params: { type: inboxTypeParam(type) }
  })
}

const handleSortChange = (order) => {
  conversationStore.setListSortField(order)
}

const loadNextPage = () => {
  conversationStore.fetchNextConversations()
}

const hasConversations = computed(() => conversationStore.conversationsList.length !== 0)
const hasErrored = computed(() => !!conversationStore.conversations.errorMessage)
const isLoading = computed(() => conversationStore.conversations.loading)

onMounted(fetchCurrentList)

watch(
  () => [route.params.type, route.params.teamID, route.params.viewID, locale.value],
  fetchCurrentList
)
</script>
