<template>
  <AuthLayout>
    <Card
      class="bg-card box border-primary/20 shadow-2xl shadow-primary/10 overflow-hidden"
      :class="{ 'animate-shake': shakeCard }"
      id="login-container"
      ref="cardRef"
    >
      <div class="h-2 w-full bg-[linear-gradient(90deg,hsl(var(--primary)),hsl(var(--accent)))]"></div>
      <CardContent class="p-6 space-y-6">
        <div class="space-y-3 text-center">
          <div class="flex flex-col items-center gap-3">
            <div class="flex h-16 w-16 items-center justify-center rounded-2xl border border-primary/20 bg-white/80 shadow-lg shadow-primary/10">
              <img
                src="/images/beta-logo.png"
                alt="Beta"
                class="h-12 w-12 object-contain"
              />
            </div>
            <img
              v-if="appLogoUrl"
              :src="appLogoUrl"
              :alt="appSettingsStore.public_config?.['app.site_name'] || 'CANALGOV'"
              class="mx-auto h-20 w-auto object-contain"
            />
          </div>
          <div class="inline-flex items-center rounded-full border border-accent/40 bg-accent/15 px-3 py-1 text-xs font-semibold uppercase tracking-[0.18em] text-primary">
            Identidade Beta
          </div>
          <CardTitle class="text-3xl font-bold text-foreground">
            {{ appSettingsStore.public_config?.['app.site_name'] || 'CANALGOV' }}
          </CardTitle>
          <p class="text-muted-foreground">{{ t('auth.signIn') }}</p>
        </div>

        <div v-if="enabledOIDCProviders.length" class="space-y-4">
          <Button
            v-for="oidcProvider in enabledOIDCProviders"
            :key="oidcProvider.id"
            variant="outline"
            type="button"
            @click="redirectToOIDC(oidcProvider)"
            class="w-full bg-card hover:bg-secondary text-foreground border-border rounded py-2 transition-all duration-200 ease-in-out transform hover:scale-105"
          >
            <img
              :src="oidcProvider.logo_url"
              :alt="oidcProvider.name"
              width="20"
              v-if="oidcProvider.logo_url"
            />
            {{ oidcProvider.name }}
          </Button>

          <div class="relative">
            <div class="absolute inset-0 flex items-center">
              <span class="w-full border-t border-border"></span>
            </div>
            <div class="relative flex justify-center text-xs uppercase">
              <span class="px-2 text-muted-foreground bg-card">{{ t('auth.orContinueWith') }}</span>
            </div>
          </div>
        </div>

        <form @submit.prevent="loginAction" class="space-y-4">
          <div class="space-y-2">
            <Label for="email" class="text-sm font-medium text-foreground">{{
              t('globals.terms.email')
            }}</Label>
            <Input
              id="email"
              type="text"
              autocomplete="username"
              :placeholder="t('auth.enterEmail')"
              v-model.trim="loginForm.email"
              :class="{ 'border-destructive': emailHasError }"
              class="w-full bg-card border-border text-foreground placeholder:text-muted-foreground rounded py-2 px-3 focus:ring-2 focus:ring-ring focus:border-ring transition-all duration-200 ease-in-out"
            />
          </div>

          <div class="space-y-2">
            <Label for="password" class="text-sm font-medium text-foreground">
              {{ t('globals.terms.password') }}
            </Label>
            <div class="relative">
              <Input
                id="password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="current-password"
                :placeholder="t('auth.enterPassword')"
                v-model="loginForm.password"
                :class="{ 'border-destructive': passwordHasError }"
                class="w-full bg-card border-border text-foreground placeholder:text-muted-foreground rounded py-2 px-3 pr-10 focus:ring-2 focus:ring-ring focus:border-ring transition-all duration-200 ease-in-out"
              />
              <button
                type="button"
                class="absolute inset-y-0 right-0 flex items-center pr-3 text-muted-foreground hover:text-foreground"
                @click="showPassword = !showPassword"
              >
                <Eye v-if="!showPassword" class="w-5 h-5" />
                <EyeOff v-else class="w-5 h-5" />
              </button>
            </div>
          </div>

          <div class="flex items-center justify-between">
            <router-link
              :to="{ name: 'reset-password' }"
              class="text-sm text-primary hover:text-primary/80 transition-all duration-200 ease-in-out"
            >
              {{ t('auth.forgotPassword') }}
            </router-link>
          </div>

          <Button
            class="w-full bg-primary hover:bg-primary/90 text-primary-foreground rounded py-2 transition-all duration-200 ease-in-out transform hover:scale-105"
            :disabled="isLoading"
            type="submit"
          >
            <span v-if="isLoading" class="flex items-center justify-center">
              <div
                class="w-5 h-5 border-2 border-primary-foreground/30 border-t-primary-foreground rounded-full animate-spin mr-3"
              ></div>
              {{ t('auth.loggingIn') }}
            </span>
            <span v-else>{{ t('auth.signInButton') }}</span>
          </Button>
        </form>

        <Error
          v-if="errorMessage"
          :errorMessage="errorMessage"
          :border="true"
          class="w-full bg-destructive/10 text-destructive border-destructive/20 p-3 rounded text-sm"
        />
      </CardContent>
    </Card>
  </AuthLayout>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { handleHTTPError } from '@/utils/http'
import api from '@/api'
import { validateEmail } from '@/utils/strings'
import { useTemporaryClass } from '@/composables/useTemporaryClass'
import { Button } from '@/components/ui/button'
import { Error } from '@/components/ui/error'
import { Card, CardContent, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useEmitter } from '@/composables/useEmitter'
import { useUserStore } from '@/stores/user'
import { useI18n } from 'vue-i18n'
import { EMITTER_EVENTS } from '@/constants/emitterEvents.js'
import { useAppSettingsStore } from '@/stores/appSettings'
import AuthLayout from '@/layouts/auth/AuthLayout.vue'
import { Eye, EyeOff } from 'lucide-vue-next'

const emitter = useEmitter()
const { t } = useI18n()
const errorMessage = ref('')
const isLoading = ref(false)
const router = useRouter()
const userStore = useUserStore()
const shakeCard = ref(false)
const showPassword = ref(false)
const loginForm = ref({
  email: '',
  password: ''
})
const oidcProviders = ref([])
const appSettingsStore = useAppSettingsStore()
const appLogoUrl = computed(() => appSettingsStore.public_config?.['app.logo_url'] || '')
// Demo build has the credentials prefilled.
const isDemoBuild = import.meta.env.VITE_DEMO_BUILD === 'true'

const demoCredentials = {
  email: 'demo@canalgov.local',
  password: 'demo@canalgov.local'
}

onMounted(async () => {
  // Prefill the login form with demo credentials if it's a demo build
  if (isDemoBuild) {
    loginForm.value.email = demoCredentials.email
    loginForm.value.password = demoCredentials.password
  }
  fetchOIDCProviders()
})

const fetchOIDCProviders = async () => {
  try {
    const config = appSettingsStore.public_config
    if (config && config['app.sso_providers']) {
      oidcProviders.value = config['app.sso_providers'] || []
    }
  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  }
}

const redirectToOIDC = (provider) => {
  // Pass the 'next' parameter to OIDC login if it exists
  const nextParam = router.currentRoute.value.query.next
  const url = nextParam
    ? `/api/v1/oidc/${provider.id}/login?next=${encodeURIComponent(nextParam)}`
    : `/api/v1/oidc/${provider.id}/login`
  window.location.href = url
}

const validateForm = () => {
  if (!validateEmail(loginForm.value.email) && loginForm.value.email !== 'System') {
    errorMessage.value = t('globals.messages.invalidEmailAddress')
    useTemporaryClass('login-container', 'animate-shake')
    return false
  }
  if (!loginForm.value.password) {
    errorMessage.value = t('globals.messages.cannotBeEmpty', {
      name: t('globals.terms.password')
    })
    useTemporaryClass('login-container', 'animate-shake')
    return false
  }
  return true
}

const loginAction = () => {
  if (!validateForm()) return

  errorMessage.value = ''
  isLoading.value = true

  api
    .login({
      email: loginForm.value.email,
      password: loginForm.value.password
    })
    .then((resp) => {
      if (resp?.data?.data) {
        userStore.setCurrentUser(resp.data.data)
      }
      // Also fetch general setting as user's logged in.
      appSettingsStore.fetchSettings('general')

      const nextParam = router.currentRoute.value.query.next
      if (typeof nextParam === 'string' && nextParam.length > 0) {
        window.location.href = nextParam
        return
      }

      router.push({ name: 'inboxes' })
    })
    .catch((error) => {
      errorMessage.value = handleHTTPError(error).message
      useTemporaryClass('login-container', 'animate-shake')
    })
    .finally(() => {
      isLoading.value = false
    })
}

const enabledOIDCProviders = computed(() => {
  return oidcProviders.value.filter((provider) => !provider.disabled)
})

const emailHasError = computed(() => {
  const email = loginForm.value.email
  return email !== 'System' && !validateEmail(email) && email !== ''
})

const passwordHasError = computed(
  () => !loginForm.value.password && loginForm.value.password !== ''
)
</script>
