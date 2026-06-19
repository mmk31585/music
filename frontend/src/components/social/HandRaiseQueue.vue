<template>
  <div v-if="pendingRaises.length" class="rounded-2xl bg-white/[0.04] p-6 ring-1 ring-white/[0.07]">
    <div class="mb-4 flex items-center justify-between">
      <h2 class="text-xs font-bold uppercase tracking-wider text-white/30">
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
        class="flex items-center gap-3 rounded-xl bg-white/[0.03] px-3 py-2.5"
      >
        <div class="flex h-9 w-9 shrink-0 items-center justify-center overflow-hidden rounded-full bg-white/10">
          <img
            v-if="raise.avatar_url"
            :src="raise.avatar_url"
            :alt="raise.username || 'User'"
            class="h-full w-full object-cover"
          />
          <span v-else class="text-xs font-bold text-white/60">
            {{ (raise.username || '?').charAt(0).toUpperCase() }}
          </span>
        </div>

        <div class="min-w-0 flex-1">
          <p class="truncate text-sm font-medium text-white">{{ raise.username || raise.user_id.slice(0, 8) }}</p>
        </div>

        <div class="flex shrink-0 gap-1">
          <button
            class="flex h-8 w-8 items-center justify-center rounded-lg bg-[#1db954]/10 text-[#1db954] transition hover:bg-[#1db954]/20"
            title="تأیید"
            aria-label="تأیید"
            @click="$emit('approve', raise.user_id)"
          >
            <i aria-hidden="true" class="pi pi-check text-sm" />
          </button>
          <button
            class="flex h-8 w-8 items-center justify-center rounded-lg bg-white/5 text-white/30 transition hover:bg-white/10 hover:text-white/60"
            title="رد"
            aria-label="رد"
            @click="$emit('deny', raise.user_id)"
          >
            <i aria-hidden="true" class="pi pi-times text-sm" />
          </button>
        </div>
      </div>
    </TransitionGroup>

  </div>

  <!-- Empty state (host sees this when no pending requests) -->
  <div
    v-else
    class="rounded-2xl bg-white/[0.04] p-6 text-center text-sm text-white/30 ring-1 ring-white/[0.07]"
  >
    درخواستی برای صحبت نیست
  </div>
</template>

<script setup lang="ts">
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
