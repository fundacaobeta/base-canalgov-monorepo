import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { useAppSettingsStore } from './stores/appSettings'
import { createAppRouter } from './router'
import i18n, { setupI18n } from './plugins/i18n'
import mitt from 'mitt'
import api from './api'
import './assets/styles/main.scss'
import './utils/strings.js'
import Root from './Root.vue'

const DEFAULT_FAVICON = '/images/beta-logo.png?v=20260331'

const setFavicon = (url) => {
  let link = document.querySelector("link[rel='icon']")
  if (!link) {
    link = document.createElement("link")
    link.rel = "icon"
    document.head.appendChild(link)
  }
  link.href = url
}

async function initApp () {
  const config = (await api.getConfig()).data.data
  const emitter = mitt()
  const lang = config['app.lang'] || 'pt-BR'
  const langMessages = await api.getLanguage(lang)

  // Set favicon.
  setFavicon(config['app.favicon_url'] || DEFAULT_FAVICON)

  // Configure the i18n singleton with the loaded messages before mounting.
  setupI18n(lang, langMessages.data)

  // Create the router with locale-aware URL paths.
  const router = createAppRouter(lang)

  const app = createApp(Root)
  const pinia = createPinia()
  app.use(pinia)

  // Fetch and store app settings in store (after pinia is initialized)
  const settingsStore = useAppSettingsStore()

  // Store the public config in the store
  settingsStore.setPublicConfig(config)

  try {
    await settingsStore.fetchSettings('general')
  } catch (error) {
    // Pass
  }

  // Add emitter to global properties.
  app.config.globalProperties.emitter = emitter

  app.use(router)
  app.use(i18n)
  app.mount('#app')
}

initApp().catch((error) => {
  console.error('Error initializing app: ', error)
})
