<template>
  <AdminPageWithHelp>
    <template #content>
      <router-view />
    </template>

    <template #help>
      <p>Configure o login único com um ou mais provedores OpenID Connect.</p>
      <a
        href="https://docs.libredesk.io/configuration/sso"
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
