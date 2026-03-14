<template>
  <div class="h-screen overflow-y-auto">
    <div class="mx-auto flex min-h-full max-w-7xl flex-col gap-6 p-6 sm:p-8 lg:flex-row">
      <aside class="w-full lg:max-w-xs">
        <div class="rounded-3xl border border-border/70 bg-background/95 p-5 shadow-sm">
          <div class="space-y-1">
            <p class="text-xs font-medium uppercase tracking-[0.18em] text-muted-foreground">
              {{ $t('globals.terms.account') }}
            </p>
            <h1 class="text-2xl font-semibold">{{ $t('account.accountArea') }}</h1>
            <p class="text-sm text-muted-foreground">
              {{ $t('account.accountAreaDescription') }}
            </p>
          </div>

          <nav class="mt-6 space-y-2">
            <router-link
              v-for="item in navItems"
              :key="item.name"
              :to="{ name: item.name }"
              class="flex items-start gap-3 rounded-2xl border px-4 py-3 transition-colors"
              :class="route.name === item.name
                ? 'border-foreground/20 bg-muted text-foreground'
                : 'border-transparent text-muted-foreground hover:border-border hover:bg-muted/40 hover:text-foreground'"
            >
              <component :is="item.icon" class="mt-0.5 h-4 w-4 shrink-0" />
              <div>
                <p class="text-sm font-medium">{{ item.label }}</p>
                <p class="text-xs text-muted-foreground">{{ item.description }}</p>
              </div>
            </router-link>
          </nav>
        </div>
      </aside>

      <div class="min-w-0 flex-1">
        <router-view class="flex-grow" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { Settings2, UserRound } from 'lucide-vue-next'

const route = useRoute()
const { t } = useI18n()

const navItems = computed(() => [
  {
    name: 'profile',
    icon: UserRound,
    label: t('account.editProfile'),
    description: t('account.nav.profile')
  },
  {
    name: 'account-preferences',
    icon: Settings2,
    label: t('account.appPreferences'),
    description: t('account.nav.preferences')
  }
])
</script>
