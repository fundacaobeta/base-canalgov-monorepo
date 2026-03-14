<template>
  <AdminPageWithHelp>
    <template #content>
      <router-view></router-view>
    </template>

    <template #help>
      <p>{{ $t('admin.team.help') }}</p>
      <p>{{ $t('admin.team.help2') }}</p>
    </template>
  </AdminPageWithHelp>
</template>

<script setup>
import AdminPageWithHelp from '@/layouts/admin/AdminPageWithHelp.vue'
</script>
