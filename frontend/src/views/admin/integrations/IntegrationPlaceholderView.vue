<template>
  <div class="space-y-6">
    <div class="box p-6">
      <h2 class="text-xl font-semibold">{{ title }}</h2>
      <p class="mt-2 text-sm text-muted-foreground">
        {{ description }}
      </p>
    </div>

    <div class="grid gap-4 md:grid-cols-2">
      <div class="box p-5">
        <h3 class="font-medium">Status da integração</h3>
        <p class="mt-2 text-sm text-muted-foreground">
          Estrutura inicial criada para a integração de {{ title.toLowerCase() }}.
        </p>
      </div>
      <div class="box p-5">
        <h3 class="font-medium">Próximos passos</h3>
        <p class="mt-2 text-sm text-muted-foreground">
          Configure credenciais, webhooks ou endpoints específicos conforme a sua operação.
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true
  },
  description: {
    type: String,
    required: true
  }
})
</script>
