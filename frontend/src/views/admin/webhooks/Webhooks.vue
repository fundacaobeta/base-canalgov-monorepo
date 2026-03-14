<template>
  <AdminPageWithHelp>
    <template #content>
      <router-view />
    </template>

    <template #help>
      <p>Configure webhooks para receber notificações em tempo real quando eventos ocorrerem no seu ambiente CanalGov.</p>
      <p>Os webhooks permitem integrar o CanalGov com serviços externos enviando requisições HTTP POST quando eventos específicos acontecerem.</p>
      <a
        href="https://docs.libredesk.io/configuration/webhooks"
        target="_blank"
        rel="noopener noreferrer"
        class="link-style"
      >
        <p>Saiba mais</p>
      </a>
    </template>
  </AdminPageWithHelp>
</template>

<script setup>
import AdminPageWithHelp from '@/layouts/admin/AdminPageWithHelp.vue'
</script>
