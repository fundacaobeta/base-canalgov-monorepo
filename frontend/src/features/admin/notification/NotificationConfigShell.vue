<template>
  <div class="space-y-6">
    <div class="box p-6">
      <div class="flex flex-col gap-6 xl:flex-row xl:items-start xl:justify-between">
        <div class="space-y-2">
          <h2 class="text-xl font-semibold">{{ title }}</h2>
          <p class="max-w-3xl text-sm text-muted-foreground">{{ description }}</p>
        </div>
        <div class="min-w-72 rounded-lg border bg-muted/30 p-4">
          <div class="text-sm font-medium">{{ statusLabel }}</div>
          <p class="mt-2 text-sm text-muted-foreground">{{ statusDescription }}</p>
        </div>
      </div>
    </div>

    <Collapsible v-model:open="helpOpen" class="box p-5">
      <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
        <div class="space-y-1">
          <h3 class="font-medium">Ajuda integrada</h3>
          <p class="text-sm text-muted-foreground">
            Referências rápidas para configurar este canal com segurança e consistência operacional.
          </p>
        </div>
        <CollapsibleTrigger as-child>
          <Button type="button" variant="outline">
            {{ helpOpen ? 'Ocultar ajuda' : 'Abrir ajuda' }}
          </Button>
        </CollapsibleTrigger>
      </div>

      <CollapsibleContent class="pt-5">
        <div class="grid gap-3 lg:grid-cols-2">
          <div
            v-for="item in helpItems"
            :key="item.title"
            class="rounded border border-border bg-card p-4"
          >
            <div class="text-sm font-medium">{{ item.title }}</div>
            <p class="mt-2 text-sm text-muted-foreground">{{ item.description }}</p>
          </div>
        </div>
      </CollapsibleContent>
    </Collapsible>

    <slot />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Button } from '@/components/ui/button'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'

const helpOpen = ref(false)

defineProps({
  title: {
    type: String,
    required: true
  },
  description: {
    type: String,
    required: true
  },
  statusLabel: {
    type: String,
    required: true
  },
  statusDescription: {
    type: String,
    required: true
  },
  helpItems: {
    type: Array,
    required: true
  }
})
</script>
