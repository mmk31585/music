<template>
  <section
    v-if="items.length"
    ref="heroRef"
    class="relative mb-12 overflow-hidden rounded-2xl"
    :style="{ height: `${heroHeight}px` }"
    role="region"
    aria-roledescription="carousel"
    :aria-label="'اسلایدر ویژه'"
    @mouseenter="pauseAutoRotate"
    @mouseleave="resumeAutoRotate"
  >
    <!-- Video background (loop) -->
    <video
      v-if="videoSrc"
      :key="activeIndex"
      :src="videoSrc"
      class="absolute inset-0 h-full w-full object-cover"
      autoplay
      muted
      loop
      playsinline
    />
    <!-- Fallback blurred image background when no video -->
    <div v-else class="absolute inset-0 overflow-hidden">
      <img
        v-for="(item, i) in items"
        :key="i"
        :src="item.image"
        class="absolute inset-0 h-full w-full object-cover transition-opacity duration-700"
        :class="i === activeIndex ? 'opacity-100' : 'opacity-0'"
        style="filter: blur(40px) saturate(1.3); transform: scale(1.1)"
        aria-hidden="true"
      />
    </div>
    <div class="absolute inset-0 bg-linear-to-l from-transparent via-black/40 to-black/80" />

    <!-- Frosted glass card overlay -->
    <div class="relative z-10 flex h-full items-center px-6 md:px-10">
      <div class="flex w-full items-center justify-between">
        <div
          class="max-w-lg rounded-2xl border border-white/10 bg-black/30 px-8 py-8 backdrop-blur-xl shadow-2xl"
          style="backdrop-filter: blur(24px); -webkit-backdrop-filter: blur(24px);"
          :style="{
            opacity: slideVisible ? 1 : 0,
            transform: slideVisible ? 'translateY(0)' : 'translateY(12px)',
            transition: 'opacity 400ms ease, transform 400ms ease',
          }"
        >
          <span
            class="mb-3 inline-block rounded-full px-3 py-1 text-[11px] font-bold"
            :class="badgeClass"
          >
            {{ currentItem.badge }}
          </span>
          <h2 class="font-hero text-4xl font-black leading-tight text-white md:text-5xl lg:text-6xl" style="letter-spacing: -0.03em">
            {{ currentItem.title }}
          </h2>
          <p class="mt-3 text-lg text-white/60 md:text-xl">
            {{ currentItem.subtitle }}
          </p>
          <div class="mt-5 flex flex-wrap gap-3">
            <button
              type="button"
              class="glow-green inline-flex h-12 cursor-pointer items-center gap-2 rounded-full bg-spotify px-8 text-base font-bold text-black transition-all hover:bg-spotify-hover hover:scale-105 active:scale-95 focus-visible:ring-2 focus-visible:ring-spotify focus-visible:ring-offset-2 focus-visible:outline-hidden"
              :aria-label="'پخش ' + currentItem.title"
              @click="$emit('play', currentItem)"
            >
              <Play aria-hidden="true" class="text-sm"  />
              پخش
            </button>
            <button
              type="button"
              class="inline-flex h-12 cursor-pointer items-center gap-2 rounded-full border border-white/15 bg-white/4 px-6 text-base font-medium text-white backdrop-blur-xs transition-all hover:bg-white/10 active:scale-95 focus-visible:ring-2 focus-visible:ring-spotify focus-visible:ring-offset-2 focus-visible:outline-hidden"
              @click="$emit('add-to-library', currentItem)"
            >
              <Plus aria-hidden="true" class="text-sm"  />
              افزودن به کتابخانه
            </button>
          </div>
        </div>

        <!-- Floating art (desktop only) -->
        <div class="hidden shrink-0 md:block">
          <img
            v-if="currentItem.image"
            :src="currentItem.image"
            :alt="currentItem.title"
            fetchpriority="high"
            class="h-44 w-44 rounded-2xl object-cover shadow-2xl md:h-52 md:w-52"
            :style="{
              transform: 'rotate(-3deg)',
              boxShadow: '0 24px 64px rgba(0,0,0,0.6)',
            }"
          />
        </div>
      </div>
    </div>

    <!-- Arrow buttons -->
    <button
      type="button"
      class="absolute top-1/2 right-4 z-20 hidden h-10 w-10 -translate-y-1/2 items-center justify-center rounded-full border border-white/10 bg-black/60 text-white backdrop-blur-xs transition hover:bg-white/10 md:flex"
      :class="showArrows ? 'opacity-100' : 'opacity-0'"
      aria-label="اسلاید قبلی"
      @click="prev"
    >
      <ChevronRight aria-hidden="true" class="text-sm"  />
    </button>
    <button
      type="button"
      class="absolute top-1/2 left-4 z-20 hidden h-10 w-10 -translate-y-1/2 items-center justify-center rounded-full border border-white/10 bg-black/60 text-white backdrop-blur-xs transition hover:bg-white/10 md:flex"
      :class="showArrows ? 'opacity-100' : 'opacity-0'"
      aria-label="اسلاید بعدی"
      @click="next"
    >
      <ChevronLeft aria-hidden="true" class="text-sm"  />
    </button>

    <!-- Indicators -->
    <div class="absolute bottom-4 left-1/2 z-20 flex -translate-x-1/2 items-center gap-2">
      <button
        v-for="(item, i) in items"
        :key="i"
        type="button"
        class="cursor-pointer rounded-full transition-all"
        :class="i === activeIndex
          ? 'w-6 bg-white'
          : 'w-2 bg-white/30'"
        :style="{ height: '8px', transition: 'width 300ms cubic-bezier(0.19, 1, 0.22, 1)' }"
        :aria-label="'برو به اسلاید ' + (i + 1)"
        @click="goTo(i)"
      />
    </div>

    <!-- Pause auto-rotate button -->
    <button
      type="button"
      class="absolute top-4 left-4 z-20 flex h-8 w-8 cursor-pointer items-center justify-center rounded-full border border-white/10 bg-black/40 text-white/60 backdrop-blur-xs transition hover:bg-white/10"
      :aria-label="autoRotating ? 'توقف چرخش خودکار' : 'شروع چرخش خودکار'"
      @click="toggleAutoRotate"
    >
      <component :is="autoRotating ? Pause : Play"<i aria-hidden="true"  class="text-xs" /> />
    </button>
  </section>
</template>

<script setup lang="ts">
import { ChevronLeft, ChevronRight, Pause, Play, Plus } from 'lucide-vue-next'
import { ref, computed, onMounted, onUnmounted } from 'vue'

export interface HeroItem {
  id: string
  title: string
  subtitle: string
  image: string
  badge: string
  badgeVariant?: 'green' | 'purple'
  type: 'album' | 'artist' | 'playlist'
  videoSrc?: string
}

const props = withDefaults(defineProps<{
  items: HeroItem[]
  interval?: number
}>(), {
  interval: 8000,
})

defineEmits<{
  play: [item: HeroItem]
  'add-to-library': [item: HeroItem]
}>()

const activeIndex = ref(0)
const autoRotating = ref(true)
const showArrows = ref(false)
const slideVisible = ref(true)
const heroRef = ref<HTMLElement | null>(null)
let timer: ReturnType<typeof setInterval> | null = null
let reducedMotion = false

const currentItem = computed((): HeroItem => props.items[activeIndex.value] ?? props.items[0]!)

const videoSrc = computed(() => currentItem.value.videoSrc || '')

const heroHeight = computed(() => {
  if (typeof window !== 'undefined' && window.innerWidth < 768) return 360
  return 520
})

const badgeClass = computed(() => {
  const v = currentItem.value.badgeVariant || 'green'
  return v === 'purple'
    ? 'bg-aurora-purple/20 text-aurora-purple'
    : 'bg-spotify/20 text-spotify'
})

function startTimer() {
  stopTimer()
  if (reducedMotion) return
  timer = setInterval(() => {
    next()
  }, props.interval)
}

function stopTimer() {
  if (timer) {
    clearInterval(timer)
    timer = null
  }
}

function pauseAutoRotate() {
  if (!reducedMotion) {
    showArrows.value = true
    stopTimer()
  }
}

function resumeAutoRotate() {
  if (autoRotating.value && !reducedMotion) {
    showArrows.value = false
    startTimer()
  }
}

function toggleAutoRotate() {
  autoRotating.value = autoRotating.value!
  if (autoRotating.value) {
    startTimer()
  } else {
    stopTimer()
  }
}

function crossfade(toIndex: number) {
  if (toIndex === activeIndex.value) return
  slideVisible.value = false
  setTimeout(() => {
    activeIndex.value = toIndex
    slideVisible.value = true
  }, 50)
}

function next() {
  const next = (activeIndex.value + 1) % props.items.length
  crossfade(next)
}

function prev() {
  const prev = (activeIndex.value - 1 + props.items.length) % props.items.length
  crossfade(prev)
}

function goTo(index: number) {
  crossfade(index)
  if (autoRotating.value && !reducedMotion) {
    startTimer()
  }
}

onMounted(() => {
  reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  if (autoRotating.value && !reducedMotion) {
    startTimer()
  }
})

onUnmounted(() => {
  stopTimer()
})
</script>
