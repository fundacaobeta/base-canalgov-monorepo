<template>
  <div class="text-center">
    <div class="text-muted-foreground text-sm">
      {{ message.content }}
      <DateTimeMeta
        v-if="message.updated_at"
        :value="message.updated_at"
        :show-absolute="false"
        inline
        compact
      />
    </div>
  </div>
</template>

<script setup>
import DateTimeMeta from '@/components/datetime/DateTimeMeta.vue'

defineProps({
  message: Object
})
</script>
