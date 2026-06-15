<template>
  <div class="mx-auto w-full max-w-7xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      eyebrow="Overview"
      title="Admin dashboard"
      description="A central place to manage uploaded media and catalog content."
    >
      <template #actions>
        <RouterLink
          to="/admin/media"
          class="inline-flex items-center gap-2 rounded-xl bg-emerald-500 px-4 py-2 text-sm font-semibold text-black transition-colors hover:bg-emerald-400"
        >
          <i class="pi pi-upload text-xs" />
          Upload media
        </RouterLink>
      </template>
    </AdminSectionHeader>

    <div
      v-if="error"
      class="mb-6 rounded-2xl border border-red-500/20 bg-red-500/10 px-5 py-4 text-sm text-red-300"
    >
      <div class="flex items-start gap-3">
        <i class="pi pi-exclamation-triangle mt-0.5 text-xs" />
        <div>
          <p class="font-semibold">Dashboard data could not be fully loaded.</p>
          <p class="mt-1 text-xs text-red-300/80">
            Some API responses may not match the expected frontend schema.
          </p>
        </div>
      </div>
    </div>

    <!-- Stats -->
    <section class="mb-8 grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
      <AdminStatCard
        label="Tracks"
        :value="trackCount"
        hint="Total in catalog"
        icon="pi pi-play-circle"
        color="emerald"
        :loading="loading"
      />

      <AdminStatCard
        label="Artists"
        :value="artistCount"
        hint="Total in catalog"
        icon="pi pi-users"
        color="blue"
        :loading="loading"
      />

      <AdminStatCard
        label="Albums"
        :value="albumCount"
        hint="Total in catalog"
        icon="pi pi-book"
        color="purple"
        :loading="loading"
      />

      <AdminStatCard
        label="Genres"
        :value="genreCount"
        hint="Total in catalog"
        icon="pi pi-tags"
        color="amber"
        :loading="loading"
      />
    </section>

    <!-- Ingestion Stats -->
    <section class="mb-8 grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
      <AdminStatCard
        label="Published This Month"
        :value="ingestionStats?.publishedThisMonth ?? 0"
        hint="Tracks published via ingestion"
        icon="pi pi-check-circle"
        color="emerald"
        :loading="loading"
      />

      <AdminStatCard
        label="Pending Review"
        :value="ingestionStats?.pendingReview ?? 0"
        hint="Drafts awaiting review"
        icon="pi pi-eye"
        color="blue"
        :loading="loading"
      />

      <AdminStatCard
        label="Total Drafts"
        :value="ingestionStats?.totalDrafts ?? 0"
        hint="All ingestion drafts"
        icon="pi pi-file"
        color="purple"
        :loading="loading"
      />

      <RouterLink
        to="/admin/ingestion"
        class="flex items-center justify-center gap-2 rounded-2xl border border-white/[0.06] bg-white/[0.02] px-5 py-8 transition-colors hover:bg-white/[0.04]"
      >
        <i class="pi pi-arrow-right text-sm text-primary" />
        <span class="text-sm font-medium text-white">Go to Ingestion</span>
      </RouterLink>
    </section>

    <!-- Quick Actions -->
    <section class="mb-8">
      <CatalogQuickActions />
    </section>

    <!-- Content Panels -->
    <section class="grid gap-6 xl:grid-cols-[1.2fr_0.8fr]">
      <!-- Recent Tracks -->
      <div class="overflow-hidden rounded-2xl border border-white/[0.06] bg-white/[0.02]">
        <div class="flex items-center justify-between border-b border-white/[0.06] px-5 py-4">
          <div class="flex items-center gap-3">
            <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-emerald-500/10">
              <i class="pi pi-play-circle text-xs text-emerald-400" />
            </div>

            <div>
              <h2 class="text-base font-semibold text-white">Recent tracks</h2>
              <p class="text-xs text-slate-500">Latest additions to the catalog</p>
            </div>
          </div>

          <RouterLink
            to="/admin/tracks"
            class="text-xs font-medium text-emerald-400 transition-colors hover:text-emerald-300"
          >
            View all →
          </RouterLink>
        </div>

        <div v-if="loading" class="divide-y divide-white/[0.04]">
          <div v-for="i in 5" :key="i" class="flex items-center gap-3 px-5 py-4">
            <div class="h-10 w-10 animate-pulse rounded-lg bg-white/[0.06]" />
            <div class="flex-1 space-y-1.5">
              <div class="h-3.5 w-32 animate-pulse rounded bg-white/[0.06]" />
              <div class="h-3 w-48 animate-pulse rounded bg-white/[0.04]" />
            </div>
            <div class="h-3.5 w-10 animate-pulse rounded bg-white/[0.04]" />
          </div>
        </div>

        <div v-else-if="recentTracks.length === 0" class="py-12 text-center">
          <i class="pi pi-play-circle text-2xl text-slate-700" />
          <p class="mt-2 text-sm text-slate-500">No tracks yet</p>
        </div>

        <div v-else class="divide-y divide-white/[0.04]">
          <div
            v-for="(track, i) in recentTracks"
            :key="track.id ?? i"
            class="group flex items-center gap-3 px-5 py-3 transition-colors hover:bg-white/[0.02]"
          >
            <span class="w-5 text-center text-xs text-slate-600 tabular-nums">
              {{ i + 1 }}
            </span>

            <div class="h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/[0.04]">
              <img
                v-if="track.cover_url"
                :src="track.cover_url"
                :alt="track.title || 'Track cover'"
                class="h-full w-full object-cover"
                @error="hideBrokenImage"
              />

              <div v-else class="flex h-full w-full items-center justify-center">
                <i class="pi pi-music text-xs text-slate-700" />
              </div>
            </div>

            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-medium text-white">
                {{ track.title || 'Untitled track' }}
              </p>

              <p class="truncate text-xs text-slate-500">
                {{ track.artist_name || 'Unknown artist' }}
              </p>
            </div>

            <span class="text-xs text-slate-600 tabular-nums">
              {{ formatDuration(track.duration_seconds) }}
            </span>
          </div>
        </div>
      </div>

      <!-- Right column -->
      <div class="space-y-6">
        <!-- Activity / Tips -->
        <div class="overflow-hidden rounded-2xl border border-white/[0.06] bg-white/[0.02]">
          <div class="border-b border-white/[0.06] px-5 py-4">
            <div class="flex items-center gap-3">
              <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-blue-500/10">
                <i class="pi pi-info-circle text-xs text-blue-400" />
              </div>

              <h2 class="text-base font-semibold text-white">Getting started</h2>
            </div>
          </div>

          <div class="p-5">
            <ul class="space-y-3">
              <li v-for="(tip, i) in tips" :key="i" class="flex items-start gap-3">
                <div
                  class="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full text-xs font-semibold"
                  :class="
                    tip.done
                      ? 'bg-emerald-500/10 text-emerald-400'
                      : 'bg-white/[0.06] text-slate-500'
                  "
                >
                  <i v-if="tip.done" class="pi pi-check text-[10px]" />
                  <span v-else>{{ i + 1 }}</span>
                </div>

                <div>
                  <p class="text-sm text-white" :class="{ 'line-through opacity-50': tip.done }">
                    {{ tip.label }}
                  </p>

                  <p class="mt-0.5 text-xs text-slate-500">
                    {{ tip.description }}
                  </p>
                </div>
              </li>
            </ul>
          </div>
        </div>

        <!-- System Health -->
        <div class="overflow-hidden rounded-2xl border border-white/[0.06] bg-white/[0.02]">
          <div class="border-b border-white/[0.06] px-5 py-4">
            <div class="flex items-center gap-3">
              <div
                class="flex h-8 w-8 items-center justify-center rounded-lg"
                :class="error ? 'bg-red-500/10' : 'bg-emerald-500/10'"
              >
                <i
                  class="pi pi-server text-xs"
                  :class="error ? 'text-red-400' : 'text-emerald-400'"
                />
              </div>

              <h2 class="text-base font-semibold text-white">System</h2>
            </div>
          </div>

          <div class="divide-y divide-white/[0.04]">
            <div class="flex items-center justify-between px-5 py-3">
              <span class="text-sm text-slate-400">API Status</span>

              <div v-if="error" class="flex items-center gap-1.5">
                <div class="h-2 w-2 rounded-full bg-red-400" />
                <span class="text-xs text-red-400">Partial error</span>
              </div>

              <div v-else class="flex items-center gap-1.5">
                <div class="h-2 w-2 rounded-full bg-emerald-400" />
                <span class="text-xs text-emerald-400">Online</span>
              </div>
            </div>

            <div class="flex items-center justify-between px-5 py-3">
              <span class="text-sm text-slate-400">Catalog data</span>

              <span class="text-xs text-slate-500"> {{ totalEntityCount }} total entities </span>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import AdminSectionHeader from '@/components/admin/AdminSectionHeader.vue'
import AdminStatCard from '@/components/admin/AdminStatCard.vue'
import CatalogQuickActions from '@/components/admin/CatalogQuickActions.vue'
import { useTracksApi, type Track } from '@/services/api/catalog/tracks'
import { useIngestionApi } from '@/services/api/ingestion/routes'
import type { IngestionStats } from '@/services/api/ingestion/types'

const { getTracks } = useTracksApi()
const ingestionApi = useIngestionApi()

const loading = ref(false)
const error = ref<unknown>(null)

const tracks = ref<Track[]>([])
const ingestionStats = ref<IngestionStats | null>(null)
const stats = ref<{ track_count: number; artist_count: number; album_count: number; genre_count: number } | null>(null)

const trackCount = computed(() => stats.value?.track_count ?? tracks.value.length)
const artistCount = computed(() => stats.value?.artist_count ?? 0)
const albumCount = computed(() => stats.value?.album_count ?? 0)
const genreCount = computed(() => stats.value?.genre_count ?? 0)

const totalEntityCount = computed(() => {
  return trackCount.value + artistCount.value + albumCount.value + genreCount.value
})

const recentTracks = computed(() => tracks.value.slice(0, 8))

const tips = computed(() => [
  {
    label: 'Upload audio files',
    description: 'Use the media upload page to add audio.',
    done: trackCount.value > 0,
  },
  {
    label: 'Add artists',
    description: 'Create artist profiles for your catalog.',
    done: artistCount.value > 0,
  },
  {
    label: 'Organize into albums',
    description: 'Group tracks into album collections.',
    done: albumCount.value > 0,
  },
  {
    label: 'Categorize with genres',
    description: 'Add genre tags for better discovery.',
    done: genreCount.value > 0,
  },
])

function toArray<T>(value: unknown): T[] {
  if (Array.isArray(value)) {
    return value as T[]
  }

  if (
    value &&
    typeof value === 'object' &&
    'data' in value &&
    Array.isArray((value as { data?: unknown }).data)
  ) {
    return (value as { data: T[] }).data
  }

  if (
    value &&
    typeof value === 'object' &&
    'items' in value &&
    Array.isArray((value as { items?: unknown }).items)
  ) {
    return (value as { items: T[] }).items
  }

  return []
}

function formatDuration(value?: number | string | null): string {
  const seconds = Number(value)

  if (!Number.isFinite(seconds) || seconds <= 0) {
    return '—'
  }

  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)

  return `${mins}:${String(secs).padStart(2, '0')}`
}

function hideBrokenImage(event: Event) {
  const image = event.target as HTMLImageElement | null

  if (image) {
    image.style.display = 'none'
  }
}

async function fetchDashboard(): Promise<void> {
  loading.value = true
  error.value = null

  try {
    const [statsResult, tracksResult, ingestionResult] = await Promise.allSettled([
      fetch('/api/v1/admin/dashboard/stats').then(r => r.json()).then(r => r.data),
      getTracks(),
      ingestionApi.getIngestionStats(),
    ])

    if (statsResult.status === 'fulfilled' && statsResult.value) {
      stats.value = statsResult.value
    } else {
      console.error('Failed to fetch dashboard stats:', statsResult.status === 'rejected' ? statsResult.reason : 'no data')
    }

    if (tracksResult.status === 'fulfilled') {
      tracks.value = toArray<Track>(tracksResult.value)
    } else {
      console.error('Failed to fetch dashboard tracks:', tracksResult.reason)
      tracks.value = []
      error.value = tracksResult.reason
    }

    if (ingestionResult.status === 'fulfilled') {
      ingestionStats.value = ingestionResult.value
    } else {
      console.error('Failed to fetch ingestion stats:', ingestionResult.reason)
      ingestionStats.value = null
    }
  } catch (err) {
    console.error('Failed to fetch dashboard:', err)
    error.value = err
    tracks.value = []
    ingestionStats.value = null
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  void fetchDashboard()
})
</script>
