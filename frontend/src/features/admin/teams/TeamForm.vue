<template>
  <form @submit="onSubmit" class="space-y-6">
    <FormField name="emoji" v-slot="{ componentField }">
      <FormItem ref="emojiPickerContainer" class="relative">
        <FormLabel>Emoji</FormLabel>
        <FormControl>
          <Input type="text" v-bind="componentField" @click="toggleEmojiPicker" />
          <div v-if="isEmojiPickerVisible" class="absolute z-10 mt-2">
            <EmojiPicker :native="true" @select="onSelectEmoji" class="w-[300px]" />
          </div>
        </FormControl>
        <FormDescription>Emoji exibido para esta equipe.</FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="name">
      <FormItem>
        <FormLabel>Nome</FormLabel>
        <FormControl>
          <Input type="text" placeholder="Nome" v-bind="componentField" />
        </FormControl>
        <FormDescription>Selecione um nome único para a equipe.</FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField name="conversation_assignment_type" v-slot="{ componentField }">
      <FormItem>
        <FormLabel>Tipo de atribuição automática</FormLabel>
        <FormControl>
          <Select v-bind="componentField">
            <SelectTrigger>
              <SelectValue placeholder="Selecione um tipo de atribuição" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectItem v-for="at in assignmentTypes" :key="at" :value="at">
                  {{ at }}
                </SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
        </FormControl>
        <FormDescription>
          Round robin: as conversas são atribuídas aos membros da equipe em rodízio. <br />
          Manual: as conversas devem ser assumidas pelos membros da equipe.
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="max_auto_assigned_conversations">
      <FormItem>
        <FormLabel>Máximo de conversas autoatribuídas</FormLabel>
        <FormControl>
          <Input type="number" placeholder="0" v-bind="componentField" />
        </FormControl>
        <FormDescription>
          Número máximo de conversas que podem ser atribuídas automaticamente a um agente.
          Conversas com status "Resolvido" ou "Fechado" não contam para esse limite. Use 0 para ilimitado.
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="timezone">
      <FormItem>
        <FormLabel>Fuso horário</FormLabel>
        <FormControl>
          <Select v-bind="componentField">
            <SelectTrigger>
              <SelectValue placeholder="Selecione um fuso horário" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectItem v-for="(value, label) in timeZones" :key="value" :value="value">
                  {{ label }}
                </SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
        </FormControl>
        <FormDescription>O fuso horário da equipe será usado para calcular o SLA.</FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="business_hours_id">
      <FormItem>
        <FormLabel>Horário comercial</FormLabel>
        <FormControl>
          <Select v-bind="componentField">
            <SelectTrigger>
              <SelectValue placeholder="Selecione o horário comercial" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectItem :value = 0>Nenhum</SelectItem>
                <SelectItem v-for="bh in businessHours" :key="bh.id" :value="bh.id">
                  {{ bh.name }}
                </SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
        </FormControl>
        <FormDescription
          >Horário comercial padrão da equipe; será usado para calcular o SLA.</FormDescription
        >
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="sla_policy_id">
      <FormItem>
        <FormLabel>Política de SLA</FormLabel>
        <FormControl>
          <Select v-bind="componentField">
            <SelectTrigger>
              <SelectValue placeholder="Selecione a política" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectItem :value= 0>Nenhuma</SelectItem>
                <SelectItem
                  v-for="sla in slaStore.options"
                  :key="sla.value"
                  :value="parseInt(sla.value)"
                >
                  {{ sla.label }}
                </SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
        </FormControl>
        <FormDescription
          >Política de SLA aplicada automaticamente às conversas quando forem atribuídas a esta equipe.</FormDescription
        >
        <FormMessage />
      </FormItem>
    </FormField>

    <Button type="submit" :isLoading="isLoading"> {{ submitLabel }} </Button>
  </form>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { onClickOutside } from '@vueuse/core'
import { Button } from '@/components/ui/button'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { teamFormSchema } from './teamFormSchema.js'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  FormDescription
} from '@/components/ui/form'
import { EMITTER_EVENTS } from '@/constants/emitterEvents.js'
import { useEmitter } from '@/composables/useEmitter'
import { Input } from '@/components/ui/input'
import EmojiPicker from 'vue3-emoji-picker'
import 'vue3-emoji-picker/css'
import { handleHTTPError } from '@/utils/http'
import { useSlaStore } from '@/stores/sla'
import { timeZones } from '@/constants/timezones.js'
import api from '@/api'

const emitter = useEmitter()
const slaStore = useSlaStore()
const assignmentTypes = ['Round robin', 'Manual']
const businessHours = ref([])

const props = defineProps({
  initialValues: { type: Object, required: false },
  submitForm: { type: Function, required: true },
  submitLabel: { type: String, default: 'Enviar' },
  isLoading: { type: Boolean }
})

const form = useForm({
  validationSchema: toTypedSchema(teamFormSchema)
})

const isEmojiPickerVisible = ref(false)
const emojiPickerContainer = ref(null)

onMounted(() => {
  fetchBusinessHours()
  onClickOutside(emojiPickerContainer, () => {
    isEmojiPickerVisible.value = false
  })
})

const fetchBusinessHours = async () => {
  try {
    const response = await api.getAllBusinessHours()
    businessHours.value = response.data.data
  } catch (error) {
    // If unauthorized (no permission), show a toast message.
    const toastPayload =
      error.response.status === 403
        ? {
            title: 'Não autorizado',
            variant: 'destructive',
            description: 'Você não tem permissão para visualizar horários comerciais.'
          }
        : {
            title: 'Não foi possível buscar horários comerciais',
            variant: 'destructive',
            description: handleHTTPError(error).message
          }
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, toastPayload)
  }
}

const onSubmit = form.handleSubmit((values) => {
  props.submitForm({
    ...values,
    business_hours_id: values.business_hours_id > 0 ? values.business_hours_id : null,
    sla_policy_id: values.sla_policy_id > 0 ? values.sla_policy_id: null
  })
})

watch(
  () => props.initialValues,
  (newValues) => {
    if (Object.keys(newValues).length === 0) return
    form.setValues(newValues)
  },
  { immediate: true }
)

function toggleEmojiPicker() {
  isEmojiPickerVisible.value = !isEmojiPickerVisible.value
}

function onSelectEmoji(emoji) {
  form.setFieldValue('emoji', emoji.i || emoji)
}
</script>
