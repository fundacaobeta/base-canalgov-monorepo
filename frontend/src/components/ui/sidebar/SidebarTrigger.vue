<script setup>
import { ViewVerticalIcon } from '@radix-icons/vue';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { useSidebar } from './utils';
import { useI18n } from 'vue-i18n';

const props = defineProps({
  class: { type: null, required: false },
});

const { toggleSidebar } = useSidebar();
const { t } = useI18n();
</script>

<template>
  <Button
    data-sidebar="trigger"
    variant="ghost"
    size="icon"
    :class="cn('h-7 w-7', props.class)"
    @click="toggleSidebar"
  >
    <ViewVerticalIcon />
    <span class="sr-only">{{ t('globals.ui.toggleSidebar') }}</span>
  </Button>
</template>
