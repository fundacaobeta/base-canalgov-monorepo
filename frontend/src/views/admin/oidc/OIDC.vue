<template>
  <div>
    <AdminPageHeader
      :title="$t('admin.oidc.title')"
      :description="$t('admin.oidc.description')"
      :breadcrumbs="[{ label: $t('globals.terms.administration'), to: '/admin' }, { label: $t('admin.oidc.title') }]"
    >
      <template v-if="isListRoute" #actions>
        <RouterLink :to="{ name: 'new-sso' }">
          <Button>
            <Plus class="h-4 w-4 mr-1.5" />
            {{ $t('globals.messages.new', { name: $t('globals.terms.sso') }) }}
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
          <p>{{ $t('admin.oidc.help') }}</p>
          <a
            href="https://docs.canalgov.io/configuration/sso"
            target="_blank"
            rel="noopener noreferrer"
            class="link-style"
          >
            {{ $t('globals.messages.learnMore') }}
          </a>
        </template>
        <template v-else>
          <p>{{ $t('admin.oidc.form.help') }}</p>
          <p>{{ $t('admin.oidc.form.help2') }}</p>
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
const isListRoute = computed(() => route.name === 'sso-list')
</script>
