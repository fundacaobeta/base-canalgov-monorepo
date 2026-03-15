import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { FIELD_OPERATORS, FIELD_TYPE } from '@/constants/filterConfig.js'

// Map FIELD_TYPE values → FIELD_OPERATORS arrays (mirrors the FIELD_OPERATORS keys)
const OPERATORS_BY_TYPE = {
  [FIELD_TYPE.TEXT]: FIELD_OPERATORS.TEXT,
  [FIELD_TYPE.SELECT]: FIELD_OPERATORS.SELECT,
  [FIELD_TYPE.MULTI_SELECT]: FIELD_OPERATORS.MULTI_SELECT,
  [FIELD_TYPE.BOOLEAN]: FIELD_OPERATORS.BOOLEAN,
  [FIELD_TYPE.DATE]: FIELD_OPERATORS.DATE,
  [FIELD_TYPE.NUMBER]: FIELD_OPERATORS.NUMBER,
}

export function useContactFilters() {
  const { t } = useI18n()

  const contactListFilters = computed(() => ({
    first_name: {
      label: t('globals.terms.firstName'),
      type: FIELD_TYPE.TEXT,
      operators: OPERATORS_BY_TYPE[FIELD_TYPE.TEXT]
    },
    last_name: {
      label: t('globals.terms.lastName'),
      type: FIELD_TYPE.TEXT,
      operators: OPERATORS_BY_TYPE[FIELD_TYPE.TEXT]
    },
    email: {
      label: t('globals.terms.email'),
      type: FIELD_TYPE.TEXT,
      operators: OPERATORS_BY_TYPE[FIELD_TYPE.TEXT]
    },
    phone_number: {
      label: t('globals.terms.phoneNumber'),
      type: FIELD_TYPE.TEXT,
      operators: OPERATORS_BY_TYPE[FIELD_TYPE.TEXT]
    },
    created_at: {
      label: t('globals.terms.createdAt'),
      type: FIELD_TYPE.DATE,
      operators: OPERATORS_BY_TYPE[FIELD_TYPE.DATE]
    }
    // Note: custom attributes are harder to filter statically, 
    // but we can add the most common ones here if needed.
  }))

  return {
    contactListFilters
  }
}
