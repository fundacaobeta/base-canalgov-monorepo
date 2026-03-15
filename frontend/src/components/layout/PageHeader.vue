<template>
  <div v-if="!isHidden">
    <div class="flex items-center space-x-4 h-12 px-2">
      <SidebarTrigger class="cursor-pointer" />
      <span class="text-xl font-semibold">
        {{ title }}
      </span>
    </div>
    <Separator />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { Separator } from '@/components/ui/separator'
import { SidebarTrigger } from '@/components/ui/sidebar'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'

const route = useRoute()
const { t } = useI18n()
const title = computed(() => (route.meta.title ? t(route.meta.title) : ''))
const isHidden = computed(() => route.meta.hidePageHeader === true)
</script>
