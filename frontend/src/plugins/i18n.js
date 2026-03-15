import { createI18n } from 'vue-i18n'

// Singleton instance created once — configured with messages in main.js before app mount.
// Exported so the router can call i18n.global.t() in beforeEach without needing the Vue app context.
const i18n = createI18n({
  legacy: false,
  locale: 'pt-BR',
  fallbackLocale: 'pt-BR',
  messages: {}
})

export function setupI18n(locale, messages) {
  i18n.global.locale.value = locale
  i18n.global.setLocaleMessage(locale, messages)
}

export default i18n
