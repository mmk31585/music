<template>
  <div class="mx-auto w-full max-w-7xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      eyebrow="Performance"
      title="Cache Management"
      description="View and manage platform caches, Redis stats, and warm-up operations."
    >
      <template #actions>
        <Button
          type="button"
          severity="warn"
          size="small"
          icon="pi pi-trash"
          :loading="clearing"
          @click="clearAll"
        >
          Clear all caches
        </Button>
      </template>
    </AdminSectionHeader>

    <!-- Cache Stats Cards -->
    <section class="mb-8 grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
      <AdminStatCard
        label="Redis Memory"
        :value="cacheStats.redisMemoryMB"
        :hint="`Used / ${cacheStats.redisMaxMemoryMB} MB`"
        icon="pi pi-database"
        color="blue"
        :loading="loading"
      />
      <AdminStatCard
        label="Cache Keys"
        :value="cacheStats.totalKeys"
        hint="Total keys in Redis"
        icon="pi pi-key"
        color="purple"
        :loading="loading"
      />
      <AdminStatCard
        label="Hit Rate"
        :value="cacheStats.hitRatePercent"
        hint="Cache hit rate (last 24h)"
        icon="pi pi-check-circle"
        color="emerald"
        :loading="loading"
      />
      <AdminStatCard
        label="Miss Rate"
        :value="cacheStats.missRatePercent"
        hint="Cache miss rate (last 24h)"
        icon="pi pi-exclamation-circle"
        color="amber"
        :loading="loading"
      />
    </section>

    <!-- Cache Sections -->
    <section class="mb-8 grid gap-6 xl:grid-cols-2">
      <!-- Cache Groups -->
      <div class="overflow-hidden rounded-2xl border border-white/6 bg-white/2">
        <div class="border-b border-white/6 px-5 py-4">
          <div class="flex items-center gap-3">
            <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-purple-500/10">
              <Layers aria-hidden="true" class="text-xs text-purple-400"  />
            </div>
            <h2 class="text-base font-semibold text-white">Cache Groups</h2>
          </div>
        </div>

        <div v-if="loading" class="divide-y divide-white/4">
          <div v-for="i in 4" :key="i" class="flex items-center justify-between px-5 py-4">
            <div class="space-y-1.5">
              <div class="h-3.5 w-28 animate-pulse rounded bg-white/6" />
              <div class="h-3 w-20 animate-pulse rounded bg-white/4" />
            </div>
            <div class="h-7 w-20 animate-pulse rounded-lg bg-white/6" />
          </div>
        </div>

        <div v-else-if="cacheGroups.length === 0" class="py-12 text-center">
          <Layers aria-hidden="true" class="text-2xl text-slate-700"  />
          <p class="mt-2 text-sm text-slate-500">No cache groups available</p>
        </div>

        <div v-else class="divide-y divide-white/4">
          <div
            v-for="group in cacheGroups"
            :key="group.key"
            class="flex items-center justify-between px-5 py-4 transition-colors hover:bg-white/2"
          >
            <div>
              <p class="text-sm font-medium text-white">{{ group.label }}</p>
              <p class="mt-0.5 text-xs text-slate-500">
                {{ group.keys.toLocaleString() }} keys
              </p>
            </div>
            <Button
              type="button"
              severity="secondary"
              size="small"
              icon="pi pi-refresh"
              :loading="clearingKey === group.key"
              @click="clearGroup(group.key)"
            >
              Clear
            </Button>
          </div>
        </div>
      </div>

      <!-- Server Info -->
      <div class="overflow-hidden rounded-2xl border border-white/6 bg-white/2">
        <div class="border-b border-white/6 px-5 py-4">
          <div class="flex items-center gap-3">
            <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-emerald-500/10">
              <Server aria-hidden="true" class="text-xs text-emerald-400"  />
            </div>
            <h2 class="text-base font-semibold text-white">Server Status</h2>
          </div>
        </div>

        <div v-if="loading" class="space-y-3 p-5">
          <div v-for="i in 4" :key="i" class="h-4 w-full animate-pulse rounded bg-white/6" />
        </div>

        <div v-else class="divide-y divide-white/4">
          <div class="flex items-center justify-between px-5 py-3">
            <span class="text-sm text-slate-400">Redis Server</span>
            <div class="flex items-center gap-1.5">
              <div
                class="h-2 w-2 rounded-full"
                :class="serverOnline ? 'bg-emerald-400' : 'bg-red-400'"
              />
              <span class="text-xs" :class="serverOnline ? 'text-emerald-400' : 'text-red-400'">
                {{ serverOnline ? 'Online' : 'Offline' }}
              </span>
            </div>
          </div>
          <div class="flex items-center justify-between px-5 py-3">
            <span class="text-sm text-slate-400">Uptime</span>
            <span class="text-xs text-slate-500 tabular-nums">{{ serverUptime }}</span>
          </div>
          <div class="flex items-center justify-between px-5 py-3">
            <span class="text-sm text-slate-400">Connected Clients</span>
            <span class="text-xs text-slate-500 tabular-nums">{{ cacheStats.connectedClients }}</span>
          </div>
          <div class="flex items-center justify-between px-5 py-3">
            <span class="text-sm text-slate-400">Cache Version</span>
            <span class="text-xs text-slate-500">{{ cacheVersion }}</span>
          </div>
        </div>

        <div class="border-t border-white/6 px-5 py-4">
          <p class="mb-3 text-xs font-medium text-slate-500 uppercase tracking-wider">Actions</p>
          <div class="flex flex-wrap gap-2">
            <Button
              type="button"
              severity="secondary"
              size="small"
              icon="pi pi-refresh"
              :loading="warming"
              @click="warmCache"
            >
              Warm cache
            </Button>
            <Button
              type="button"
              severity="danger"
              size="small"
              icon="pi pi-trash"
              :loading="flushing"
              @click="flushRedis"
            >
              Flush Redis
            </Button>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { Layers, Server } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import Button from 'primevue/button'
import AdminSectionHeader from '@/components/admin/AdminSectionHeader.vue'
import AdminStatCard from '@/components/admin/AdminStatCard.vue'
import { client } from '@/composables/useRequest'
import { useToast } from 'primevue/usetoast'

interface CacheGroup {
  key: string
  label: string
  keys: number
}

interface CacheStats {
  redisMemoryMB: number
  redisMaxMemoryMB: number
  totalKeys: number
  hitRatePercent: number
  missRatePercent: number
  connectedClients: number
}

const loading = ref(false)
const clearing = ref(false)
const clearingKey = ref<string | null>(null)
const warming = ref(false)
const flushing = ref(false)
const serverOnline = ref(true)
const serverUptime = ref('--')
const cacheVersion = ref('1.0')

const toast = useToast()

const cacheStats = ref<CacheStats>({
  redisMemoryMB: 0,
  redisMaxMemoryMB: 0,
  totalKeys: 0,
  hitRatePercent: 0,
  missRatePercent: 0,
  connectedClients: 0,
})

const cacheGroups = ref<CacheGroup[]>([
  { key: 'catalog', label: 'Catalog (tracks, albums, artists)', keys: 0 },
  { key: 'session', label: 'User sessions', keys: 0 },
  { key: 'search', label: 'Search results', keys: 0 },
  { key: 'media', label: 'Media metadata', keys: 0 },
  { key: 'player', label: 'Player queue & state', keys: 0 },
])

async function fetchCacheStats() {
  loading.value = true
  try {
    const res = await client.get('/admin/cache/stats').then(r => r.data)
    if (res.data) {
      cacheStats.value = { ...cacheStats.value, ...res.data }
    }
    const groupsRes = await client.get('/admin/cache/groups').then(r => r.data)
    if (groupsRes.data && Array.isArray(groupsRes.data)) {
      cacheGroups.value = groupsRes.data
    }
    const infoRes = await client.get('/admin/cache/info').then(r => r.data)
    if (infoRes.data) {
      serverOnline.value = infoRes.data.online ?? true
      serverUptime.value = infoRes.data.uptime ?? '--'
      cacheVersion.value = infoRes.data.version ?? '1.0'
    }
  } catch {
    serverOnline.value = false
    toast.add({ severity: 'warn', summary: 'Cache server unavailable', detail: 'Could not reach Redis. Showing placeholder data.', life: 5000 })
  } finally {
    loading.value = false
  }
}

async function clearGroup(key: string) {
  clearingKey.value = key
  try {
    await client.post(`/admin/cache/clear/${key}`)
    toast.add({ severity: 'success', summary: `Cleared "${key}" cache`, life: 3000 })
    await fetchCacheStats()
  } catch {
    toast.add({ severity: 'error', summary: 'Failed to clear cache group', life: 3000 })
  } finally {
    clearingKey.value = null
  }
}

async function clearAll() {
  clearing.value = true
  try {
    await client.post('/admin/cache/clear-all')
    toast.add({ severity: 'success', summary: 'All caches cleared', life: 3000 })
    await fetchCacheStats()
  } catch {
    toast.add({ severity: 'error', summary: 'Failed to clear all caches', life: 3000 })
  } finally {
    clearing.value = false
  }
}

async function warmCache() {
  warming.value = true
  try {
    await client.post('/admin/cache/warm')
    toast.add({ severity: 'success', summary: 'Cache warm initiated', life: 3000 })
  } catch {
    toast.add({ severity: 'error', summary: 'Failed to warm cache', life: 3000 })
  } finally {
    warming.value = false
  }
}

async function flushRedis() {
  flushing.value = true
  try {
    await client.post('/admin/cache/flush')
    toast.add({ severity: 'success', summary: 'Redis flushed successfully', life: 3000 })
    await fetchCacheStats()
  } catch {
    toast.add({ severity: 'error', summary: 'Failed to flush Redis', life: 3000 })
  } finally {
    flushing.value = false
  }
}

onMounted(() => {
  void fetchCacheStats()
})
</script>