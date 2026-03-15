<template>
  <div>
    <AdminPageHeader
      :title="$t('admin.activityLog.title')"
      :description="$t('admin.activityLog.description')"
      :breadcrumbs="[{ label: $t('globals.terms.administration'), to: '/admin' }, { label: $t('admin.activityLog.title') }]"
    />

    <AdminPageWithHelp>
      <template #content>
        <ActivityLog />
      </template>

      <template #help>
        <p>{{ $t('admin.activityLog.help') }}</p>
        <p>{{ $t('admin.activityLog.help2') }}</p>
      </template>
    </AdminPageWithHelp>
  </div>
</template>

<script setup>
import ActivityLog from '@/features/admin/activity-log/ActivityLog.vue'
import AdminPageWithHelp from '@/layouts/admin/AdminPageWithHelp.vue'
import AdminPageHeader from '@/components/layout/AdminPageHeader.vue'
</script>
