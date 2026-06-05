import { ref } from 'vue'
import {
  useGenresApi,
  type Genre,
  type GenreCreatePayload,
  type GenreUpdatePayload,
} from '@/services/api/catalog/genres'

export type GenreFormPayload = GenreCreatePayload & GenreUpdatePayload

export function useAdminGenres() {
  const {
    adminGetGenres,
    adminCreateGenre,
    adminUpdateGenre,
    adminDeleteGenre,
  } = useGenresApi()

  const genres = ref<Genre[]>([])
  const loading = ref(false)
  const saving = ref(false)
  const deleting = ref(false)
  const error = ref<unknown>(null)

  async function fetchGenres() {
    loading.value = true
    error.value = null

    try {
      const response = await adminGetGenres()
      genres.value = response
      return response
    } catch (err) {
      error.value = err
      throw err
    } finally {
      loading.value = false
    }
  }

  async function createGenre(payload: GenreFormPayload) {
    saving.value = true
    error.value = null

    try {
      const created = await adminCreateGenre(payload)
      genres.value = [created, ...genres.value]
      return created
    } catch (err) {
      error.value = err
      throw err
    } finally {
      saving.value = false
    }
  }

  async function updateGenre(id: string | number, payload: GenreFormPayload) {
    saving.value = true
    error.value = null

    try {
      const updated = await adminUpdateGenre(id, payload)
      genres.value = genres.value.map((genre) =>
        String(genre.id) === String(id) ? updated : genre,
      )
      return updated
    } catch (err) {
      error.value = err
      throw err
    } finally {
      saving.value = false
    }
  }

  async function deleteGenre(id: string | number) {
    deleting.value = true
    error.value = null

    try {
      await adminDeleteGenre(id)
      genres.value = genres.value.filter((genre) => String(genre.id) !== String(id))
    } catch (err) {
      error.value = err
      throw err
    } finally {
      deleting.value = false
    }
  }

  function setGenres(value: Genre[]) {
    genres.value = value
  }

  function resetGenres() {
    genres.value = []
    error.value = null
  }

  return {
    genres,
    loading,
    saving,
    deleting,
    error,

    fetchGenres,
    createGenre,
    updateGenre,
    deleteGenre,

    setGenres,
    resetGenres,
  }
}
