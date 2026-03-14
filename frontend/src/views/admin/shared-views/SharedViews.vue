<template>
  <AdminPageWithHelp>
    <template #content>
      <router-view></router-view>
    </template>
    <template #help>
      <p>Crie visões compartilhadas visíveis para todos os agentes ou para equipes específicas.</p>
    </template>
  </AdminPageWithHelp>
</template>

<script setup>
import AdminPageWithHelp from '@/layouts/admin/AdminPageWithHelp.vue'
</script>
