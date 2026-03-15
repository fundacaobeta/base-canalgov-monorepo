<template>
  <div>
    <AdminPageHeader
      :title="t('globals.terms.tag', 2)"
      :description="t('admin.tags.description', 'Categorize conversas para organizar e filtrar atendimentos.')"
      :breadcrumbs="[{ label: $t('globals.terms.administration'), to: '/admin' }, { label: t('globals.terms.tag', 2) }]"
    >
      <template #actions>
        <Dialog v-model:open="dialogOpen">
          <DialogTrigger as-child>
            <Button>
              <Plus class="h-4 w-4 mr-1.5" aria-hidden="true" />
              {{ t('globals.messages.new', { name: t('globals.terms.tag') }) }}
            </Button>
          </DialogTrigger>
          <DialogContent class="sm:max-w-[440px]">
            <DialogHeader>
              <DialogTitle>{{ t('globals.messages.new', { name: t('globals.terms.tag') }) }}</DialogTitle>
              <DialogDescription>{{ t('admin.conversationTags.new.description') }}</DialogDescription>
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
    </AdminPageHeader>

    <AdminPageWithHelp>
      <template #content>
        <DataTable :columns="createColumns(t)" :data="tags" :loading="isLoading" />
      </template>
      <template #help>
        <p>{{ t('admin.tags.help') }}</p>
      </template>
    </AdminPageWithHelp>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Plus } from 'lucide-vue-next'
import DataTable from '@/components/datatable/DataTable.vue'
import AdminPageWithHelp from '@/layouts/admin/AdminPageWithHelp.vue'
import AdminPageHeader from '@/components/layout/AdminPageHeader.vue'
import { createColumns } from '@/features/admin/tags/dataTableColumns.js'
import { Button } from '@/components/ui/button'
import TagsForm from '@/features/admin/tags/TagsForm.vue'
import {
  Dialog, DialogContent, DialogDescription,
  DialogFooter, DialogHeader, DialogTitle, DialogTrigger
} from '@/components/ui/dialog'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { createFormSchema } from '@/features/admin/tags/formSchema.js'
import { useAdminListRefresh } from '@/composables/useAdminListRefresh'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'
import api from '@/api'

const { t } = useI18n()
const isLoading = ref(false)
const tags = ref([])
const dialogOpen = ref(false)
const { showErrorToast, showSuccessToast } = useAdminErrorToast()

const getTags = async () => {
  isLoading.value = true
  try {
    const resp = await api.getTags()
    tags.value = resp.data.data
  } finally {
    isLoading.value = false
  }
}

useAdminListRefresh('tags', getTags)

const form = useForm({ validationSchema: toTypedSchema(createFormSchema(t)) })

const onSubmit = form.handleSubmit(async (values) => {
  isLoading.value = true
  try {
    await api.createTag(values)
    dialogOpen.value = false
    showSuccessToast(t('globals.messages.createdSuccessfully', { name: t('globals.terms.tag') }))
    await getTags()
  } catch (error) {
    showErrorToast(error)
  } finally {
    isLoading.value = false
  }
})
</script>
