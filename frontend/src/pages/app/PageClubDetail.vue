<template>
  <div class="mx-auto max-w-5xl space-y-6 px-4 pt-20 pb-24 md:px-8">
    <button
      class="inline-flex items-center gap-1.5 text-sm text-white/40 transition hover:text-white/70"
      @click="goBack"
    >
      &larr; Back to Social
    </button>

    <div v-if="loading" class="space-y-6">
      <SkeletonLoader variant="card" />
      <SkeletonLoader variant="card" class="h-48" />
    </div>

    <div v-else-if="error" class="rounded-2xl bg-white/[0.03] p-12 text-center">
      <p class="text-sm text-white/40">{{ error }}</p>
    </div>

    <template v-else-if="club">
      <div class="grid gap-6 lg:grid-cols-3">

        <!-- Main column -->
        <div class="space-y-6 lg:col-span-2">
          <!-- Club header -->
          <div class="glass-strong rounded-2xl p-6">
            <div class="flex items-start gap-4">
              <div
                class="flex h-16 w-16 shrink-0 items-center justify-center rounded-xl bg-[#1db954]/10 text-2xl"
                :style="club.cover_url ? { backgroundImage: `url(${club.cover_url})`, backgroundSize: 'cover' } : {}"
              >
                <span v-if="!club.cover_url">🏛</span>
              </div>
              <div class="min-w-0 flex-1">
                <h1 class="truncate text-2xl font-black text-white">{{ club.name }}</h1>
                <p v-if="club.description" class="mt-1.5 line-clamp-2 text-sm text-white/40">
                  {{ club.description }}
                </p>
                <div class="mt-3 flex items-center gap-4 text-xs text-white/30">
                  <span>{{ club.member_count }} / {{ club.max_members }} members</span>
                  <span>Created by {{ userName(club.created_by) }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Posts -->
          <div class="glass-strong rounded-2xl p-6">
            <h2 class="mb-4 text-sm font-bold uppercase tracking-wider text-white/30">Posts</h2>

            <!-- Create post -->
            <form class="mb-6 flex gap-2" @submit.prevent="createPost">
              <input
                v-model="postInput"
                type="text"
                placeholder="Write something..."
                class="min-w-0 flex-1 rounded-xl border border-white/10 bg-white/5 px-4 py-2.5 text-sm text-white placeholder-white/20 outline-none transition focus:border-white/20"
              />
              <button
                type="submit"
                class="rounded-xl bg-[#1db954]/10 px-4 py-2.5 text-sm font-semibold text-[#1db954] transition hover:bg-[#1db954]/20 disabled:opacity-40"
                :disabled="!postInput.trim()"
              >
                Post
              </button>
            </form>

            <!-- Posts list -->
            <div class="space-y-3">
              <div
                v-for="post in posts"
                :key="post.id"
                class="rounded-xl bg-white/[0.03] px-4 py-3"
              >
                <div class="flex items-center gap-2 text-xs text-white/30">
                  <span class="font-semibold text-[#1db954]">{{ userName(post.user_id) }}</span>
                  <span>{{ timeAgo(post.created_at) }}</span>
                </div>
                <p class="mt-1 text-sm text-white/70">{{ post.content }}</p>
              </div>
              <div v-if="!posts.length" class="flex flex-col items-center gap-3 py-16 text-center">
                <i class="pi pi-inbox text-4xl text-slate-500" />
                <p class="text-sm text-slate-400">No posts yet</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Sidebar -->
        <div class="space-y-6">
          <!-- Members -->
          <div class="glass-strong rounded-2xl p-6">
            <h2 class="mb-4 text-sm font-bold uppercase tracking-wider text-white/30">
              Members ({{ members.length }})
            </h2>
            <div class="space-y-2">
              <div
                v-for="m in members"
                :key="m.id"
                class="flex items-center gap-3 rounded-lg bg-white/[0.03] px-3 py-2"
              >
                <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-[#1db954]/20 text-xs font-bold text-[#1db954]">
                  {{ userName(m.user_id)?.charAt(0).toUpperCase() || '?' }}
                </div>
                <div class="min-w-0 flex-1">
                  <p class="truncate text-xs font-medium text-white">{{ userName(m.user_id) || 'Unknown' }}</p>
                </div>
                <span
                  v-if="m.role === 'admin'"
                  class="rounded bg-yellow-500/10 px-2 py-0.5 text-[10px] font-semibold text-yellow-400"
                >Admin</span>
                <span
                  v-else-if="m.role === 'moderator'"
                  class="rounded bg-blue-500/10 px-2 py-0.5 text-[10px] font-semibold text-blue-400"
                >Mod</span>
              </div>
              <div v-if="!members.length" class="flex flex-col items-center gap-3 py-16 text-center">
                <i class="pi pi-inbox text-4xl text-slate-500" />
                <p class="text-sm text-slate-400">No members yet</p>
              </div>
            </div>
          </div>

          <!-- Join/Leave button -->
          <button
            class="w-full rounded-xl py-3 text-sm font-bold transition"
            :class="isMember
              ? 'border border-red-500/20 text-red-400 hover:bg-red-500/10'
              : 'bg-[#1db954] text-black hover:bg-[#1db954]/90'"
            @click="isMember ? handleLeave() : handleJoin()"
          >
            {{ isMember ? 'Leave Club' : 'Join Club' }}
          </button>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { SkeletonLoader } from '@/components/common'
import { useSocialApi } from '@/services/api/social'
import { useUserApi } from '@/services/api/users'
import { useUserAuthStore } from '@/stores'
import type { MusicClub, MusicClubMember, MusicClubPost } from '@/services/api/social'

const router = useRouter()
const route = useRoute()
const api = useSocialApi()
const auth = useUserAuthStore()

const clubId = route.params.id as string

const loading = ref(true)
const error = ref('')
const club = ref<MusicClub | null>(null)
const members = ref<MusicClubMember[]>([])
const posts = ref<MusicClubPost[]>([])
const postInput = ref('')
const userNames = ref<Record<string, string>>({})

const isMember = ref(false)

async function fetchUserName(userId: string) {
  if (userNames.value[userId]) return
  try {
    const profile = await useUserApi().getPublicUserProfile(userId)
    const p = profile as any
    userNames.value[userId] = p.full_name || p.username || userId.slice(0, 8)
  } catch {
    userNames.value[userId] = userId.slice(0, 8)
  }
}

function userName(userId: string): string {
  return userNames.value[userId] || userId?.slice(0, 8) || ''
}

async function loadClub() {
  try {
    const [clubData, membersData, postsData] = await Promise.all([
      api.getClub(clubId),
      api.getClubMembers(clubId),
      api.getClubPosts(clubId, { limit: 50 }),
    ])
    club.value = clubData
    members.value = Array.isArray(membersData) ? membersData : []
    posts.value = Array.isArray(postsData) ? postsData : []

    const userIds = new Set<string>()
    userIds.add(clubData.created_by)
    if (Array.isArray(membersData)) {
      membersData.forEach((m: MusicClubMember) => userIds.add(m.user_id))
    }
    if (Array.isArray(postsData)) {
      postsData.forEach((p: MusicClubPost) => userIds.add(p.user_id))
    }
    await Promise.all(Array.from(userIds).map(fetchUserName))

    isMember.value = Array.isArray(membersData) && membersData.some(
      (m: MusicClubMember) => m.user_id === auth.user?.id
    )
  } catch (e: any) {
    error.value = e?.message || 'Failed to load club'
  } finally {
    loading.value = false
  }
}

async function handleJoin() {
  try {
    await api.joinClub(clubId)
    isMember.value = true
    loadClub()
  } catch { /* ignore */ }
}

async function handleLeave() {
  try {
    await api.leaveClub(clubId)
    isMember.value = false
    loadClub()
  } catch { /* ignore */ }
}

async function createPost() {
  if (!postInput.value.trim()) return
  try {
    await api.createClubPost(clubId, postInput.value.trim())
    postInput.value = ''
    loadClub()
  } catch { /* ignore */ }
}

function goBack() {
  router.push({ name: 'social' })
}

function timeAgo(dateStr: string): string {
  const now = Date.now()
  const then = new Date(dateStr).getTime()
  const diff = Math.max(0, Math.floor((now - then) / 1000))
  if (diff < 60) return 'just now'
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`
  return `${Math.floor(diff / 86400)}d ago`
}

onMounted(loadClub)
</script>
