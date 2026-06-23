<template>
  <slot :proceed="handleAction" />
  <GuestUpgradePrompt v-model:visible="showPrompt" />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useUserAuthStore } from '@/stores'
import { useGuestSession } from '@/composables/useGuestSession'
import GuestUpgradePrompt from './GuestUpgradePrompt.vue'

const props = defineProps<{
  action: 'play' | 'like' | 'playlist' | 'follow'
}>()

const emit = defineEmits<{
  proceed: []
}>()

const { isAuthenticated } = useUserAuthStore()
const { hasReachedLimit, incrementGuestPlay } = useGuestSession()

const showPrompt = ref(false)

function handleAction() {
  if (isAuthenticated) {
    emit('proceed')
    return
  }

  if (props.action === 'play' && hasReachedLimit.value!) {
    incrementGuestPlay()
    emit('proceed')
    return
  }

  showPrompt.value = true
}
</script>
