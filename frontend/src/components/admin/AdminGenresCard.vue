<template>
  <div class="rounded-3xl border border-white/10 bg-white/5 p-6 shadow-lg">
    <div class="mb-6 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      <div>
        <h2 class="text-xl font-bold text-white">Genres</h2>
        <p class="mt-2 text-sm text-slate-400">Manage genre names used by tracks.</p>
      </div>

      <Button
        label="Add genre"
        icon="pi pi-plus"
        class="border-0 bg-[#1db954] text-black"
        @click="openCreate"
      />
    </div>

    <DataTable
      :value="genres"
      :loading="loading"
      data-key="id"
      responsive-layout="scroll"
      class="overflow-hidden rounded-2xl"
    >
      <Column field="id" header="ID" sortable />

      <Column field="name" header="Name" sortable />

      <Column header="Actions">
        <template #body="{ data }">
          <div class="flex gap-2">
            <Button icon="pi pi-pencil" severity="secondary" text rounded @click="openEdit(data)" />

            <Button
              icon="pi pi-trash"
              severity="danger"
              text
              rounded
              :loading="deleting"
              @click="handleDelete(data.id)"
            />
          </div>
        </template>
      </Column>

      <template #empty>
        <div class="py-8 text-center text-sm text-slate-400">No genres found.</div>
      </template>
    </DataTable>

    <Dialog
      v-model:visible="dialogVisible"
      modal
      :header="editingGenre ? 'Edit genre' : 'Add genre'"
      class="w-full max-w-lg"
    >
      <div class="flex flex-col gap-4">
        <div>
          <label class="mb-2 block text-sm font-medium text-slate-300"> Name </label>

          <InputText v-model="form.name" class="w-full" placeholder="Genre name" autofocus />
        </div>

        <div
          v-if="errorMessage"
          class="rounded-2xl border border-red-500/20 bg-red-500/10 px-4 py-3 text-sm text-red-300"
        >
          {{ errorMessage }}
        </div>
      </div>

      <template #footer>
        <Button
          label="Cancel"
          severity="secondary"
          outlined
          :disabled="saving"
          @click="dialogVisible = false"
        />

        <Button
          :label="editingGenre ? 'Save changes' : 'Create genre'"
          icon="pi pi-check"
          :loading="saving"
          :disabled="!canSubmit"
          class="border-0 bg-[#1db954] text-black"
          @click="submit"
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import { useToast } from 'primevue/usetoast'

import { useAdminGenres } from '@/composables/admin'
import type { Genre } from '@/services/api/catalog'

const toast = useToast()

const { genres, loading, saving, deleting, fetchGenres, createGenre, updateGenre, deleteGenre } =
  useAdminGenres()

const dialogVisible = ref(false)
const editingGenre = ref<Genre | null>(null)
const errorMessage = ref('')

const form = reactive({
  name: '',
})

const canSubmit = computed(() => {
  return form.name.trim().length > 0 && !saving.value
})

onMounted(() => {
  fetchGenres()
})

function resetForm() {
  form.name = ''
  errorMessage.value = ''
  editingGenre.value = null
}

function openCreate() {
  resetForm()
  dialogVisible.value = true
}

function openEdit(genre: Genre) {
  editingGenre.value = genre
  form.name = genre.name
  errorMessage.value = ''
  dialogVisible.value = true
}

async function submit() {
  if (!canSubmit.value) return

  errorMessage.value = ''

  try {
    if (editingGenre.value) {
      const updated = await updateGenre(editingGenre.value.id, {
        name: form.name.trim(),
      })

      toast.add({
        severity: 'success',
        summary: 'Genre updated',
        detail: updated.name,
        life: 2500,
      })
    } else {
      const created = await createGenre({
        name: form.name.trim(),
      })

      toast.add({
        severity: 'success',
        summary: 'Genre created',
        detail: created.name,
        life: 2500,
      })
    }

    dialogVisible.value = false
    resetForm()
  } catch {
    errorMessage.value = 'Could not save genre. Please try again.'
  }
}

async function handleDelete(id: string | number) {
  const confirmed = window.confirm('Delete this genre?')
  if (!confirmed) return

  try {
    await deleteGenre(id)

    toast.add({
      severity: 'success',
      summary: 'Genre deleted',
      life: 2500,
    })
  } catch {
    toast.add({
      severity: 'error',
      summary: 'Delete failed',
      detail: 'Could not delete genre.',
      life: 3000,
    })
  }
}
</script>
