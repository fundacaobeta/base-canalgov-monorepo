function getFallbackLocale() {
  if (typeof document !== 'undefined' && document.documentElement.lang) {
    return document.documentElement.lang
  }
  if (typeof navigator !== 'undefined' && navigator.language) {
    return navigator.language
  }
  return 'pt-BR'
}

function getFallbackLabels(locale) {
  const isPortuguese = locale?.toLowerCase().startsWith('pt')
  return isPortuguese
    ? {
        Open: 'Aberto',
        Snoozed: 'Adiado',
        Resolved: 'Resolvido',
        Closed: 'Fechado',
        Replied: 'Respondido'
      }
    : {
        Open: 'Open',
        Snoozed: 'Snoozed',
        Resolved: 'Resolved',
        Closed: 'Closed',
        Replied: 'Replied'
      }
}

export function translateConversationStatus(status, t, locale = getFallbackLocale()) {
  if (!status) return ''

  if (typeof t === 'function') {
    const map = {
      Open: t('globals.terms.open'),
      Snoozed: t('globals.terms.snoozed'),
      Resolved: t('globals.terms.resolved'),
      Closed: t('globals.terms.closed'),
      Replied: t('globals.terms.replied')
    }

    return map[status] || status
  }

  return getFallbackLabels(locale)[status] || status
}
