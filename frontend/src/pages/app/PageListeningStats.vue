<template>
  <div class="mx-auto max-w-4xl px-4 pt-20 pb-32 md:px-8">
    <!-- Header -->
    <div class="text-center">
      <h1 class="text-4xl font-black text-white md:text-5xl">آمار شنیدن تو</h1>
      <p class="mt-2 text-sm text-white/40">Your personalized listening statistics</p>
    </div>

    <!-- Period Toggle -->
    <div
      class="mt-8 flex justify-center gap-2 rounded-xl bg-white/4 p-1"
      role="tablist"
    >
      <button
        v-for="p in periods"
        :key="p.key"
        role="tab"
        :aria-selected="activePeriod === p.key"
        class="rounded-lg px-5 py-2 text-sm font-medium transition-all duration-200"
        :class="
          activePeriod === p.key
            ? 'bg-white/10 text-white shadow-lg'
            : 'text-white/30 hover:text-white/50'
        "
        @click="setPeriod(p.key)"
      >
        {{ p.label }}
      </button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="mt-10 space-y-6">
      <div class="grid gap-6 sm:grid-cols-2">
        <SkeletonLoader v-for="i in 2" :key="i" variant="card" class="h-32" />
      </div>
      <SkeletonLoader variant="card" class="h-64" />
      <SkeletonLoader variant="card" class="h-48" />
    </div>

    <!-- Stats Content -->
    <template v-else-if="stats">
      <!-- Big Number Cards -->
      <section class="mt-10 grid gap-5 sm:grid-cols-2">
        <div
          class="group rounded-3xl border border-white/6 bg-linear-to-br from-purple-600/20 to-[#0C0C14] p-8 text-center transition hover:-translate-y-0.5 hover:border-purple-500/30"
        >
          <p class="text-5xl font-black text-white tabular-nums">
            <CountUp :to="stats.total_minutes_listened" />
          </p>
          <p class="mt-2 text-sm text-white/50">دقیقه گوش دادی</p>
        </div>

        <div
          class="group rounded-3xl border border-white/6 bg-linear-to-br from-pink-600/20 to-[#0C0C14] p-8 text-center transition hover:-translate-y-0.5 hover:border-pink-500/30"
        >
          <p class="text-5xl font-black text-white tabular-nums">
            <CountUp :to="stats.total_tracks_played" />
          </p>
          <p class="mt-2 text-sm text-white/50">آهنگ پخش شده</p>
        </div>
      </section>

      <!-- Streak + Discovery -->
      <section class="mt-6 grid gap-5 sm:grid-cols-2">
        <div class="rounded-3xl border border-white/6 bg-white/2 p-6">
          <div class="flex items-center gap-3">
            <span class="text-3xl">🔥</span>
            <div>
              <p class="text-2xl font-black text-white tabular-nums">
                <CountUp :to="stats.longest_listening_streak_days" />
              </p>
              <p class="text-sm text-white/50">روز پشت سر هم گوش دادی!</p>
            </div>
          </div>
        </div>

        <div class="rounded-3xl border border-white/6 bg-white/2 p-6">
          <div class="flex items-center gap-3">
            <span class="text-3xl">🎯</span>
            <div>
              <p class="text-2xl font-black text-white tabular-nums">
                {{ (stats.discovery_score * 100).toFixed(0) }}%
              </p>
              <p class="text-sm text-white/50">امتیاز کاوشگری</p>
              <p class="text-xs text-white/30">
                {{ stats.unique_artists_count }} هنرمند مختلف
              </p>
            </div>
          </div>
        </div>
      </section>

      <!-- Top Tracks -->
      <section v-if="stats.top_tracks.length" class="mt-10">
        <h2 class="mb-5 text-2xl font-black text-white">پرشنیده‌ترین‌های تو</h2>
        <div class="overflow-hidden rounded-3xl border border-white/10 bg-black/20 p-2 backdrop-blur-xs">
          <div
            v-for="(track, i) in stats.top_tracks"
            :key="track.track_id"
            class="flex items-center gap-4 rounded-2xl px-4 py-3 transition hover:bg-white/4"
          >
            <span class="w-8 text-center text-sm font-bold text-slate-500">{{ i + 1 }}</span>
            <div class="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-purple-500/20 text-lg text-purple-400">
              <i aria-hidden="true" class="pi pi-music" />
            </div>
            <div class="min-w-0 flex-1">
              <p class="truncate font-semibold text-white">{{ track.title }}</p>
              <p class="truncate text-sm text-slate-400">{{ track.artist_name }}</p>
            </div>
            <span class="shrink-0 text-xs text-slate-500">{{ track.play_count }} plays</span>
          </div>
        </div>
      </section>

      <!-- Top Artists -->
      <section v-if="stats.top_artists.length" class="mt-10">
        <h2 class="mb-5 text-2xl font-black text-white">خواننده‌های محبوبت</h2>
        <div class="overflow-hidden rounded-3xl border border-white/10 bg-black/20 p-2 backdrop-blur-xs">
          <div
            v-for="(artist, i) in stats.top_artists"
            :key="artist.artist_id"
            class="flex items-center gap-4 rounded-2xl px-4 py-3 transition hover:bg-white/4"
          >
            <span class="w-8 text-center text-sm font-bold text-slate-500">{{ i + 1 }}</span>
            <div class="flex h-12 w-12 shrink-0 items-center justify-center rounded-full bg-pink-500/20 text-lg text-pink-400">
              <i aria-hidden="true" class="pi pi-user" />
            </div>
            <div class="min-w-0 flex-1">
              <p class="truncate font-semibold text-white">{{ artist.artist_name }}</p>
            </div>
            <span class="shrink-0 text-xs text-slate-500">{{ artist.play_count }} plays</span>
          </div>
        </div>
      </section>

      <!-- Top Genres -->
      <section v-if="stats.top_genres.length" class="mt-10">
        <h2 class="mb-5 text-2xl font-black text-white">ژانرهای مورد علاقه‌ت</h2>
        <div class="flex flex-wrap gap-3">
          <div
            v-for="genre in stats.top_genres"
            :key="genre.genre_name"
            class="rounded-2xl border border-white/6 bg-white/4 px-5 py-4 text-center transition hover:-translate-y-0.5 hover:border-emerald-500/30 hover:bg-emerald-500/10"
          >
            <p class="text-sm font-bold text-white">{{ genre.genre_name }}</p>
            <p class="mt-1 text-xs text-slate-500">{{ genre.play_count }} plays</p>
          </div>
        </div>
      </section>

      <!-- Share -->
      <section class="mt-12 text-center">
        <button
          class="inline-flex items-center gap-2 rounded-full bg-sky-500/20 px-6 py-3 text-sm font-bold text-sky-300 transition hover:bg-sky-500/30"
          @click="shareStats"
        >
          <i aria-hidden="true" class="pi pi-telegram" />
          اشتراک‌گذاری در تلگرام
        </button>
      </section>
    </template>

    <!-- Empty -->
    <div
      v-else-if="!loading"
      class="mt-10 flex flex-col items-center gap-4 rounded-3xl border border-white/10 px-6 py-20 text-center"
    >
      <div class="flex h-16 w-16 items-center justify-center rounded-full bg-white/10">
        <i aria-hidden="true" class="pi pi-chart-bar text-2xl text-slate-400" />
      </div>
      <h3 class="text-xl font-bold text-white">No stats yet</h3>
      <p class="text-sm text-slate-400">Start listening to see your listening stats.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRecommendationsApi } from '@/services/api/recommendation'
import { SkeletonLoader } from '@/components/common'
import CountUp from '@/components/common/CountUp.vue'
import type { ListeningStats } from '@/services/api/recommendation'

const api = useRecommendationsApi()
const loading = ref(false)
const stats = ref<ListeningStats | null>(null)
const activePeriod = ref('month')

const periods = [
  { key: 'month', label: 'این ماه' },
  { key: 'year', label: 'امسال' },
  { key: 'all_time', label: 'همیشه' },
]

async function fetchStats(period: string) {
  loading.value = true
  try {
    const response = await api.getListeningStats({ period })
    stats.value = response
  } catch (err) {
    console.error('Failed to fetch stats:', err)
    stats.value = null
  } finally {
    loading.value = false
  }
}

function setPeriod(key: string) {
  activePeriod.value = key
  void fetchStats(key)
}

function shareStats() {
  if (stats.value!) return
  const s = stats.value
  const text = [
    `🎵 *آمار شنیدن من در ${s.period_label}*`,
    '',
    `📊 ${s.total_minutes_listened} دقیقه گوش دادم`,
    `🎶 ${s.total_tracks_played} آهنگ پخش شد`,
    `🔥 ${s.longest_listening_streak_days} روز پشت سر هم`,
    `🎯 امتیاز کاوشگری: ${(s.discovery_score * 100).toFixed(0)}%`,
    '',
    s.top_artists?.length
      ? `🎤 خواننده‌ی محبوب: ${s.top_artists[0]!.artist_name}`
      : '',
    '',
    'via Moja',
  ]
    .filter(Boolean)
    .join('\n')

  const url = `https://t.me/share/url?url=${encodeURIComponent(window.location.href)}&text=${encodeURIComponent(text)}`
  window.open(url, '_blank', 'noopener')
}

onMounted(() => fetchStats('month'))
</script>
