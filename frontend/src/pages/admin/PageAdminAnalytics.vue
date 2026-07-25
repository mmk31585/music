<template>
  <div class="mx-auto w-full max-w-7xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      eyebrow="Insights"
      title="Analytics"
      description="Platform metrics, trends, and performance data."
    >
      <template #actions>
        <div class="flex items-center gap-2">
          <Select
            v-model="period"
            :options="periodOptions"
            option-label="label"
            option-value="value"
            class="w-40"
            size="small"
          />
          <Button
            icon="pi pi-refresh"
            severity="secondary"
            size="small"
            :loading="loading"
            @click="fetchAnalytics"
            aria-label="RefreshCw analytics"
          />
        </div>
      </template>
    </AdminSectionHeader>

    <!-- KPI Cards -->
    <section class="mb-8 grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
      <AdminStatCard
        label="Total Plays"
        :value="stats.totalPlays"
        hint="All-time track plays"
        icon="pi pi-play"
        color="emerald"
        :loading="loading"
      />
      <AdminStatCard
        label="Active Users"
        :value="stats.activeUsers"
        hint="Unique users this period"
        icon="pi pi-users"
        color="blue"
        :loading="loading"
      />
      <AdminStatCard
        label="New Tracks"
        :value="stats.newTracks"
        hint="Tracks added this period"
        icon="pi pi-music"
        color="purple"
        :loading="loading"
      />
      <AdminStatCard
        label="Revenue"
        :value="stats.revenue"
        hint="Estimated revenue (USD)"
        icon="pi pi-dollar"
        color="amber"
        :loading="loading"
      />
    </section>

    <!-- Charts Row -->
    <section class="mb-8 grid gap-6 xl:grid-cols-2">
      <!-- Plays Over Time -->
      <div class="overflow-hidden rounded-2xl border border-white/6 bg-white/2">
        <div class="flex items-center justify-between border-b border-white/6 px-5 py-4">
          <div class="flex items-center gap-3">
            <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-emerald-500/10">
              <TrendingUp aria-hidden="true" class="text-xs text-emerald-400"  />
            </div>
            <h2 class="text-base font-semibold text-white">Plays Over Time</h2>
          </div>
        </div>
        <div class="p-5">
          <div v-if="loading" class="flex items-end gap-2" style="height: 160px">
            <div
              v-for="i in 12"
              :key="i"
              class="flex-1 animate-pulse rounded-t bg-white/6"
              :style="{ height: `${20 + Math.random() * 80}%` }"
            />
          </div>
          <div v-else-if="chartData.dailyPlays.length === 0" class="flex flex-col items-center py-8 text-center">
            <BarChart3 aria-hidden="true" class="text-2xl text-slate-700"  />
            <p class="mt-2 text-sm text-slate-500">No play data for this period</p>
          </div>
          <div v-else class="flex items-end gap-2" style="height: 160px">
            <div
              v-for="(day, i) in chartData.dailyPlays"
              :key="i"
              class="group relative flex-1"
            >
              <div
                class="w-full rounded-t bg-emerald-500/60 transition-all hover:bg-emerald-400/80"
                :style="{ height: `${Math.max(4, (day.count / maxPlayCount) * 100)}%` }"
              >
                <div
                  class="absolute -top-8 left-1/2 -translate-x-1/2 whitespace-nowrap rounded-md bg-slate-800 px-2 py-1 text-[10px] text-white opacity-0 shadow-lg transition-opacity group-hover:opacity-100"
                >
                  {{ day.count.toLocaleString() }} plays
                </div>
              </div>
              <p class="mt-1.5 text-center text-[10px] text-slate-600">
                {{ day.label }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- User Growth -->
      <div class="overflow-hidden rounded-2xl border border-white/6 bg-white/2">
        <div class="flex items-center justify-between border-b border-white/6 px-5 py-4">
          <div class="flex items-center gap-3">
            <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-blue-500/10">
              <Users aria-hidden="true" class="text-xs text-blue-400"  />
            </div>
            <h2 class="text-base font-semibold text-white">User Growth</h2>
          </div>
        </div>
        <div class="p-5">
          <div v-if="loading" class="flex items-end gap-2" style="height: 160px">
            <div
              v-for="i in 12"
              :key="i"
              class="flex-1 animate-pulse rounded-t bg-white/6"
              :style="{ height: `${20 + Math.random() * 60}%` }"
            />
          </div>
          <div v-else-if="chartData.userGrowth.length === 0" class="flex flex-col items-center justify-center py-8 text-center">
            <BarChart3 aria-hidden="true" class="text-2xl text-slate-700"  />
            <p class="mt-2 text-sm text-slate-500">No user data for this period</p>
          </div>
          <div v-else class="flex items-end gap-2" style="height: 160px">
            <div
              v-for="(month, i) in chartData.userGrowth"
              :key="i"
              class="group relative flex-1"
            >
              <div
                class="w-full rounded-t bg-blue-500/60 transition-all hover:bg-blue-400/80"
                :style="{ height: `${Math.max(4, (month.count / maxUserGrowth) * 100)}%` }"
              >
                <div
                  class="absolute -top-8 left-1/2 -translate-x-1/2 whitespace-nowrap rounded bg-slate-800 px-2 py-1 text-[10px] text-white opacity-0 shadow-lg transition-opacity group-hover:opacity-100"
                >
                  +{{ month.count.toLocaleString() }}
                </div>
              </div>
              <p class="mt-1.5 text-center text-[10px] text-slate-600">
                {{ month.label }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Top Content Table -->
    <section class="overflow-hidden rounded-2xl border border-white/6 bg-white/2">
      <div class="flex items-center justify-between border-b border-white/6 px-5 py-4">
        <div class="flex items-center gap-3">
          <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-purple-500/10">
            <Star aria-hidden="true" class="text-xs text-purple-400"  />
          </div>
          <h2 class="text-base font-semibold text-white">Top Played Tracks</h2>
        </div>
      </div>

      <div v-if="loading" class="divide-y divide-white/4">
        <div v-for="i in 5" :key="i" class="flex items-center gap-3 px-5 py-4">
          <div class="h-3 w-6 animate-pulse rounded bg-white/6" />
          <div class="h-10 w-10 animate-pulse rounded-lg bg-white/6" />
          <div class="flex-1 space-y-1.5">
            <div class="h-3.5 w-32 animate-pulse rounded bg-white/6" />
            <div class="h-3 w-48 animate-pulse rounded bg-white/4" />
          </div>
          <div class="h-3.5 w-12 animate-pulse rounded bg-white/4" />
        </div>
      </div>

      <div v-else-if="topTracks.length === 0" class="py-12 text-center">
        <Star aria-hidden="true" class="text-2xl text-slate-700"  />
        <p class="mt-2 text-sm text-slate-500">No track data available yet</p>
      </div>

      <div v-else class="divide-y divide-white/4">
        <div
          v-for="(track, i) in topTracks"
          :key="track.id ?? i"
          class="flex items-center gap-3 px-5 py-3 transition-colors hover:bg-white/2"
        >
          <span class="w-5 text-center text-xs font-medium tabular-nums" :class="i < 3 ? 'text-emerald-400' : 'text-slate-500'">
            {{ i + 1 }}
          </span>
          <div class="h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/4">
            <img
              v-if="track.cover_url"
              :src="track.cover_url"
              :alt="track.title || 'Cover'"
              class="h-full w-full object-cover"
              @error="hideBrokenImage"
            />
            <div v-else class="flex h-full w-full items-center justify-center">
              <Music aria-hidden="true" class="text-xs text-slate-700"  />
            </div>
          </div>
          <div class="min-w-0 flex-1">
            <p class="truncate text-sm font-medium text-white">{{ track.title || 'Untitled' }}</p>
            <p class="truncate text-xs text-slate-500">{{ track.artist_name || 'Unknown' }}</p>
          </div>
          <span class="text-xs tabular-nums text-slate-500">
            {{ (track.play_count ?? 0).toLocaleString() }} plays
          </span>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { BarChart3, Music, Star, TrendingUp, Users } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import Select from 'primevue/select'
import Button from 'primevue/button'
import AdminSectionHeader from '@/components/admin/AdminSectionHeader.vue'
import AdminStatCard from '@/components/admin/AdminStatCard.vue'
import { client } from '@/composables/useRequest'

interface DailyPlay {
  label: string
  count: number
}

interface MonthlyGrowth {
  label: string
  count: number
}

interface TopTrack {
  id: number
  title: string
  artist_name: string
  cover_url: string | null
  play_count: number
}

interface ChartData {
  dailyPlays: DailyPlay[]
  userGrowth: MonthlyGrowth[]
}

interface AnalyticsStats {
  totalPlays: number
  activeUsers: number
  newTracks: number
  revenue: number
}

const loading = ref(false)
const period = ref<'7d' | '30d' | '90d' | '1y'>('30d')

const periodOptions = [
  { label: 'Last 7 days', value: '7d' },
  { label: 'Last 30 days', value: '30d' },
  { label: 'Last 90 days', value: '90d' },
  { label: 'Last year', value: '1y' },
]

const stats = ref<AnalyticsStats>({
  totalPlays: 0,
  activeUsers: 0,
  newTracks: 0,
  revenue: 0,
})

const chartData = ref<ChartData>({
  dailyPlays: [],
  userGrowth: [],
})

const topTracks = ref<TopTrack[]>([])

const maxPlayCount = computed(() => Math.max(...chartData.value.dailyPlays.map(d => d.count), 1))
const maxUserGrowth = computed(() => Math.max(...chartData.value.userGrowth.map(m => m.count), 1))

function hideBrokenImage(event: Event) {
  const image = event.target as HTMLImageElement | null
  if (image) image.style.display = 'none'
}

async function fetchAnalytics() {
  loading.value = true
  try {
    const [statsRes, chartRes, topRes] = await Promise.allSettled([
      client.get('/admin/analytics/stats', { params: { period: period.value } }).then(r => r.data.data),
      client.get('/admin/analytics/charts', { params: { period: period.value } }).then(r => r.data.data),
      client.get('/admin/analytics/top-tracks', { params: { period: period.value } }).then(r => r.data.data),
    ])

    if (statsRes.status === 'fulfilled' && statsRes.value) {
      stats.value = statsRes.value
    }
    if (chartRes.status === 'fulfilled' && chartRes.value) {
      chartData.value = chartRes.value
    }
    if (topRes.status === 'fulfilled' && topRes.value) {
      topTracks.value = Array.isArray(topRes.value) ? topRes.value : []
    }
  } catch {
    // Silent — charts will show empty state
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  void fetchAnalytics()
})
</script>