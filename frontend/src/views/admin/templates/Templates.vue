<template>
  <div>
    <AdminPageHeader
      :title="$t('admin.templates.title')"
      :description="$t('admin.templates.description')"
      :breadcrumbs="[{ label: $t('globals.terms.administration'), to: '/admin' }, { label: $t('admin.templates.title') }]"
    >
      <template v-if="isListRoute" #actions>
        <Button @click="navigateToNewTemplate" :disabled="templateType === 'email_notification'">
          <Plus class="h-4 w-4 mr-1.5" aria-hidden="true" />
          {{ $t('globals.messages.new', { name: $t('globals.terms.template') }) }}
        </Button>
      </template>
    </AdminPageHeader>

    <AdminPageWithHelp>
      <template #content>
        <template v-if="isListRoute">
          <Tabs default-value="response" v-model="templateType">
            <TabsList class="grid w-full grid-cols-3 mb-5">
              <TabsTrigger value="response">{{ $t('admin.template.responseTemplates') }}</TabsTrigger>
              <TabsTrigger value="email_outgoing">{{ $t('admin.template.emailLayouts') }}</TabsTrigger>
              <TabsTrigger value="email_notification">{{ $t('admin.template.emailNotifications') }}</TabsTrigger>
            </TabsList>
            <TabsContent value="response">
              <DataTable :columns="createResponseTemplateColumns(t)" :data="templates" :loading="isLoading" />
            </TabsContent>
            <TabsContent value="email_outgoing">
              <DataTable :columns="createOutgoingEmailTableColumns(t)" :data="templates" :loading="isLoading" />
            </TabsContent>
            <TabsContent value="email_notification">
              <DataTable :columns="createEmailNotificationTableColumns(t)" :data="templates" :loading="isLoading" />
            </TabsContent>
          </Tabs>
        </template>
        <template v-else>
          <router-view />
        </template>
      </template>

      <template #help>
        <template v-if="isListRoute">
          <p>{{ $t('admin.template.help') }}</p>
          <p>{{ $t('admin.template.help2') }}</p>
          <p>{{ $t('admin.template.help3') }}</p>
        </template>
        <template v-else>
          <p>{{ $t('admin.template.form.help') }}</p>
          <p>{{ $t('admin.template.form.help2') }}</p>
        </template>
      </template>
    </AdminPageWithHelp>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { Plus } from 'lucide-vue-next'
import DataTable from '@/components/datatable/DataTable.vue'
import {
  createResponseTemplateColumns,
  createOutgoingEmailTableColumns,
  createEmailNotificationTableColumns
} from '@/features/admin/templates/dataTableColumns.js'
import { Button } from '@/components/ui/button'
import { useRouter, useRoute } from 'vue-router'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useStorage } from '@vueuse/core'
import AdminPageWithHelp from '@/layouts/admin/AdminPageWithHelp.vue'
import AdminPageHeader from '@/components/layout/AdminPageHeader.vue'
import { useAdminListRefresh } from '@/composables/useAdminListRefresh'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'
import api from '@/api'

const templateType = useStorage('templateType', 'response')
const { t } = useI18n()
const templates = ref([])
const isLoading = ref(false)
const router = useRouter()
const route = useRoute()
const isListRoute = computed(() => route.name === 'templates')
const { showErrorToast } = useAdminErrorToast()

const fetchAll = async () => {
  isLoading.value = true
  try {
    const resp = await api.getTemplates(templateType.value)
    templates.value = resp.data.data
  } catch (error) {
    showErrorToast(error)
  } finally {
    isLoading.value = false
  }
}

useAdminListRefresh('templates', fetchAll)

const navigateToNewTemplate = () => {
  router.push({ name: 'new-template', query: { type: templateType.value } })
}

watch(templateType, () => {
  templates.value = []
  fetchAll()
})

watch(
  () => route.name,
  () => {
    if (route.name === 'templates') fetchAll()
  }
)
</script>
