<template>
  <div class="mx-auto max-w-5xl px-4 pt-16 pb-24 md:px-8">
    <!-- Back -->
    <button
      class="mb-6 inline-flex items-center gap-1.5 text-sm text-white/40 transition hover:text-white/70"
      @click="router.push({ name: 'social' })"
    >
      &larr; بازگشت
    </button>

    <SkeletonLoader v-if="loading" variant="card" class="h-64" />

    <div v-else-if="error" class="rounded-2xl bg-white/3 p-12 text-center">
      <p class="text-sm text-white/40">{{ error }}</p>
    </div>

    <template v-else-if="detail">
      <!-- Cover Hero -->
      <div
        class="relative mb-8 overflow-hidden rounded-2xl border border-white/6"
        :style="coverBg"
      >
        <div class="absolute inset-0 bg-linear-to-t from-surface-base via-surface-base/70 to-transparent" />
        <div class="relative z-10 flex flex-col gap-6 p-8 pt-48">
          <div>
            <div class="flex items-center gap-3">
              <h1 class="text-3xl font-black text-white md:text-4xl">{{ detail.club.name }}</h1>
              <span
                v-if="detail.club.genre"
                class="rounded-full bg-white/10 px-3 py-0.5 text-[11px] font-medium text-white/60"
              >
                {{ detail.club.genre }}
              </span>
            </div>
            <div class="mt-2 flex items-center gap-4 text-sm text-white/40">
              <span>{{ detail.club.member_count }} عضو</span>
              <span>{{ detail.track_count }} آهنگ</span>
            </div>
          </div>

          <div class="flex flex-wrap items-center gap-3">
            <button
              v-if="!detail.is_member"
              class="inline-flex items-center gap-2 rounded-xl bg-spotify px-6 py-2.5 text-sm font-bold text-black transition hover:bg-spotify-hover"
              @click="handleJoin"
            >
              عضو شدم ✓
            </button>

            <template v-else>
              <button
                class="inline-flex items-center gap-2 rounded-xl border border-white/10 bg-white/5 px-6 py-2.5 text-sm font-bold text-white transition hover:bg-white/10"
                @click="handleLaunchParty"
              >
                شروع گوش دادن گروهی
              </button>
              <button
                class="rounded-xl border border-red-500/20 px-4 py-2.5 text-xs font-medium text-red-400 transition hover:bg-red-500/10"
                @click="showLeaveConfirm = true"
              >
                خروج از کلاب
              </button>
            </template>
          </div>
        </div>
      </div>

      <div class="grid gap-8 lg:grid-cols-3">
        <!-- Main content -->
        <div class="space-y-8 lg:col-span-2">
          <!-- About -->
          <section class="rounded-2xl bg-white/3 p-6">
            <h2 class="mb-3 text-sm font-bold uppercase tracking-wider text-white/30">درباره</h2>
            <p class="text-sm leading-relaxed text-white/60">
              {{ detail.club.description || 'هنوز توضیحی ثبت نشده.' }}
            </p>
          </section>

          <!-- Shared Playlist -->
          <section class="rounded-2xl bg-white/3 p-6">
            <div class="mb-4 flex items-center justify-between">
              <h2 class="text-sm font-bold uppercase tracking-wider text-white/30">
                پلی‌لیست مشترک
              </h2>
              <button
                v-if="detail.is_member && playlistId"
                class="inline-flex items-center gap-1.5 rounded-lg bg-spotify/10 px-3 py-1.5 text-xs font-semibold text-spotify transition hover:bg-spotify/20"
                @click="showTrackPicker = true"
              >
                افزودن آهنگ
              </button>
            </div>
    <TrackList
      v-if="mappedTracks.length > 0"
      :tracks="mappedTracks"
    />
            <div
              v-else
              class="flex flex-col items-center gap-3 py-12 text-center"
            >
              <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-white/4">
                <Music aria-hidden="true" class="text-xl text-slate-500"  />
              </div>
              <p class="text-sm text-white/40">هنوز آهنگی به پلی‌لیست اضافه نشده.</p>
            </div>
          </section>

          <!-- Members -->
          <section class="rounded-2xl bg-white/3 p-6">
            <h2 class="mb-4 text-sm font-bold uppercase tracking-wider text-white/30">
              اعضا ({{ detail.members.length }})
            </h2>
            <div class="flex flex-wrap gap-3">
              <div
                v-for="m in visibleMembers"
                :key="m.id"
                class="group relative"
              >
                <div
                  class="flex h-10 w-10 items-center justify-center rounded-full bg-spotify/20 text-xs font-bold text-spotify transition hover:bg-spotify/30"
                  :title="m.user_id"
                >
                  {{ initials(m.user_id) }}
                </div>
              </div>
              <div
                v-if="overflowCount > 0"
                class="flex h-10 w-10 items-center justify-center rounded-full bg-white/5 text-[11px] font-medium text-white/40"
                :title="`+${overflowCount} more`"
              >
                +{{ overflowCount }}
              </div>
            </div>
          </section>

          <!-- Discussions (Phase 6) -->
          <section class="rounded-2xl bg-white/3 p-6">
            <div class="mb-4 flex items-center justify-between">
              <h2 class="text-sm font-bold uppercase tracking-wider text-white/30">
                بحث و گفتگو
              </h2>
              <button
                v-if="detail.is_member"
                class="inline-flex items-center gap-1.5 rounded-lg bg-spotify/10 px-3 py-1.5 text-xs font-semibold text-spotify transition hover:bg-spotify/20"
                @click="showCreateDiscussion = true"
              >
                بحث جدید
              </button>
            </div>
            <div v-if="discussionsLoading" class="flex items-center justify-center py-8">
              <span class="inline-block h-5 w-5 animate-spin rounded-full border-2 border-white/20 border-t-accent" />
            </div>
            <div
              v-else-if="discussionsError"
              class="rounded-xl bg-red-500/10 px-4 py-3 text-sm text-red-400"
            >
              {{ discussionsError }}
              <button
                class="mr-2 underline hover:text-red-300"
                @click="retryLoadDiscussions"
              >Retry</button>
            </div>
            <ClubDiscussionThread
              v-else
              :discussions="discussions"
              :current-user-id="authUserId"
              :current-user-role="detail?.member_role ?? ''"
              @delete="handleDeleteDiscussion"
            />
            <button
              v-if="discussions.length > 0 && discussions.length >= discussionLimit"
              class="mt-4 w-full rounded-lg bg-white/5 py-2.5 text-xs font-medium text-white/40 transition hover:bg-white/10"
              @click="loadMoreDiscussions"
            >
              بیشتر
            </button>
          </section>
        </div>

        <!-- Sidebar -->
        <div class="space-y-6">
          <div class="rounded-2xl bg-white/3 p-6">
            <h3 class="mb-3 text-xs font-bold uppercase tracking-wider text-white/30">اطلاعات</h3>
            <div class="space-y-3 text-sm">
              <div class="flex justify-between">
                <span class="text-white/40">وضعیت</span>
                <span class="text-white/70">{{ detail.club.is_public ? 'عمومی' : 'خصوصی' }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-white/40">اعضا</span>
                <span class="text-white/70">{{ detail.club.member_count }} / {{ detail.club.max_members }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-white/40">آهنگ‌ها</span>
                <span class="text-white/70">{{ detail.track_count }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-white/40">پست‌ها</span>
                <span class="text-white/70">{{ detail.post_count }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- Track Picker Dialog -->
    <TrackPickerDialog
      :visible="showTrackPicker"
      title="افزودن آهنگ به پلی‌لیست کلاب"
      @select="handleAddTrack"
      @update:visible="showTrackPicker = $event"
    />

    <!-- Leave Confirm Dialog -->
    <Teleport to="body">
      <div
        v-if="showLeaveConfirm"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs"
        @click.self="showLeaveConfirm = false"
      >
        <div class="glass-strong mx-4 w-full max-w-sm rounded-2xl p-8 text-center">
          <p class="mb-6 text-sm text-white/60">
            مطمئنی می‌خوای از این کلاب خارج بشی؟
          </p>
          <div class="flex gap-3">
            <button
              class="flex-1 rounded-xl bg-white/5 py-3 text-sm font-medium text-white/50 transition hover:bg-white/10"
              @click="showLeaveConfirm = false"
            >
              انصراف
            </button>
            <button
              class="flex-1 rounded-xl bg-red-500 py-3 text-sm font-bold text-white transition hover:bg-red-600"
              @click="handleLeave"
            >
              خروج
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <CreateDiscussionDialog
      :visible="showCreateDiscussion"
      :club-id="clubId"
      @close="showCreateDiscussion = false"
      @created="onDiscussionCreated"
    />
  </div>
</template>

<script setup lang="ts">
import { Music } from 'lucide-vue-next'
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { SkeletonLoader } from '@/components/common'
import { TrackList } from '@/components/music'
import TrackPickerDialog from '@/components/social/TrackPickerDialog.vue'
import { useSocialApi } from '@/services/api/social'
import { usePlaylistsApi } from '@/services/api/playlist'
import { useCollaborativePlaylist } from '@/composables/useCollaborativePlaylist'
import { useUserAuthStore } from '@/stores'
import { useAppToast } from '@/composables/useAppToast'
import type { ClubDetailResponse } from '@/services/api/social'
import type { ClubDiscussion } from '@/services/api/social'
import type { PlaylistTrackItem } from '@/services/api/playlist'
import type { Track } from '@/services/api/catalog/tracks'
import ClubDiscussionThread from '@/components/social/ClubDiscussionThread.vue'
import CreateDiscussionDialog from '@/components/social/CreateDiscussionDialog.vue'

const router = useRouter()
const route = useRoute()
const socialApi = useSocialApi()
const playlistApi = usePlaylistsApi()
const toast = useAppToast()

const auth = useUserAuthStore()

const clubId = route.params.id as string
const authUserId = computed(() => String(auth.user?.id ?? ''))

const loading = ref(true)
const error = ref('')
const detail = ref<ClubDetailResponse | null>(null)
const showTrackPicker = ref(false)
const showLeaveConfirm = ref(false)

const playlistTracks = ref<PlaylistTrackItem[]>([])

const collab = useCollaborativePlaylist('', () => loadDetail(), () => loadDetail())

const mappedTracks = computed(() => playlistTracks.value.map(t => ({
  id: t.track_id,
  title: t.title,
  duration_seconds: t.duration_seconds ?? 0,
  audio_url: t.audio_url ?? null,
  cover_url: t.cover_url ?? null,
  audio_media_id: null,
  cover_media_id: null,
  artist_id: null,
  album_id: null,
  genre_id: null,
  artist_name: t.artist_name ?? null,
  album_title: null,
  genres: [],
  play_count: 0,
  track_number: null,
  explicit: false,
  is_public: true,
  created_at: null,
  updated_at: null,
})))

const coverBg = computed(() => {
  if (detail.value?.club.cover_url) {
    return { background: 'linear-gradient(135deg, #1a1a2e 0%, #16213e 100%)' }
  }
  return {
    backgroundImage: `url(${detail.value?.club.cover_url})`,
    backgroundSize: 'cover',
    backgroundPosition: 'center',
  }
})

const playlistId = computed(() => detail.value?.club.playlist_id ?? '')

const visibleMembers = computed(() => detail.value?.members.slice(0, 12) ?? [])
const overflowCount = computed(() => {
  const total = detail.value?.members.length ?? 0
  return Math.max(0, total - 12)
})

function initials(userId: string): string {
  if (!userId) return '?'
  return userId.charAt(0).toUpperCase()
}

const discussions = ref<ClubDiscussion[]>([])
const showCreateDiscussion = ref(false)
const discussionLimit = ref(20)
const discussionOffset = ref(0)
const discussionsLoading = ref(false)
const discussionsError = ref('')

async function loadDiscussions() {
  discussionsLoading.value = true
  discussionsError.value = ''
  try {
    const data = await socialApi.listClubDiscussions(clubId, { limit: discussionLimit.value, offset: discussionOffset.value })
    const items = Array.isArray(data) ? data : []
    if (discussionOffset.value === 0) {
      discussions.value = items
    } else {
      discussions.value = [...discussions.value, ...items]
    }
  } catch (err: any) {
    discussionsError.value = err?.response?.data?.message || err.message || 'Failed to load discussions'
    toast.error(discussionsError.value)
  } finally {
    discussionsLoading.value = false
  }
}

function retryLoadDiscussions() {
  discussionOffset.value = 0
  loadDiscussions()
}

function loadMoreDiscussions() {
  discussionOffset.value += discussionLimit.value
  loadDiscussions()
}

async function handleDeleteDiscussion(discussionId: string) {
  try {
    await socialApi.deleteClubDiscussion(discussionId)
    discussions.value = discussions.value.filter(d => d.id !== discussionId)
    toast.success('Discussion deleted')
  } catch (err: any) {
    toast.apiError(err, 'Failed to delete discussion')
  }
}

async function loadDetail() {
  try {
    const res = await socialApi.getClubDetail(clubId)
    detail.value = res

    if (res.club.playlist_id) {
      const playlistDetail = await playlistApi.getPlaylist(res.club.playlist_id)
      playlistTracks.value = playlistDetail.tracks ?? []
    }
  } catch (err: any) {
    error.value = err?.response?.data?.message || err.message || 'خطا در بارگذاری کلاب'
    toast.error(error.value)
  } finally {
    loading.value = false
  }
}

async function handleJoin() {
  try {
    await socialApi.joinClub(clubId)
    toast.success('Joined club!')
    await loadDetail()
    discussionOffset.value = 0
    loadDiscussions()
  } catch (err: any) {
    toast.apiError(err, 'Failed to join club')
  }
}

async function handleLeave() {
  try {
    showLeaveConfirm.value = false
    await socialApi.leaveClub(clubId)
    toast.info('Left club')
    await loadDetail()
  } catch (err: any) {
    toast.apiError(err, 'Failed to leave club')
  }
}

async function handleLaunchParty() {
  try {
    const party = await socialApi.launchParty(clubId, { title: detail.value?.club.name, is_public: true })
    if (party?.id) {
      toast.success('Party launched!')
      router.push({ name: 'social.party', params: { id: party.id } })
    }
  } catch (err: any) {
    toast.apiError(err, 'Failed to launch party')
  }
}

async function handleAddTrack(track: Track) {
  try {
    const pid = playlistId.value
    if (!pid) return
    toast.info('Adding track...')
    await playlistApi.addTrack(pid, { track_id: String(track.id) })
    showTrackPicker.value = false
    toast.success('Track added to playlist')
    await loadDetail()
  } catch (err: any) {
    toast.apiError(err, 'Failed to add track')
  }
}

function onDiscussionCreated() {
  showCreateDiscussion.value = false
  discussionOffset.value = 0
  loadDiscussions()
}

onMounted(() => {
  loadDetail()
  loadDiscussions()
})

onUnmounted(() => {
  collab.teardown()
})
</script>
