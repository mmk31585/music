<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useUploadsApi } from '@/services/api/uploads'
import type { Draft } from '@/services/api/uploads/types'

const uploadsApi = useUploadsApi()
const pending = ref<Draft[]>([])
const loading = ref(true)
const selectedDraft = ref<Draft | null>(null)
const reviewNotes = ref('')
const processing = ref(false)

async function loadPending() {
  loading.value = true
  try {
    const result = await uploadsApi.listPendingDrafts()
    pending.value = Array.isArray(result) ? result : [result]
  } finally {
    loading.value = false
  }
}

async function handleReview(action: 'accept' | 'reject') {
  if (!selectedDraft.value) return
  processing.value = true
  try {
    await uploadsApi.reviewDraft(
      selectedDraft.value.id,
      action,
      reviewNotes.value || undefined,
    )
    pending.value = pending.value.filter((d) => d.id !== selectedDraft.value!.id)
    selectedDraft.value = null
    reviewNotes.value = ''
  } finally {
    processing.value = false
  }
}

onMounted(loadPending)
</script>

<template>
  <div class="mx-auto w-full max-w-6xl px-4 pb-32 pt-8 md:px-6">
    <h1 class="mb-8 text-3xl font-black text-white">Upload Review Queue</h1>

    <div v-if="loading" class="space-y-4">
      <div v-for="i in 5" :key="i" class="h-16 animate-pulse rounded-xl bg-white/5" />
    </div>

    <template v-else>
      <p class="mb-4 text-sm text-white/40">
        {{ pending.length }} pending upload{{ pending.length === 1 ? '' : 's' }} awaiting review
      </p>

      <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <!-- Draft list -->
        <div class="space-y-2 lg:col-span-1">
          <button
            v-for="d in pending"
            :key="d.id"
            class="w-full rounded-xl p-3 text-left transition"
            :class="
              selectedDraft?.id === d.id
                ? 'bg-white/10 ring-1 ring-white/20'
                : 'bg-white/5 hover:bg-white/8'
            "
            @click="selectedDraft = d"
          >
            <p class="text-sm font-medium text-white truncate">{{ d.originalFilename }}</p>
            <p class="text-xs text-white/40">
              {{ d.uploadSource }} &middot; {{ new Date(d.createdAt).toLocaleDateString() }}
            </p>
          </button>
          <p v-if="!pending.length" class="py-8 text-center text-sm text-white/30">
            All caught up! No pending uploads.
          </p>
        </div>

        <!-- Detail + Review -->
        <div v-if="selectedDraft" class="space-y-6 lg:col-span-2">
          <div class="rounded-2xl bg-white/5 p-6">
            <h2 class="mb-4 text-lg font-bold text-white">{{ selectedDraft.originalFilename }}</h2>
            <dl class="grid grid-cols-2 gap-4 text-sm">
              <div>
                <dt class="text-white/40">Status</dt>
                <dd class="font-medium text-white">{{ selectedDraft.status }}</dd>
              </div>
              <div>
                <dt class="text-white/40">Source</dt>
                <dd class="font-medium text-white">{{ selectedDraft.uploadSource }}</dd>
              </div>
              <div>
                <dt class="text-white/40">Size</dt>
                <dd class="font-medium text-white">{{ (selectedDraft.fileSize / 1e6).toFixed(1) }} MB</dd>
              </div>
              <div>
                <dt class="text-white/40">Submitted</dt>
                <dd class="font-medium text-white">{{ new Date(selectedDraft.createdAt).toLocaleDateString() }}</dd>
              </div>
            </dl>
          </div>

          <!-- Review form -->
          <div class="rounded-2xl bg-white/5 p-6">
            <h3 class="mb-3 text-sm font-bold text-white/60 uppercase tracking-wider">Review Decision</h3>
            <textarea
              v-model="reviewNotes"
              placeholder="Add review notes (optional)..."
              class="mb-4 w-full rounded-xl border border-white/10 bg-white/5 p-3 text-sm text-white placeholder-white/30 outline-none transition focus:border-spotify/50"
              rows="3"
            />
            <div class="flex gap-3">
              <button
                class="rounded-full bg-emerald-600 px-8 py-2.5 text-sm font-bold text-white transition hover:bg-emerald-500 disabled:opacity-50"
                :disabled="processing"
                @click="handleReview('accept')"
              >
                {{ processing ? 'Processing...' : 'Accept' }}
              </button>
              <button
                class="rounded-full bg-red-600 px-8 py-2.5 text-sm font-bold text-white transition hover:bg-red-500 disabled:opacity-50"
                :disabled="processing"
                @click="handleReview('reject')"
              >
                {{ processing ? 'Processing...' : 'Reject' }}
              </button>
            </div>
          </div>
        </div>

        <div v-else class="flex items-center justify-center py-24 lg:col-span-2">
          <p class="text-sm text-white/30">Select a draft to review</p>
        </div>
      </div>
    </template>
  </div>
</template>
