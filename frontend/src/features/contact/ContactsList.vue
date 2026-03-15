<template>
  <div class="flex flex-col gap-4 p-4 h-full overflow-hidden">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold tracking-tight">{{ t('globals.terms.citizen', 2) }}</h1>
      <div class="flex items-center gap-2">
        <!-- Segment Selector -->
        <Select v-model="selectedSegmentId" @update:model-value="handleSegmentChange">
          <SelectTrigger class="w-[200px] h-10">
            <SelectValue :placeholder="t('globals.terms.allGroups')" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem :value="null">{{ t('globals.terms.allGroups') }}</SelectItem>
            <SelectItem v-for="seg in segments" :key="seg.id" :value="seg.id.toString()">
              {{ seg.name }}
            </SelectItem>
          </SelectContent>
        </Select>

        <Input
          id="search-contacts"
          name="search-contacts"
          type="text"
          v-model="searchTerm"
          :placeholder="t('globals.terms.searchPlaceholder')"
          class="w-[300px]"
          @input="fetchContactsDebounced"
        />
        <Popover>
          <PopoverTrigger as-child>
            <Button variant="outline" size="sm" class="h-10">
              <Filter class="mr-2 h-4 w-4" />
              {{ t('globals.terms.filter', 2) }}
            </Button>
          </PopoverTrigger>
          <PopoverContent class="w-[240px] p-4 flex flex-col gap-4" align="end">
            <div class="space-y-2">
              <label class="text-xs font-medium uppercase text-muted-foreground tracking-wider">{{ t('globals.terms.sortBy') }}</label>
              <Select v-model="orderByField" @update:model-value="fetchContacts">
                <SelectTrigger class="h-9 w-full">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="users.created_at">{{ t('globals.terms.registeredAt') }}</SelectItem>
                  <SelectItem value="users.first_name">{{ t('globals.terms.firstName') }}</SelectItem>
                  <SelectItem value="users.email">{{ t('globals.terms.email') }}</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div class="space-y-2">
              <label class="text-xs font-medium uppercase text-muted-foreground tracking-wider">{{ t('globals.terms.order') }}</label>
              <Select v-model="orderByDirection" @update:model-value="fetchContacts">
                <SelectTrigger class="h-9 w-full">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="asc">{{ t('globals.terms.ascending') }}</SelectItem>
                  <SelectItem value="desc">{{ t('globals.terms.descending') }}</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </PopoverContent>
        </Popover>
      </div>
    </div>

    <div class="rounded-md border bg-card overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead class="w-[80px]"></TableHead>
            <TableHead>{{ t('globals.terms.name') }}</TableHead>
            <TableHead>{{ t('globals.terms.contact') }}</TableHead>
            <TableHead>{{ t('globals.terms.document', 2) }} (SUS/NIS)</TableHead>
            <TableHead>{{ t('globals.terms.registeredAt') }}</TableHead>
            <TableHead class="text-right">{{ t('globals.terms.actions') }}</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <template v-if="loading">
            <TableRow v-for="i in 10" :key="i">
              <TableCell><Skeleton class="h-10 w-10 rounded-full" /></TableCell>
              <TableCell><Skeleton class="h-4 w-[150px]" /></TableCell>
              <TableCell><Skeleton class="h-4 w-[200px]" /></TableCell>
              <TableCell><Skeleton class="h-4 w-[100px]" /></TableCell>
              <TableCell><Skeleton class="h-4 w-[120px]" /></TableCell>
              <TableCell class="text-right"><Skeleton class="h-8 w-8 ml-auto rounded-md" /></TableCell>
            </TableRow>
          </template>
          <template v-else-if="contacts.length > 0">
            <TableRow 
              v-for="contact in contacts" 
              :key="contact.id" 
              class="cursor-pointer hover:bg-muted/50 transition-colors"
              @click="$router.push({ name: 'contact-detail', params: { id: contact.id } })"
            >
              <TableCell>
                <Avatar class="h-9 w-9 border shadow-sm">
                  <AvatarImage :src="contact.avatar_url || ''" />
                  <AvatarFallback class="text-xs bg-primary/5 text-primary">
                    {{ getInitials(contact.first_name, contact.last_name) }}
                  </AvatarFallback>
                </Avatar>
              </TableCell>
              <TableCell>
                <div class="font-medium text-foreground">{{ contact.first_name }} {{ contact.last_name }}</div>
                <div v-if="!contact.enabled" class="text-[10px] uppercase font-bold text-destructive">{{ t('globals.terms.blocked') }}</div>
              </TableCell>
              <TableCell>
                <div class="flex flex-col text-sm">
                  <span class="text-muted-foreground flex items-center gap-1.5">
                    <Mail class="h-3 w-3" /> {{ contact.email }}
                  </span>
                  <span v-if="contact.phone_number" class="text-muted-foreground flex items-center gap-1.5">
                    <Phone class="h-3 w-3" /> {{ contact.phone_number }}
                  </span>
                </div>
              </TableCell>
              <TableCell>
                <div class="flex gap-1.5 flex-wrap">
                  <Badge v-if="contact.custom_attributes?.cartao_sus" variant="outline" class="text-[10px] px-1.5 py-0">
                    {{ t('globals.terms.sus') }}: {{ contact.custom_attributes.cartao_sus }}
                  </Badge>
                  <Badge v-if="contact.custom_attributes?.nis" variant="secondary" class="text-[10px] px-1.5 py-0 bg-blue-50 text-blue-700 border-blue-100">
                    {{ t('globals.terms.nis') }}: {{ contact.custom_attributes.nis }}
                  </Badge>
                  <span v-if="!contact.custom_attributes?.cartao_sus && !contact.custom_attributes?.nis" class="text-muted-foreground text-xs italic">
                    -
                  </span>
                </div>
              </TableCell>
              <TableCell class="text-sm text-muted-foreground">
                {{ formatDate(contact.created_at) }}
              </TableCell>
              <TableCell class="text-right">
                <Button variant="ghost" size="icon" class="h-8 w-8">
                  <ChevronRight class="h-4 w-4" />
                </Button>
              </TableCell>
            </TableRow>
          </template>
          <TableRow v-else>
            <TableCell colspan="6" class="h-32 text-center">
              <div class="flex flex-col items-center justify-center gap-1 text-muted-foreground">
                <Users class="h-8 w-8 mb-2 opacity-20" />
                <p class="font-medium">{{ t('globals.terms.noCitizenFound') }}</p>
                <p class="text-sm opacity-70">{{ t('globals.terms.adjustFilters') }}</p>
              </div>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>

    <!-- Pagination -->
    <div class="flex items-center justify-between mt-auto pt-4 border-t border-dashed">
      <div class="flex items-center gap-4">
        <span class="text-xs font-medium text-muted-foreground uppercase tracking-wider">
          {{ total }} {{ t('globals.terms.result', 2).toLowerCase() }} • {{ t('globals.terms.page') }} {{ page }} de {{ totalPages }}
        </span>
        <Select v-model="perPage" @update:model-value="handlePerPageChange">
          <SelectTrigger class="h-8 w-[75px] text-xs">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem :value="15">15</SelectItem>
            <SelectItem :value="30">30</SelectItem>
            <SelectItem :value="50">50</SelectItem>
            <SelectItem :value="100">100</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <Pagination v-if="totalPages > 1">
        <PaginationList class="flex items-center gap-1">
          <Button variant="outline" size="icon" class="h-8 w-8" :disabled="page === 1" @click="goToPage(1)">
            <ChevronsLeft class="h-4 w-4" />
          </Button>
          <Button variant="outline" size="icon" class="h-8 w-8" :disabled="page === 1" @click="goToPage(page - 1)">
            <ChevronLeft class="h-4 w-4" />
          </Button>
          
          <template v-for="p in visiblePages" :key="p">
            <Button v-if="p !== '...'" :variant="p === page ? 'default' : 'outline'" size="sm" class="h-8 w-8 text-xs" @click="goToPage(p)">
              {{ p }}
            </Button>
            <span v-else class="px-2 text-muted-foreground">...</span>
          </template>

          <Button variant="outline" size="icon" class="h-8 w-8" :disabled="page === totalPages" @click="goToPage(page + 1)">
            <ChevronRight class="h-4 w-4" />
          </Button>
          <Button variant="outline" size="icon" class="h-8 w-8" :disabled="page === totalPages" @click="goToPage(totalPages)">
            <ChevronsRight class="h-4 w-4" />
          </Button>
        </PaginationList>
      </Pagination>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Card } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from '@/components/ui/table'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Pagination, PaginationList } from '@/components/ui/pagination'
import { 
  Filter, Mail, Phone, Users, ChevronRight, 
  ChevronLeft, ChevronsLeft, ChevronsRight 
} from 'lucide-vue-next'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { useDebounceFn } from '@vueuse/core'
import { EMITTER_EVENTS } from '@/constants/emitterEvents.js'
import { useEmitter } from '@/composables/useEmitter'
import { useContactSegments } from '@/composables/useContactSegments'
import { handleHTTPError } from '@/utils/http'
import { getVisiblePages } from '@/utils/pagination'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import api from '@/api'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const contacts = ref([])
const loading = ref(false)
const page = ref(1)
const perPage = ref(15)
const totalPages = ref(0)
const searchTerm = ref('')
const selectedSegmentId = ref(route.query.segment_id ? route.query.segment_id.toString() : null)
const orderByField = ref('users.created_at')
const orderByDirection = ref('desc')
const total = ref(0)
const emitter = useEmitter()
const { segments, fetchSegments } = useContactSegments()

const visiblePages = computed(() => getVisiblePages(page.value, totalPages.value))

watch(() => route.query.segment_id, (newId) => {
  selectedSegmentId.value = newId ? newId.toString() : null
  fetchContacts()
})

const fetchContactsDebounced = useDebounceFn(() => {
  page.value = 1
  fetchContacts()
}, 400)

const handleSegmentChange = (val) => {
  router.push({ name: 'contacts', query: { segment_id: val || undefined } })
}

const fetchContacts = async () => {
  loading.value = true
  let filters = []
  
  if (searchTerm.value) {
    filters.push({
      model: 'users',
      field: 'first_name',
      operator: 'ilike',
      value: `%${searchTerm.value}%`
    })
  }

  if (selectedSegmentId.value) {
    filters.push({
      field: 'segment_id',
      operator: 'eq',
      value: selectedSegmentId.value.toString()
    })
  }

  try {
    const response = await api.getContacts({
      page: page.value,
      page_size: perPage.value,
      filters: filters.length > 0 ? JSON.stringify(filters) : '',
      order: orderByDirection.value,
      order_by: orderByField.value
    })
    contacts.value = response.data.data.results
    totalPages.value = response.data.data.total_pages
    total.value = response.data.data.total
  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  } finally {
    loading.value = false
  }
}

const formatDate = (date) => {
  if (!date) return '-'
  return format(new Date(date), 'dd/MM/yyyy', { locale: ptBR })
}

const getInitials = (firstName, lastName) => {
  return `${firstName?.[0] || ''}${lastName?.[0] || ''}`.toUpperCase()
}

const goToPage = (newPage) => {
  page.value = newPage
  fetchContacts()
}

const handlePerPageChange = (newPerPage) => {
  page.value = 1
  perPage.value = newPerPage
  fetchContacts()
}

onMounted(() => {
  fetchSegments()
  fetchContacts()
})
</script>
