<template>
  <div class="flex flex-col gap-6 p-6 h-full overflow-y-auto">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-4">
        <Button variant="outline" size="icon" @click="$router.push({ name: 'contacts' })">
          <ArrowLeft class="h-4 w-4" />
        </Button>
        <div>
          <h1 class="text-2xl font-bold tracking-tight text-foreground">
            {{ contact?.first_name }} {{ contact?.last_name }}
          </h1>
          <div class="flex items-center gap-2 text-muted-foreground">
            <Mail class="h-3 w-3" />
            <span class="text-sm">{{ contact?.email }}</span>
          </div>
        </div>
      </div>
      <div class="flex items-center gap-2">
        <Badge v-if="!contact?.enabled" variant="destructive">{{ t('globals.terms.blocked') }}</Badge>
        <Badge v-else variant="outline" class="bg-green-50 text-green-700 border-green-200 uppercase text-[10px] font-bold">
          {{ t('globals.terms.active') }}
        </Badge>
      </div>
    </div>

    <!-- Stats Grid -->
    <div class="grid gap-4 md:grid-cols-4">
      <Card class="bg-card/50 shadow-none border-dashed">
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-[10px] font-bold uppercase tracking-widest text-muted-foreground">
            {{ t('reports.conversations.total') }}
          </CardTitle>
          <MessageSquare class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.total_conversations || 0 }}</div>
        </CardContent>
      </Card>
      
      <Card class="bg-card/50 shadow-none border-dashed">
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-[10px] font-bold uppercase tracking-widest text-muted-foreground">
            {{ t('globals.status.open') }}
          </CardTitle>
          <CircleDot class="h-4 w-4 text-orange-500" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.open_conversations || 0 }}</div>
        </CardContent>
      </Card>

      <Card class="bg-card/50 shadow-none border-dashed">
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-[10px] font-bold uppercase tracking-widest text-muted-foreground">
            {{ t('globals.status.resolved') }}
          </CardTitle>
          <CheckCircle2 class="h-4 w-4 text-green-500" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.resolved_conversations || 0 }}</div>
        </CardContent>
      </Card>

      <Card class="bg-card/50 shadow-none border-dashed">
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-[10px] font-bold uppercase tracking-widest text-muted-foreground">
            {{ t('globals.terms.phoneNumber') }}
          </CardTitle>
          <Phone class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-md font-bold truncate">{{ contact?.phone_number || '-' }}</div>
        </CardContent>
      </Card>
    </div>

    <!-- Tabs Container -->
    <Tabs defaultValue="conversations" class="w-full flex flex-col flex-1">
      <TabsList class="w-full justify-start bg-transparent border-b rounded-none h-auto p-0 gap-6">
        <TabsTrigger 
          value="conversations" 
          class="data-[state=active]:border-primary data-[state=active]:bg-transparent border-b-2 border-transparent rounded-none px-2 pb-3 pt-0 shadow-none"
        >
          {{ t('globals.terms.conversations') }}
        </TabsTrigger>
        <TabsTrigger 
          value="info" 
          class="data-[state=active]:border-primary data-[state=active]:bg-transparent border-b-2 border-transparent rounded-none px-2 pb-3 pt-0 shadow-none"
        >
          {{ t('globals.terms.information') }}
        </TabsTrigger>
        <TabsTrigger 
          value="notes" 
          class="data-[state=active]:border-primary data-[state=active]:bg-transparent border-b-2 border-transparent rounded-none px-2 pb-3 pt-0 shadow-none"
        >
          {{ t('globals.terms.notes') }}
        </TabsTrigger>
      </TabsList>

      <!-- Content -->
      <div class="mt-6 flex-1">
        <!-- Tab: Conversas -->
        <TabsContent value="conversations" class="space-y-4 outline-none">
          <div v-if="conversations.length > 0" class="grid gap-3">
            <router-link 
              v-for="conv in conversations" 
              :key="conv.uuid" 
              :to="{ name: 'inbox-conversation', params: { type: inboxTypeParam('all'), uuid: conv.uuid } }"
              class="group block no-underline"
            >
              <div class="flex items-center justify-between p-4 rounded-lg border bg-card hover:bg-accent/50 hover:border-primary/30 transition-all">
                <div class="space-y-1">
                  <div class="flex items-center gap-2">
                    <div class="h-2 w-2 rounded-full" :class="conv.status === 'Open' ? 'bg-orange-500' : 'bg-green-500'"></div>
                    <span class="font-semibold text-foreground group-hover:text-primary">
                      {{ conv.subject || '(' + t('globals.messages.noSubject') + ')' }}
                    </span>
                  </div>
                  <div class="flex items-center gap-2 text-xs text-muted-foreground">
                    <span class="font-medium">#{{ conv.reference_number }}</span>
                    <span>•</span>
                    <span class="truncate max-w-[400px]">{{ conv.last_message }}</span>
                    <span>•</span>
                    <span>{{ formatDate(conv.last_message_at) }}</span>
                  </div>
                </div>
                <ChevronRight class="h-4 w-4 text-muted-foreground group-hover:text-primary transition-colors" />
              </div>
            </router-link>
          </div>
          <div v-else class="flex flex-col items-center justify-center py-20 border border-dashed rounded-xl bg-muted/5">
            <MessageSquare class="h-10 w-10 text-muted-foreground/20 mb-3" />
            <p class="text-muted-foreground text-sm">{{ t('globals.messages.noData') }}</p>
          </div>
        </TabsContent>

        <!-- Tab: Informações -->
        <TabsContent value="info" class="space-y-6 outline-none">
          <div class="grid gap-6 md:grid-cols-2">
            <!-- Basic Data Card -->
            <Card class="shadow-sm">
              <CardHeader>
                <CardTitle class="text-lg flex items-center gap-2">
                  <User class="h-5 w-5 text-primary" />
                  {{ t('globals.terms.citizenData') }}
                </CardTitle>
              </CardHeader>
              <CardContent class="grid gap-4">
                <div class="grid grid-cols-2 gap-4 border-b pb-4 border-dashed">
                  <div>
                    <p class="text-[10px] font-bold uppercase text-muted-foreground tracking-widest mb-1">{{ t('globals.terms.email') }}</p>
                    <p class="text-sm font-medium">{{ contact?.email }}</p>
                  </div>
                  <div>
                    <p class="text-[10px] font-bold uppercase text-muted-foreground tracking-widest mb-1">{{ t('globals.terms.phoneNumber') }}</p>
                    <p class="text-sm font-medium">{{ contact?.phone_number || '-' }}</p>
                  </div>
                </div>
                <div class="grid grid-cols-2 gap-4">
                  <div>
                    <p class="text-[10px] font-bold uppercase text-muted-foreground tracking-widest mb-1">{{ t('globals.terms.registration') }}</p>
                    <p class="text-sm font-medium">{{ formatDate(contact?.created_at) }}</p>
                  </div>
                  <div>
                    <p class="text-[10px] font-bold uppercase text-muted-foreground tracking-widest mb-1">{{ t('globals.terms.update') }}</p>
                    <p class="text-sm font-medium">{{ formatDate(contact?.updated_at) }}</p>
                  </div>
                </div>
              </CardContent>
            </Card>

            <!-- Custom Attributes Card -->
            <Card class="shadow-sm">
              <CardHeader class="flex flex-row items-center justify-between space-y-0">
                <CardTitle class="text-lg flex items-center gap-2">
                  <Settings class="h-5 w-5 text-primary" />
                  {{ t('globals.terms.customAttributes') }}
                </CardTitle>
                <router-link :to="{ name: 'custom-attributes' }">
                  <Button variant="ghost" size="sm" class="h-8 text-primary">
                    <PlusCircle class="h-4 w-4 mr-1" />
                    {{ t('globals.terms.manage') }}
                  </Button>
                </router-link>
              </CardHeader>
              <CardContent>
                <div v-if="hasAttributes" class="grid gap-3">
                  <div v-for="(value, key) in parsedAttributes" :key="key" class="flex items-center justify-between p-2 rounded bg-muted/30 border border-transparent hover:border-border transition-colors">
                    <span class="text-xs font-semibold text-muted-foreground uppercase">{{ key.replace(/_/g, ' ') }}</span>
                    <span class="text-sm font-bold text-foreground">{{ value }}</span>
                  </div>
                </div>
                <div v-else class="text-center py-10">
                  <Settings class="h-10 w-10 text-muted-foreground/10 mx-auto mb-2" />
                  <p class="text-xs text-muted-foreground italic">{{ t('globals.terms.noAttributesFilled') }}</p>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <!-- Tab: Notas -->
        <TabsContent value="notes" class="space-y-4 outline-none">
          <div class="space-y-4 mb-6">
            <div class="flex items-center justify-between">
              <h3 class="text-sm font-semibold text-muted-foreground uppercase tracking-wider">{{ t('globals.terms.newObservation') }}</h3>
              
              <!-- Template Selector for Notes -->
              <Popover>
                <PopoverTrigger as-child>
                  <Button variant="ghost" size="sm" class="h-8 text-primary">
                    <FileText class="h-4 w-4 mr-1" />
                    {{ t('globals.terms.useTemplate') }}
                  </Button>
                </PopoverTrigger>
                <PopoverContent class="w-[300px] p-0" align="end">
                  <div class="p-2 border-b bg-muted/20">
                    <Input :placeholder="t('globals.messages.filter', { name: t('globals.terms.template', 2).toLowerCase() })" class="h-8 text-xs" v-model="templateSearch" />
                  </div>
                  <div class="max-h-[300px] overflow-y-auto">
                    <div v-for="category in filteredCategories" :key="category.id" class="p-2">
                      <p class="text-[10px] font-bold uppercase text-muted-foreground px-2 py-1">{{ category.name }}</p>
                      <div class="space-y-1">
                        <Button 
                          v-for="tpl in category.templates" 
                          :key="tpl.id"
                          variant="ghost" 
                          class="w-full justify-start text-xs h-auto py-2 px-2 text-left block"
                          @click="applyTemplate(tpl.body)"
                        >
                          <div class="font-medium truncate">{{ tpl.name }}</div>
                          <div class="text-[10px] text-muted-foreground truncate opacity-70">{{ tpl.body }}</div>
                        </Button>
                      </div>
                    </div>
                    <div v-if="filteredCategories.length === 0" class="p-8 text-center text-muted-foreground text-xs italic">
                      {{ t('globals.messages.noData') }}
                    </div>
                  </div>
                </PopoverContent>
              </Popover>
            </div>

            <div class="flex gap-2">
              <Textarea 
                v-model="newNote" 
                :placeholder="t('globals.terms.internalNotePlaceholder')" 
                class="min-h-[100px] resize-none text-sm"
              />
            </div>
            <div class="flex justify-end">
              <Button @click="handleAddNote" :disabled="!newNote.trim()">{{ t('globals.messages.add', { name: '' }) }}</Button>
            </div>
          </div>
          
          <div v-if="notes.length > 0" class="space-y-4 pb-10">
            <Card v-for="note in notes" :key="note.id" class="shadow-none border-l-4 border-l-muted-foreground/20">
              <CardContent class="p-4 space-y-2">
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-2">
                    <Avatar class="h-6 w-6">
                      <AvatarImage :src="note.avatar_url" />
                      <AvatarFallback>{{ note.first_name[0] }}</AvatarFallback>
                    </Avatar>
                    <span class="text-xs font-bold">{{ note.first_name }} {{ note.last_name }}</span>
                    <span class="text-[10px] text-muted-foreground">{{ formatDate(note.created_at) }}</span>
                  </div>
                  <Button variant="ghost" size="icon" class="h-7 w-7 text-destructive hover:bg-destructive/10" @click="handleDeleteNote(note.id)">
                    <Trash2 class="h-3.5 w-3.5" />
                  </Button>
                </div>
                <p class="text-sm text-foreground/80 leading-relaxed">{{ note.note }}</p>
              </CardContent>
            </Card>
          </div>
          <div v-else class="flex flex-col items-center justify-center py-20 border border-dashed rounded-xl bg-muted/5 text-muted-foreground text-sm">
            {{ t('globals.terms.noNotesFound') }}
          </div>
        </TabsContent>
      </div>
    </Tabs>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useContact } from '@/composables/useContact'
import { useInboxTypes } from '@/composables/useInboxTypes'
import { useUserStore } from '@/stores/user'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import { 
  ArrowLeft, ArrowRight, MessageSquare, CircleDot, 
  CheckCircle2, Trash2, User, Phone, Mail, 
  PlusCircle, Settings, ChevronRight, FileText
} from 'lucide-vue-next'

import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { useTemplateCategories } from '@/composables/useTemplateCategories'
import api from '@/api'

const route = useRoute()
const { t } = useI18n()
const userStore = useUserStore()
const { getContact, getConversations, getStats, getNotes, createNote, deleteNote } = useContact()
const { toParam: inboxTypeParam } = useInboxTypes()
const { categories, fetchCategories } = useTemplateCategories()

const contact = ref(null)
const conversations = ref([])
const stats = ref({})
const notes = ref([])
const newNote = ref('')
const allTemplates = ref([])
const templateSearch = ref('')

const parsedAttributes = computed(() => {
  if (!contact.value?.custom_attributes) return {}
  let attrs = contact.value.custom_attributes
  if (typeof attrs === 'string') {
    try {
      attrs = JSON.parse(attrs)
    } catch (e) {
      return {}
    }
  }
  return attrs
})

const hasAttributes = computed(() => {
  return Object.keys(parsedAttributes.value).length > 0
})

const filteredCategories = computed(() => {
  const userTeamIds = userStore.user?.teams?.map(team => team.id) || []
  
  return categories.value
    .map(category => {
      const categoryTemplates = allTemplates.value?.filter(tpl => 
        tpl.category_id === category.id && 
        tpl.type === 'note' &&
        (tpl.name.toLowerCase().includes(templateSearch.value.toLowerCase()) || 
         tpl.body.toLowerCase().includes(templateSearch.value.toLowerCase()))
      ) || []

      const hasTeamAccess = !category.team_ids || category.team_ids.length === 0 || 
                           category.team_ids.some(id => userTeamIds.includes(id))

      if (categoryTemplates.length > 0 && hasTeamAccess) {
        return { ...category, templates: categoryTemplates }
      }
      return null
    })
    .filter(Boolean)
})

const applyTemplate = (body) => {
  newNote.value = body
}

const formatDate = (date) => {
  if (!date) return ''
  return format(new Date(date), 'dd/MM/yyyy HH:mm', { locale: ptBR })
}

const loadData = async () => {
  const id = route.params.id
  try {
    const [c, convs, s, n, cats, tpls] = await Promise.all([
      getContact(id),
      getConversations(id),
      getStats(id),
      getNotes(id),
      fetchCategories(),
      api.getTemplates('note')
    ])
    contact.value = c
    conversations.value = convs
    stats.value = s
    notes.value = n
    allTemplates.value = tpls.data.data
  } catch (err) {
    console.error('Error loading data', err)
  }
}

const handleAddNote = async () => {
  if (!newNote.value.trim()) return
  try {
    await createNote(contact.value.id, newNote.value)
    newNote.value = ''
    notes.value = await getNotes(contact.value.id)
  } catch (err) {
    console.error('Error adding note', err)
  }
}

const handleDeleteNote = async (noteId) => {
  if (!confirm(t('globals.messages.confirmDelete'))) return
  try {
    await deleteNote(contact.value.id, noteId)
    notes.value = await getNotes(contact.value.id)
  } catch (err) {
    console.error('Error deleting note', err)
  }
}

onMounted(loadData)
</script>
