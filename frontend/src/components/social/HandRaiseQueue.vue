<template>
  <div v-if="pendingRaises.length" class="rounded-2xl bg-surface-overlay/60 p-6 ring-1 ring-border-subtle">
    <div class="mb-4 flex items-center justify-between">
      <h2 class="text-xs font-bold uppercase tracking-wider text-muted">
        درخواست‌های صحبت ({{ pendingRaises.length }})
      </h2>
    </div>

    <TransitionGroup
      name="hand-raise-item"
      tag="div"
      class="space-y-2"
    >
      <div
        v-for="raise in pendingRaises"
        :key="raise.user_id"
        class="flex items-center gap-3 rounded-xl bg-surface-overlay/50 px-3 py-2.5"
      >
        <div class="flex h-9 w-9 shrink-0 items-center justify-center overflow-hidden rounded-full bg-surface-active">
          <img
            v-if="raise.avatar_url"
            :src="raise.avatar_url"
            :alt="raise.username || 'User'"
            class="h-full w-full object-cover"
          />
          <span v-else class="text-xs font-bold text-secondary">
            {{ (raise.username || '?').charAt(0).toUpperCase() }}
          </span>
        </div>

        <div class="min-w-0 flex-1">
          <p class="truncate text-sm font-medium text-primary">{{ raise.username || raise.user_id.slice(0, 8) }}</p>
        </div>

        <div class="flex shrink-0 gap-1">
          <button
            class="flex h-8 w-8 items-center justify-center rounded-lg bg-accent-subtle text-accent transition hover:bg-accent-subtle"
            title="تأیید"
            aria-label="تأیید"
            @click="$emit('approve', raise.user_id)"
          >
            <Check aria-hidden="true" class="text-sm"  />
          </button>
          <button
            class="flex h-8 w-8 items-center justify-center rounded-lg bg-surface-overlay text-muted transition hover:bg-surface-active hover:text-secondary"
            title="رد"
            aria-label="رد"
            @click="$emit('deny', raise.user_id)"
          >
            <X aria-hidden="true" class="text-sm"  />
          </button>
        </div>
      </div>
    </TransitionGroup>

  </div>

  <!-- Empty state (host sees this when no pending requests) -->
  <div
    v-else
    class="rounded-2xl bg-surface-overlay/60 p-6 text-center text-sm text-muted ring-1 ring-border-subtle"
  >
    درخواستی برای صحبت نیست
  </div>
</template>

<script setup lang="ts">
import { Check, X } from 'lucide-vue-next'
defineProps<{
  pendingRaises: { user_id: string; username?: string; avatar_url?: string }[]
}>()

defineEmits<{
  approve: [userId: string]
  deny: [userId: string]
}>()
</script>

<style scoped>
.hand-raise-item-leave-active {
  transition: all 200ms cubic-bezier(0.4, 0, 0.2, 1);
}

.hand-raise-item-enter-active {
  transition: all 200ms cubic-bezier(0.4, 0, 0.2, 1);
}

.hand-raise-item-leave-to {
  opacity: 0;
  transform: translateX(20px);
  max-height: 0;
}

.hand-raise-item-enter-from {
  opacity: 0;
  transform: translateX(-20px);
}

@media (prefers-reduced-motion: reduce) {
  .hand-raise-item-enter-active,
  .hand-raise-item-leave-active {
    transition: none !important;
  }
  .hand-raise-item-enter-from,
  .hand-raise-item-leave-to {
    opacity: 1 !important;
    transform: none !important;
  }
}
</style>
