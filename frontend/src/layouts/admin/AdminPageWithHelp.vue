<template>
  <div class="grid grid-cols-1 lg:grid-cols-[1fr_280px] gap-8 items-start">
    <div class="min-w-0">
      <slot name="content" />
    </div>
    <aside class="help-panel self-start lg:sticky lg:top-4">
      <p class="help-panel-title">{{ $t('globals.terms.help') }}</p>
      <div class="space-y-3 text-sm text-muted-foreground">
        <slot name="help" />
      </div>
    </aside>
  </div>
</template>
