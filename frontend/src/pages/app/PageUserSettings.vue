<template>
  <div class="mx-auto w-full max-w-4xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <div v-if="loading" class="space-y-6">
      <SkeletonLoader variant="lines" :lines="1" class="w-48" />
      <SkeletonLoader variant="lines" :lines="8" />
    </div>

    <template v-else>
      <div class="mb-8">
        <p class="text-xs font-bold tracking-[0.25em] text-spotify uppercase">Settings</p>
        <h1 class="mt-1 text-3xl font-black text-white">Account Settings</h1>
      </div>

      <div class="mb-8 flex gap-1 rounded-2xl border border-white/6 bg-white/3 p-1">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        type="button"
        class="flex-1 rounded-xl px-4 py-2.5 text-sm font-bold transition"
        :class="
          activeTab === tab.key ? 'bg-white/10 text-white' : 'text-slate-400 hover:text-white'
        "
        @click="activeTab = tab.key"
      >
        <i aria-hidden="true" :class="tab.icon" class="mr-2" />
        {{ tab.label }}
      </button>
    </div>

    <!-- Profile Tab -->
    <div v-show="activeTab === 'profile'" class="space-y-6">
      <div class="flex items-center gap-6">
        <div class="relative shrink-0">
          <div
            class="flex h-20 w-20 items-center justify-center overflow-hidden rounded-full bg-white/10 text-3xl font-black text-white"
          >
            <img v-if="previewAvatar" :src="previewAvatar" alt="Avatar preview" class="h-full w-full object-cover" />
            <span v-else>{{ initials }}</span>
          </div>
          <button
            type="button"
            aria-label="Edit avatar"
            class="absolute -right-1 -bottom-1 flex h-7 w-7 items-center justify-center rounded-full border-2 border-black bg-spotify text-[10px] text-black transition hover:bg-spotify-hover"
            @click="triggerAvatarUpload"
          >
            <i aria-hidden="true" class="pi pi-pencil" />
          </button>
          <input
            ref="avatarInput"
            type="file"
            accept="image/*"
            class="hidden"
            @change="handleAvatar"
          />
        </div>
        <div>
          <h2 class="text-xl font-bold text-white">{{ form.displayName || 'User' }}</h2>
          <p class="text-sm text-slate-400">@{{ form.username || 'username' }}</p>
        </div>
      </div>

      <div class="grid gap-6 md:grid-cols-2">
        <div>
          <label class="mb-1.5 block text-xs font-semibold tracking-wider text-slate-400 uppercase"
            >Display Name</label
          >
          <input
            v-model="form.displayName"
            type="text"
            aria-label="Display name"
            class="w-full rounded-xl border border-white/8 bg-white/4 px-4 py-2.5 text-sm text-white transition outline-hidden placeholder:text-slate-600 focus:border-spotify/50 focus:bg-white/6"
            placeholder="Your display name"
          />
        </div>
        <div>
          <label class="mb-1.5 block text-xs font-semibold tracking-wider text-slate-400 uppercase"
            >Username</label
          >
          <input
            v-model="form.username"
            type="text"
            aria-label="Username"
            class="w-full rounded-xl border border-white/8 bg-white/4 px-4 py-2.5 text-sm text-white transition outline-hidden placeholder:text-slate-600 focus:border-spotify/50 focus:bg-white/6"
            placeholder="username"
          />
        </div>
        <div class="md:col-span-2">
          <label class="mb-1.5 block text-xs font-semibold tracking-wider text-slate-400 uppercase"
            >Bio</label
          >
          <textarea
            v-model="form.bio"
            rows="3"
            aria-label="Bio"
            class="w-full resize-none rounded-xl border border-white/8 bg-white/4 px-4 py-2.5 text-sm text-white transition outline-hidden placeholder:text-slate-600 focus:border-spotify/50 focus:bg-white/6"
            placeholder="Tell us about yourself"
          />
        </div>
        <div>
          <label class="mb-1.5 block text-xs font-semibold tracking-wider text-slate-400 uppercase"
            >Location</label
          >
          <input
            v-model="form.location"
            type="text"
            aria-label="Location"
            class="w-full rounded-xl border border-white/8 bg-white/4 px-4 py-2.5 text-sm text-white transition outline-hidden placeholder:text-slate-600 focus:border-spotify/50 focus:bg-white/6"
            placeholder="Tehran, Iran"
          />
        </div>
        <div>
          <label class="mb-1.5 block text-xs font-semibold tracking-wider text-slate-400 uppercase"
            >Website</label
          >
          <input
            v-model="form.website"
            type="url"
            aria-label="Website"
            class="w-full rounded-xl border border-white/8 bg-white/4 px-4 py-2.5 text-sm text-white transition outline-hidden placeholder:text-slate-600 focus:border-spotify/50 focus:bg-white/6"
            placeholder="https://example.com"
          />
        </div>
      </div>

      <div class="flex items-center gap-3 pt-2">
        <button
          type="button"
          :disabled="saving"
          class="rounded-full bg-spotify px-6 py-2.5 text-sm font-bold text-black transition hover:bg-spotify-hover disabled:opacity-50"
          @click="saveProfile"
        >
          {{ saving ? 'Saving...' : 'Save changes' }}
        </button>
        <p
          v-if="saveMessage"
          class="text-sm"
          :class="saveError ? 'text-red-400' : 'text-spotify'"
          :role="saveError ? 'alert' : undefined"
          aria-live="polite"
        >
          {{ saveMessage }}
        </p>
      </div>
    </div>

    <!-- Account Tab -->
    <div v-show="activeTab === 'account'" class="space-y-8">
      <div>
        <h3 class="mb-1 text-lg font-bold text-white">Email</h3>
        <p class="mb-3 text-sm text-slate-400">{{ auth.user?.email }}</p>
        <p class="text-xs text-slate-500">Email cannot be changed at this time.</p>
      </div>

      <div class="rounded-2xl border border-white/6 bg-white/2 p-6">
        <h3 class="mb-4 text-lg font-bold text-white">Change Password</h3>
        <div class="space-y-4">
          <div>
            <label
              class="mb-1.5 block text-xs font-semibold tracking-wider text-slate-400 uppercase"
              >Current Password</label
            >
            <input
              v-model="passwordForm.currentPassword"
              type="password"
              aria-label="Current Password"
              class="w-full rounded-xl border border-white/8 bg-white/4 px-4 py-2.5 text-sm text-white transition outline-hidden placeholder:text-slate-600 focus:border-spotify/50 focus:bg-white/6"
            />
          </div>
          <div>
            <label
              class="mb-1.5 block text-xs font-semibold tracking-wider text-slate-400 uppercase"
              >New Password</label
            >
            <input
              v-model="passwordForm.newPassword"
              type="password"
              aria-label="New Password"
              class="w-full rounded-xl border border-white/8 bg-white/4 px-4 py-2.5 text-sm text-white transition outline-hidden placeholder:text-slate-600 focus:border-spotify/50 focus:bg-white/6"
            />
          </div>
          <button
            type="button"
            :disabled="passwordSaving"
            class="rounded-full bg-white/10 px-6 py-2.5 text-sm font-bold text-white transition hover:bg-white/15 disabled:opacity-50"
            @click="savePassword"
          >
            {{ passwordSaving ? 'Updating...' : 'Update Password' }}
          </button>
          <p
            v-if="passwordMessage"
            class="text-sm"
            :class="passwordError ? 'text-red-400' : 'text-spotify'"
            :role="passwordError ? 'alert' : undefined"
            aria-live="polite"
          >
            {{ passwordMessage }}
          </p>
        </div>
      </div>

      <div class="rounded-2xl border border-red-500/20 bg-red-500/3 p-6">
        <h3 class="mb-2 text-lg font-bold text-red-400">Delete Account</h3>
        <p class="mb-4 text-sm text-slate-400">
          Permanently delete your account and all associated data. This action cannot be undone.
        </p>
        <button
          type="button"
          class="rounded-full border border-red-500/30 px-6 py-2.5 text-sm font-bold text-red-400 transition hover:bg-red-500/10"
          @click="confirmDelete"
        >
          Delete Account
        </button>
      </div>
    </div>

    <!-- Preferences Tab -->
    <div v-show="activeTab === 'preferences'" class="space-y-6">
      <div class="rounded-2xl border border-white/6 bg-white/2 p-6">
        <h3 class="mb-4 text-lg font-bold text-white">Display</h3>
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-white">Persian (RTL) Layout</p>
              <p class="text-xs text-slate-500">Enable right-to-left layout for Persian language</p>
            </div>
            <input
              v-model="preferences.rtl"
              type="checkbox"
              aria-label="Persian (RTL) Layout"
              class="h-5 w-5 rounded border-white/20 bg-white/10 accent-[#1db954]"
            />
          </div>
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-white">Lyrics Autoscroll</p>
              <p class="text-xs text-slate-500">Automatically scroll lyrics during playback</p>
            </div>
            <input
              v-model="preferences.lyricsAutoscroll"
              type="checkbox"
              aria-label="Lyrics Autoscroll"
              class="h-5 w-5 rounded border-white/20 bg-white/10 accent-[#1db954]"
            />
          </div>
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-white">Explicit Content</p>
              <p class="text-xs text-slate-500">Allow playback of explicit tracks</p>
            </div>
            <input
              v-model="preferences.explicitContent"
              type="checkbox"
              aria-label="Explicit Content"
              class="h-5 w-5 rounded border-white/20 bg-white/10 accent-[#1db954]"
            />
          </div>
        </div>
      </div>

      <div class="rounded-2xl border border-white/6 bg-white/2 p-6">
        <h3 class="mb-4 text-lg font-bold text-white">Notifications</h3>
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-white">New Releases</p>
              <p class="text-xs text-slate-500">
                Get notified when followed artists release new music
              </p>
            </div>
            <input
              v-model="preferences.notifyReleases"
              type="checkbox"
              aria-label="New Releases"
              class="h-5 w-5 rounded border-white/20 bg-white/10 accent-[#1db954]"
            />
          </div>
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-white">Social Activity</p>
              <p class="text-xs text-slate-500">Get notified about follows, shares, and likes</p>
            </div>
            <input
              v-model="preferences.notifySocial"
              type="checkbox"
              aria-label="Social Activity"
              class="h-5 w-5 rounded border-white/20 bg-white/10 accent-[#1db954]"
            />
          </div>
        </div>
      </div>

      <button
        type="button"
        :disabled="prefSaving"
        class="rounded-full bg-spotify px-6 py-2.5 text-sm font-bold text-black transition hover:bg-spotify-hover disabled:opacity-50"
        @click="savePreferences"
      >
        {{ prefSaving ? 'Saving...' : 'Save Preferences' }}
      </button>
      <p v-if="prefMessage" class="text-sm text-spotify" aria-live="polite">{{ prefMessage }}</p>
    </div>
  </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useToast } from 'primevue/usetoast'
import { SkeletonLoader } from '@/components/common'
import { useUserAuthStore } from '@/stores'
import { useRequest, useRTL } from '@/composables'

const auth = useUserAuthStore()
const toast = useToast()
const { isRTL, setRTL } = useRTL()

const activeTab = ref('profile')
const tabs = [
  { key: 'profile', label: 'Profile', icon: 'pi pi-user' },
  { key: 'account', label: 'Account', icon: 'pi pi-lock' },
  { key: 'preferences', label: 'Preferences', icon: 'pi pi-cog' },
]

const saving = ref(false)
const saveMessage = ref('')
const saveError = ref(false)
const passwordSaving = ref(false)
const passwordMessage = ref('')
const passwordError = ref(false)
const prefSaving = ref(false)
const prefMessage = ref('')
const avatarInput = ref<HTMLInputElement | null>(null)
const previewAvatar = ref('')
const loading = ref(true)

const form = reactive({
  displayName: '',
  username: '',
  bio: '',
  location: '',
  website: '',
})

const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
})

const preferences = reactive({
  rtl: false,
  lyricsAutoscroll: true,
  explicitContent: true,
  notifyReleases: true,
  notifySocial: true,
})

// Sync RTL preference with composable
watch(() => preferences.rtl, (val) => {
  setRTL(val)
}, { immediate: true })

// Sync initial value from composable to form
if (isRTL.value) {
  preferences.rtl = true
}

const initials = computed(() => {
  const name = auth.user?.displayName || auth.user?.username || auth.user?.name || '?'
  const words = name.split(/\s+/).filter(Boolean)
  return words.length >= 2
    ? (words[0]![0]! + words[words.length - 1]![0]!).toUpperCase()
    : name.slice(0, 2).toUpperCase()
})

onMounted(() => {
  if (auth.user) {
    form.displayName = auth.user.displayName || ''
    form.username = auth.user.username || ''
    form.bio = (auth.user as any as { bio?: string }).bio || ''
    form.location = (auth.user as any as { location?: string }).location || ''
    form.website = (auth.user as any as { website?: string }).website || ''
  }
  loading.value = false
})

function triggerAvatarUpload() {
  avatarInput.value?.click()
}

function handleAvatar(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  previewAvatar.value = URL.createObjectURL(file)
  toast.add({ severity: 'info', summary: 'Avatar upload coming soon', life: 2000 })
}

async function saveProfile() {
  saving.value = true
  saveMessage.value = ''
  saveError.value = false
  try {
    await useRequest('/users/me/profile', {
      method: 'PUT',
      data: {
        displayName: form.displayName || undefined,
        username: form.username || undefined,
        bio: form.bio || undefined,
        location: form.location || undefined,
        website: form.website || undefined,
      },
    })
    await auth.me()
    saveMessage.value = 'Profile updated'
    setTimeout(() => {
      saveMessage.value = ''
    }, 3000)
  } catch (e) {
    saveError.value = true
    saveMessage.value = e instanceof Error ? e.message : 'Failed to update profile'
  } finally {
    saving.value = false
  }
}

async function savePassword() {
  if (!passwordForm.currentPassword || !passwordForm.newPassword) {
    passwordMessage.value = 'Both fields are required'
    passwordError.value = true
    return
  }
  if (passwordForm.newPassword.length < 8) {
    passwordMessage.value = 'New password must be at least 8 characters'
    passwordError.value = true
    return
  }
  passwordSaving.value = true
  passwordMessage.value = ''
  passwordError.value = false
  try {
    await useRequest('/users/me/password', {
      method: 'PUT',
      data: {
        currentPassword: passwordForm.currentPassword,
        newPassword: passwordForm.newPassword,
      },
    })
    passwordForm.currentPassword = ''
    passwordForm.newPassword = ''
    passwordMessage.value = 'Password updated successfully'
  } catch (e) {
    passwordError.value = true
    passwordMessage.value = e instanceof Error ? e.message : 'Failed to update password'
  } finally {
    passwordSaving.value = false
  }
}

async function savePreferences() {
  prefSaving.value = true
  prefMessage.value = ''
  try {
    await useRequest('/users/me/profile', {
      method: 'PUT',
      data: { preferences: JSON.stringify({ ...preferences }) },
    })
    prefMessage.value = 'Preferences saved'
    setTimeout(() => {
      prefMessage.value = ''
    }, 3000)
  } catch (err) {
    console.error('Failed to save preferences:', err)
    prefMessage.value = 'Failed to save preferences'
  } finally {
    prefSaving.value = false
  }
}

function confirmDelete() {
  toast.add({
    severity: 'warn',
    summary: 'Coming soon',
    detail: 'Account deletion is not yet implemented',
    life: 3000,
  })
}
</script>
