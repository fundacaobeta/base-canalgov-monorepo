<template>
  <div class="min-h-screen flex flex-col">
    <div class="flex flex-wrap gap-4 pb-4">
      <div class="flex items-center gap-2 mb-4">
        <!-- Filter Popover -->
        <Popover :open="filtersOpen" @update:open="filtersOpen = $event">
          <PopoverTrigger @click="filtersOpen = !filtersOpen">
            <Button variant="outline" size="sm" class="flex items-center gap-2 h-8">
              <ListFilter size="14" />
              <span>{{ t('globals.terms.filter') }}</span>
              <span
                v-if="filters.length > 0"
                class="flex items-center justify-center bg-primary text-primary-foreground rounded-full size-4 text-xs"
              >
                {{ filters.length }}
              </span>
            </Button>
          </PopoverTrigger>
          <PopoverContent class="w-full p-4 flex flex-col gap-4">
            <div class="w-[32rem]">
              <FilterBuilder
                :fields="filterFields"
                :showButtons="true"
                v-model="filters"
                @apply="fetchActivityLogs"
                @clear="fetchActivityLogs"
              />
            </div>
          </PopoverContent>
        </Popover>

        <!-- Order By Popover -->
        <Popover>
          <PopoverTrigger>
            <Button variant="outline" size="sm" class="flex items-center h-8">
              <ArrowDownWideNarrow size="18" class="text-muted-foreground cursor-pointer" />
            </Button>
          </PopoverTrigger>
          <PopoverContent class="w-[200px] p-4 flex flex-col gap-4">
            <Select v-model="orderByField" @update:model-value="fetchActivityLogs">
              <SelectTrigger class="h-8 w-full">
                <SelectValue :placeholder="orderByField" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem :value="'activity_logs.created_at'">
                  {{ t('globals.terms.createdAt') }}
                </SelectItem>
              </SelectContent>
            </Select>

            <Select v-model="orderByDirection" @update:model-value="fetchActivityLogs">
              <SelectTrigger class="h-8 w-full">
                <SelectValue :placeholder="orderByDirection" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem :value="'asc'">{{ t('globals.messages.ascending') }}</SelectItem>
                <SelectItem :value="'desc'">{{ t('globals.messages.descending') }}</SelectItem>
              </SelectContent>
            </Select>
          </PopoverContent>
        </Popover>
      </div>

      <div class="w-full overflow-x-auto">
        <SimpleTable
          :headers="[
            t('globals.terms.name'),
            t('globals.terms.timestamp'),
            t('globals.terms.ipAddress')
          ]"
          :keys="['activity_description', 'created_at', 'ip']"
          :data="activityLogs"
          :showDelete="false"
          :loading="loading"
          :skeletonRows="15"
        />
      </div>
    </div>

    <AdminPagination
      :page="page"
      :per-page="perPage"
      :total-pages="totalPages"
      @go-to="goToPage"
      @per-page-change="handlePerPageChange"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import SimpleTable from '@/components/table/SimpleTable.vue'
import AdminPagination from '@/components/admin/AdminPagination.vue'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import FilterBuilder from '@/components/filter/FilterBuilder.vue'
import { Button } from '@/components/ui/button'
import { ListFilter, ArrowDownWideNarrow } from 'lucide-vue-next'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { useActivityLogFilters } from '@/composables/useActivityLogFilters'
import { usePagination } from '@/composables/usePagination'
import { useI18n } from 'vue-i18n'
import { format } from 'date-fns'
import api from '@/api'

const { t } = useI18n()
const activityLogs = ref([])
const loading = ref(true)
const orderByField = ref('activity_logs.created_at')
const orderByDirection = ref('desc')
const filters = ref([])
const filtersOpen = ref(false)
const { activityLogListFilters } = useActivityLogFilters()

const filterFields = computed(() =>
  Object.entries(activityLogListFilters.value).map(([field, value]) => ({
    model: 'activity_logs',
    label: value.label,
    field,
    type: value.type,
    operators: value.operators,
    options: value.options ?? []
  }))
)

async function fetchActivityLogs() {
  filtersOpen.value = false
  loading.value = true
  try {
    const resp = await api.getActivityLogs({
      page: page.value,
      page_size: perPage.value,
      filters: JSON.stringify(filters.value),
      order: orderByDirection.value,
      order_by: orderByField.value
    })
    const data = resp.data.data
    activityLogs.value = data.results.map((log) => ({
      ...log,
      created_at: format(new Date(log.created_at), 'PPpp')
    }))
    setPaginationMeta({ count: data.count, total_pages: data.total_pages })
  } catch (err) {
    console.error('Error fetching activity logs:', err)
    activityLogs.value = []
  } finally {
    loading.value = false
  }
}

const { page, perPage, totalPages, goToPage, handlePerPageChange, setPaginationMeta } =
  usePagination(fetchActivityLogs)

onMounted(fetchActivityLogs)
</script>
