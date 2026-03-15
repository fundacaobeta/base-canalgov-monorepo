<template>
  <div>
    <AdminPageHeader
      :title="$t('admin.customAttributes.title')"
      :description="$t('admin.customAttributes.description')"
      :breadcrumbs="[{ label: $t('globals.terms.administration'), to: '/admin' }, { label: $t('admin.customAttributes.title') }]"
    >
      <template #actions>
        <Dialog v-model:open="dialogOpen">
          <DialogTrigger as-child>
            <Button @click="newCustomAttribute">
              <Plus class="h-4 w-4 mr-1.5" aria-hidden="true" />
              {{ $t('globals.messages.new', { name: $t('globals.terms.customAttribute').toLowerCase() }) }}
            </Button>
          </DialogTrigger>
          <DialogContent class="sm:max-w-[600px]">
            <DialogHeader>
              <DialogTitle>
                {{
                  isEditing
                    ? $t('globals.messages.edit', { name: $t('globals.terms.customAttribute').toLowerCase() })
                    : $t('globals.messages.new', { name: $t('globals.terms.customAttribute').toLowerCase() })
                }}
              </DialogTitle>
              <DialogDescription />
            </DialogHeader>
            <CustomAttributesForm @submit.prevent="onSubmit" :form="form">
              <template #footer>
                <DialogFooter class="mt-6 gap-2">
                  <Button variant="outline" type="button" @click="dialogOpen = false">
                    {{ $t('globals.messages.cancel') }}
                  </Button>
                  <Button type="submit" :disabled="isLoading">
                    {{ isEditing ? $t('globals.messages.update') : $t('globals.messages.create') }}
                  </Button>
                </DialogFooter>
              </template>
            </CustomAttributesForm>
          </DialogContent>
        </Dialog>
      </template>
    </AdminPageHeader>

    <AdminPageWithHelp>
      <template #content>
        <Tabs default-value="contact" v-model="appliesTo">
          <TabsList class="grid w-full grid-cols-2 mb-5">
            <TabsTrigger value="contact">{{ $t('globals.terms.contact') }}</TabsTrigger>
            <TabsTrigger value="conversation">{{ $t('globals.terms.conversation') }}</TabsTrigger>
          </TabsList>
          <TabsContent value="contact">
            <DataTable :columns="createColumns(t)" :data="customAttributes" :loading="isLoading" />
          </TabsContent>
          <TabsContent value="conversation">
            <DataTable :columns="createColumns(t)" :data="customAttributes" :loading="isLoading" />
          </TabsContent>
        </Tabs>
      </template>

      <template #help>
        <p>{{ $t('admin.customAttributes.help') }}</p>
        <p>{{ $t('admin.customAttributes.help2') }}</p>
      </template>
    </AdminPageWithHelp>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { Plus } from 'lucide-vue-next'
import DataTable from '@/components/datatable/DataTable.vue'
import { createColumns } from '@/features/admin/custom-attributes/dataTableColumns.js'
import CustomAttributesForm from '@/features/admin/custom-attributes/CustomAttributesForm.vue'
import { Button } from '@/components/ui/button'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { createFormSchema } from '@/features/admin/custom-attributes/formSchema.js'
import {
  Dialog, DialogContent, DialogDescription,
  DialogFooter, DialogHeader, DialogTitle, DialogTrigger
} from '@/components/ui/dialog'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useStorage } from '@vueuse/core'
import AdminPageWithHelp from '@/layouts/admin/AdminPageWithHelp.vue'
import AdminPageHeader from '@/components/layout/AdminPageHeader.vue'
import { useEmitter } from '@/composables/useEmitter'
import { EMITTER_EVENTS } from '@/constants/emitterEvents.js'
import { useAdminListRefresh } from '@/composables/useAdminListRefresh'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'
import api from '@/api'

const appliesTo = useStorage('appliesTo', 'contact')
const { t } = useI18n()
const customAttributes = ref([])
const isLoading = ref(false)
const emitter = useEmitter()
const dialogOpen = ref(false)
const isEditing = ref(false)
const { showErrorToast, showSuccessToast } = useAdminErrorToast()

const fetchAll = async () => {
  if (!appliesTo.value) return
  isLoading.value = true
  try {
    const resp = await api.getCustomAttributes(appliesTo.value)
    customAttributes.value = resp.data.data
  } catch (error) {
    showErrorToast(error)
  } finally {
    isLoading.value = false
  }
}

useAdminListRefresh('custom-attributes', fetchAll)

onMounted(() => {
  emitter.on(EMITTER_EVENTS.EDIT_MODEL, (data) => {
    if (data?.model === 'custom-attributes') {
      form.setValues(data.data)
      form.setErrors({})
      isEditing.value = true
      dialogOpen.value = true
    }
  })
})

onUnmounted(() => {
  emitter.off(EMITTER_EVENTS.EDIT_MODEL)
})

const newCustomAttribute = () => {
  form.resetForm()
  form.setErrors({})
  isEditing.value = false
}

const form = useForm({
  validationSchema: toTypedSchema(createFormSchema(t)),
  initialValues: { id: 0, name: '', data_type: 'text', applies_to: appliesTo.value, values: [] }
})

const onSubmit = form.handleSubmit(async (values) => {
  isLoading.value = true
  try {
    if (values.id) {
      await api.updateCustomAttribute(values.id, values)
      showSuccessToast(t('globals.messages.updatedSuccessfully', { name: t('globals.terms.customAttribute') }))
    } else {
      await api.createCustomAttribute(values)
      showSuccessToast(t('globals.messages.createdSuccessfully', { name: t('globals.terms.customAttribute') }))
    }
    dialogOpen.value = false
    fetchAll()
  } catch (error) {
    showErrorToast(error)
  } finally {
    isLoading.value = false
  }
})

watch(appliesTo, (newVal) => {
  form.resetForm({ values: { ...form.values, applies_to: newVal } })
  fetchAll()
}, { immediate: true })
</script>
