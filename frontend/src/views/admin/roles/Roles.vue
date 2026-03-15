<template>
  <div>
    <AdminPageHeader
      :title="$t('admin.roles.title')"
      :description="$t('admin.roles.description')"
      :breadcrumbs="[{ label: $t('globals.terms.administration'), to: '/admin' }, { label: $t('admin.roles.title') }]"
    >
      <template v-if="isListRoute" #actions>
        <RouterLink :to="{ name: 'new-role' }">
          <Button>
            <Plus class="h-4 w-4 mr-1.5" />
            {{ $t('globals.messages.new', { name: $t('globals.terms.role') }) }}
          </Button>
        </RouterLink>
      </template>
    </AdminPageHeader>

    <AdminPageWithHelp>
      <template #content>
        <router-view></router-view>
      </template>

      <template #help>
        <template v-if="isListRoute">
          <p>{{ $t('admin.role.help') }}</p>
        </template>
        <template v-else>
          <p>{{ $t('admin.role.form.help') }}</p>
          <p>{{ $t('admin.role.form.help2') }}</p>
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
const isListRoute = computed(() => route.name === 'role-list')
</script>
