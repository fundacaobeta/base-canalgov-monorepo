<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button variant="ghost" class="w-8 h-8 p-0">
        <span class="sr-only">{{ t('globals.messages.openMenu') }}</span>
        <MoreHorizontal class="w-4 h-4" />
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent>
      <DropdownMenuItem v-if="showEdit" @click="$emit('edit')">
        {{ t('globals.messages.edit') }}
      </DropdownMenuItem>
      <slot name="extra-items" />
      <DropdownMenuItem v-if="showDelete" class="text-destructive" @click="alertOpen = true">
        {{ t('globals.messages.delete') }}
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>

  <AlertDialog :open="alertOpen" @update:open="alertOpen = $event">
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>{{ deleteTitle || t('globals.messages.areYouAbsolutelySure') }}</AlertDialogTitle>
        <AlertDialogDescription>
          {{ deleteDescription || t('globals.messages.deletionConfirmation', { name: entityName }) }}
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel>{{ t('globals.messages.cancel') }}</AlertDialogCancel>
        <AlertDialogAction @click="confirmDelete">
          {{ t('globals.messages.delete') }}
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>

<script setup>
import { ref } from 'vue'
import { MoreHorizontal } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
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
import { useI18n } from 'vue-i18n'

const props = defineProps({
  entityName: { type: String, default: '' },
  deleteTitle: { type: String, default: '' },
  deleteDescription: { type: String, default: '' },
  showEdit: { type: Boolean, default: true },
  showDelete: { type: Boolean, default: true }
})

const emit = defineEmits(['edit', 'delete'])

const { t } = useI18n()
const alertOpen = ref(false)

const confirmDelete = () => {
  alertOpen.value = false
  emit('delete')
}
</script>
