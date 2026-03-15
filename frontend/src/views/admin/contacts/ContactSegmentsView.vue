<template>
  <div>
    <AdminPageHeader
      :title="t('admin.contactSegment.title')"
      :description="t('admin.contactSegment.description')"
      :breadcrumbs="[{ label: $t('globals.terms.administration'), to: '/admin' }, { label: t('admin.contactSegment.title') }]"
    >
      <template #actions>
        <Button @click="openCreateDialog">
          <Plus class="h-4 w-4 mr-1.5" aria-hidden="true" />
          {{ t('globals.messages.new', { name: t('globals.terms.group') }) }}
        </Button>
      </template>
    </AdminPageHeader>

    <AdminPageWithHelp>
      <template #content>
        <!-- Loading skeletons -->
        <div v-if="loading && !segments?.length" class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          <Card v-for="i in 3" :key="i" class="flex flex-col">
            <CardHeader>
              <Skeleton class="h-6 w-3/4" />
              <Skeleton class="h-4 w-1/2" />
            </CardHeader>
            <CardContent><Skeleton class="h-20 w-full" /></CardContent>
          </Card>
        </div>

        <!-- Empty state -->
        <div
          v-else-if="!loading && !segments?.length"
          class="flex flex-col items-center justify-center py-20 border-2 border-dashed rounded-xl"
        >
          <Users class="h-10 w-10 text-muted-foreground/30 mb-3" />
          <p class="text-sm font-medium text-muted-foreground">{{ t('admin.contactSegment.empty') }}</p>
          <Button variant="link" class="mt-1" @click="openCreateDialog">
            {{ t('admin.contactSegment.createFirst') }}
          </Button>
        </div>

        <!-- Card grid -->
        <div v-else class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          <Card
            v-for="segment in segments"
            :key="segment.id"
            class="flex flex-col group hover:border-primary/40 transition-colors"
          >
            <CardHeader>
              <div class="flex items-center justify-between">
                <CardTitle class="text-base">{{ segment.name }}</CardTitle>
                <Users class="h-4 w-4 text-muted-foreground group-hover:text-primary transition-colors" />
              </div>
              <CardDescription class="line-clamp-2 min-h-[40px]">
                {{ segment.description || t('globals.messages.noData') }}
              </CardDescription>
            </CardHeader>
            <CardContent class="flex-1">
              <div class="flex flex-wrap gap-1">
                <Badge
                  v-for="(f, i) in getFiltersList(segment.filters)"
                  :key="i"
                  variant="secondary"
                  class="text-[10px]"
                >
                  {{ f.label }}: {{ f.operator }}
                </Badge>
              </div>
            </CardContent>
            <CardFooter class="flex justify-end gap-2 border-t bg-muted/5 pt-4">
              <Button variant="ghost" size="sm" @click="openEditDialog(segment)" :aria-label="t('globals.messages.edit', { name: segment.name })">
                <Edit class="h-4 w-4" />
              </Button>
              <Button
                variant="ghost"
                size="sm"
                class="text-destructive hover:bg-destructive/10"
                @click="confirmDelete(segment)"
                :aria-label="t('globals.messages.delete', { name: segment.name })"
              >
                <Trash2 class="h-4 w-4" />
              </Button>
            </CardFooter>
          </Card>
        </div>
      </template>

      <template #help>
        <p>{{ t('admin.contactSegment.help') }}</p>
        <p>{{ t('admin.contactSegment.help2') }}</p>
      </template>
    </AdminPageWithHelp>

    <!-- Create / Edit Dialog -->
    <Dialog v-model:open="isDialogOpen">
      <DialogContent class="sm:max-w-[650px] max-h-[90vh] flex flex-col p-0">
        <DialogHeader class="p-6 pb-0">
          <DialogTitle>{{ editMode ? t('admin.contactSegment.editTitle') : t('admin.contactSegment.newTitle') }}</DialogTitle>
          <DialogDescription>
            {{ t('admin.contactSegment.filterDescription') }}
          </DialogDescription>
        </DialogHeader>

        <div class="flex-1 overflow-y-auto p-6 space-y-6">
          <div class="space-y-4">
            <div class="grid gap-2">
              <Label for="name">{{ t('admin.contactSegment.form.name') }}</Label>
              <Input id="name" v-model="form.name" :placeholder="t('admin.contactSegment.form.namePlaceholder')" />
            </div>
            <div class="grid gap-2">
              <Label for="desc">{{ t('admin.contactSegment.form.description') }}</Label>
              <Input id="desc" v-model="form.description" :placeholder="t('admin.contactSegment.form.description')" />
            </div>
          </div>

          <div class="space-y-4 pt-4 border-t">
            <Label class="text-sm font-bold uppercase tracking-widest text-primary flex items-center gap-2">
              <Filter class="h-4 w-4" />
              {{ t('admin.contactSegment.filterLabel') }}
            </Label>
            <FilterBuilder
              :fields="filterFields"
              v-model="form.filters"
              class="border rounded-lg p-4 bg-muted/10"
            />
          </div>
        </div>

        <DialogFooter class="p-6 pt-0 border-t mt-0">
          <Button variant="ghost" @click="isDialogOpen = false" :disabled="saving">
            {{ t('globals.messages.cancel') }}
          </Button>
          <Button @click="handleSave" :disabled="saving || !form.name">
            <Loader2 v-if="saving" class="mr-2 h-4 w-4 animate-spin" />
            {{ t('globals.messages.save') }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Delete Confirmation -->
    <AlertDialog v-model:open="deleteDialogOpen">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{{ t('globals.messages.areYouAbsolutelySure') }}</AlertDialogTitle>
          <AlertDialogDescription>
            {{ t('globals.messages.deletionConfirmation', { name: segmentToDelete?.name || t('globals.terms.group') }) }}
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>{{ t('globals.messages.cancel') }}</AlertDialogCancel>
          <AlertDialogAction
            class="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            @click="handleDelete"
          >
            {{ t('globals.messages.delete') }}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive, computed } from 'vue'
import { Plus, Users, Edit, Trash2, Filter, Loader2 } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import AdminPageHeader from '@/components/layout/AdminPageHeader.vue'
import AdminPageWithHelp from '@/layouts/admin/AdminPageWithHelp.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { useContactSegments } from '@/composables/useContactSegments'
import { useContactFilters } from '@/composables/useContactFilters'
import FilterBuilder from '@/components/filter/FilterBuilder.vue'
import {
  Dialog, DialogContent, DialogDescription,
  DialogFooter, DialogHeader, DialogTitle
} from '@/components/ui/dialog'
import {
  AlertDialog, AlertDialogAction, AlertDialogCancel,
  AlertDialogContent, AlertDialogDescription,
  AlertDialogFooter, AlertDialogHeader, AlertDialogTitle
} from '@/components/ui/alert-dialog'

const { t } = useI18n()
const { segments, loading, fetchSegments, createSegment, updateSegment, deleteSegment } = useContactSegments()
const { contactListFilters } = useContactFilters()

const isDialogOpen = ref(false)
const deleteDialogOpen = ref(false)
const editMode = ref(false)
const selectedId = ref(null)
const segmentToDelete = ref(null)
const saving = ref(false)

const filterFields = computed(() =>
  Object.entries(contactListFilters.value).map(([field, value]) => ({
    model: 'users',
    label: value.label,
    field,
    type: value.type,
    operators: value.operators,
    options: value.options ?? []
  }))
)

const form = reactive({ name: '', description: '', filters: [] })

const openCreateDialog = () => {
  editMode.value = false
  form.name = ''
  form.description = ''
  form.filters = []
  isDialogOpen.value = true
}

const openEditDialog = (segment) => {
  editMode.value = true
  selectedId.value = segment.id
  form.name = segment.name
  form.description = segment.description
  form.filters = Array.isArray(segment.filters) ? segment.filters : []
  isDialogOpen.value = true
}

const confirmDelete = (segment) => {
  segmentToDelete.value = segment
  deleteDialogOpen.value = true
}

const getFiltersList = (filters) => {
  if (!Array.isArray(filters)) return []
  return filters.map((f) => ({
    label: contactListFilters.value[f.field]?.label || f.field,
    operator: f.operator
  }))
}

const handleSave = async () => {
  try {
    saving.value = true
    const data = { name: form.name, description: form.description, filters: form.filters }
    if (editMode.value) {
      await updateSegment(selectedId.value, data)
    } else {
      await createSegment(data)
    }
    isDialogOpen.value = false
    await fetchSegments()
  } catch (err) {
    console.error(err)
  } finally {
    saving.value = false
  }
}

const handleDelete = async () => {
  if (!segmentToDelete.value) return
  await deleteSegment(segmentToDelete.value.id)
  segmentToDelete.value = null
  await fetchSegments()
}

onMounted(fetchSegments)
</script>
