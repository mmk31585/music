import { ref } from 'vue'
import {
  useArtistsApi,
  type Artist,
  type ArtistCreatePayload,
  type ArtistUpdatePayload,
} from '@/services/api/catalog/artists'

export type ArtistFormPayload = ArtistCreatePayload & ArtistUpdatePayload

export function useAdminArtists() {
  const {
    adminGetArtists,
    adminCreateArtist,
    adminUpdateArtist,
    adminDeleteArtist,
  } = useArtistsApi()

  const artists = ref<Artist[]>([])
  const loading = ref(false)
  const saving = ref(false)
  const deleting = ref(false)
  const error = ref<unknown>(null)

  async function fetchArtists() {
    loading.value = true
    error.value = null

    try {
      const response = await adminGetArtists()
      artists.value = response
      return response
    } catch (err) {
      error.value = err
      throw err
    } finally {
      loading.value = false
    }
  }

  async function createArtist(payload: ArtistFormPayload) {
    saving.value = true
    error.value = null

    try {
      const created = await adminCreateArtist(payload)
      artists.value = [created, ...artists.value]
      return created
    } catch (err) {
      error.value = err
      throw err
    } finally {
      saving.value = false
    }
  }

  async function updateArtist(id: string | number, payload: ArtistFormPayload) {
    saving.value = true
    error.value = null

    try {
      const updated = await adminUpdateArtist(id, payload)
      artists.value = artists.value.map((artist) =>
        String(artist.id) === String(id) ? updated : artist,
      )
      return updated
    } catch (err) {
      error.value = err
      throw err
    } finally {
      saving.value = false
    }
  }

  async function deleteArtist(id: string | number) {
    deleting.value = true
    error.value = null

    try {
      await adminDeleteArtist(id)
      artists.value = artists.value.filter((artist) => String(artist.id) !== String(id))
    } catch (err) {
      error.value = err
      throw err
    } finally {
      deleting.value = false
    }
  }

  function setArtists(value: Artist[]) {
    artists.value = value
  }

  function resetArtists() {
    artists.value = []
    error.value = null
  }

  return {
    artists,
    loading,
    saving,
    deleting,
    error,

    fetchArtists,
    createArtist,
    updateArtist,
    deleteArtist,

    setArtists,
    resetArtists,
  }
}
