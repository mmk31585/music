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
          : 'text-secondary hover:bg-surface-overlay hover:text-primary',
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
  if (type === 'like') return 'bg-accent-subtle text-accent'
  if (type === 'love') return 'bg-danger-subtle text-danger'
  return 'bg-warning-subtle text-warning'
}
</script>
