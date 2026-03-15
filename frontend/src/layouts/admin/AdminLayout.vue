<template>
  <div class="overflow-y-auto h-screen bg-muted/30">
    <div class="px-6 py-8 sm:px-10 sm:py-10 min-h-full flex flex-col max-w-6xl w-full mx-auto">
      <div
        v-if="$route.name === 'admin'"
        class="flex flex-col items-center justify-center flex-grow gap-3 text-center"
      >
        <div class="rounded-full bg-primary/8 p-4">
          <Settings class="h-8 w-8 text-primary/60" aria-hidden="true" />
        </div>
        <p class="text-sm font-medium text-foreground">{{ $t('admin.empty') }}</p>
        <p class="text-xs text-muted-foreground max-w-xs">{{ $t('admin.emptyHint') }}</p>
      </div>

      <RouterView v-slot="{ Component }">
        <Transition name="page" mode="out-in" appear>
          <component :is="Component" class="flex-grow" />
        </Transition>
      </RouterView>
    </div>
  </div>
</template>

<script setup>
import { Settings } from 'lucide-vue-next'
</script>

<style scoped>
.page-enter-active,
.page-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}
.page-enter-from {
  opacity: 0;
  transform: translateY(4px);
}
.page-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
