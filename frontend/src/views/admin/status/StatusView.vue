<template>
  <div>
    <AdminPageHeader
      :title="t('globals.terms.status', 2)"
      :description="t('admin.status.description', 'Defina os estados possíveis de uma conversa no seu fluxo de atendimento.')"
      :breadcrumbs="[{ label: $t('globals.terms.administration'), to: '/admin' }, { label: t('globals.terms.status', 2) }]"
    >
      <template #actions>
        <Dialog v-model:open="dialogOpen">
          <DialogTrigger as-child>
            <Button>
              <Plus class="h-4 w-4 mr-1.5" aria-hidden="true" />
              {{ t('globals.messages.new', { name: t('globals.terms.status') }) }}
            </Button>
          </DialogTrigger>
          <DialogContent class="sm:max-w-[425px]">
            <DialogHeader>
              <DialogTitle>{{ t('globals.messages.new', { name: t('globals.terms.status') }) }}</DialogTitle>
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
    </AdminPageHeader>

    <AdminPageWithHelp>
      <template #content>
        <DataTable :columns="createColumns(t)" :data="statuses" :loading="isLoading" />
      </template>
      <template #help>
        <p>{{ t('admin.status.help') }}</p>
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
import { createColumns } from '@/features/admin/status/dataTableColumns.js'
import { Button } from '@/components/ui/button'
import StatusForm from '@/features/admin/status/StatusForm.vue'
import {
  Dialog, DialogContent, DialogDescription,
  DialogFooter, DialogHeader, DialogTitle, DialogTrigger
} from '@/components/ui/dialog'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { createFormSchema } from '@/features/admin/status/formSchema.js'
import { useAdminListRefresh } from '@/composables/useAdminListRefresh'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'
import api from '@/api'

const { t } = useI18n()
const isLoading = ref(false)
const statuses = ref([])
const dialogOpen = ref(false)
const { showErrorToast } = useAdminErrorToast()

const getStatuses = async () => {
  isLoading.value = true
  try {
    const resp = await api.getStatuses()
    statuses.value = resp.data.data
  } finally {
    isLoading.value = false
  }
}

useAdminListRefresh('status', getStatuses)

const form = useForm({ validationSchema: toTypedSchema(createFormSchema(t)) })

const onSubmit = form.handleSubmit(async (values) => {
  isLoading.value = true
  try {
    await api.createStatus(values)
    dialogOpen.value = false
    await getStatuses()
  } catch (error) {
    showErrorToast(error)
  } finally {
    isLoading.value = false
  }
})
</script>
