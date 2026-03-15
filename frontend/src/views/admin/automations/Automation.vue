<template>
  <div>
    <AdminPageHeader
      :title="$t('admin.automation.title')"
      :description="$t('admin.automation.description')"
      :breadcrumbs="[{ label: $t('globals.terms.administration'), to: '/admin' }, { label: $t('admin.automation.title') }]"
    >
      <template v-if="isListRoute" #actions>
        <Button @click="newRule">
          <Plus class="h-4 w-4 mr-1.5" aria-hidden="true" />
          {{
            $t('globals.messages.new', {
              name: $t('globals.terms.rule')
            })
          }}
        </Button>
      </template>
    </AdminPageHeader>

    <AdminPageWithHelp>
      <template #content>
        <div v-if="isListRoute">
          <div v-if="selectedTab">
            <AutomationTabs v-model:automationsTab="selectedTab" />
          </div>
        </div>
        <router-view />
      </template>

      <template #help>
        <p>{{ $t('admin.automation.help') }}</p>
        <p>{{ $t('admin.automation.help2') }}</p>
      </template>
    </AdminPageWithHelp>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { Plus } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { useRoute, useRouter } from 'vue-router'
import { useStorage } from '@vueuse/core'
import AutomationTabs from '@/features/admin/automation/AutomationTabs.vue'
import AdminPageWithHelp from '@/layouts/admin/AdminPageWithHelp.vue'
import AdminPageHeader from '@/components/layout/AdminPageHeader.vue'

const route = useRoute()
const router = useRouter()
const isListRoute = computed(() => route.name === 'automations')
const selectedTab = useStorage('automationsTab', 'new_conversation')
const newRule = () => {
  router.push({ name: 'new-automation', query: { type: selectedTab.value } })
}
</script>
