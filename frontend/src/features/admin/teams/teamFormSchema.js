import * as z from 'zod'

export const teamFormSchema = z.object({
  name: z
    .string({
      required_error: 'O nome da equipe é obrigatório.'
    })
    .min(2, {
      message: 'O nome da equipe deve ter pelo menos 2 caracteres.'
    }),
  emoji: z.string({ required_error: 'O emoji é obrigatório.' }),
  conversation_assignment_type: z.string({ required_error: 'O tipo de atribuição da conversa é obrigatório.' }),
  max_auto_assigned_conversations: z.coerce.number().optional().default(0),
  timezone: z.string({ required_error: 'O fuso horário é obrigatório.' }),
  business_hours_id: z.number().optional().nullable(),
  sla_policy_id: z.number().optional().nullable(),
})
