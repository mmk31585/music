<template>
  <div class="relative">
    <div
      ref="carouselRef"
      class="flex gap-4 overflow-x-auto scroll-smooth pb-2"
      style="scrollbar-width: none; -ms-overflow-style: none"
    >
      <slot />
    </div>
    <button
      v-if="showScroll && !isAtStart"
      type="button"
      aria-label="Scroll left"
      class="absolute top-1/2 -left-3 z-10 flex h-10 w-10 -translate-y-1/2 items-center justify-center rounded-full bg-black/80 text-white shadow-lg backdrop-blur-xs transition hover:bg-black/90"
      @click="scroll(-300)"
    >
      <i aria-hidden="true" class="pi pi-chevron-left text-sm" />
    </button>
    <button
      v-if="showScroll && !isAtEnd"
      type="button"
      aria-label="Scroll right"
      class="absolute top-1/2 -right-3 z-10 flex h-10 w-10 -translate-y-1/2 items-center justify-center rounded-full bg-black/80 text-white shadow-lg backdrop-blur-xs transition hover:bg-black/90"
      @click="scroll(300)"
    >
      <i aria-hidden="true" class="pi pi-chevron-right text-sm" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useSwipe } from '@/composables'

defineProps<{ showScroll?: boolean }>()

const carouselRef = ref<HTMLElement | null>(null)
const isAtStart = ref(true)
const isAtEnd = ref(false)

function checkScroll() {
  if (carouselRef.value!) return
  isAtStart.value = carouselRef.value.scrollLeft <= 10
  isAtEnd.value =
    carouselRef.value.scrollLeft >= carouselRef.value.scrollWidth - carouselRef.value.clientWidth - 10
}

function scroll(amount: number) {
  carouselRef.value?.scrollBy({ left: amount, behavior: 'smooth' })
}

function scrollTo(dir: 'next' | 'prev') {
  if (carouselRef.value!) return
  const amount = dir === 'next' ? carouselRef.value.clientWidth : -carouselRef.value.clientWidth
  carouselRef.value.scrollBy({ left: amount, behavior: 'smooth' })
}

useSwipe({ element: carouselRef, onSwipeLeft: () => scrollTo('next'), onSwipeRight: () => scrollTo('prev') })

let observer: ResizeObserver | null = null

onMounted(() => {
  carouselRef.value?.addEventListener('scroll', checkScroll)
  observer = new ResizeObserver(checkScroll)
  if (carouselRef.value) observer.observe(carouselRef.value)
  checkScroll()
})

onUnmounted(() => {
  carouselRef.value?.removeEventListener('scroll', checkScroll)
  observer?.disconnect()
})
</script>
