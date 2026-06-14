import { computed, ref } from 'vue'
import { safeLocalStorage } from '@/services/storage'

const STORAGE_KEY = 'moja_guest_plays'
export const GUEST_PLAY_LIMIT = 3

const guestPlayCount = ref<number>((safeLocalStorage.getItem<number>(STORAGE_KEY) as number) || 0)

export function useGuestSession() {
  function incrementGuestPlay() {
    guestPlayCount.value++
    safeLocalStorage.setItem(STORAGE_KEY, guestPlayCount.value)
  }

  const hasReachedLimit = computed(() => guestPlayCount.value >= GUEST_PLAY_LIMIT)

  return {
    guestPlayCount,
    incrementGuestPlay,
    hasReachedLimit,
    GUEST_PLAY_LIMIT,
  }
}
