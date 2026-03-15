<template>
  <form @submit="onSubmit" class="max-w-2xl space-y-0">
    <!-- Identidade -->
    <div class="form-section">
      <h2 class="form-section-title">{{ t('admin.general.section.identity') }}</h2>

      <FormField v-slot="{ field }" name="site_name">
        <FormItem>
          <FormLabel>{{ t('admin.general.siteName') }}</FormLabel>
          <FormControl>
            <Input type="text" placeholder="ex: Prefeitura de São Paulo" v-bind="field" />
          </FormControl>
          <FormDescription>
            {{ t('admin.general.siteName.description') }}
          </FormDescription>
          <FormMessage />
        </FormItem>
      </FormField>

      <FormField v-slot="{ field }" name="logo_url" :value="props.initialValues.logo_url">
        <FormItem>
          <FormLabel>{{ t('admin.general.logoURL') }}</FormLabel>
          <FormControl>
            <Input type="text" placeholder="https://..." v-bind="field" />
          </FormControl>
          <FormDescription>{{ t('admin.general.logoURL.description') }}</FormDescription>
          <FormMessage />
        </FormItem>
      </FormField>

      <FormField v-slot="{ field }" name="favicon_url" :value="props.initialValues.favicon_url">
        <FormItem>
          <FormLabel>{{ t('admin.general.faviconURL') }}</FormLabel>
          <FormControl>
            <Input type="text" placeholder="https://..." v-bind="field" />
          </FormControl>
          <FormDescription>{{ t('admin.general.faviconURL.description') }}</FormDescription>
          <FormMessage />
        </FormItem>
      </FormField>
    </div>

    <!-- Localização -->
    <div class="form-section">
      <h2 class="form-section-title">{{ t('admin.general.section.localization') }}</h2>

      <FormField v-slot="{ componentField }" name="lang">
        <FormItem>
          <FormLabel>{{ t('globals.terms.language') }}</FormLabel>
          <FormControl>
            <Select v-bind="componentField" :modelValue="componentField.modelValue">
              <SelectTrigger>
                <SelectValue :placeholder="t('admin.general.language.placeholder')" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectItem value="pt-BR">{{ t('globals.languages.pt-BR') }}</SelectItem>
                  <SelectItem value="da">{{ t('globals.languages.da') }}</SelectItem>
                  <SelectItem value="de">{{ t('globals.languages.de') }}</SelectItem>
                  <SelectItem value="en">{{ t('globals.languages.en') }}</SelectItem>
                  <SelectItem value="es">{{ t('globals.languages.es') }}</SelectItem>
                  <SelectItem value="fa">{{ t('globals.languages.fa') }}</SelectItem>
                  <SelectItem value="fr">{{ t('globals.languages.fr') }}</SelectItem>
                  <SelectItem value="it">{{ t('globals.languages.it') }}</SelectItem>
                  <SelectItem value="ja">{{ t('globals.languages.ja') }}</SelectItem>
                  <SelectItem value="mr">{{ t('globals.languages.mr') }}</SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
          </FormControl>
          <FormDescription>
            {{ t('admin.general.language.description') }}
          </FormDescription>
          <FormMessage />
        </FormItem>
      </FormField>

      <FormField v-slot="{ componentField }" name="timezone">
        <FormItem>
          <FormLabel>
            {{ t('globals.terms.timezone') }}
          </FormLabel>
          <FormControl>
            <Select v-bind="componentField">
              <SelectTrigger>
                <SelectValue :placeholder="t('admin.general.timezone.placeholder')" />
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
          <FormDescription>
            {{ t('admin.general.timezone.description') }}
          </FormDescription>
          <FormMessage />
        </FormItem>
      </FormField>
    </div>

    <!-- Atendimento -->
    <div class="form-section">
      <h2 class="form-section-title">{{ t('admin.general.section.service') }}</h2>

      <FormField v-slot="{ componentField }" name="business_hours_id">
        <FormItem>
          <FormLabel>
            {{ t('globals.terms.businessHour', 2) }}
          </FormLabel>
          <FormControl>
            <Select v-bind="componentField">
              <SelectTrigger>
                <SelectValue :placeholder="t('admin.general.businessHours.placeholder')" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectItem v-for="bh in businessHours" :key="bh.id" :value="bh.id">
                    {{ bh.name }}
                  </SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
          </FormControl>
          <FormDescription>
            {{ t('admin.general.businessHours.description') }}
          </FormDescription>
          <FormMessage />
        </FormItem>
      </FormField>

      <FormField v-slot="{ field }" name="root_url">
        <FormItem>
          <FormLabel>
            {{ t('globals.terms.rootURL') }}
          </FormLabel>
          <FormControl>
            <Input type="text" placeholder="https://atendimento.prefeitura.gov.br" v-bind="field" />
          </FormControl>
          <FormDescription>
            {{ t('admin.general.rootURL.description') }}
          </FormDescription>
          <FormMessage />
        </FormItem>
      </FormField>
    </div>

    <!-- Arquivos e Uploads -->
    <div class="form-section">
      <h2 class="form-section-title">{{ t('admin.general.section.uploads') }}</h2>

      <FormField
        v-slot="{ field }"
        name="max_file_upload_size"
        :value="props.initialValues.max_file_upload_size"
      >
        <FormItem>
          <FormLabel>
            {{ t('admin.general.maxAllowedFileUploadSize') }}
          </FormLabel>
          <FormControl>
            <div class="flex items-center gap-2">
              <Input type="number" placeholder="10" class="w-32" v-bind="field" />
              <span>{{ t('globals.units.mb') }}</span>
            </div>
          </FormControl>
          <FormDescription>
            {{ t('admin.general.maxAllowedFileUploadSize.description') }}
          </FormDescription>
          <FormMessage />
        </FormItem>
      </FormField>

      <FormField name="allowed_file_upload_extensions" v-slot="{ componentField, handleChange }">
        <FormItem>
          <FormLabel>
            {{ t('admin.general.allowedFileUploadExtensions') }}
          </FormLabel>
          <FormControl>
            <TagsInput :modelValue="componentField.modelValue" @update:modelValue="handleChange">
              <TagsInputItem v-for="item in componentField.modelValue" :key="item" :value="item">
                <TagsInputItemText />
                <TagsInputItemDelete />
              </TagsInputItem>
              <TagsInputInput placeholder="jpg" />
            </TagsInput>
          </FormControl>
          <FormDescription>
            {{ t('admin.general.allowedFileUploadExtensions.description') }}
          </FormDescription>
          <FormMessage />
        </FormItem>
      </FormField>
    </div>

    <div class="pt-6 border-t border-border">
      <Button type="submit" :isLoading="formLoading">{{ submitLabel }}</Button>
    </div>
  </form>
</template>

<script setup>
import { watch, ref, onMounted } from 'vue'
import { Button } from '@/components/ui/button'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { createFormSchema } from './formSchema.js'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  FormDescription
} from '@/components/ui/form'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import {
  TagsInput,
  TagsInputInput,
  TagsInputItem,
  TagsInputItemDelete,
  TagsInputItemText
} from '@/components/ui/tags-input'
import { Input } from '@/components/ui/input'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { timeZones } from '@/constants/timezones.js'
import { useI18n } from 'vue-i18n'
import api from '@/api'

const { showErrorToast, showSuccessToast } = useAdminErrorToast()
const { t } = useI18n()
const businessHours = ref({})
const formLoading = ref(false)
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
const form = useForm({
  validationSchema: toTypedSchema(createFormSchema(t))
})

onMounted(() => {
  fetchBusinessHours()
})

const fetchBusinessHours = async () => {
  try {
    const response = await api.getAllBusinessHours()
    // Convert business hours id to string
    response.data.data.forEach((bh) => {
      bh.id = bh.id.toString()
    })
    businessHours.value = response.data.data
  } catch (error) {
    showErrorToast(error)
  }
}

const onSubmit = form.handleSubmit(async (values) => {
  try {
    formLoading.value = true
    await props.submitForm(values)
    showSuccessToast(t('globals.messages.updatedSuccessfully', { name: t('globals.terms.setting', 2) }))
  } catch (error) {
    showErrorToast(error)
  } finally {
    formLoading.value = false
  }
})

// Watch for changes in initialValues and update the form.
watch(
  () => props.initialValues,
  (newValues) => {
    if (Object.keys(newValues).length === 0) {
      return
    }
    // Convert business hours id to string
    if (newValues.business_hours_id)
      newValues.business_hours_id = newValues.business_hours_id.toString()
    form.setValues(newValues)
  },
  { deep: true }
)
</script>
