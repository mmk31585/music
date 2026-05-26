<template>
  <div class="mx-auto w-full max-w-7xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      eyebrow="Catalog"
      title="Tracks"
      description="Create, edit, and delete tracks in the catalog."
    >
      <template #actions>
        <Button
          label="Add track"
          icon="pi pi-plus"
          class="border-0 bg-[#1db954] text-black"
          @click="openCreate"
        />
      </template>
    </AdminSectionHeader>

    <div class="mb-4 flex items-center justify-between gap-3">
      <span class="text-sm text-slate-400"> {{ tracks.length }} tracks </span>
    </div>

    <div class="overflow-hidden rounded-3xl border border-white/10 bg-white/5">
      <DataTable :value="tracks" :loading="loading" data-key="id" class="admin-tracks-table">
        <Column field="title" header="Title" sortable />
        <Column field="artist_name" header="Artist" sortable />
        <Column field="album_title" header="Album" sortable />
        <Column field="duration_seconds" header="Duration" :body="durationBody" />

        <Column header="Actions" :body="actionsBody" style="width: 150px" />
      </DataTable>
    </div>

    <TrackFormDialog v-model="showForm" :track="selectedTrack" @submit="handleSubmitTrack" />
  </div>
</template>

<script setup lang="ts">
import { h, onMounted, ref } from 'vue'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import { AdminSectionHeader } from '@/components/admin'
import TrackFormDialog from '@/components/admin/TrackFormDialog.vue'
import { useAdminTracks } from '@/composables/admin/useAdminTracks'
import type { Track, TrackFormPayload } from '@/services/api/catalog'

const { tracks, loading, fetchTracks, createTrack, updateTrack, deleteTrack } = useAdminTracks()
const confirm = useConfirm()
const toast = useToast()

const showForm = ref(false)
const selectedTrack = ref<Track | null>(null)

onMounted(fetchTracks)

function openCreate() {
  selectedTrack.value = null
  showForm.value = true
}

function openEdit(track: Track) {
  selectedTrack.value = track
  showForm.value = true
}

async function handleSubmitTrack(payload: TrackFormPayload) {
  try {
    if (payload.id) {
      await updateTrack(payload.id, payload)
      toast.add({
        severity: 'success',
        summary: 'Track updated',
        life: 2000,
      })
    } else {
      await createTrack(payload)
      toast.add({
        severity: 'success',
        summary: 'Track created',
        life: 2000,
      })
    }
  } catch {
    toast.add({
      severity: 'error',
      summary: 'Operation failed',
      life: 2500,
    })
  }
}

function confirmDelete(track: Track) {
  confirm.require({
    message: `Delete "${track.title}"?`,
    header: 'Confirm delete',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: async () => {
      try {
        await deleteTrack(track.id)
        toast.add({
          severity: 'success',
          summary: 'Track deleted',
          life: 2000,
        })
      } catch {
        toast.add({
          severity: 'error',
          summary: 'Delete failed',
          life: 2500,
        })
      }
    },
  })
}

function durationBody(rowData: Track) {
  const value = rowData.duration_seconds
  if (!value) return '--:--'
  const mins = Math.floor(value / 60)
  const secs = value % 60
  return `${mins}:${String(secs).padStart(2, '0')}`
}

function actionsBody(rowData: Track) {
  return h('div', { class: 'flex justify-end gap-2' }, [
    h(Button, {
      icon: 'pi pi-pencil',
      text: true,
      rounded: true,
      class: 'text-emerald-300',
      onClick: () => openEdit(rowData),
    }),
    h(Button, {
      icon: 'pi pi-trash',
      text: true,
      rounded: true,
      class: 'text-red-300',
      onClick: () => confirmDelete(rowData),
    }),
  ])
}
</script>

<style scoped>
.admin-tracks-table :deep(.p-datatable-header),
.admin-tracks-table :deep(.p-datatable-thead > tr > th),
.admin-tracks-table :deep(.p-datatable-tbody > tr > td) {
  border-color: rgb(255 255 255 / 0.1);
  background: transparent;
  color: rgb(226 232 240);
}
.admin-tracks-table :deep(.p-datatable-wrapper) {
  background: transparent;
}
</style>
