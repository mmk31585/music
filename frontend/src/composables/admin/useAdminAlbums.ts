import { ref } from 'vue'
import {
  useAlbumsApi,
  type Album,
  type AlbumCreatePayload,
  type AlbumUpdatePayload,
} from '@/services/api/catalog/albums'

export type AlbumFormPayload = AlbumCreatePayload & AlbumUpdatePayload

export function useAdminAlbums() {
  const {
    adminGetAlbums,
    adminCreateAlbum,
    adminUpdateAlbum,
    adminDeleteAlbum,
  } = useAlbumsApi()

  const albums = ref<Album[]>([])
  const loading = ref(false)
  const saving = ref(false)
  const deleting = ref(false)
  const error = ref<unknown>(null)

  async function fetchAlbums() {
    loading.value = true
    error.value = null

    try {
      const response = await adminGetAlbums()
      albums.value = response
      return response
    } catch (err) {
      error.value = err
      throw err
    } finally {
      loading.value = false
    }
  }

  async function createAlbum(payload: AlbumFormPayload) {
    saving.value = true
    error.value = null

    try {
      const created = await adminCreateAlbum(payload)
      albums.value = [created, ...albums.value]
      return created
    } catch (err) {
      error.value = err
      throw err
    } finally {
      saving.value = false
    }
  }

  async function updateAlbum(id: string | number, payload: AlbumFormPayload) {
    saving.value = true
    error.value = null

    try {
      const updated = await adminUpdateAlbum(id, payload)
      albums.value = albums.value.map((album) =>
        String(album.id) === String(id) ? updated : album,
      )
      return updated
    } catch (err) {
      error.value = err
      throw err
    } finally {
      saving.value = false
    }
  }

  async function deleteAlbum(id: string | number) {
    deleting.value = true
    error.value = null

    try {
      await adminDeleteAlbum(id)
      albums.value = albums.value.filter((album) => String(album.id) !== String(id))
    } catch (err) {
      error.value = err
      throw err
    } finally {
      deleting.value = false
    }
  }

  function setAlbums(value: Album[]) {
    albums.value = value
  }

  function resetAlbums() {
    albums.value = []
    error.value = null
  }

  return {
    albums,
    loading,
    saving,
    deleting,
    error,

    fetchAlbums,
    createAlbum,
    updateAlbum,
    deleteAlbum,

    setAlbums,
    resetAlbums,
  }
}
