<template>
  <div>
    <AdminPageHeader
      :title="$t('admin.businessHours.title')"
      :description="$t('admin.businessHours.description')"
      :breadcrumbs="[{ label: $t('globals.terms.administration'), to: '/admin' }, { label: $t('admin.businessHours.title') }]"
    >
      <template v-if="isListRoute" #actions>
        <RouterLink :to="{ name: 'new-business-hours' }">
          <Button>
            <Plus class="h-4 w-4 mr-1.5" />
            {{ $t('globals.messages.new', { name: $t('globals.terms.businessHour') }) }}
          </Button>
        </RouterLink>
      </template>
    </AdminPageHeader>

    <AdminPageWithHelp>
      <template #content>
        <router-view />
      </template>

      <template #help>
        <template v-if="isListRoute">
          <p>{{ $t('admin.businessHours.help') }}</p>
          <p>{{ $t('admin.businessHours.help2') }}</p>
        </template>
        <template v-else>
          <p>{{ $t('admin.businessHours.form.help') }}</p>
          <p>{{ $t('admin.businessHours.form.help2') }}</p>
        </template>
      </template>
    </AdminPageWithHelp>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { Plus } from 'lucide-vue-next'
import { useRoute } from 'vue-router'
import AdminPageWithHelp from '@/layouts/admin/AdminPageWithHelp.vue'
import AdminPageHeader from '@/components/layout/AdminPageHeader.vue'
import { Button } from '@/components/ui/button'

const route = useRoute()
const isListRoute = computed(() => route.name === 'business-hours-list')
</script>
