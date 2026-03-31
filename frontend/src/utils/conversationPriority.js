function normalizePriority(priority) {
  return String(priority || '').trim()
}

export function translateConversationPriority(priority) {
  const normalizedPriority = normalizePriority(priority)
  if (!normalizedPriority) return ''

  const priorityMap = {
    low: 'Baixa',
    medium: 'Média',
    high: 'Alta',
    urgent: 'Urgente'
  }

  return priorityMap[normalizedPriority.toLowerCase()] || normalizedPriority
}

export function getConversationPriorityBadgeVariant(priority) {
  const normalizedPriority = normalizePriority(priority).toLowerCase()

  const variantMap = {
    low: 'neutral',
    medium: 'info',
    high: 'warning',
    urgent: 'destructive'
  }

  return variantMap[normalizedPriority] || 'outline'
}
