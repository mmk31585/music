import { ref } from 'vue'
import { useLibraryApi } from '@/services/api/library'
import { useToast } from 'primevue/usetoast'
import type { LibraryAlbum, LibraryArtist, LibraryTrack } from '@/services/api/library'

export const useLibrary = () => {
  const libraryApi = useLibraryApi()
  const toast = useToast()

  const isLoading = ref(false)
  const error = ref('')

  const likedTracks = ref<LibraryTrack[]>([])
  const likedAlbums = ref<LibraryAlbum[]>([])
  const followedArtists = ref<LibraryArtist[]>([])

  const fetchLibrary = async () => {
    isLoading.value = true
    error.value = ''

    try {
      const [tracksRes, albumsRes, artistsRes] = await Promise.all([
        libraryApi.getLikedTracks(),
        libraryApi.getLikedAlbums(),
        libraryApi.getFollowedArtists(),
      ])

      likedTracks.value = tracksRes ?? []
      likedAlbums.value = albumsRes ?? []
      followedArtists.value = artistsRes ?? []
    } catch (err: any) {
      const e = err as Record<string, any>
      const msg = (e.response as Record<string, any> | undefined)?.data?.message || e?.message || 'Failed to load library.'
      error.value = msg
      toast.add({ severity: 'error', summary: 'Library Error', detail: msg, life: 5000 })
    } finally {
      isLoading.value = false
    }
  }

  return {
    isLoading,
    error,
    likedTracks,
    likedAlbums,
    followedArtists,
    fetchLibrary,
  }
}
