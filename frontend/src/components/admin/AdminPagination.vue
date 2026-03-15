<template>
  <div class="sticky bottom-0 bg-background border-t px-1 py-3 mt-auto">
    <div class="flex flex-col sm:flex-row items-center justify-between gap-4">
      <div class="flex items-center gap-2">
        <span class="text-sm text-muted-foreground">
          {{ $t('globals.terms.page') }} {{ page }} / {{ totalPages }}
        </span>
        <Select :model-value="String(perPage)" @update:model-value="(v) => $emit('per-page-change', Number(v))">
          <SelectTrigger class="h-8 w-[70px]">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="15">15</SelectItem>
            <SelectItem value="30">30</SelectItem>
            <SelectItem value="50">50</SelectItem>
            <SelectItem value="100">100</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <Pagination>
        <PaginationList class="flex items-center gap-1">
          <PaginationListItem>
            <PaginationFirst
              :class="{ 'cursor-not-allowed opacity-50': page === 1 }"
              @click.prevent="page > 1 ? $emit('go-to', 1) : null"
            />
          </PaginationListItem>
          <PaginationListItem>
            <PaginationPrev
              :class="{ 'cursor-not-allowed opacity-50': page === 1 }"
              @click.prevent="page > 1 ? $emit('go-to', page - 1) : null"
            />
          </PaginationListItem>
          <template v-for="pageNumber in visiblePages" :key="pageNumber">
            <PaginationListItem v-if="pageNumber === '...'">
              <PaginationEllipsis />
            </PaginationListItem>
            <PaginationListItem v-else>
              <Button
                :is-active="pageNumber === page"
                @click.prevent="$emit('go-to', pageNumber)"
                :variant="pageNumber === page ? 'default' : 'outline'"
              >
                {{ pageNumber }}
              </Button>
            </PaginationListItem>
          </template>
          <PaginationListItem>
            <PaginationNext
              :class="{ 'cursor-not-allowed opacity-50': page === totalPages }"
              @click.prevent="page < totalPages ? $emit('go-to', page + 1) : null"
            />
          </PaginationListItem>
          <PaginationListItem>
            <PaginationLast
              :class="{ 'cursor-not-allowed opacity-50': page === totalPages }"
              @click.prevent="page < totalPages ? $emit('go-to', totalPages) : null"
            />
          </PaginationListItem>
        </PaginationList>
      </Pagination>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { getVisiblePages } from '@/utils/pagination'
import { Button } from '@/components/ui/button'
import {
  Pagination,
  PaginationEllipsis,
  PaginationFirst,
  PaginationLast,
  PaginationList,
  PaginationListItem,
  PaginationNext,
  PaginationPrev
} from '@/components/ui/pagination'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'

const props = defineProps({
  page: { type: Number, required: true },
  perPage: { type: Number, required: true },
  totalPages: { type: Number, required: true }
})

defineEmits(['go-to', 'per-page-change'])

const visiblePages = computed(() => getVisiblePages(props.page, props.totalPages))
</script>
