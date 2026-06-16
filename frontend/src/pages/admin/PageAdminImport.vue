<template>
  <div class="mx-auto w-full max-w-5xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      eyebrow="Import"
      title="Import from Internet"
      description="Search YouTube and import tracks directly into the ingestion pipeline."
    />

    <div class="mt-6">
      <div class="flex gap-3">
        <IconField class="flex-1">
          <InputIcon><i class="pi pi-search"></i></InputIcon>
          <InputText
            v-model="query"
            placeholder="Search for a track — e.g. 'feel it from d4vd'"
            class="w-full"
            @keydown.enter="doSearch"
          />
        </IconField>
        <Button
          label="Search"
          icon="pi pi-search"
          :loading="searching"
          :disabled="!query.trim()"
          @click="doSearch"
        />
      </div>
    </div>

    <div v-if="searching" class="mt-12 text-center">
      <i class="pi pi-spin pi-spinner text-3xl text-slate-400"></i>
      <p class="mt-3 text-sm text-slate-500">Searching...</p>
    </div>

    <div v-else-if="searchError" class="mt-6">
      <Message severity="error" :closable="false">{{ searchError }}</Message>
    </div>

    <div v-else-if="results.length > 0" class="mt-6 space-y-3">
      <p class="text-sm text-slate-500">{{ results.length }} result{{ results.length !== 1 ? 's' : '' }}</p>

      <div
        v-for="(r, i) in results"
        :key="r.url"
        class="group flex items-center gap-4 rounded-xl border border-white/[0.06] bg-white/[0.03] p-4 transition hover:border-white/[0.12]"
      >
        <img
          v-if="r.thumbnail"
          :src="r.thumbnail"
          alt=""
          class="h-16 w-16 flex-shrink-0 rounded-lg object-cover"
        />
        <div v-else class="flex h-16 w-16 flex-shrink-0 items-center justify-center rounded-lg bg-white/[0.06]">
          <i class="pi pi-music text-xl text-slate-500"></i>
        </div>

        <div class="min-w-0 flex-1">
          <p class="truncate text-sm font-medium text-white">{{ r.title }}</p>
          <p class="mt-0.5 text-xs text-slate-400">{{ r.artist }}</p>
          <p class="mt-0.5 text-xs text-slate-500">{{ formatDuration(r.duration) }} &bull; {{ r.source }}</p>
        </div>

        <Button
          :label="importingUrl === r.url ? 'Importing...' : 'Import'"
          :icon="importingUrl === r.url ? 'pi pi-spin pi-spinner' : 'pi pi-download'"
          :loading="importingUrl === r.url"
          :disabled="!!importingUrl"
          size="small"
          @click="doImport(r)"
        />
      </div>
    </div>

    <div v-else-if="searched" class="mt-12 text-center">
      <i class="pi pi-search text-3xl text-slate-500"></i>
      <p class="mt-3 text-sm text-slate-500">No results found. Try a different search term.</p>
    </div>

    <div v-else class="mt-12 text-center">
      <i class="pi pi-cloud-download text-3xl text-slate-500"></i>
      <p class="mt-3 text-sm text-slate-500">Search for a track to get started.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputIcon from 'primevue/inputicon'
import IconField from 'primevue/iconfield'
import Message from 'primevue/message'
import { AdminSectionHeader } from '@/components/admin'
import { useImportApi } from '@/services/api/importcmd'
import type { SearchResult } from '@/services/api/importcmd'

const router = useRouter()
const toast = useToast()
const importApi = useImportApi()

const query = ref('')
const results = ref<SearchResult[]>([])
const searching = ref(false)
const searchError = ref('')
const searched = ref(false)
const importingUrl = ref('')

async function doSearch() {
  const q = query.value.trim()
  if (!q) return

  searching.value = true
  searchError.value = ''
  results.value = []
  searched.value = false

  try {
    results.value = await importApi.search(q)
    searched.value = true
  } catch (err: any) {
    searchError.value = err?.message || 'Search failed.'
  } finally {
    searching.value = false
  }
}

async function doImport(r: SearchResult) {
  importingUrl.value = r.url
  try {
    const res = await importApi.importTrack(r.url)
    toast.add({
      severity: 'success',
      summary: 'Import started',
      detail: `${res.title} added to ingestion.`,
      life: 4000,
    })
    router.push({ name: 'admin.ingestion.review', params: { id: res.draftId } })
  } catch (err: any) {
    toast.add({
      severity: 'error',
      summary: 'Import failed',
      detail: err?.message || 'Could not import track.',
      life: 5000,
    })
  } finally {
    importingUrl.value = ''
  }
}

function formatDuration(seconds: number): string {
  if (!seconds || seconds <= 0) return '—'
  const m = Math.floor(seconds / 60)
  const s = seconds % 60
  return `${m}:${s.toString().padStart(2, '0')}`
}
</script>
