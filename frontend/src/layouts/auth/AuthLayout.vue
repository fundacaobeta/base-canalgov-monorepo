<template>
  <div class="min-h-screen flex flex-col bg-[#F9F9F9]">
    <main class="flex-1 flex items-center justify-center p-4">
      <div class="w-full max-w-[450px]">
        <slot></slot>
      </div>
    </main>
    <div class="flex justify-center mb-6">
      <span class="text-[0.75rem] text-center text-[#aaa]">
        Desenvolvido por CanalGov
      </span>
    </div>
  </div>
</template>

<script setup></script>
