<template>
  <div>
    <!-- Toolbar -->
    <div class="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div class="flex items-center gap-3">
        <span class="text-sm tabular-nums text-slate-500">
          {{ genres.length }} genre{{ genres.length !== 1 ? 's' : '' }}
        </span>
      </div>

      <Button
        label="Add genre"
        icon="pi pi-plus"
        size="small"
        class="!rounded-xl !bg-emerald-500 !px-4 !text-black hover:!bg-emerald-400"
        @click="openCreate"
      />
    </div>

    <!-- Content -->
    <div class="overflow-hidden rounded-2xl border border-white/[0.06] bg-white/[0.02]">
      <!-- Loading State -->
      <div v-if="loading" class="p-6">
        <div class="flex flex-wrap gap-3">
          <div
            v-for="i in 8"
            :key="i"
            class="h-9 animate-pulse rounded-full bg-white/[0.06]"
            :style="{ width: `${60 + Math.random() * 60}px` }"
          />
        </div>
      </div>

      <!-- Empty State -->
      <AdminEmptyState
        v-else-if="genres.length === 0"
        icon="pi pi-tags"
        title="No genres yet"
        description="Create your first music genre to categorize tracks."
      >
        <template #action>
          <Button
            label="Add genre"
            icon="pi pi-plus"
            size="small"
            class="!rounded-xl !bg-emerald-500 !px-4 !text-black hover:!bg-emerald-400"
            @click="openCreate"
          />
        </template>
      </AdminEmptyState>

      <!-- Genre Grid -->
      <div v-else class="p-5">
        <div class="flex flex-wrap gap-2">
          <div
            v-for="genre in genres"
            :key="genre.id"
            class="group flex items-center gap-2 rounded-full border border-white/[0.08] bg-white/[0.03] py-1.5 pl-4 pr-2 transition-all hover:border-emerald-500/20 hover:bg-emerald-500/5"
          >
            <span class="text-sm text-slate-300 group-hover:text-white">{{ genre.name }}</span>

            <div class="flex items-center gap-0.5 opacity-0 transition-opacity group-hover:opacity-100">
              <Button
                icon="pi pi-pencil"
                text
                rounded
                size="small"
                class="!h-6 !w-6 !text-xs !text-slate-500 hover:!text-emerald-400"
                @click="openEdit(genre)"
              />
              <Button
                icon="pi pi-times"
                text
                rounded
                size="small"
                class="!h-6 !w-6 !text-xs !text-slate-500 hover:!text-red-400"
                @click="openDeleteConfirm(genre)"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Inline Create/Edit -->
    <Dialog
      v-model:visible="showForm"
      modal
      :draggable="false"
      :style="{ width: '400px' }"
      :pt="{
        root: { class: '!border-white/[0.06] !bg-[#141414] !rounded-2xl !shadow-2xl' },
        header: { class: '!bg-transparent !border-0 !pb-2' },
        content: { class: '!bg-transparent !px-6 !pt-0 !pb-2' },
        footer: { class: '!bg-transparent !border-0' },
        mask: { class: '!backdrop-blur-sm' },
      }"
    >
      <template #header>
        <div class="flex items-center gap-3">
          <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-purple-500/10">
            <i aria-hidden="true" class="pi pi-tag text-purple-400" />
          </div>
          <h3 class="text-base font-semibold text-white">
            {{ selectedGenre ? 'Edit genre' : 'New genre' }}
          </h3>
        </div>
      </template>

      <div class="mt-4">
        <label class="mb-1.5 block text-xs font-medium text-slate-400">
          Name <span class="text-red-400">*</span>
        </label>
        <InputText
          v-model="genreName"
          placeholder="e.g. Hip-Hop, Jazz, Electronic"
          class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white placeholder:!text-slate-600 focus:!border-emerald-500/40"
          autofocus
          @keydown.enter="handleSubmitGenre"
        />
      </div>

      <template #footer>
        <div class="flex justify-end gap-2">
          <Button
            label="Cancel"
            text
            :disabled="saving"
            class="!text-slate-400 hover:!text-white"
            @click="showForm = false"
          />
          <Button
            :label="selectedGenre ? 'Save' : 'Create'"
            :loading="saving"
            class="!rounded-xl !bg-emerald-500 !text-black hover:!bg-emerald-400"
            @click="handleSubmitGenre"
          />
        </div>
      </template>
    </Dialog>

    <!-- Delete Confirm -->
    <AdminDeleteConfirm
      v-model="showDelete"
      title="Delete genre"
      :item-name="deleteTarget?.name"
      :deleting="deleting"
      @confirm="handleDelete"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import { useToast } from 'primevue/usetoast'
import AdminEmptyState from './AdminEmptyState.vue'
import AdminDeleteConfirm from './AdminDeleteConfirm.vue'
import { useAdminGenres } from '@/composables/admin/useAdminGenres'
import type { Genre } from '@/services/api/catalog/genres'

const toast = useToast()
const {
  genres,
  loading,
  saving,
  deleting,
  fetchGenres,
  createGenre,
  updateGenre,
  deleteGenre,
} = useAdminGenres()

const showForm = ref(false)
const showDelete = ref(false)
const selectedGenre = ref<Genre | null>(null)
const deleteTarget = ref<Genre | null>(null)
const genreName = ref('')

onMounted(fetchGenres)

function openCreate() {
  selectedGenre.value = null
  genreName.value = ''
  showForm.value = true
}

function openEdit(genre: Genre) {
  selectedGenre.value = genre
  genreName.value = genre.name
  showForm.value = true
}

function openDeleteConfirm(genre: Genre) {
  deleteTarget.value = genre
  showDelete.value = true
}

async function handleSubmitGenre() {
  if (!genreName.value.trim()) return

  try {
    if (selectedGenre.value) {
      await updateGenre(selectedGenre.value.id, { name: genreName.value.trim() })
      toast.add({ severity: 'success', summary: 'Genre updated', life: 2500 })
    } else {
      await createGenre({ name: genreName.value.trim() })
      toast.add({ severity: 'success', summary: 'Genre created', life: 2500 })
    }
    showForm.value = false
  } catch {
    toast.add({ severity: 'error', summary: 'Operation failed', life: 3000 })
  }
}

async function handleDelete() {
  if (!deleteTarget.value) return
  try {
    await deleteGenre(deleteTarget.value.id)
    toast.add({ severity: 'success', summary: 'Genre deleted', life: 2500 })
    showDelete.value = false
    deleteTarget.value = null
  } catch {
    toast.add({ severity: 'error', summary: 'Delete failed', life: 3000 })
  }
}
</script>
