<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { usePermissionsApi } from '@/services/api/permissions'
import { useAuthApi } from '@/services/api/auth'
import type { RoleWithPermissions, Permission, UserPermissionOverride, UserAccess } from '@/services/api/permissions/types'
import type { AdminUser } from '@/services/api/auth/types'

const api = usePermissionsApi()
const authApi = useAuthApi()
const roles = ref<RoleWithPermissions[]>([])
const allPermissions = ref<Permission[]>([])
const loading = ref(true)
const activeSection = ref<'roles' | 'users'>('roles')

// Role editing
const editingRole = ref<RoleWithPermissions | null>(null)
const editingPermissions = ref<string[]>([])
const saving = ref(false)

// User override
const searchQuery = ref('')
const userSuggestions = ref<AdminUser[]>([])
const showSuggestions = ref(false)
const selectedUser = ref<AdminUser | null>(null)
const userAccess = ref<UserAccess | null>(null)
const userOverrides = ref<UserPermissionOverride[]>([])
const searchLoading = ref(false)

async function loadData() {
  loading.value = true
  try {
    const [perms, roleList] = await Promise.all([
      api.listPermissions(),
      api.listRoles(),
    ])
    allPermissions.value = Array.isArray(perms) ? perms : [perms]
    roles.value = Array.isArray(roleList) ? roleList : [roleList]
  } finally {
    loading.value = false
  }
}

function startEditRole(role: RoleWithPermissions) {
  editingRole.value = role
  editingPermissions.value = [...role.permissions]
}

function cancelEdit() {
  editingRole.value = null
  editingPermissions.value = []
}

async function saveRolePermissions() {
  if (!editingRole.value) return
  saving.value = true
  try {
    await api.setRolePermissions(editingRole.value.slug, editingPermissions.value)
    editingRole.value.permissions = [...editingPermissions.value]
    cancelEdit()
  } finally {
    saving.value = false
  }
}

function togglePermission(slug: string) {
  const idx = editingPermissions.value.indexOf(slug)
  if (idx >= 0) editingPermissions.value.splice(idx, 1)
  else editingPermissions.value.push(slug)
}

watch(searchQuery, async (val) => {
  if (val.trim().length < 2) {
    userSuggestions.value = []
    showSuggestions.value = false
    return
  }
  try {
    const res = await authApi.adminListUsers({ search: val.trim() })
    const data = res as any
    userSuggestions.value = data?.items ?? []
    showSuggestions.value = userSuggestions.value.length > 0
  } catch {
    userSuggestions.value = []
    showSuggestions.value = false
  }
})

function selectUser(user: AdminUser) {
  selectedUser.value = user
  searchQuery.value = `${user.display_name || user.username} (${user.email})`
  showSuggestions.value = false
  loadUserPerms(user.id)
}

async function loadUserPerms(userId: string) {
  searchLoading.value = true
  try {
    const [access, overrides] = await Promise.all([
      api.getUserAccess(userId),
      api.getUserOverrides(userId),
    ])
    userAccess.value = access
    userOverrides.value = Array.isArray(overrides) ? overrides : [overrides]
  } finally {
    searchLoading.value = false
  }
}

async function grantOverride(permSlug: string) {
  if (!userAccess.value) return
  await api.grantPermission(userAccess.value.userId, permSlug, 'Manual override by admin')
  await loadUserPerms(userAccess.value.userId)
}

function closeSuggestions() {
  setTimeout(() => { showSuggestions.value = false }, 200)
}

async function revokeOverride(permSlug: string) {
  if (!userAccess.value) return
  await api.revokePermission(userAccess.value.userId, permSlug)
  await loadUserPerms(userAccess.value.userId)
}

onMounted(loadData)
</script>

<template>
  <div class="mx-auto w-full max-w-6xl px-4 pb-32 pt-8 md:px-6">
    <h1 class="mb-8 text-3xl font-black text-white">Permissions Management</h1>

    <!-- Section tabs -->
    <div class="mb-6 flex gap-1 rounded-xl bg-white/4 p-1">
      <button
        class="flex-1 rounded-lg py-2.5 text-sm font-medium transition"
        :class="activeSection === 'roles' ? 'bg-white/10 text-white' : 'text-white/30 hover:text-white/50'"
        @click="activeSection = 'roles'"
      >Roles & Permissions</button>
      <button
        class="flex-1 rounded-lg py-2.5 text-sm font-medium transition"
        :class="activeSection === 'users' ? 'bg-white/10 text-white' : 'text-white/30 hover:text-white/50'"
        @click="activeSection = 'users'"
      >User Overrides</button>
    </div>

    <div v-if="loading" class="space-y-4">
      <div v-for="i in 3" :key="i" class="h-16 animate-pulse rounded-xl bg-white/5" />
    </div>

    <!-- Roles Section -->
    <div v-show="activeSection === 'roles'" class="space-y-6">
      <div v-for="role in roles" :key="role.id" class="rounded-2xl bg-white/5 p-6">
        <div class="mb-4 flex items-center justify-between">
          <div>
            <h2 class="text-lg font-bold text-white">{{ role.label }}</h2>
            <p class="text-xs text-white/40">
              {{ role.slug }} &middot; Level {{ role.hierarchyLevel }}
              <span v-if="role.isSystemRole" class="ml-2 rounded bg-amber-500/20 px-2 py-0.5 text-amber-400">System</span>
            </p>
          </div>
          <button
            class="rounded-lg bg-white/10 px-4 py-2 text-sm text-white transition hover:bg-white/20"
            @click="startEditRole(role)"
          >Edit</button>
        </div>

        <div class="flex flex-wrap gap-2">
          <span
            v-for="perm in role.permissions"
            :key="perm"
            class="rounded-full bg-white/8 px-3 py-1 text-xs text-white/60"
          >{{ perm }}</span>
        </div>
      </div>
    </div>

    <!-- Edit role modal -->
    <Teleport to="body">
      <div
        v-if="editingRole"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
        @click.self="cancelEdit"
      >
        <div class="mx-4 w-full max-w-2xl rounded-2xl bg-slate-900 p-6 shadow-2xl">
          <h2 class="mb-1 text-lg font-bold text-white">{{ editingRole.label }}</h2>
          <p class="mb-6 text-sm text-white/40">Select permissions for this role</p>

          <div class="mb-6 max-h-80 space-y-1 overflow-y-auto">
            <div v-for="cat in [...new Set(allPermissions.map(p => p.category))]" :key="cat">
              <p class="mb-1 mt-3 text-xs font-bold uppercase tracking-wider text-white/30">{{ cat }}</p>
              <button
                v-for="perm in allPermissions.filter(p => p.category === cat)"
                :key="perm.slug"
                class="mr-1 mb-1 inline-flex items-center gap-1.5 rounded-full px-3 py-1.5 text-xs transition"
                :class="editingPermissions.includes(perm.slug) ? 'bg-spotify/30 text-white ring-1 ring-spotify/50' : 'bg-white/8 text-white/50 hover:bg-white/12'"
                @click="togglePermission(perm.slug)"
              >
                {{ perm.slug }}
              </button>
            </div>
          </div>

          <div class="flex justify-end gap-3">
            <button
              class="rounded-full bg-white/10 px-6 py-2 text-sm text-white transition hover:bg-white/20"
              @click="cancelEdit"
            >Cancel</button>
            <button
              class="rounded-full bg-spotify px-6 py-2 text-sm font-bold text-black transition hover:bg-spotify-hover disabled:opacity-50"
              :disabled="saving"
              @click="saveRolePermissions"
            >{{ saving ? 'Saving...' : 'Save' }}</button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- User Overrides Section -->
    <div v-show="activeSection === 'users'" class="space-y-6">
      <div class="relative">
        <input
          v-model="searchQuery"
          placeholder="Search users by name, username, or email..."
          class="w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2.5 text-sm text-white placeholder-white/30 outline-none transition focus:border-spotify/50"
          @focus="showSuggestions = userSuggestions.length > 0"
          @blur="closeSuggestions"
        />
        <div
          v-if="showSuggestions && userSuggestions.length > 0"
          class="absolute z-10 mt-1 w-full rounded-xl border border-white/10 bg-slate-900 py-1 shadow-2xl"
        >
          <button
            v-for="u in userSuggestions"
            :key="u.id"
            class="flex w-full items-center gap-3 px-4 py-2.5 text-left text-sm text-white/80 transition hover:bg-white/5"
            @click="selectUser(u)"
          >
            <span class="flex h-8 w-8 items-center justify-center rounded-full bg-white/10 text-xs font-bold text-white/50">
              {{ (u.display_name || u.username || '?')[0] }}
            </span>
            <div class="flex-1">
              <p class="font-medium text-white">{{ u.display_name || u.username }}</p>
              <p class="text-xs text-white/40">{{ u.email }} &middot; {{ u.role }}</p>
            </div>
          </button>
        </div>
      </div>

      <div v-if="searchLoading" class="flex items-center justify-center py-12">
        <span class="text-sm text-white/40">Loading user permissions...</span>
      </div>

      <div v-else-if="selectedUser && userAccess" class="rounded-2xl bg-white/5 p-6">
        <div class="mb-4 flex items-center gap-3">
          <span class="flex h-10 w-10 items-center justify-center rounded-full bg-white/10 text-sm font-bold text-white">
            {{ (selectedUser.display_name || selectedUser.username || '?')[0] }}
          </span>
          <div>
            <h3 class="text-lg font-bold text-white">{{ selectedUser.display_name || selectedUser.username }}</h3>
            <p class="text-xs text-white/40">{{ selectedUser.email }}</p>
          </div>
        </div>
        <div class="mb-4 grid grid-cols-3 gap-4 text-sm">
          <div><span class="text-white/40">Role:</span> <span class="text-white">{{ userAccess.roleSlug }}</span></div>
          <div><span class="text-white/40">Level:</span> <span class="text-white">{{ userAccess.roleLevel }}</span></div>
          <div><span class="text-white/40">Status:</span> <span class="text-white" :class="selectedUser.is_active ? 'text-emerald-400' : 'text-red-400'">{{ selectedUser.is_active ? 'Active' : 'Inactive' }}</span></div>
        </div>

        <h4 class="mb-2 text-sm font-bold text-white/60 uppercase tracking-wider">Permissions ({{ userAccess.permissions.length }})</h4>
        <div class="mb-6 flex flex-wrap gap-2">
          <span
            v-for="p in userAccess.permissions"
            :key="p"
            class="rounded-full bg-white/8 px-3 py-1 text-xs text-white/60"
          >{{ p }}</span>
          <span v-if="userAccess.permissions.length === 0" class="text-xs text-white/30 italic">No permissions assigned</span>
        </div>

        <h4 class="mb-2 text-sm font-bold text-white/60 uppercase tracking-wider">Overrides</h4>
        <div v-for="o in userOverrides" :key="o.id" class="mb-2 flex items-center gap-3 rounded-lg bg-white/5 p-3">
          <span class="text-sm text-white">{{ o.permSlug }}</span>
          <span
            class="rounded-full px-2 py-0.5 text-xs font-bold"
            :class="o.granted ? 'bg-emerald-500/20 text-emerald-400' : 'bg-red-500/20 text-red-400'"
          >{{ o.granted ? 'Granted' : 'Revoked' }}</span>
          <span class="text-xs text-white/30">{{ o.reason }}</span>
          <button
            class="ml-auto rounded px-3 py-1 text-xs text-red-400 transition hover:bg-red-500/20"
            v-if="o.granted"
            @click="revokeOverride(o.permSlug)"
          >Revoke</button>
        </div>
        <div v-if="userOverrides.length === 0" class="rounded-lg bg-white/5 p-4 text-center text-xs text-white/30 italic">
          No overrides for this user
        </div>
      </div>
    </div>
  </div>
</template>
