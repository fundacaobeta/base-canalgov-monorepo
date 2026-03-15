<template>
  <div>
    <AdminPageHeader
      :title="$t('admin.sla.title')"
      :description="$t('admin.sla.description')"
      :breadcrumbs="[{ label: $t('globals.terms.administration'), to: '/admin' }, { label: $t('admin.sla.title') }]"
    >
      <template v-if="isListRoute" #actions>
        <RouterLink :to="{ name: 'new-sla' }">
          <Button>
            <Plus class="h-4 w-4 mr-1.5" />
            {{ $t('globals.messages.new', { name: $t('globals.terms.slaPolicy') }) }}
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
          <p>{{ $t('admin.sla.help') }}</p>
          <p>{{ $t('admin.sla.help2') }}</p>
        </template>
        <template v-else>
          <p>{{ $t('admin.sla.form.help') }}</p>
          <p>{{ $t('admin.sla.form.help2') }}</p>
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
const isListRoute = computed(() => route.name === 'sla-list')
</script>
