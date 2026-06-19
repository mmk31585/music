<template>
  <div>
    <!-- Toolbar -->
    <div class="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div class="flex items-center gap-3">
        <span class="text-sm tabular-nums text-slate-500">
          {{ filteredAlbums.length }} album{{ filteredAlbums.length !== 1 ? 's' : '' }}
        </span>
        <div class="h-4 w-px bg-white/10" />
        <InputText
          v-model="searchQuery"
          placeholder="Search albums..."
          class="!h-9 !w-56 !rounded-lg !border-white/[0.08] !bg-white/[0.03] !text-sm !text-white placeholder:!text-slate-600"
        />
      </div>

      <Button
        label="Add album"
        icon="pi pi-plus"
        size="small"
        class="!rounded-xl !bg-emerald-500 !px-4 !text-black hover:!bg-emerald-400"
        @click="openCreate"
      />
    </div>

    <!-- Content -->
    <div class="overflow-hidden rounded-2xl border border-white/[0.06] bg-white/[0.02]">
      <!-- Loading -->
      <div v-if="loading" class="grid grid-cols-2 gap-4 p-5 sm:grid-cols-3 lg:grid-cols-4">
        <div v-for="i in 8" :key="i" class="space-y-3">
          <div class="aspect-square animate-pulse rounded-xl bg-white/[0.06]" />
          <div class="h-4 w-3/4 animate-pulse rounded bg-white/[0.06]" />
          <div class="h-3 w-1/2 animate-pulse rounded bg-white/[0.04]" />
        </div>
      </div>

      <!-- Empty -->
      <AdminEmptyState
        v-else-if="filteredAlbums.length === 0 && !searchQuery"
        icon="pi pi-book"
        title="No albums yet"
        description="Create your first album to organize tracks."
      >
        <template #action>
          <Button
            label="Add album"
            icon="pi pi-plus"
            size="small"
            class="!rounded-xl !bg-emerald-500 !px-4 !text-black hover:!bg-emerald-400"
            @click="openCreate"
          />
        </template>
      </AdminEmptyState>

      <AdminEmptyState
        v-else-if="filteredAlbums.length === 0 && searchQuery"
        icon="pi pi-search"
        title="No results found"
        :description="`No albums matching &quot;${searchQuery}&quot;`"
      />

      <!-- Album Grid -->
      <div v-else class="grid grid-cols-2 gap-4 p-5 sm:grid-cols-3 lg:grid-cols-4">
        <div
          v-for="album in filteredAlbums"
          :key="album.id"
          class="group cursor-default"
        >
          <!-- Cover -->
          <div
            class="relative aspect-square overflow-hidden rounded-xl border border-white/[0.06] bg-white/[0.04]"
          >
            <img
              v-if="album.cover_url"
              :src="album.cover_url"
              :alt="album.title"
              class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
              @error="($event.target as HTMLImageElement).style.display = 'none'"
            />
            <div
              v-else
              class="flex h-full w-full items-center justify-center"
            >
              <i aria-hidden="true" class="pi pi-image text-3xl text-slate-700" />
            </div>

            <!-- Overlay actions -->
            <div
              class="absolute inset-0 flex items-end justify-end gap-1 bg-gradient-to-t from-black/60 via-transparent p-3 opacity-0 transition-opacity group-hover:opacity-100"
            >
              <Button
                icon="pi pi-eye"
                rounded
                size="small"
                class="!h-8 !w-8 !bg-white/20 !text-white !backdrop-blur-sm hover:!bg-white/30"
                @click="router.push({ name: 'admin.album.detail', params: { id: album.id } })"
              />
              <Button
                icon="pi pi-pencil"
                rounded
                size="small"
                class="!h-8 !w-8 !bg-white/20 !text-white !backdrop-blur-sm hover:!bg-white/30"
                @click="openEdit(album)"
              />
              <Button
                icon="pi pi-trash"
                rounded
                size="small"
                class="!h-8 !w-8 !bg-white/20 !text-white !backdrop-blur-sm hover:!bg-red-500/60"
                @click="openDeleteConfirm(album)"
              />
            </div>
          </div>

          <!-- Info -->
          <div class="mt-2.5 min-w-0">
            <p class="truncate text-sm font-medium text-white">{{ album.title }}</p>
            <p class="mt-0.5 truncate text-xs text-slate-500">
              {{ album.artist_name || 'Unknown artist' }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Form Dialog -->
    <AlbumFormDialog
      v-model="showForm"
      :album="selectedAlbum"
      :saving="saving"
      @submit="handleSubmit"
    />

    <!-- Delete Confirm -->
    <AdminDeleteConfirm
      v-model="showDelete"
      title="Delete album"
      :item-name="deleteTarget?.title"
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
import AlbumFormDialog from './AlbumFormDialog.vue'
import AdminDeleteConfirm from './AdminDeleteConfirm.vue'
import { useAdminAlbums, type AlbumFormPayload } from '@/composables/admin/useAdminAlbums'
import type { Album } from '@/services/api/catalog/albums'

const router = useRouter()
const toast = useToast()
const {
  albums,
  loading,
  saving,
  deleting,
  fetchAlbums,
  createAlbum,
  updateAlbum,
  deleteAlbum,
} = useAdminAlbums()

const searchQuery = ref('')
const showForm = ref(false)
const showDelete = ref(false)
const selectedAlbum = ref<Album | null>(null)
const deleteTarget = ref<Album | null>(null)

const filteredAlbums = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return albums.value
  return albums.value.filter(
    (a) =>
      a.title.toLowerCase().includes(q) ||
      a.artist_name?.toLowerCase().includes(q),
  )
})

onMounted(fetchAlbums)

function openCreate() {
  selectedAlbum.value = null
  showForm.value = true
}

function openEdit(album: Album) {
  selectedAlbum.value = album
  showForm.value = true
}

function openDeleteConfirm(album: Album) {
  deleteTarget.value = album
  showDelete.value = true
}

async function handleSubmit(payload: AlbumFormPayload) {
  try {
    if (selectedAlbum.value) {
      await updateAlbum(selectedAlbum.value.id, payload)
      toast.add({ severity: 'success', summary: 'Album updated', life: 2500 })
    } else {
      await createAlbum(payload)
      toast.add({ severity: 'success', summary: 'Album created', life: 2500 })
    }
    showForm.value = false
  } catch {
    toast.add({ severity: 'error', summary: 'Operation failed', life: 3000 })
  }
}

async function handleDelete() {
  if (!deleteTarget.value) return
  try {
    await deleteAlbum(deleteTarget.value.id)
    toast.add({ severity: 'success', summary: 'Album deleted', life: 2500 })
    showDelete.value = false
    deleteTarget.value = null
  } catch {
    toast.add({ severity: 'error', summary: 'Delete failed', life: 3000 })
  }
}
</script>
