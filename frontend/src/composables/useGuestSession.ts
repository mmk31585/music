import { computed, ref } from 'vue'
import { safeLocalStorage } from '@/services/storage'

const COUNT_KEY = 'moja_guest_play_count'
const DATE_KEY = 'moja_guest_play_date'
export const GUEST_PLAY_LIMIT = 5

function getToday(): string {
  return new Date().toISOString().slice(0, 10)
}

function initCount(): number {
  const storedDate = safeLocalStorage.getItem<string>(DATE_KEY)
  const today = getToday()
  if (storedDate !== today) {
    safeLocalStorage.setItem(DATE_KEY, today)
    safeLocalStorage.removeItem(COUNT_KEY)
    return 0
  }
  return safeLocalStorage.getItem<number>(COUNT_KEY) ?? 0
}

const guestPlayCount = ref<number>(initCount())

export function useGuestSession() {
  const remainingPlays = computed(() => Math.max(0, GUEST_PLAY_LIMIT - guestPlayCount.value))
  const hasReachedLimit = computed(() => guestPlayCount.value >= GUEST_PLAY_LIMIT)

  function incrementPlay() {
    if (hasReachedLimit.value) return
    guestPlayCount.value++
    safeLocalStorage.setItem(COUNT_KEY, guestPlayCount.value)
  }

  function canPlay(): boolean {
    return !hasReachedLimit.value
  }

  function getRemainingPlays(): number {
    return remainingPlays.value
  }

  return {
    guestPlayCount,
    remainingPlays,
    hasReachedLimit,
    incrementPlay,
    canPlay,
    getRemainingPlays,
    GUEST_PLAY_LIMIT,
  }
}
