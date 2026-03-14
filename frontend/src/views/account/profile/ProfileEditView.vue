<template>
  <div class="mx-auto max-w-6xl space-y-6 p-6">
    <div class="grid gap-6 xl:grid-cols-[minmax(0,1.7fr)_minmax(320px,1fr)]">
      <Card class="border-border/70 shadow-sm">
        <CardHeader class="space-y-4">
          <div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
            <div class="space-y-2">
              <div class="flex items-center gap-2">
                <CardTitle class="text-2xl">{{ t('account.editProfile') }}</CardTitle>
                <Badge variant="secondary">{{ availabilityLabel }}</Badge>
              </div>
              <CardDescription class="max-w-2xl">
                {{ t('account.profileDescription') }}
              </CardDescription>
            </div>

            <div class="flex flex-wrap gap-2">
              <Button
                variant="outline"
                @click="selectAvatar"
              >
                {{ t('account.chooseAFile') }}
              </Button>
              <Button
                variant="destructive"
                @click="removeAvatar"
                :disabled="!hasAvatar && !hasPendingAvatar"
              >
                {{ t('account.removeAvatar') }}
              </Button>
              <Button
                @click="saveUser"
                :isLoading="isSaving"
                :disabled="!hasPendingAvatar"
              >
                {{ t('account.saveAvatarChanges') }}
              </Button>
            </div>
          </div>
        </CardHeader>

        <CardContent class="space-y-6">
          <input
            ref="uploadInput"
            type="file"
            hidden
            accept="image/jpg, image/jpeg, image/png"
            @change="selectFile"
          />

          <div class="grid gap-6 lg:grid-cols-[auto_minmax(0,1fr)]">
            <div class="flex flex-col items-center gap-3">
              <Avatar class="size-32 border border-border/60 shadow-sm">
                <AvatarImage :src="userStore.avatar" alt="Avatar do usuário" />
                <AvatarFallback class="text-2xl font-semibold">
                  {{ userStore.getInitials }}
                </AvatarFallback>
              </Avatar>
              <p class="text-center text-xs text-muted-foreground">
                {{ t('account.avatarHint') }}
              </p>
            </div>

            <div class="grid gap-4 md:grid-cols-2">
              <div class="rounded-2xl border border-border/70 bg-muted/20 p-4">
                <p class="text-xs font-medium uppercase tracking-[0.18em] text-muted-foreground">
                  {{ t('globals.terms.name') }}
                </p>
                <p class="mt-2 text-base font-semibold">
                  {{ userStore.getFullName || t('globals.terms.notAvailable') }}
                </p>
              </div>

              <div class="rounded-2xl border border-border/70 bg-muted/20 p-4">
                <p class="text-xs font-medium uppercase tracking-[0.18em] text-muted-foreground">
                  {{ t('globals.terms.email') }}
                </p>
                <p class="mt-2 text-base font-semibold break-all">
                  {{ userStore.email || t('globals.terms.notAvailable') }}
                </p>
              </div>

              <div class="rounded-2xl border border-border/70 bg-muted/20 p-4">
                <p class="text-xs font-medium uppercase tracking-[0.18em] text-muted-foreground">
                  {{ t('account.teams') }}
                </p>
                <div class="mt-3 flex flex-wrap gap-2">
                  <Badge
                    v-for="team in userStore.teams"
                    :key="team.id || team.name"
                    variant="outline"
                  >
                    {{ team.name }}
                  </Badge>
                  <span v-if="!userStore.teams.length" class="text-sm text-muted-foreground">
                    {{ t('account.noTeams') }}
                  </span>
                </div>
              </div>

              <div class="rounded-2xl border border-border/70 bg-muted/20 p-4">
                <p class="text-xs font-medium uppercase tracking-[0.18em] text-muted-foreground">
                  {{ t('account.roles') }}
                </p>
                <div class="mt-3 flex flex-wrap gap-2">
                  <Badge
                    v-for="role in userStore.roles"
                    :key="role"
                    variant="secondary"
                  >
                    {{ role }}
                  </Badge>
                  <span v-if="!userStore.roles.length" class="text-sm text-muted-foreground">
                    {{ t('account.noRoles') }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <Separator />

          <Alert>
            <ShieldCheck class="h-4 w-4" />
            <AlertTitle>{{ t('account.avatarSecurityTitle') }}</AlertTitle>
            <AlertDescription>
              {{ t('account.avatarSecurityDescription') }}
            </AlertDescription>
          </Alert>
        </CardContent>
      </Card>

      <div class="space-y-6">
        <Card class="border-border/70 shadow-sm">
          <CardHeader>
            <CardTitle>{{ t('account.publicAvatar') }}</CardTitle>
            <CardDescription>{{ t('account.changeAvatar') }}</CardDescription>
          </CardHeader>
          <CardContent class="space-y-4">
            <div class="rounded-2xl border border-dashed border-border/80 bg-muted/20 p-4">
              <p class="text-sm font-medium">{{ t('account.recommendedFormat') }}</p>
              <p class="mt-2 text-sm text-muted-foreground">
                {{ t('account.recommendedFormatDescription') }}
              </p>
            </div>

            <div class="rounded-2xl border border-border/70 bg-muted/20 p-4">
              <p class="text-sm font-medium">{{ t('account.currentStatus') }}</p>
              <p class="mt-2 text-sm text-muted-foreground">
                {{ availabilityDescription }}
              </p>
            </div>
          </CardContent>
        </Card>

        <Card class="border-border/70 shadow-sm">
          <CardHeader>
            <CardTitle>{{ t('account.quickActions') }}</CardTitle>
            <CardDescription>{{ t('account.quickActionsDescription') }}</CardDescription>
          </CardHeader>
          <CardContent class="space-y-3">
            <Button class="w-full justify-start" variant="outline" @click="selectAvatar">
              {{ t('account.chooseAFile') }}
            </Button>
            <Button
              class="w-full justify-start"
              :disabled="!hasPendingAvatar"
              @click="saveUser"
              :isLoading="isSaving"
            >
              {{ t('account.saveAvatarChanges') }}
            </Button>
            <Button
              class="w-full justify-start"
              variant="destructive"
              :disabled="!hasAvatar && !hasPendingAvatar"
              @click="removeAvatar"
            >
              {{ t('account.removeAvatar') }}
            </Button>
          </CardContent>
        </Card>
      </div>
    </div>

    <Dialog :open="showCropper">
      <DialogContent class="sm:max-w-2xl">
        <DialogHeader>
          <DialogTitle class="text-xl">{{ t('account.cropAvatar') }}</DialogTitle>
          <DialogDescription>{{ t('account.cropAvatarDescription') }}</DialogDescription>
        </DialogHeader>

        <VuePictureCropper
          :boxStyle="{
            width: '100%',
            height: '420px',
            backgroundColor: '#f8f8f8',
            margin: 'auto'
          }"
          :img="newUserAvatar"
          :options="{ viewMode: 1, dragMode: 'crop', aspectRatio: 1 }"
        />
        <DialogFooter class="sm:justify-end">
          <Button variant="secondary" @click="closeDialog">
            {{ t('globals.messages.close') }}
          </Button>
          <Button @click="getResult">{{ t('globals.messages.save') }}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { ShieldCheck } from 'lucide-vue-next'
import VuePictureCropper, { cropper } from 'vue-picture-cropper'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/stores/user'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { useEmitter } from '@/composables/useEmitter'
import { handleHTTPError } from '@/utils/http'
import { EMITTER_EVENTS } from '@/constants/emitterEvents.js'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogDescription
} from '@/components/ui/dialog'
import api from '@/api'

const emitter = useEmitter()
const { t } = useI18n()
const isSaving = ref(false)
const userStore = useUserStore()
const uploadInput = ref(null)
const newUserAvatar = ref('')
const showCropper = ref(false)
let croppedBlob = null
let avatarFile = null

const hasAvatar = computed(() => Boolean(userStore.avatar))
const hasPendingAvatar = computed(() => Boolean(croppedBlob))

const availabilityLabel = computed(() => {
  const labels = {
    online: t('account.availability.online'),
    away: t('account.availability.away'),
    away_manual: t('account.availability.away'),
    offline: t('account.availability.offline')
  }

  return labels[userStore.user.availability_status] || t('account.availability.offline')
})

const availabilityDescription = computed(() => {
  const descriptions = {
    online: t('account.availabilityDescription.online'),
    away: t('account.availabilityDescription.away'),
    away_manual: t('account.availabilityDescription.away'),
    offline: t('account.availabilityDescription.offline')
  }

  return descriptions[userStore.user.availability_status] || t('account.availabilityDescription.offline')
})

const selectAvatar = () => {
  uploadInput.value.click()
}

const selectFile = (event) => {
  newUserAvatar.value = ''
  const { files } = event.target
  if (!files || !files.length) return
  avatarFile = files[0]
  const reader = new FileReader()
  reader.readAsDataURL(avatarFile)
  reader.onload = () => {
    newUserAvatar.value = String(reader.result)
    showCropper.value = true
    uploadInput.value.value = ''
  }
}

const closeDialog = () => {
  showCropper.value = false
}

const getResult = async () => {
  if (!cropper) return
  croppedBlob = await cropper.getBlob()
  if (!croppedBlob) return
  userStore.setAvatar(URL.createObjectURL(croppedBlob))
  showCropper.value = false
}

const saveUser = async () => {
  if (!croppedBlob) return

  const formData = new FormData()
  formData.append('files', croppedBlob, 'avatar.png')
  try {
    isSaving.value = true
    await api.updateCurrentUser(formData)
    croppedBlob = null
    await userStore.getCurrentUser()
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      description: t('globals.messages.updatedSuccessfully', {
        name: t('globals.terms.profile')
      })
    })
  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  } finally {
    isSaving.value = false
  }
}

const removeAvatar = async () => {
  croppedBlob = null
  avatarFile = null
  try {
    await api.deleteUserAvatar()
    userStore.clearAvatar()
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      description: t('account.avatarRemoved')
    })
  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  }
}
</script>
