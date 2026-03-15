import * as z from 'zod'

export const createFormSchema = (t) => z.object({
  name: z
    .string({
      required_error: t('admin.inbox.form.validation.required')
    })
    .min(2, {
      message: t('admin.team.form.validation.nameMinLength')
    }),
  emoji: z.string({ required_error: t('admin.inbox.form.validation.required') }),
  auto_assign_type: z.string({ required_error: t('admin.inbox.form.validation.required') }),
  max_auto_assign_conversations: z.coerce.number().optional().default(0),
  timezone: z.string({ required_error: t('admin.inbox.form.validation.required') }),
  business_hours_id: z.number().optional().nullable(),
  sla_policy_id: z.number().optional().nullable(),
})
