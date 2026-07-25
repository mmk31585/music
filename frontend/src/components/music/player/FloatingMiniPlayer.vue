<template>
  <Teleport to="body">
    <Transition name="float-up">
      <div v-if="visible && currentTrack" class="fixed z-[150]">
        <FloatingMiniPlayerContent
          :visible="visible"
          @update:visible="$emit('update:visible', $event)"
          @open-fullscreen="$emit('open-fullscreen')"
          @toggle-pip="$emit('toggle-pip')"
        />
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { toRef } from 'vue'
import { usePlayerControls } from '@/composables/player'
import FloatingMiniPlayerContent from './FloatingMiniPlayerContent.vue'

defineProps<{ visible: boolean }>()

defineEmits<{
  'update:visible': [value: boolean]
  'open-fullscreen': []
  'toggle-pip': []
}>()

const pc = usePlayerControls()
const currentTrack = toRef(pc.currentTrack)
</script>

<style scoped>
.float-up-enter-active {
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.float-up-leave-active {
  transition: all 0.2s ease-in;
}
.float-up-enter-from {
  opacity: 0;
  transform: translateY(20px) scale(0.9);
}
.float-up-leave-to {
  opacity: 0;
  transform: translateY(20px) scale(0.9);
}
</style>
