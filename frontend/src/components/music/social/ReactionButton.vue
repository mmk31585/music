<template>
  <div class="flex items-center gap-1">
    <button
      v-for="reactionType in ['like', 'love', 'dislike']"
      :key="reactionType"
      type="button"
      :class="[
        'inline-flex items-center gap-1 rounded-full px-3 py-1.5 text-xs font-medium transition',
        currentReaction === reactionType
          ? activeClass(reactionType)
          : 'text-slate-400 hover:bg-white/5 hover:text-white',
      ]"
      @click="$emit('react', reactionType)"
    >
      <i aria-hidden="true" :class="iconClass(reactionType)" />
      {{ reactionType === 'like' ? 'Like' : reactionType === 'love' ? 'Love' : 'Dislike' }}
    </button>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  currentReaction: string | null
}>()

defineEmits<{
  react: [type: string]
}>()

function iconClass(type: string) {
  if (type === 'like') return 'pi pi-thumbs-up text-sm'
  if (type === 'love') return 'pi pi-heart text-sm'
  return 'pi pi-thumbs-down text-sm'
}

function activeClass(type: string) {
  if (type === 'like') return 'bg-spotify/10 text-spotify'
  if (type === 'love') return 'bg-red-500/10 text-red-400'
  return 'bg-yellow-500/10 text-yellow-400'
}
</script>
