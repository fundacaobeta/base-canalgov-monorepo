<template>
  <AdminPageWithHelp>
    <template #content>
      <router-view />
    </template>

    <template #help>
      <div class="native-html">
        <p>Conecte sua conta de e-mail Google ou Microsoft, ou configure IMAP e SMTP manualmente.</p>
        <p>Cada caixa de entrada adicionada cria um novo canal de e-mail para receber mensagens dos clientes.</p>
        <p>
          <a href="https://docs.libredesk.io/configuration/connecting-inboxes" target="_blank">
            Saiba mais
          </a>
        </p>
      </div>
    </template>
  </AdminPageWithHelp>
</template>

<script setup>
import AdminPageWithHelp from '@/layouts/admin/AdminPageWithHelp.vue'
</script>
