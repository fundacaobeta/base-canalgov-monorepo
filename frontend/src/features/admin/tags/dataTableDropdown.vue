<template>
  <Dialog v-model:open="dialogOpen">
    <DataTableRowActions
      :entity-name="t('globals.terms.tag')"
      :delete-description="t('admin.tags.deleteConfirmation')"
      @edit="dialogOpen = true"
      @delete="deleteTag"
    />
    <DialogContent class="sm:max-w-[425px]">
      <DialogHeader>
        <DialogTitle>{{ t('globals.messages.edit', { name: t('globals.terms.tag') }) }}</DialogTitle>
        <DialogDescription>{{ t('admin.conversationTags.edit.description') }}</DialogDescription>
      </DialogHeader>
      <TagsForm @submit.prevent="onSubmit">
        <template #footer>
          <DialogFooter class="mt-6 gap-2">
            <Button variant="outline" type="button" @click="dialogOpen = false">
              {{ t('globals.messages.cancel') }}
            </Button>
            <Button type="submit">{{ t('globals.messages.save') }}</Button>
          </DialogFooter>
        </template>
      </TagsForm>
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
import TagsForm from './TagsForm.vue'
import DataTableRowActions from '@/components/admin/DataTableRowActions.vue'
import { useAdminListRefresh } from '@/composables/useAdminListRefresh'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'
import api from '@/api/index.js'

const { t } = useI18n()
const dialogOpen = ref(false)
const { emitRefresh } = useAdminListRefresh('tags', () => {})
const { showErrorToast, showSuccessToast } = useAdminErrorToast()

const props = defineProps({
  tag: { type: Object, required: true, default: () => ({ id: '', name: '' }) }
})

const form = useForm({ validationSchema: toTypedSchema(createFormSchema(t)) })

const onSubmit = form.handleSubmit(async (values) => {
  try {
    await api.updateTag(props.tag.id, values)
    showSuccessToast(t('globals.messages.updatedSuccessfully', { name: t('globals.terms.tag') }))
    dialogOpen.value = false
    emitRefresh()
  } catch (error) {
    showErrorToast(error)
  }
})

const deleteTag = async () => {
  try {
    await api.deleteTag(props.tag.id)
    emitRefresh()
  } catch (error) {
    showErrorToast(error)
  }
}

watch(() => props.tag, (val) => form.setValues(val), { immediate: true, deep: true })
</script>
