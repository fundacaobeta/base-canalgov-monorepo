<template>
  <div>
    <AdminPageHeader
      :title="$t('admin.inbox.title')"
      :description="$t('admin.inbox.description')"
      :breadcrumbs="[{ label: $t('globals.terms.administration'), to: '/admin' }, { label: $t('admin.inbox.title') }]"
    >
      <template v-if="route.name === 'inbox-list'" #actions>
        <RouterLink :to="{ name: 'new-inbox' }">
          <Button>
            <Plus class="h-4 w-4 mr-1.5" aria-hidden="true" />
            {{ $t('globals.messages.new', { name: $t('globals.terms.inbox') }) }}
          </Button>
        </RouterLink>
      </template>
    </AdminPageHeader>

    <AdminPageWithHelp>
      <template #content>
        <router-view />
      </template>

      <template #help>
        <p>{{ $t('admin.inbox.help') }}</p>
        <p>{{ $t('admin.inbox.help2') }}</p>
        <a
          href="https://docs.canalgov.io/configuration/connecting-inboxes"
          target="_blank"
          rel="noopener noreferrer"
          class="link-style"
        >
          Saiba mais
        </a>
      </template>
    </AdminPageWithHelp>
  </div>
</template>

<script setup>
import { Plus } from 'lucide-vue-next'
import { useRoute } from 'vue-router'
import AdminPageWithHelp from '@/layouts/admin/AdminPageWithHelp.vue'
import AdminPageHeader from '@/components/layout/AdminPageHeader.vue'
import { Button } from '@/components/ui/button'

const route = useRoute()
</script>
