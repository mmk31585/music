<script setup lang="ts">
import { onBeforeUnmount, ref } from 'vue'
import random from 'lodash.random'
import { usePageLoaderStore } from '@/stores'
import { watchImmediate } from '@vueuse/core'

/* ---------- Props ---------- */
const props = defineProps({
  /** Maximum time (ms) the loading simulation may run */
  defaultDuration: {
    type: Number,
    default: 8000,
    validator: (v: number) => v >= 0,
  },
  /** Interval (ms) between progress updates */
  defaultInterval: {
    type: Number,
    default: 1000,
    validator: (v: number) => v >= 0,
  },
  /** Random variation (0‑1) applied to each increment */
  variation: {
    type: Number,
    default: 0.5,
    validator: (v: number) => v >= 0 && v <= 1,
  },
  /** Starting percentage of the bar */
  startingPoint: {
    type: Number,
    default: 5,
    validator: (v: number) => v >= 0 && v <= 100,
  },
  /** Maximum percentage before the bar is considered “done” */
  endingPoint: {
    type: Number,
    default: 90,
    validator: (v: number) => v >= 0 && v <= 100,
  },
})

/* ---------- Store & Refs ---------- */
const loadingStore = usePageLoaderStore()

const loaderRef = ref<HTMLElement | null>(null)
const isLoading = ref(true)
const isVisible = ref(false)
const progress = ref<number>(props.startingPoint)
const timeoutId = ref<ReturnType<typeof setTimeout> | null>(null)

/* ---------- Helper Functions ---------- */
function clearTimeoutId(): void {
  if (timeoutId.value) {
    clearTimeout(timeoutId.value)
    timeoutId.value = null
  }
}

function changeLoaderDelay(delay: number): void {
  if (loaderRef.value?.style) {
    loaderRef.value.style.setProperty('--delay', `${delay}ms`)
  }
}

/* ---------- Core Logic ---------- */
function start(): void {
  isLoading.value = true
  isVisible.value = true
  progress.value = props.startingPoint
  loop()
}

async function loop(): Promise<void> {
  clearTimeoutId()

  if (progress.value >= props.endingPoint) return

  changeLoaderDelay(2200)

  const stepSize =
    (props.endingPoint - props.startingPoint) / (props.defaultDuration / props.defaultInterval)
  const randomInc = random(stepSize * (1 - props.variation), stepSize * (1 + props.variation))
  const next = Math.ceil(progress.value + randomInc)

  progress.value = Math.min(next, props.endingPoint)

  timeoutId.value = setTimeout(
    loop,
    random(
      props.defaultInterval * (1 - props.variation),
      props.defaultInterval * (1 + props.variation),
    ),
  )
}

function stop(): void {
  changeLoaderDelay(200)
  progress.value = 100
  clearTimeoutId()

  setTimeout(() => (isLoading.value = false), 300)
  setTimeout(() => {
    if (!isLoading.value) isVisible.value = false
  }, 450)
}

/* ---------- Reactivity ---------- */
watchImmediate(
  () => loadingStore.loading,
  (newVal) => (newVal ? !isLoading.value && start() : stop()),
)

onBeforeUnmount(() => {
  clearTimeoutId()
})
</script>

<template>
  <div :class="{ 'loading-container': true, loading: isLoading, visible: isVisible }" dir="ltr">
    <div ref="loaderRef" :style="{ width: progress + '%' }" class="loader">
      <div class="light"></div>
    </div>
    <div class="glow"></div>
  </div>
</template>

<style scoped>
.loading-container {
  font-size: 0;
  position: fixed;
  top: 0;
  inset-inline-start: 0;
  height: 5px;
  width: 100%;
  opacity: 0;
  display: none;
  z-index: 100;
  transition: opacity 200ms;
}

.loading-container.visible {
  display: block;
}

.loading-container.loading {
  opacity: 1;
}

/* Formal gradient: deep navy → cool gray */
.loader {
  --delay: 200ms;
  background: linear-gradient(90deg, #314f8a 15%, #96a9be 100%);
  display: inline-block;
  height: 100%;
  width: 50%;
  overflow: hidden;
  border-end-end-radius: 5px;
  transition: var(--delay) width ease-out;
}

.loader > .light {
  float: inline-end;
  height: 100%;
  width: 180px;
  background-image: linear-gradient(to right, transparent, #a5c8ec, transparent);
  animation: loading-animation 3.5s ease-in infinite;
}

.glow {
  display: inline-block;
  height: 100%;
  width: 30px;
  margin-inline-start: -30px;
  border-end-end-radius: 5px;
  box-shadow: 0 0 10px #6c7a89;
}

@keyframes loading-animation {
  0%,
  50% {
    margin-inline-start: 100%;
  }
  100% {
    margin-inline-start: -10%;
  }
}
</style>
