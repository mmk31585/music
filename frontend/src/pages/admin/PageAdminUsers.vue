<template>
  <div class="mx-auto w-full max-w-7xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      eyebrow="Account management"
      title="Users"
      description="View and manage platform users."
    >
      <template #actions>
        <span
          class="rounded-full bg-blue-500/10 px-3 py-1 text-xs font-medium text-blue-400 tabular-nums"
        >
          {{ total }} users
        </span>
        <Button
          icon="pi pi-refresh"
          size="small"
          severity="secondary"
          class="!rounded-xl !border-white/[0.08] !bg-white/[0.04] !text-white/60 hover:!bg-white/[0.08]"
          :loading="loading"
          @click="fetchUsers()"
        />
      </template>
    </AdminSectionHeader>

    <div class="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div class="flex flex-wrap items-center gap-3">
        <div class="relative">
          <i aria-hidden="true" class="pi pi-search absolute left-3 top-1/2 -translate-y-1/2 text-xs text-slate-500" />
          <InputText
            v-model="searchQuery"
            placeholder="Search email, username, display name..."
            class="!h-9 !w-full !rounded-lg !border-white/[0.08] !bg-white/[0.03] !pl-9 !text-sm !text-white placeholder:!text-slate-600 sm:!w-72"
            @input="onSearchInput"
          />
        </div>

        <div
          class="flex items-center gap-1 rounded-lg border border-white/[0.06] bg-white/[0.02] p-0.5"
        >
          <button
            v-for="opt in roleOptions"
            :key="opt.value"
            class="rounded-md px-2.5 py-1 text-xs font-medium transition"
            :class="
              roleFilter === opt.value
                ? 'bg-white/10 text-white'
                : 'text-slate-500 hover:text-slate-300'
            "
            @click="roleFilter = opt.value; fetchUsers()"
          >
            {{ opt.label }}
          </button>
        </div>

        <div
          class="flex items-center gap-1 rounded-lg border border-white/[0.06] bg-white/[0.02] p-0.5"
        >
          <button
            v-for="opt in statusOptions"
            :key="opt.value"
            class="rounded-md px-2.5 py-1 text-xs font-medium transition"
            :class="
              statusFilter === opt.value
                ? 'bg-white/10 text-white'
                : 'text-slate-500 hover:text-slate-300'
            "
            @click="statusFilter = opt.value; fetchUsers()"
          >
            {{ opt.label }}
          </button>
        </div>
      </div>
    </div>

    <div class="overflow-hidden rounded-2xl border border-white/[0.06] bg-white/[0.02]">
      <div v-if="loading" class="divide-y divide-white/[0.04]">
        <div v-for="i in 6" :key="i" class="flex items-center gap-4 px-5 py-4">
          <div class="h-9 w-9 animate-pulse rounded-full bg-white/[0.06]" />
          <div class="flex-1 space-y-2">
            <div class="h-4 w-44 animate-pulse rounded bg-white/[0.06]" />
            <div class="h-3 w-32 animate-pulse rounded bg-white/[0.04]" />
          </div>
          <div class="h-4 w-16 animate-pulse rounded bg-white/[0.04]" />
          <div class="h-4 w-12 animate-pulse rounded bg-white/[0.04]" />
        </div>
      </div>

      <AdminEmptyState
        v-else-if="users.length === 0 && !searchQuery && !roleFilter && !statusFilter"
        icon="pi pi-users"
        title="No users found"
        description="No users registered yet."
      />

      <AdminEmptyState
        v-else-if="users.length === 0"
        icon="pi pi-search"
        title="No results found"
        :description="filterDescription"
      />

      <template v-else>
        <div class="overflow-x-auto">
          <table class="admin-table min-w-[800px]">
            <thead>
              <tr>
                <th
                  class="cursor-pointer select-none hover:text-white"
                  @click="toggleSort('display_name')"
                >
                  <span class="inline-flex items-center gap-1">
                    User
                    <i aria-hidden="true" v-if="sortBy === 'display_name'" :class="sortIcon" class="text-[10px]" />
                  </span>
                </th>
                <th
                  class="cursor-pointer select-none hover:text-white"
                  @click="toggleSort('email')"
                >
                  <span class="inline-flex items-center gap-1">
                    Email
                    <i aria-hidden="true" v-if="sortBy === 'email'" :class="sortIcon" class="text-[10px]" />
                  </span>
                </th>
                <th
                  class="hidden cursor-pointer select-none hover:text-white md:table-cell"
                  @click="toggleSort('username')"
                >
                  <span class="inline-flex items-center gap-1">
                    Username
                    <i aria-hidden="true" v-if="sortBy === 'username'" :class="sortIcon" class="text-[10px]" />
                  </span>
                </th>
                <th
                  class="cursor-pointer select-none hover:text-white"
                  @click="toggleSort('role')"
                >
                  <span class="inline-flex items-center gap-1">
                    Role
                    <i aria-hidden="true" v-if="sortBy === 'role'" :class="sortIcon" class="text-[10px]" />
                  </span>
                </th>
                <th
                  class="cursor-pointer select-none hover:text-white"
                  @click="toggleSort('is_active')"
                >
                  <span class="inline-flex items-center gap-1">
                    Status
                    <i aria-hidden="true" v-if="sortBy === 'is_active'" :class="sortIcon" class="text-[10px]" />
                  </span>
                </th>
                <th
                  class="hidden cursor-pointer select-none hover:text-white lg:table-cell"
                  @click="toggleSort('created_at')"
                >
                  <span class="inline-flex items-center gap-1">
                    Created
                    <i aria-hidden="true" v-if="sortBy === 'created_at'" :class="sortIcon" class="text-[10px]" />
                  </span>
                </th>
                <th class="text-right">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-white/[0.04]">
              <tr
                v-for="u in users"
                :key="u.id"
                class="group transition-colors"
              >
                <td>
                  <div class="flex items-center gap-3">
                    <div
                      class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-white/10 text-sm font-bold text-white"
                    >
                      {{ initials(u) }}
                    </div>
                    <div class="min-w-0">
                      <p class="truncate text-sm font-medium text-white">
                        {{ u.display_name || u.displayName || u.username || u.email }}
                      </p>
                      <div class="mt-0.5 flex items-center gap-1.5">
                        <span
                          v-if="u.email_verified"
                          class="flex h-4 w-4 items-center justify-center rounded-full bg-blue-500/20 text-[8px] text-blue-400"
                          title="Verified"
                        >
                          <i aria-hidden="true" class="pi pi-check" />
                        </span>
                        <span v-else class="text-[10px] text-slate-600" title="Not verified"
                          >unverified</span
                        >
                      </div>
                    </div>
                  </div>
                </td>
                <td>
                  {{ u.email }}
                </td>
                <td class="hidden md:table-cell">
                  {{ u.username }}
                </td>
                <td>
                  <span
                    class="rounded-full px-2 py-0.5 text-xs font-medium"
                    :class="roleClass(u.role)"
                  >
                    {{ u.role || 'user' }}
                  </span>
                </td>
                <td>
                  <span
                    class="rounded-full px-2 py-0.5 text-[10px] font-medium"
                    :class="
                      u.is_active
                        ? 'bg-emerald-500/10 text-emerald-400'
                        : 'bg-red-500/10 text-red-400'
                    "
                  >
                    {{ u.is_active ? 'Active' : 'Inactive' }}
                  </span>
                </td>
                <td class="hidden tabular-nums lg:table-cell">
                  {{ formatDate(u.created_at) }}
                </td>
                <td class="text-right">
                  <div class="flex items-center justify-end gap-1">
                    <Button
                      icon="pi pi-eye"
                      size="small"
                      severity="secondary"
                      class="!h-8 !w-8 !rounded-lg !border-white/[0.06] !bg-transparent !text-slate-500 hover:!bg-white/[0.06] hover:!text-white"
                      @click="viewUser(u)"
                    />
                    <Button
                      icon="pi pi-pencil"
                      size="small"
                      severity="secondary"
                      class="!h-8 !w-8 !rounded-lg !border-white/[0.06] !bg-transparent !text-slate-500 hover:!bg-white/[0.06] hover:!text-white"
                      @click="openEdit(u)"
                    />
                    <Button
                      icon="pi pi-trash"
                      size="small"
                      severity="danger"
                      class="!h-8 !w-8 !rounded-lg !border-white/[0.06] !bg-transparent !text-red-400/60 hover:!bg-red-500/10 hover:!text-red-400"
                      @click="confirmDelete(u)"
                    />
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="flex items-center justify-between border-t border-white/[0.06] px-5 py-3">
          <span class="text-xs text-slate-500 tabular-nums">
            Page {{ page }} of {{ totalPages }}
          </span>
          <div class="flex items-center gap-1">
            <Button
              icon="pi pi-chevron-left"
              size="small"
              severity="secondary"
              class="!h-8 !w-8 !rounded-lg !border-white/[0.06] !bg-transparent !text-slate-500 hover:!bg-white/[0.06] hover:!text-white"
              :disabled="page <= 1"
              @click="page > 1 && goToPage(page - 1)"
            />
            <button
              v-for="p in visiblePages"
              :key="p"
              class="flex h-8 w-8 items-center justify-center rounded-lg text-xs font-medium transition"
              :class="
                p === page
                  ? 'bg-white/10 text-white'
                  : 'text-slate-500 hover:bg-white/[0.04] hover:text-slate-300'
              "
              @click="goToPage(p)"
            >
              {{ p }}
            </button>
            <Button
              icon="pi pi-chevron-right"
              size="small"
              severity="secondary"
              class="!h-8 !w-8 !rounded-lg !border-white/[0.06] !bg-transparent !text-slate-500 hover:!bg-white/[0.06] hover:!text-white"
              :disabled="page >= totalPages"
              @click="page < totalPages && goToPage(page + 1)"
            />
          </div>
        </div>
      </template>
    </div>

    <Dialog
      v-model:visible="showEditDialog"
      :modal="true"
      :closable="true"
      :draggable="false"
      :pt="{
        root: { class: '!border-white/[0.06] !bg-[#141414] !rounded-2xl !shadow-2xl' },
        header: { class: '!bg-transparent !border-0 !pb-0' },
        content: { class: '!bg-transparent !p-0' },
        footer: { class: '!bg-transparent !border-0' },
        mask: { class: '!backdrop-blur-sm' },
      }"
    >
      <template #header>
        <div class="flex items-center gap-3">
          <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-blue-500/10">
            <i aria-hidden="true" class="pi pi-user text-blue-400" />
          </div>
          <div>
            <h3 class="text-base font-semibold text-white">
              Edit User — {{ editTarget?.display_name || editTarget?.username || '' }}
            </h3>
          </div>
        </div>
      </template>
      <div class="space-y-4 p-6">
        <div>
          <label class="mb-1.5 block text-xs font-medium text-slate-400">Role</label>
          <div class="flex items-center gap-2">
            <button
              v-for="r in ['user', 'creator', 'admin']"
              :key="r"
              class="rounded-lg px-3 py-1.5 text-xs font-medium capitalize transition"
              :class="
                editRole === r
                  ? 'bg-white/15 text-white'
                  : 'bg-white/[0.04] text-slate-400 hover:text-white'
              "
              @click="editRole = r"
            >
              {{ r }}
            </button>
          </div>
        </div>

        <div class="flex items-center justify-between rounded-xl bg-white/[0.03] px-4 py-3">
          <div>
            <p class="text-sm font-medium text-white">Active</p>
            <p class="text-xs text-slate-500">Allow this user to log in</p>
          </div>
          <label class="relative inline-flex h-5 w-9 cursor-pointer items-center">
            <input type="checkbox" v-model="editActive" class="peer sr-only" />
            <span
              class="absolute inset-0 rounded-full bg-white/[0.08] transition peer-checked:bg-emerald-500"
            />
            <span
              class="absolute left-0.5 h-4 w-4 rounded-full bg-white transition peer-checked:translate-x-4"
            />
          </label>
        </div>

        <div class="flex items-center justify-between rounded-xl bg-white/[0.03] px-4 py-3">
          <div>
            <p class="text-sm font-medium text-white">Email verified</p>
            <p class="text-xs text-slate-500">Mark email as confirmed</p>
          </div>
          <label class="relative inline-flex h-5 w-9 cursor-pointer items-center">
            <input type="checkbox" v-model="editEmailVerified" class="peer sr-only" />
            <span
              class="absolute inset-0 rounded-full bg-white/[0.08] transition peer-checked:bg-blue-500"
            />
            <span
              class="absolute left-0.5 h-4 w-4 rounded-full bg-white transition peer-checked:translate-x-4"
            />
          </label>
        </div>
      </div>

      <div class="flex justify-end gap-2 border-t border-white/[0.06] px-6 py-4">
        <Button
          label="Cancel"
          severity="secondary"
          size="small"
          class="!rounded-xl !border-white/[0.08] !bg-white/[0.04] !text-slate-300 hover:!bg-white/[0.08]"
          @click="showEditDialog = false"
        />
        <Button
          label="Save"
          size="small"
          :loading="saving"
          class="!rounded-xl !bg-emerald-500 !px-5 !text-black hover:!bg-emerald-400"
          @click="saveEdit"
        />
      </div>
    </Dialog>

    <Dialog
      v-model:visible="showDetailDialog"
      :modal="true"
      :closable="true"
      :draggable="false"
      :pt="{
        root: { class: '!border-white/[0.06] !bg-[#141414] !rounded-2xl !shadow-2xl' },
        header: { class: '!bg-transparent !border-0 !pb-0' },
        content: { class: '!bg-transparent !p-0' },
        footer: { class: '!bg-transparent !border-0' },
        mask: { class: '!backdrop-blur-sm' },
      }"
    >
      <template #header>
        <div class="flex items-center gap-3">
          <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-sky-500/10">
            <i aria-hidden="true" class="pi pi-eye text-sky-400" />
          </div>
          <div>
            <h3 class="text-base font-semibold text-white">{{ userDetail?.display_name || userDetail?.username || 'User Detail' }}</h3>
          </div>
        </div>
      </template>
      <div v-if="detailLoading" class="flex items-center justify-center p-12">
        <i aria-hidden="true" class="pi pi-spin pi-spinner text-2xl text-slate-500" />
      </div>
      <div v-else-if="userDetail" class="divide-y divide-white/[0.06]">
        <div class="flex items-center gap-4 p-6">
          <div
            class="flex h-14 w-14 shrink-0 items-center justify-center rounded-full bg-white/10 text-xl font-bold text-white"
          >
            {{ initials(userDetail) }}
          </div>
          <div>
            <p class="text-lg font-bold text-white">{{ userDetail.display_name || userDetail.displayName || userDetail.username || 'Unknown' }}</p>
            <p class="text-sm text-slate-400">{{ userDetail.email }}</p>
          </div>
        </div>
        <div class="space-y-4 p-6">
          <div class="flex items-center justify-between">
            <span class="text-sm text-slate-400">Username</span>
            <span class="text-sm text-white/80">{{ userDetail.username || '-' }}</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-sm text-slate-400">Role</span>
            <span
              class="rounded-full px-2.5 py-0.5 text-xs font-medium"
              :class="roleClass(userDetail.role)"
            >
              {{ userDetail.role || 'user' }}
            </span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-sm text-slate-400">Status</span>
            <span
              class="rounded-full px-2.5 py-0.5 text-xs font-medium"
              :class="
                userDetail.is_active
                  ? 'bg-emerald-500/10 text-emerald-400'
                  : 'bg-red-500/10 text-red-400'
              "
            >
              {{ userDetail.is_active ? 'Active' : 'Inactive' }}
            </span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-sm text-slate-400">Email verified</span>
            <span
              v-if="userDetail.email_verified"
              class="flex items-center gap-1.5 rounded-full bg-blue-500/10 px-2.5 py-0.5 text-xs font-medium text-blue-400"
            >
              <i aria-hidden="true" class="pi pi-check-circle" /> Verified
            </span>
            <span v-else class="text-xs text-slate-600">Not verified</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-sm text-slate-400">Created</span>
            <span class="text-sm text-white/60 tabular-nums">{{ formatDate(userDetail.created_at) }}</span>
          </div>
          <div v-if="userDetail.updated_at" class="flex items-center justify-between">
            <span class="text-sm text-slate-400">Updated</span>
            <span class="text-sm text-white/60 tabular-nums">{{ formatDate(userDetail.updated_at) }}</span>
          </div>
        </div>
      </div>
      <div v-else class="p-6 text-center text-sm text-slate-500">Could not load user details.</div>
      <div class="flex justify-end border-t border-white/[0.06] px-6 py-4">
        <Button
          label="Close"
          severity="secondary"
          size="small"
          class="!rounded-xl !border-white/[0.08] !bg-white/[0.04] !text-slate-300 hover:!bg-white/[0.08]"
          @click="showDetailDialog = false"
        />
      </div>
    </Dialog>

    <AdminDeleteConfirm
      v-model="showDeleteDialog"
      title="Delete User"
      :item-name="deleteTarget?.display_name || deleteTarget?.username || deleteTarget?.email"
      :deleting="deleting"
      @confirm="deleteUser"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Dialog from 'primevue/dialog'
import { useToast } from 'primevue/usetoast'
import { useAuthApi } from '@/services/api'
import { AdminSectionHeader, AdminEmptyState } from '@/components/admin'
import AdminDeleteConfirm from '@/components/admin/AdminDeleteConfirm.vue'

interface AdminUser {
  id: string
  email: string
  username?: string
  display_name?: string
  displayName?: string
  role: string
  is_active: boolean
  email_verified: boolean
  avatar_url?: string
  created_at: string
  updated_at?: string
}

const { adminListUsers, adminGetUser, adminUpdateUser, adminDeleteUser } = useAuthApi()
const toast = useToast()

const users = ref<AdminUser[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const searchQuery = ref('')
const roleFilter = ref('')
const statusFilter = ref('')
const sortBy = ref('created_at')
const sortOrder = ref('desc')
let searchTimer: ReturnType<typeof setTimeout> | null = null

const roleOptions = [
  { label: 'All', value: '' },
  { label: 'User', value: 'user' },
  { label: 'Creator', value: 'creator' },
  { label: 'Admin', value: 'admin' },
]

const statusOptions = [
  { label: 'All', value: '' },
  { label: 'Active', value: 'active' },
  { label: 'Inactive', value: 'inactive' },
]

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize.value)))

const visiblePages = computed(() => {
  const pages: number[] = []
  const tp = totalPages.value
  const cp = page.value
  let start = Math.max(1, cp - 2)
  let end = Math.min(tp, cp + 2)
  if (end - start < 4) {
    if (start === 1) end = Math.min(tp, start + 4)
    else start = Math.max(1, end - 4)
  }
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

const sortIcon = computed(() =>
  sortOrder.value === 'desc' ? 'pi pi-sort-amount-down' : 'pi pi-sort-amount-up-alt',
)

const filterDescription = computed(() => {
  const parts: string[] = []
  if (searchQuery.value) parts.push(`"${searchQuery.value}"`)
  if (roleFilter.value) parts.push(`role: ${roleFilter.value}`)
  if (statusFilter.value) parts.push(`status: ${statusFilter.value}`)
  return `No users matching ${parts.join(', ') || 'criteria'}.`
})

function initials(u: AdminUser) {
  const name = u.display_name || u.displayName || u.username || u.email || '?'
  return name.slice(0, 2).toUpperCase()
}

function formatDate(dateStr: string) {
  try {
    return new Date(dateStr).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    })
  } catch {
    return dateStr
  }
}

function roleClass(role: string): string {
  if (role === 'admin') return 'bg-yellow-500/10 text-yellow-400'
  if (role === 'creator') return 'bg-purple-500/10 text-purple-400'
  return 'bg-white/5 text-slate-400'
}

function toggleSort(col: string) {
  if (sortBy.value === col) {
    sortOrder.value = sortOrder.value === 'desc' ? 'asc' : 'desc'
  } else {
    sortBy.value = col
    sortOrder.value = 'desc'
  }
  fetchUsers()
}

function goToPage(p: number) {
  if (p < 1 || p > totalPages.value) return
  page.value = p
  fetchUsers()
}

function onSearchInput() {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    page.value = 1
    fetchUsers()
  }, 300)
}

async function fetchUsers() {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: page.value,
      page_size: pageSize.value,
      sort_by: sortBy.value,
      sort_order: sortOrder.value,
    }
    if (searchQuery.value) params.search = searchQuery.value
    if (roleFilter.value) params.role = roleFilter.value
    if (statusFilter.value === 'active') params.is_active = true
    else if (statusFilter.value === 'inactive') params.is_active = false
    const res = await adminListUsers(params)
    users.value = res?.items ?? []
    total.value = res?.total ?? 0
    page.value = res?.page ?? 1
  } catch {
    toast.add({ severity: 'error', summary: 'Failed to load users', life: 3000 })
    users.value = []
  } finally {
    loading.value = false
  }
}

const showEditDialog = ref(false)
const editTarget = ref<AdminUser | null>(null)
const editRole = ref('user')
const editActive = ref(true)
const editEmailVerified = ref(false)

const showDetailDialog = ref(false)
const userDetail = ref<AdminUser | null>(null)
const detailLoading = ref(false)

function viewUser(u: AdminUser) {
  showDetailDialog.value = true
  detailLoading.value = true
  adminGetUser(u.id)
    .then((res: Record<string, unknown>) => {
      if (res) userDetail.value = res
    })
    .catch(() => {
      toast.add({ severity: 'error', summary: 'Failed to load user details', life: 3000 })
    })
    .finally(() => {
      detailLoading.value = false
    })
}

function openEdit(u: AdminUser) {
  editTarget.value = u
  editRole.value = u.role || 'user'
  editActive.value = u.is_active !== false
  editEmailVerified.value = u.email_verified === true
  showEditDialog.value = true
}

async function saveEdit() {
  if (!editTarget.value) return
  saving.value = true
  try {
    await adminUpdateUser(editTarget.value.id, {
      role: editRole.value,
      is_active: editActive.value,
      email_verified: editEmailVerified.value,
    })
    toast.add({ severity: 'success', summary: 'User updated', life: 3000 })
    showEditDialog.value = false
    fetchUsers()
  } catch {
    toast.add({ severity: 'error', summary: 'Failed to update user', life: 3000 })
  } finally {
    saving.value = false
  }
}

const showDeleteDialog = ref(false)
const deleteTarget = ref<AdminUser | null>(null)

function confirmDelete(u: AdminUser) {
  deleteTarget.value = u
  showDeleteDialog.value = true
}

async function deleteUser() {
  if (!deleteTarget.value) return
  deleting.value = true
  try {
    await adminDeleteUser(deleteTarget.value.id)
    toast.add({ severity: 'success', summary: 'User deleted', life: 3000 })
    showDeleteDialog.value = false
    fetchUsers()
  } catch {
    toast.add({ severity: 'error', summary: 'Failed to delete user', life: 3000 })
  } finally {
    deleting.value = false
  }
}

onMounted(fetchUsers)
</script>
