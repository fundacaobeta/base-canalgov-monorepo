<template>
  <form @submit="onSubmit" class="space-y-6">
    <div class="box p-5 space-y-6">
      <h3 class="text-base font-medium">{{ t('admin.template.form.header') }}</h3>

      <FormField v-slot="{ componentField }" name="name">
        <FormItem>
          <FormLabel>{{ t('globals.terms.name') }}</FormLabel>
          <FormControl>
            <Input type="text" placeholder="ex: Saudação Inicial" v-bind="componentField" />
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>

      <FormField v-slot="{ componentField }" name="content">
        <FormItem>
          <FormLabel>{{ t('globals.terms.content') }}</FormLabel>
          <FormControl>
            <Textarea
              class="min-h-32"
              placeholder="ex: Olá, como posso ajudar hoje?"
              v-bind="componentField"
            />
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>

      <div class="grid gap-6 md:grid-cols-2">
        <FormField v-slot="{ componentField }" name="category_id">
          <FormItem>
            <FormLabel>{{ t('globals.terms.category') }}</FormLabel>
            <FormControl>
              <Select v-bind="componentField" v-model="componentField.modelValue">
                <SelectTrigger>
                  <SelectValue :placeholder="t('globals.messages.select', { name: t('globals.terms.category') })" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem v-for="cat in categories" :key="cat.id" :value="cat.id">
                    {{ cat.name }}
                  </SelectItem>
                </SelectContent>
              </Select>
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="team_id">
          <FormItem>
            <FormLabel>{{ t('admin.template.form.team') }}</FormLabel>
            <FormControl>
              <Select v-bind="componentField" v-model="componentField.modelValue">
                <SelectTrigger>
                  <SelectValue :placeholder="t('globals.messages.select', { name: t('globals.terms.team') })" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem :value="0">{{ t('admin.template.form.global') }}</SelectItem>
                  <SelectItem v-for="team in teams" :key="team.id" :value="team.id">
                    {{ team.name }}
                  </SelectItem>
                </SelectContent>
              </Select>
            </FormControl>
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
import { createFormSchema } from './formSchema.js'
import {
  FormControl,
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
import { Textarea } from '@/components/ui/textarea'
import { Button } from '@/components/ui/button'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'
import api from '@/api'

const { showErrorToast } = useAdminErrorToast()
const { t } = useI18n()
const categories = ref([])
const teams = ref([])

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

const templateForm = useForm({
  validationSchema: toTypedSchema(createFormSchema(t)),
  initialValues: {
    name: '',
    content: '',
    category_id: 0,
    team_id: 0
  }
})

onMounted(() => {
  fetchCategories()
  fetchTeams()
})

const fetchCategories = async () => {
  try {
    const response = await api.getTemplateCategories()
    categories.value = response.data.data
  } catch (error) {
    showErrorToast(error)
  }
}

const fetchTeams = async () => {
  try {
    const response = await api.getTeams()
    teams.value = response.data.data
  } catch (error) {
    showErrorToast(error)
  }
}

const onSubmit = templateForm.handleSubmit(async (values) => {
  await props.submitForm(values)
})

watch(
  () => props.initialValues,
  (newValues) => {
    if (!newValues || Object.keys(newValues).length === 0) return
    templateForm.setValues(newValues)
  },
  { deep: true, immediate: true }
)
</script>
