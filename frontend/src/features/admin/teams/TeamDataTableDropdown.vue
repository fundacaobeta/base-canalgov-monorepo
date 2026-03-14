<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button variant="ghost" class="w-8 h-8 p-0">
        <span class="sr-only">{{ t('globals.messages.openMenu') }}</span>
        <MoreHorizontal class="w-4 h-4" />
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent>
      <DropdownMenuItem @click="editTeam(props.team.id)">{{ t('globals.messages.editTeam') }}</DropdownMenuItem>
      <DropdownMenuItem @click="() => (alertOpen = true)">{{ t('globals.messages.deleteTeam') }}</DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>

  <AlertDialog :open="alertOpen" @update:open="alertOpen = $event">
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>{{ t('globals.messages.deleteTeam') }}</AlertDialogTitle>
        <AlertDialogDescription>
          {{ t('admin.team.deleteConfirmation') }}
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel>{{ t('globals.messages.cancel') }}</AlertDialogCancel>
        <AlertDialogAction @click="handleDelete">{{ t('globals.messages.deleteTeam') }}</AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>

<script setup>
import { ref } from 'vue'
import { MoreHorizontal } from 'lucide-vue-next'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { useRouter } from 'vue-router'
import { useEmitter } from '@/composables/useEmitter'
import { EMITTER_EVENTS } from '@/constants/emitterEvents.js'
import { handleHTTPError } from '@/utils/http'
import { useI18n } from 'vue-i18n'
import api from '@/api'

const alertOpen = ref(false)
const router = useRouter()
const emit = useEmitter()
const { t } = useI18n()

const props = defineProps({
  team: {
    type: Object,
    required: true,
    default: () => ({
      id: ''
    })
  }
})

function editTeam(id) {
  router.push({ path: `/admin/teams/teams/${id}/edit` })
}

async function handleDelete() {
  try {
    await api.deleteTeam(props.team.id)
    alertOpen.value = false
    emitRefreshTeamList()
  } catch (error) {
    emit.emit(EMITTER_EVENTS.SHOW_TOAST, {
      title: t('globals.terms.error'),
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  }
}

const emitRefreshTeamList = () => {
  emit.emit(EMITTER_EVENTS.REFRESH_LIST, {
    model: 'team'
  })
}
</script>
