<template>
  <div class="glass-strong rounded-2xl p-6 md:p-8">
    <h3 class="mb-6 text-lg font-bold text-white">New Contribution</h3>

    <div v-if="step === 1" class="space-y-4">
      <p class="text-sm text-white/40">What would you like to contribute?</p>
      <div class="grid grid-cols-2 gap-3 md:grid-cols-3">
        <button
          v-for="option in contributionTypes"
          :key="option.value"
          type="button"
          class="spring flex flex-col items-center gap-2 rounded-xl border border-white/6 p-4 text-center transition-all hover:border-white/20 hover:bg-white/4"
          :class="{ 'border-spotify! bg-spotify/10!': selectedType === option.value }"
          @click="selectedType = option.value; step = 2"
        >
          <i aria-hidden="true" :class="option.icon" class="text-xl" :style="{ color: option.color }" />
          <span class="text-xs font-medium text-white/70">{{ option.label }}</span>
        </button>
      </div>
    </div>

    <div v-else-if="step === 2" class="space-y-5">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <button
            type="button"
            aria-label="Back"
            class="text-xs text-white/40 transition-colors hover:text-white/60"
            @click="step = 1; selectedType = null"
          >
            <i aria-hidden="true" class="pi pi-arrow-left mr-1" /> Back
          </button>
          <span class="text-white/20">|</span>
          <span class="text-sm font-medium text-white/70">{{ typeLabel }}</span>
        </div>
      </div>

      <div class="grid grid-cols-3 gap-3">
        <button
          v-for="t in targetOptions"
          :key="t.value"
          type="button"
          class="spring rounded-lg border border-white/6 px-3 py-2 text-xs font-medium transition-all"
          :class="
            form.target_type === t.value
              ? 'border-spotify! bg-spotify/10! text-spotify!'
              : 'text-white/40 hover:border-white/20 hover:text-white/60'
          "
          @click="form.target_type = t.value"
        >
          {{ t.label }}
        </button>
      </div>

      <input
        v-model="form.target_id"
        type="text"
        placeholder="Target ID (UUID)"
        aria-label="Target ID"
        class="w-full rounded-xl bg-white/5 px-4 py-3 text-sm text-white placeholder-white/20 ring-1 ring-white/6 transition-all outline-hidden focus:ring-spotify/50"
      />

      <div v-if="selectedType === 'lyrics' || selectedType === 'translation'">
        <input
          v-model="form.locale"
          type="text"
          placeholder="Language code (e.g., fa, en, ar)"
          aria-label="Language code"
          class="w-full rounded-xl bg-white/5 px-4 py-3 text-sm text-white placeholder-white/20 ring-1 ring-white/6 transition-all outline-hidden focus:ring-spotify/50"
        />
      </div>

      <textarea
        v-model="rawData"
        :rows="8"
        aria-label="Contribution content"
        placeholder="Paste your contribution content here...
For lyrics: paste the full lyrics text
For LRC format: [00:00.00]Line 1&#10;[00:05.00]Line 2"
        class="w-full resize-y rounded-xl bg-white/5 px-4 py-3 font-mono text-sm text-white placeholder-white/20 ring-1 ring-white/6 transition-all outline-hidden focus:ring-spotify/50"
      />

      <input
        v-model="form.summary"
        type="text"
        placeholder="Brief summary of your change (optional)"
        aria-label="Summary"
        class="w-full rounded-xl bg-white/5 px-4 py-3 text-sm text-white placeholder-white/20 ring-1 ring-white/6 transition-all outline-hidden focus:ring-spotify/50"
      />

      <div class="flex items-center gap-2">
        <input
          id="is-minor"
          v-model="form.is_minor"
          type="checkbox"
          class="rounded-sm border-white/20 bg-white/5 text-spotify focus:ring-spotify"
        />
        <label for="is-minor" class="text-xs text-white/40">This is a minor edit</label>
      </div>

      <div class="flex gap-3">
        <button
          type="button"
          class="spring flex-1 rounded-xl bg-white/5 py-3 text-sm font-medium text-white/50 transition-all hover:bg-white/10 hover:text-white/70"
          @click="reset"
        >
          Cancel
        </button>
        <button
          type="button"
          :disabled="canSubmit! || submitting"
          class="spring flex-1 rounded-xl bg-spotify py-3 text-sm font-bold text-black transition-all hover:bg-spotify-hover disabled:opacity-40"
          @click="submit"
        >
          <i aria-hidden="true" v-if="submitting" class="pi pi-spin pi-spinner mr-2" />
          Submit {{ typeLabel }}
        </button>
      </div>

      <p v-if="error" class="text-xs text-red-400">{{ error }}</p>
      <p v-if="success" class="text-xs text-spotify">Submitted! Status: {{ successStatus }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useContributionApi } from '@/services/api/contribution'
import type { ContributionType, TargetType } from '@/services/api/contribution'

const emit = defineEmits<{ submitted: [] }>()

const api = useContributionApi()

const contributionTypes = [
  {
    value: 'lyrics' as ContributionType,
    label: 'Lyrics',
    icon: 'pi pi-align-left',
    color: '#1db954',
  },
  {
    value: 'translation' as ContributionType,
    label: 'Translation',
    icon: 'pi pi-language',
    color: '#60a5fa',
  },
  { value: 'credits' as ContributionType, label: 'Credits', icon: 'pi pi-users', color: '#a855f7' },
  { value: 'metadata' as ContributionType, label: 'Metadata', icon: 'pi pi-tag', color: '#f59e0b' },
  {
    value: 'album_art' as ContributionType,
    label: 'Album Art',
    icon: 'pi pi-image',
    color: '#f472b6',
  },
  {
    value: 'bio' as ContributionType,
    label: 'Artist Bio',
    icon: 'pi pi-info-circle',
    color: '#34d399',
  },
]

const targetOptions = [
  { value: 'track' as TargetType, label: 'Track' },
  { value: 'album' as TargetType, label: 'Album' },
  { value: 'artist' as TargetType, label: 'Artist' },
]

const step = ref(1)
const selectedType = ref<ContributionType | null>(null)
const rawData = ref('')
const submitting = ref(false)
const error = ref('')
const success = ref(false)
const successStatus = ref('')

const form = ref({
  target_type: 'track' as TargetType,
  target_id: '',
  locale: '',
  summary: '',
  is_minor: false,
})

const typeLabel = computed(
  () => contributionTypes.find((t) => t.value === selectedType.value)?.label || '',
)
const canSubmit = computed(() => form.value.target_id && rawData.value.trim())

function reset() {
  step.value = 1
  selectedType.value = null
  rawData.value = ''
  submitting.value = false
  error.value = ''
  success.value = false
  successStatus.value = ''
  form.value = { target_type: 'track', target_id: '', locale: '', summary: '', is_minor: false }
}

async function submit() {
  if (canSubmit.value! || selectedType.value!) return
  submitting.value = true
  error.value = ''
  success.value = false

  try {
    let data: Record<string, any>
    if (selectedType.value === 'lyrics') {
      const lines = rawData.value.trim().split('\n')
      const hasTimestamps = lines.some((l) => /^\[\d{2}:\d{2}(\.\d+)?\]/.test(l.trim()))
      data = hasTimestamps
        ? { format: 'lrc', content: rawData.value.trim() }
        : { format: 'plain', content: rawData.value.trim() }
    } else if (selectedType.value === 'translation') {
      data = { language: form.value.locale || 'en', content: rawData.value.trim() }
    } else {
      try {
        data = JSON.parse(rawData.value)
      } catch {
        data = { text: rawData.value.trim() }
      }
    }

    const res = await api.create({
      contribution_type: selectedType.value,
      target_type: form.value.target_type,
      target_id: form.value.target_id,
      locale: form.value.locale || undefined,
      data,
      summary: form.value.summary || undefined,
      is_minor: form.value.is_minor,
    })

    success.value = true
    successStatus.value = res.status || 'pending'
    emit('submitted')
    setTimeout(reset, 2000)
  } catch (err) {
    console.error('Failed to submit contribution:', err)
    error.value = (err as Record<string, any>)?.message || 'Failed to submit'
  } finally {
    submitting.value = false
  }
}
</script>
