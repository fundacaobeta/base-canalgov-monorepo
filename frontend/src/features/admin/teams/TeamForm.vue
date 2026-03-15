<template>
  <form @submit="onSubmit" class="space-y-6">
    <div class="box p-5 space-y-6">
      <FormField v-slot="{ componentField }" name="emoji">
        <FormItem>
          <FormLabel>{{ t('admin.team.form.emoji') }}</FormLabel>
          <FormControl>
            <div class="flex items-center space-x-2">
              <Input class="w-20 text-2xl text-center" v-bind="componentField" />
              <p class="text-sm text-muted-foreground">{{ t('admin.team.form.emojiDescription') }}</p>
            </div>
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>

      <FormField v-slot="{ componentField }" name="name">
        <FormItem>
          <FormLabel>{{ t('admin.team.form.name') }}</FormLabel>
          <FormControl>
            <Input type="text" placeholder="ex: Ouvidoria Geral" v-bind="componentField" />
          </FormControl>
          <FormDescription>{{ t('admin.team.form.nameDescription') }}</FormDescription>
          <FormMessage />
        </FormItem>
      </FormField>

      <div class="grid gap-6 md:grid-cols-2">
        <FormField v-slot="{ componentField }" name="auto_assign_type">
          <FormItem>
            <FormLabel>{{ t('admin.team.form.autoAssignType') }}</FormLabel>
            <FormControl>
              <Select v-bind="componentField" v-model="componentField.modelValue">
                <SelectTrigger>
                  <SelectValue :placeholder="t('globals.messages.select', { name: t('admin.team.form.autoAssignType') })" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="round_robin">Round Robin</SelectItem>
                  <SelectItem value="load_balanced">Load Balanced</SelectItem>
                  <SelectItem value="none">{{ t('globals.auth.none') }}</SelectItem>
                </SelectContent>
              </Select>
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="max_auto_assign_conversations">
          <FormItem>
            <FormLabel>{{ t('admin.team.form.maxAutoAssignConversations') }}</FormLabel>
            <FormControl>
              <Input type="number" v-bind="componentField" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>

      <div class="grid gap-6 md:grid-cols-2">
        <FormField v-slot="{ componentField }" name="timezone">
          <FormItem>
            <FormLabel>{{ t('admin.team.form.timezone') }}</FormLabel>
            <FormControl>
              <Select v-bind="componentField" v-model="componentField.modelValue">
                <SelectTrigger>
                  <SelectValue :placeholder="t('admin.general.timezone.placeholder')" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem v-for="(value, label) in timeZones" :key="value" :value="value">
                    {{ label }}
                  </SelectItem>
                </SelectContent>
              </Select>
            </FormControl>
            <FormDescription>{{ t('admin.team.form.timezoneDescription') }}</FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="business_hours_id">
          <FormItem>
            <FormLabel>{{ t('admin.team.form.businessHours') }}</FormLabel>
            <FormControl>
              <Select v-bind="componentField" v-model="componentField.modelValue">
                <SelectTrigger>
                  <SelectValue :placeholder="t('admin.general.businessHours.placeholder')" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem :value="0">{{ t('globals.auth.none') }}</SelectItem>
                  <SelectItem v-for="bh in businessHours" :key="bh.id" :value="bh.id">
                    {{ bh.name }}
                  </SelectItem>
                </SelectContent>
              </Select>
            </FormControl>
            <FormDescription>{{ t('admin.team.form.businessHoursDescription') }}</FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="sla_policy_id">
          <FormItem>
            <FormLabel>{{ t('admin.team.form.slaPolicy') }}</FormLabel>
            <FormControl>
              <Select v-bind="componentField" v-model="componentField.modelValue">
                <SelectTrigger>
                  <SelectValue :placeholder="t('globals.messages.select', { name: t('admin.team.form.slaPolicy') })" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem :value="0">{{ t('globals.auth.none') }}</SelectItem>
                  <SelectItem v-for="sla in slaPolicies" :key="sla.id" :value="sla.id">
                    {{ sla.name }}
                  </SelectItem>
                </SelectContent>
              </Select>
            </FormControl>
            <FormDescription>{{ t('admin.team.form.slaPolicyDescription') }}</FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>
    </div>

    <div class="flex items-center justify-end space-x-3">
      <Button type="submit" :isLoading="isLoading">{{ submitLabel }}</Button>
    </div>
  </form>
</template>

<script setup>
import { onMounted, ref, watch } from 'vue'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { createFormSchema } from './teamFormSchema.js'
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from '@/components/ui/form'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { timeZones } from '@/constants/timezones.js'
import { useI18n } from 'vue-i18n'
import api from '@/api'

const { showErrorToast } = useAdminErrorToast()
const { t } = useI18n()
const businessHours = ref([])
const slaPolicies = ref([])

const props = defineProps({
  initialValues: {
    type: Object,
    required: false
  },
  submitForm: {
    type: Function,
    required: true
  },
  submitLabel: {
    type: String,
    required: false,
    default: ''
  },
  isLoading: {
    type: Boolean,
    default: false
  }
})

const submitLabel = props.submitLabel || t('globals.messages.save')

const teamForm = useForm({
  validationSchema: toTypedSchema(createFormSchema(t)),
  initialValues: {
    emoji: '🏢',
    name: '',
    auto_assign_type: 'none',
    max_auto_assign_conversations: 0,
    timezone: 'America/Sao_Paulo',
    business_hours_id: 0,
    sla_policy_id: 0
  }
})

onMounted(() => {
  fetchBusinessHours()
  fetchSLAPolicies()
})

const fetchBusinessHours = async () => {
  try {
    const response = await api.getAllBusinessHours()
    businessHours.value = response.data.data
  } catch (error) {
    showErrorToast(error)
  }
}

const fetchSLAPolicies = async () => {
  try {
    const response = await api.getSLAPolicies()
    slaPolicies.value = response.data.data
  } catch (error) {
    showErrorToast(error)
  }
}

const onSubmit = teamForm.handleSubmit(async (values) => {
  await props.submitForm(values)
})

watch(
  () => props.initialValues,
  (newValues) => {
    if (!newValues || Object.keys(newValues).length === 0) return
    teamForm.setValues(newValues)
  },
  { deep: true, immediate: true }
)
</script>
