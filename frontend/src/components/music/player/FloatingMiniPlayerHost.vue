<template>
  <FloatingMiniPlayer
    v-if="!isPiPOpen"
    v-model:visible="visibleProxy"
    @open-fullscreen="$emit('open-fullscreen')"
    @toggle-pip="togglePiP"
  />

  <Teleport v-else-if="pipMountEl" :to="pipMountEl">
    <PiPPlayerContent @close="close" />
  </Teleport>

  <div v-if="showUnsupportedMessage" class="hidden">Document Picture-in-Picture not supported.</div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { usePlayerPiPController } from '@/composables/usePlayerPiPController'
import FloatingMiniPlayer from './FloatingMiniPlayer.vue'
import PiPPlayerContent from './PiPPlayerContent.vue'

const props = defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  'open-fullscreen': []
  'pip-unsupported': []
}>()

const showUnsupportedMessage = ref(false)

const visibleProxy = computed({
  get: () => props.visible,
  set: (value: boolean) => emit('update:visible', value),
})

const {
  isSupported,
  isOpen: isPiPOpen,
  mountEl: pipMountEl,
  open,
  close,
} = usePlayerPiPController()

async function togglePiP() {
  if (!isSupported.value) {
    showUnsupportedMessage.value = true
    emit('pip-unsupported')
    window.setTimeout(() => {
      showUnsupportedMessage.value = false
    }, 2500)
    return
  }

  if (isPiPOpen.value) {
    close()
    return
  }

  await open()
}

defineExpose({ togglePiP })

watch(
  () => props.visible,
  (visible) => {
    if (!visible && isPiPOpen.value) {
      close()
    }
  },
)
</script>
