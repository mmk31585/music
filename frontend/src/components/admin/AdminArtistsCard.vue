<template>
  <div>
    <!-- Toolbar -->
    <div class="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div class="flex items-center gap-3">
        <span class="text-sm tabular-nums text-slate-500">
          {{ filteredArtists.length }} artist{{ filteredArtists.length !== 1 ? 's' : '' }}
        </span>
        <div class="h-4 w-px bg-white/10" />
        <IconField class="!w-56">
          <InputIcon><i aria-hidden="true" class="pi pi-search text-xs text-slate-500" /></InputIcon>
          <InputText
            v-model="searchQuery"
            placeholder="Search artists..."
            class="!h-9 !w-full !rounded-lg !border-white/[0.08] !bg-white/[0.03] !text-sm !text-white placeholder:!text-slate-600"
          />
        </IconField>
      </div>

      <Button
        label="Add artist"
        icon="pi pi-plus"
        size="small"
        class="!rounded-xl !bg-emerald-500 !px-4 !text-black hover:!bg-emerald-400"
        @click="openCreate"
      />
    </div>

    <!-- Content -->
    <div class="overflow-hidden rounded-2xl border border-white/[0.06] bg-white/[0.02]">
      <!-- Loading State -->
      <div v-if="loading" class="divide-y divide-white/[0.04]">
        <div v-for="i in 5" :key="i" class="flex items-center gap-4 px-5 py-4">
          <div class="h-10 w-10 animate-pulse rounded-full bg-white/[0.06]" />
          <div class="flex-1 space-y-2">
            <div class="h-4 w-32 animate-pulse rounded bg-white/[0.06]" />
            <div class="h-3 w-48 animate-pulse rounded bg-white/[0.04]" />
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <AdminEmptyState
        v-else-if="filteredArtists.length === 0 && !searchQuery"
        icon="pi pi-users"
        title="No artists yet"
        description="Add your first artist to the catalog to get started."
      >
        <template #action>
          <Button
            label="Add artist"
            icon="pi pi-plus"
            size="small"
            class="!rounded-xl !bg-emerald-500 !px-4 !text-black hover:!bg-emerald-400"
            @click="openCreate"
          />
        </template>
      </AdminEmptyState>

      <!-- No Results -->
      <AdminEmptyState
        v-else-if="filteredArtists.length === 0 && searchQuery"
        icon="pi pi-search"
        title="No results found"
        :description="`No artists matching &quot;${searchQuery}&quot;`"
      />

      <!-- Artist List -->
      <div v-else class="divide-y divide-white/[0.04]">
        <div
          v-for="artist in filteredArtists"
          :key="artist.id"
          class="group flex items-center gap-4 px-5 py-4 transition-colors hover:bg-white/[0.02]"
        >
          <!-- Avatar -->
          <div class="relative h-10 w-10 shrink-0 overflow-hidden rounded-full bg-white/[0.06]">
            <img
              v-if="artist.image_url"
              :src="artist.image_url"
              :alt="artist.name"
              class="h-full w-full object-cover"
              @error="($event.target as HTMLImageElement).style.display = 'none'"
            />
            <div
              v-else
              class="flex h-full w-full items-center justify-center text-sm font-semibold text-slate-500"
            >
              {{ artist.name.charAt(0).toUpperCase() }}
            </div>
          </div>

          <!-- Info -->
          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-2">
              <p class="truncate text-sm font-medium text-white">{{ artist.name }}</p>
              <i
                v-if="artist.is_verified"
                class="pi pi-verified text-xs text-emerald-400"
                v-tooltip.top="'Verified'"
              />
            </div>
            <p class="mt-0.5 truncate text-xs text-slate-500">
              <template v-if="artist.bio">{{ artist.bio }}</template>
              <template v-else>No biography</template>
              <template v-if="artist.monthly_listeners">
                <span class="mx-1.5 text-slate-700">·</span>
                {{ formatListeners(artist.monthly_listeners) }} listeners
              </template>
            </p>
          </div>

          <!-- Actions -->
          <div
            class="flex shrink-0 items-center gap-1 opacity-0 transition-opacity group-hover:opacity-100"
          >
            <Button
              icon="pi pi-eye"
              text
              rounded
              size="small"
              class="!h-8 !w-8 !text-slate-400 hover:!text-white"
              v-tooltip.top="'View'"
              @click="router.push({ name: 'admin.artist.detail', params: { id: artist.id } })"
            />
            <Button
              icon="pi pi-refresh"
              text
              rounded
              size="small"
              :loading="enrichingId === artist.id"
              class="!h-8 !w-8 !text-slate-400 hover:!text-amber-400"
              v-tooltip.top="'Enrich'"
              @click="handleEnrich(artist)"
            />
            <Button
              icon="pi pi-pencil"
              text
              rounded
              size="small"
              class="!h-8 !w-8 !text-slate-400 hover:!text-emerald-400"
              v-tooltip.top="'Edit'"
              @click="openEdit(artist)"
            />
            <Button
              icon="pi pi-trash"
              text
              rounded
              size="small"
              class="!h-8 !w-8 !text-slate-400 hover:!text-red-400"
              v-tooltip.top="'Delete'"
              @click="openDeleteConfirm(artist)"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Form Dialog -->
    <ArtistFormDialog
      v-model="showForm"
      :artist="selectedArtist"
      :saving="saving"
      @submit="handleSubmit"
    />

    <!-- Delete Confirm -->
    <AdminDeleteConfirm
      v-model="showDelete"
      title="Delete artist"
      :item-name="deleteTarget?.name"
      :deleting="deleting"
      @confirm="handleDelete"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import { useToast } from 'primevue/usetoast'
import AdminEmptyState from './AdminEmptyState.vue'
import ArtistFormDialog from './ArtistFormDialog.vue'
import AdminDeleteConfirm from './AdminDeleteConfirm.vue'
import { useAdminArtists, type ArtistFormPayload } from '@/composables/admin/useAdminArtists'
import { useArtistsApi } from '@/services/api/catalog/artists'
import type { Artist } from '@/services/api/catalog/artists'

const router = useRouter()
const toast = useToast()
const artistsApi = useArtistsApi()
const {
  artists,
  loading,
  saving,
  deleting,
  fetchArtists,
  createArtist,
  updateArtist,
  deleteArtist,
} = useAdminArtists()

const searchQuery = ref('')
const showForm = ref(false)
const showDelete = ref(false)
const selectedArtist = ref<Artist | null>(null)
const deleteTarget = ref<Artist | null>(null)
const enrichingId = ref<string | number | null>(null)

const filteredArtists = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return artists.value
  return artists.value.filter(
    (a) =>
      a.name.toLowerCase().includes(q) ||
      a.bio?.toLowerCase().includes(q),
  )
})

onMounted(fetchArtists)

function openCreate() {
  selectedArtist.value = null
  showForm.value = true
}

function openEdit(artist: Artist) {
  selectedArtist.value = artist
  showForm.value = true
}

function openDeleteConfirm(artist: Artist) {
  deleteTarget.value = artist
  showDelete.value = true
}

async function handleEnrich(artist: Artist) {
  enrichingId.value = artist.id
  try {
    await artistsApi.adminEnrichArtist(artist.id)
    toast.add({ severity: 'success', summary: 'Artist enriched', detail: 'Data fetched and updated from external sources', life: 3000 })
    await fetchArtists()
  } catch {
    toast.add({ severity: 'error', summary: 'Enrich failed', detail: 'Could not enrich artist. Check external service connectivity.', life: 4000 })
  } finally {
    enrichingId.value = null
  }
}

async function handleSubmit(payload: ArtistFormPayload) {
  try {
    if (selectedArtist.value) {
      await updateArtist(selectedArtist.value.id, payload)
      toast.add({ severity: 'success', summary: 'Artist updated', life: 2500 })
    } else {
      await createArtist(payload)
      toast.add({ severity: 'success', summary: 'Artist created', life: 2500 })
    }
    showForm.value = false
  } catch {
    toast.add({ severity: 'error', summary: 'Operation failed', life: 3000 })
  }
}

async function handleDelete() {
  if (!deleteTarget.value) return
  try {
    await deleteArtist(deleteTarget.value.id)
    toast.add({ severity: 'success', summary: 'Artist deleted', life: 2500 })
    showDelete.value = false
    deleteTarget.value = null
  } catch {
    toast.add({ severity: 'error', summary: 'Delete failed', life: 3000 })
  }
}

function formatListeners(count: number): string {
  if (count >= 1_000_000) return `${(count / 1_000_000).toFixed(1)}M`
  if (count >= 1_000) return `${(count / 1_000).toFixed(1)}K`
  return String(count)
}
</script>
