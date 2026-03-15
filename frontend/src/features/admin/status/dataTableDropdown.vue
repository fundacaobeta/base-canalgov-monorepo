<template>
  <Dialog v-model:open="dialogOpen">
    <DataTableRowActions
      v-if="!CONVERSATION_DEFAULT_STATUSES_LIST.includes(props.status.name)"
      :entity-name="t('globals.terms.status')"
      :delete-description="t('globals.messages.deletionConfirmation', { name: t('globals.terms.status').toLowerCase() })"
      @edit="dialogOpen = true"
      @delete="handleDelete"
    />
    <div v-else class="w-8 h-8 p-0 invisible" />

    <DialogContent class="sm:max-w-[425px]">
      <DialogHeader>
        <DialogTitle>{{ t('globals.messages.edit', { name: t('globals.terms.status') }) }}</DialogTitle>
        <DialogDescription>{{ t('admin.conversationStatus.name.description') }}</DialogDescription>
      </DialogHeader>
      <StatusForm @submit.prevent="onSubmit">
        <template #footer>
          <DialogFooter class="mt-6 gap-2">
            <Button variant="outline" type="button" @click="dialogOpen = false">
              {{ t('globals.messages.cancel') }}
            </Button>
            <Button type="submit" :disabled="isLoading">{{ t('globals.messages.save') }}</Button>
          </DialogFooter>
        </template>
      </StatusForm>
    </DialogContent>
  </Dialog>
</template>

<script setup>
import { watch, ref } from 'vue'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { createFormSchema } from './formSchema.js'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import StatusForm from './StatusForm.vue'
import DataTableRowActions from '@/components/admin/DataTableRowActions.vue'
import { CONVERSATION_DEFAULT_STATUSES_LIST } from '@/constants/conversation.js'
import { useAdminListRefresh } from '@/composables/useAdminListRefresh'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'
import api from '@/api/index.js'

const { t } = useI18n()
const isLoading = ref(false)
const dialogOpen = ref(false)
const { emitRefresh } = useAdminListRefresh('status', () => {})
const { showErrorToast } = useAdminErrorToast()

const props = defineProps({
  status: { type: Object, required: true }
})

const form = useForm({ validationSchema: toTypedSchema(createFormSchema(t)) })

const onSubmit = form.handleSubmit(async (values) => {
  isLoading.value = true
  try {
    await api.updateStatus(props.status.id, values)
    dialogOpen.value = false
    emitRefresh()
  } catch (error) {
    showErrorToast(error)
  } finally {
    isLoading.value = false
  }
})

const handleDelete = async () => {
  isLoading.value = true
  try {
    await api.deleteStatus(props.status.id)
    emitRefresh()
  } catch (error) {
    showErrorToast(error)
  } finally {
    isLoading.value = false
  }
}

watch(() => props.status, (val) => form.setValues(val), { immediate: true, deep: true })
</script>
