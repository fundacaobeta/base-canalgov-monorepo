<template>
  <div
    class="min-h-screen flex flex-col bg-background text-foreground"
    style="background-image: radial-gradient(circle at top left, hsl(var(--primary) / 0.16), transparent 30%), radial-gradient(circle at top right, hsl(var(--accent) / 0.18), transparent 28%), linear-gradient(180deg, hsl(var(--background)), hsl(var(--secondary) / 0.75));"
  >
    <main class="flex-1 flex items-center justify-center p-4">
      <div class="w-full max-w-[450px]">
        <slot></slot>
      </div>
    </main>
    <div class="flex justify-center mb-6">
      <span class="text-[0.75rem] text-center text-muted-foreground">
        Desenvolvido por CanalGov
      </span>
    </div>
  </div>
</template>

<script setup></script>
