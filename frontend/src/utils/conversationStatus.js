function getFallbackLocale() {
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

  return getFallbackLabels(locale)[status] || status
}

export function getConversationStatusBadgeVariant(status) {
  const variantMap = {
    Open: 'destructive',
    Snoozed: 'warning',
    Resolved: 'success',
    Closed: 'neutral',
    Replied: 'info'
  }

  return variantMap[status] || 'outline'
}
