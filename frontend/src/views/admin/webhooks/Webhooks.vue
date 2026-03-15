<template>
  <div>
    <AdminPageHeader
      :title="$t('admin.webhook.title')"
      :description="$t('admin.webhook.description')"
      :breadcrumbs="[{ label: $t('globals.terms.administration'), to: '/admin' }, { label: $t('admin.webhook.title') }]"
    >
      <template v-if="route.name === 'webhook-list'" #actions>
        <RouterLink :to="{ name: 'new-webhook' }">
          <Button>
            <Plus class="h-4 w-4 mr-1.5" aria-hidden="true" />
            {{ $t('globals.messages.new', { name: $t('globals.terms.webhook') }) }}
          </Button>
        </RouterLink>
      </template>
    </AdminPageHeader>

    <AdminPageWithHelp>
      <template #content>
        <router-view />
      </template>

      <template #help>
        <template v-if="route.name === 'webhook-list'">
          <p>{{ $t('admin.webhook.help') }}</p>
          <p>{{ $t('admin.webhook.help2') }}</p>
        </template>
        <template v-else>
          <p>{{ $t('admin.webhook.form.help') }}</p>
          <p>{{ $t('admin.webhook.form.help2') }}</p>
        </template>
        <a
          href="https://docs.canalgov.io/configuration/webhooks"
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
