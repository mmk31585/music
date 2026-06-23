<template>
  <div class="relative mx-auto min-h-screen w-full pb-36">
    <!-- ── Ambient Aurora ── -->
    <div class="aurora-bg pointer-events-none fixed inset-0" aria-hidden="true">
      <div class="aurora-spot-1" />
      <div class="aurora-spot-2" />
    </div>

    <div class="relative z-10 px-4 pt-8 md:px-6 lg:px-8">
      <div class="mx-auto max-w-7xl">
        <!-- ════════════════════════════════════════ -->
        <!-- HERO — Your Library                      -->
        <!-- ════════════════════════════════════════ -->
        <section
          class="relative overflow-hidden rounded-[2.5rem] border border-white/[0.06] bg-gradient-to-br from-[#1db954]/10 via-[#0C0C14] to-black/60 p-8 backdrop-blur-2xl md:p-12"
        >
          <div class="pointer-events-none absolute inset-0 overflow-hidden rounded-[2.5rem]">
            <div class="absolute -top-1/2 -right-1/4 h-80 w-80 rounded-full bg-[#1db954]/8 blur-[120px]" />
            <div class="absolute -bottom-1/2 -left-1/4 h-64 w-64 rounded-full bg-[#a855f7]/8 blur-[100px]" />
          </div>
          <div class="relative">
            <p class="text-[10px] font-bold tracking-[0.35em] text-white/30 uppercase">Your collection</p>
            <h1 class="mt-2 text-4xl font-black text-white md:text-6xl">Library</h1>
            <p class="mt-3 max-w-2xl text-sm text-white/50">
              Liked songs, saved albums, followed artists, playlists, and listening history — all in one place.
            </p>
            <!-- Quick stats -->
            <div v-if="!loading" class="mt-5 flex flex-wrap items-center gap-3">
              <div class="flex items-center gap-1.5 rounded-full border border-white/[0.06] bg-white/[0.04] px-3 py-1.5 text-xs text-white/50">
                <i aria-hidden="true" class="pi pi-heart text-[10px]" /> {{ likedTracks.length }} tracks
              </div>
              <div class="flex items-center gap-1.5 rounded-full border border-white/[0.06] bg-white/[0.04] px-3 py-1.5 text-xs text-white/50">
                <i aria-hidden="true" class="pi pi-images text-[10px]" /> {{ likedAlbums.length }} albums
              </div>
              <div class="flex items-center gap-1.5 rounded-full border border-white/[0.06] bg-white/[0.04] px-3 py-1.5 text-xs text-white/50">
                <i aria-hidden="true" class="pi pi-users text-[10px]" /> {{ followedArtists.length }} artists
              </div>
              <div class="flex items-center gap-1.5 rounded-full border border-white/[0.06] bg-white/[0.04] px-3 py-1.5 text-xs text-white/50">
                <i aria-hidden="true" class="pi pi-list text-[10px]" /> {{ playlists.length }} playlists
              </div>
            </div>
          </div>
        </section>

        <!-- ════════════════════════════════════════ -->
        <!-- TAB NAVIGATION — Pill Tabs               -->
        <!-- ════════════════════════════════════════ -->
        <div class="mt-8">
          <div class="no-scrollbar flex gap-2 overflow-x-auto" role="tablist" aria-label="Library sections">
            <button
              v-for="tab in tabs"
              :key="tab.id"
              role="tab"
              :aria-selected="activeTab === tab.id"
              type="button"
              class="group relative flex shrink-0 items-center gap-2 rounded-full border px-5 py-2.5 text-xs font-bold transition-all duration-300"
              :class="
                activeTab === tab.id
                  ? 'border-[#1db954]/40 bg-[#1db954]/15 text-white shadow-lg shadow-[#1db954]/5'
                  : 'border-white/[0.06] bg-white/[0.03] text-white/50 hover:border-white/[0.12] hover:bg-white/[0.06] hover:text-white'
              "
              @click="activeTab = tab.id"
            >
              <i aria-hidden="true" :class="tab.icon" class="text-sm" />
              {{ tab.label }}
              <span
                class="flex h-4 min-w-[16px] items-center justify-center rounded-full px-1 text-[9px] font-bold"
                :class="activeTab === tab.id ? 'bg-[#1db954]/20 text-[#1db954]' : 'bg-white/[0.06] text-white/30'"
              >
                {{ tab.count }}
              </span>
            </button>
          </div>
        </div>

        <!-- ════════════════════════════════════════ -->
        <!-- TAB CONTENT                              -->
        <!-- ════════════════════════════════════════ -->
        <div class="mt-8">
          <div v-if="loading" class="space-y-6">
            <div v-for="i in 3" :key="i">
              <div class="mb-4 shimmer h-5 w-32 rounded-lg" />
              <div class="flex gap-4">
                <div v-for="j in 4" :key="j" class="shimmer h-44 w-40 shrink-0 rounded-2xl" />
              </div>
            </div>
          </div>

          <!-- ──── TRACKS TAB ──── -->
          <div v-if="activeTab === 'tracks'" role="tabpanel" aria-live="polite">
            <!-- Filter/search -->
            <div class="mb-4 flex items-center justify-between">
              <h2 class="text-lg font-bold text-white">Liked Tracks</h2>
              <div class="relative">
                <i aria-hidden="true" class="pi pi-search absolute left-3 top-1/2 -translate-y-1/2 text-xs text-white/20" />
                <input
                  v-model="trackFilter"
                  type="text"
                  placeholder="Filter tracks..."
                  class="w-48 rounded-full border border-white/[0.06] bg-white/[0.03] py-2 pl-9 pr-4 text-xs text-white outline-none transition placeholder:text-white/20 focus:border-[#1db954]/30 focus:bg-white/[0.06]"
                />
              </div>
            </div>

            <div
              v-if="filteredTracks.length"
              class="overflow-hidden rounded-2xl border border-white/[0.06] bg-white/[0.02] backdrop-blur-sm"
            >
              <TrackRow
                v-for="(item, index) in filteredTracks"
                :key="item.track_id"
                :track="item"
                :index="index"
                :queue="filteredTracks"
              />
            </div>
            <div
              v-else
              class="flex flex-col items-center gap-3 rounded-2xl border border-dashed border-white/[0.06] bg-white/[0.02] px-6 py-16 text-center"
            >
              <div class="flex h-14 w-14 items-center justify-center rounded-2xl bg-white/[0.04]">
                <i aria-hidden="true" class="pi pi-heart text-xl text-white/20" />
              </div>
              <h3 class="text-base font-bold text-white">No liked tracks yet</h3>
              <p class="text-sm text-white/40">
                {{ trackFilter ? 'No tracks match your filter.' : 'Tap the heart icon on any track to save it here.' }}
              </p>
            </div>
          </div>

          <!-- ──── ALBUMS TAB ──── -->
          <div v-if="activeTab === 'albums'" role="tabpanel" aria-live="polite">
            <div class="mb-4 flex items-center justify-between">
              <h2 class="text-lg font-bold text-white">Saved Albums</h2>
              <span class="text-xs text-white/30 tabular-nums">{{ likedAlbums.length }} albums</span>
            </div>

            <div
              v-if="likedAlbums.length"
              class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5"
            >
              <RouterLink
                v-for="album in likedAlbums"
                :key="album.album_id"
                :to="`/album/${album.album_id}`"
                class="group"
              >
                <div class="relative mb-3 aspect-square overflow-hidden rounded-2xl bg-white/5 shadow-lg ring-1 ring-white/[0.06] transition-all duration-300 group-hover:scale-[1.02] group-hover:ring-[#1db954]/30">
                  <img
                    v-if="album.cover_url"
                    :src="album.cover_url"
                    :alt="album.title"
                    loading="lazy"
                    class="h-full w-full object-cover transition duration-300 group-hover:scale-105"
                    @error="onImgError"
                  />
                  <div v-else class="flex h-full items-center justify-center">
                    <i aria-hidden="true" class="pi pi-images text-3xl text-white/20" />
                  </div>
                  <div class="absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 backdrop-blur-sm transition group-hover:opacity-100">
                    <div class="flex h-12 w-12 items-center justify-center rounded-full bg-[#1db954]/90 text-black shadow-xl">
                      <i aria-hidden="true" class="pi pi-play-fill text-lg" />
                    </div>
                  </div>
                </div>
                <p class="truncate text-sm font-semibold text-white">{{ album.title }}</p>
                <p class="truncate text-xs text-white/40">{{ album.artist_name || 'Unknown artist' }}</p>
              </RouterLink>
            </div>
            <div
              v-else
              class="flex flex-col items-center gap-3 rounded-2xl border border-dashed border-white/[0.06] bg-white/[0.02] px-6 py-16 text-center"
            >
              <div class="flex h-14 w-14 items-center justify-center rounded-2xl bg-white/[0.04]">
                <i aria-hidden="true" class="pi pi-images text-xl text-white/20" />
              </div>
              <h3 class="text-base font-bold text-white">No saved albums</h3>
              <p class="text-sm text-white/40">Save albums to your library to find them quickly.</p>
            </div>
          </div>

          <!-- ──── ARTISTS TAB ──── -->
          <div v-if="activeTab === 'artists'" role="tabpanel" aria-live="polite">
            <div class="mb-4 flex items-center justify-between">
              <h2 class="text-lg font-bold text-white">Followed Artists</h2>
              <span class="text-xs text-white/30 tabular-nums">{{ followedArtists.length }} artists</span>
            </div>

            <div
              v-if="followedArtists.length"
              class="grid grid-cols-2 gap-6 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5"
            >
              <RouterLink
                v-for="artist in followedArtists"
                :key="artist.artist_id"
                :to="`/artist/${artist.artist_id}`"
                class="group block text-center"
              >
                <div class="mx-auto mb-3 h-36 w-36 overflow-hidden rounded-full bg-white/[0.04] ring-1 ring-white/[0.06] transition-all duration-300 group-hover:ring-[#1db954]/30">
                  <img
                    v-if="artist.cover_url"
                    :src="artist.cover_url"
                    :alt="artist.name"
                    loading="lazy"
                    class="h-full w-full object-cover transition duration-300 group-hover:scale-105"
                    @error="onImgError"
                  />
                  <div v-else class="flex h-full items-center justify-center">
                    <i aria-hidden="true" class="pi pi-user text-3xl text-white/20" />
                  </div>
                </div>
                <p class="truncate text-sm font-bold text-white">{{ artist.name }}</p>
                <p class="text-xs text-white/40">Artist</p>
              </RouterLink>
            </div>
            <div
              v-else
              class="flex flex-col items-center gap-3 rounded-2xl border border-dashed border-white/[0.06] bg-white/[0.02] px-6 py-16 text-center"
            >
              <div class="flex h-14 w-14 items-center justify-center rounded-2xl bg-white/[0.04]">
                <i aria-hidden="true" class="pi pi-users text-xl text-white/20" />
              </div>
              <h3 class="text-base font-bold text-white">No followed artists</h3>
              <p class="text-sm text-white/40">Follow artists to keep up with their latest releases.</p>
            </div>
          </div>

          <!-- ──── PLAYLISTS TAB ──── -->
          <div v-if="activeTab === 'playlists'" role="tabpanel" aria-live="polite">
            <div class="mb-4 flex items-center justify-between">
              <h2 class="text-lg font-bold text-white">My Playlists</h2>
              <button
                type="button"
                class="inline-flex items-center gap-2 rounded-full bg-[#1db954] px-5 py-2.5 text-xs font-bold text-black transition hover:bg-[#1ed760] hover:scale-105"
                @click="showCreate = true"
              >
                <i aria-hidden="true" class="pi pi-plus text-[10px]" />
                Create
              </button>
            </div>

            <div
              v-if="playlists.length"
              class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4"
            >
              <RouterLink
                v-for="playlist in playlists"
                :key="playlist.id"
                :to="`/playlist/${playlist.id}`"
                class="group rounded-2xl border border-white/[0.06] bg-white/[0.02] p-4 transition-all duration-300 hover:-translate-y-0.5 hover:bg-white/[0.06] hover:border-white/[0.12]"
              >
                <div class="relative mb-3 aspect-square overflow-hidden rounded-xl bg-white/5 ring-1 ring-white/[0.06]">
                  <img
                    v-if="playlist.cover_url"
                    :src="playlist.cover_url"
                    :alt="playlist.name"
                    loading="lazy"
                    class="h-full w-full object-cover transition duration-300 group-hover:scale-105"
                    @error="onImgError"
                  />
                  <PlaylistCoverGrid
                    v-else
                    :covers="[]"
                    :track-count="playlist.track_count"
                    class="h-full w-full"
                  />
                  <div class="absolute inset-0 flex items-center justify-center bg-black/30 opacity-0 backdrop-blur-sm transition group-hover:opacity-100">
                    <div class="flex h-12 w-12 items-center justify-center rounded-full bg-[#1db954]/90 text-black shadow-xl">
                      <i aria-hidden="true" class="pi pi-play-fill text-lg" />
                    </div>
                  </div>
                </div>
                <p class="truncate text-sm font-semibold text-white">{{ playlist.name }}</p>
                <p class="text-xs text-white/40">{{ playlist.track_count }} tracks</p>
              </RouterLink>
            </div>
            <div
              v-else
              class="flex flex-col items-center gap-3 rounded-2xl border border-dashed border-white/[0.06] bg-white/[0.02] px-6 py-16 text-center"
            >
              <div class="flex h-14 w-14 items-center justify-center rounded-2xl bg-white/[0.04]">
                <i aria-hidden="true" class="pi pi-list text-xl text-white/20" />
              </div>
              <h3 class="text-base font-bold text-white">No playlists yet</h3>
              <p class="text-sm text-white/40">Create your first playlist to start organizing your music.</p>
              <button
                type="button"
                class="mt-2 rounded-full bg-[#1db954] px-6 py-2.5 text-sm font-bold text-black transition hover:bg-[#1ed760]"
                @click="showCreate = true"
              >
                <i aria-hidden="true" class="pi pi-plus mr-1 text-xs" />
                Create Playlist
              </button>
            </div>
          </div>

          <!-- ──── HISTORY TAB ──── -->
          <div v-if="activeTab === 'history'" role="tabpanel" aria-live="polite">
            <div class="mb-4 flex items-center justify-between">
              <h2 class="text-lg font-bold text-white">Recently Played</h2>
              <span class="text-xs text-white/30 tabular-nums">{{ recentTracks.length }} tracks</span>
            </div>

            <div
              v-if="recentTracks.length"
              class="overflow-hidden rounded-2xl border border-white/[0.06] bg-white/[0.02] backdrop-blur-sm"
            >
              <TrackRow
                v-for="(item, index) in recentTracks"
                :key="item.track_id"
                :track="item"
                :index="index"
                :queue="recentTracks"
              />
            </div>
            <div
              v-else
              class="flex flex-col items-center gap-3 rounded-2xl border border-dashed border-white/[0.06] bg-white/[0.02] px-6 py-16 text-center"
            >
              <div class="flex h-14 w-14 items-center justify-center rounded-2xl bg-white/[0.04]">
                <i aria-hidden="true" class="pi pi-history text-xl text-white/20" />
              </div>
              <h3 class="text-base font-bold text-white">No history yet</h3>
              <p class="text-sm text-white/40">Start playing tracks and your history will appear here.</p>
              <RouterLink
                to="/discover"
                class="mt-2 inline-flex rounded-full bg-[#1db954] px-6 py-2.5 text-sm font-bold text-black transition hover:bg-[#1ed760]"
              >
                Discover music
              </RouterLink>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ════════════════════════════════════════ -->
    <!-- CREATE PLAYLIST DIALOG                   -->
    <!-- ════════════════════════════════════════ -->
    <Dialog
      v-model:visible="showCreate"
      :modal="true"
      :draggable="false"
      :style="{ maxWidth: '440px', width: '90vw' }"
      :pt="{
        root: 'border-none',
        mask: 'backdrop-blur-sm bg-black/60',
        header: 'border-b border-white/5',
        title: 'text-white text-sm font-bold',
        content: 'p-0',
      }"
    >
      <template #header>
        <div class="flex items-center gap-2 px-1">
          <i aria-hidden="true" class="pi pi-plus text-sm text-[#1db954]" />
          <span>Create Playlist</span>
        </div>
      </template>

      <div class="space-y-5 p-6">
        <div>
          <label class="mb-2 block text-xs font-medium text-white/40">Name</label>
          <InputText
            v-model="newName"
            placeholder="My awesome playlist"
            class="w-full"
          />
        </div>
        <div>
          <label class="mb-2 block text-xs font-medium text-white/40">Description</label>
          <Textarea
            v-model="newDescription"
            placeholder="Optional description"
            rows="3"
            class="w-full"
          />
        </div>
        <div class="flex items-center gap-2">
          <Checkbox v-model="newIsPublic" :binary="true" input-id="playlist-public" />
          <label for="playlist-public" class="text-sm text-white/60">Public playlist</label>
        </div>
        <div class="flex justify-end gap-3 pt-2">
          <Button
            label="Cancel"
            severity="secondary"
            text
            class="text-sm"
            @click="showCreate = false"
          />
          <Button
            label="Create"
            severity="success"
            class="text-sm"
            :disabled="!newName.trim() || creating"
            :loading="creating"
            @click="handleCreate"
          />
        </div>
      </div>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Checkbox from 'primevue/checkbox'
import Button from 'primevue/button'
import { TrackRow, PlaylistCoverGrid } from '@/components/music'
import { useLibraryApi } from '@/services/api/library'
import { usePlaylistsApi } from '@/services/api/playlist'
import type { LibraryTrack, LibraryAlbum, LibraryArtist } from '@/services/api/library'
import type { PlaylistListItem } from '@/services/api/playlist'
import { onImgError } from '@/utils/helpers'

const router = useRouter()
const libraryApi = useLibraryApi()
const playlistsApi = usePlaylistsApi()

const loading = ref(false)
const activeTab = ref('tracks')

// Library data
const likedTracks = ref<LibraryTrack[]>([])
const likedAlbums = ref<LibraryAlbum[]>([])
const followedArtists = ref<LibraryArtist[]>([])
const playlists = ref<PlaylistListItem[]>([])
const recentTracks = ref<LibraryTrack[]>([])
const trackFilter = ref('')

// Create playlist
const showCreate = ref(false)
const creating = ref(false)
const newName = ref('')
const newDescription = ref('')
const newIsPublic = ref(true)

// ── Filtered tracks ──
const filteredTracks = computed(() => {
  if (!trackFilter.value.trim()) return likedTracks.value
  const q = trackFilter.value.toLowerCase()
  return likedTracks.value.filter(
    (t) =>
      t.title.toLowerCase().includes(q) ||
      (t.artist_name || '').toLowerCase().includes(q) ||
      (t.album_title || '').toLowerCase().includes(q),
  )
})

// ── Tab definitions ──
const tabs = computed(() => [
  { id: 'tracks', label: 'Tracks', icon: 'pi pi-heart', count: likedTracks.value.length },
  { id: 'albums', label: 'Albums', icon: 'pi pi-images', count: likedAlbums.value.length },
  { id: 'artists', label: 'Artists', icon: 'pi pi-users', count: followedArtists.value.length },
  { id: 'playlists', label: 'Playlists', icon: 'pi pi-list', count: playlists.value.length },
  { id: 'history', label: 'History', icon: 'pi pi-history', count: recentTracks.value.length },
])

// ── Fetch all data ──
async function fetchAll() {
  loading.value = true
  try {
    const [tracks, albums, artists, plists, recent] = await Promise.all([
      libraryApi.getLikedTracks().catch(() => []),
      libraryApi.getLikedAlbums().catch(() => []),
      libraryApi.getFollowedArtists().catch(() => []),
      playlistsApi.getMyPlaylists().catch(() => []),
      libraryApi.getRecentlyPlayed().catch(() => []),
    ])
    likedTracks.value = Array.isArray(tracks) ? tracks : []
    likedAlbums.value = Array.isArray(albums) ? albums : []
    followedArtists.value = Array.isArray(artists) ? artists : []
    playlists.value = Array.isArray(plists) ? plists : []
    recentTracks.value = Array.isArray(recent) ? recent : []
  } catch (err) {
    console.error('Failed to fetch library data:', err)
  } finally {
    loading.value = false
  }
}

// ── Create playlist ──
async function handleCreate() {
  if (!newName.value.trim() || creating.value) return
  creating.value = true
  try {
    const result = await playlistsApi.createPlaylist({
      name: newName.value.trim(),
      description: newDescription.value.trim() || undefined,
      is_public: newIsPublic.value,
    })
    showCreate.value = false
    newName.value = ''
    newDescription.value = ''
    if (result?.id) {
      await router.push(`/playlist/${result.id}`)
    } else {
      await fetchAll()
    }
  } catch {
    const toast = useToast()
    toast.add({ severity: 'error', summary: 'Failed to create playlist', life: 3000 })
  } finally {
    creating.value = false
  }
}

onMounted(() => {
  void fetchAll()
})
</script>
