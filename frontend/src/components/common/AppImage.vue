<template>
  <div
    v-if="hasError || !src"
    role="img"
    :aria-label="alt"
    class="flex items-center justify-center bg-white/5"
    :class="[containerClass, containerStyle]"
  >
    <i aria-hidden="true" :class="fallbackIcon || 'pi pi-music'" class="text-slate-500" :style="{ fontSize: iconSize }" />
  </div>
  <img
    v-else
    ref="imgRef"
    :src="src"
    :alt="alt"
    :loading="lazy ? 'lazy' : undefined"
    class="h-full w-full object-cover"
    :class="imgClass"
    @error="onError"
    @load="onLoad"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface Props {
  src?: string | null
  alt?: string
  lazy?: boolean
  fallbackIcon?: string
  containerClass?: string
  containerStyle?: string
  imgClass?: string
  iconSize?: string
}

withDefaults(defineProps<Props>(), {
  src: '',
  alt: '',
  lazy: true,
  fallbackIcon: 'pi pi-music',
  containerClass: '',
  containerStyle: '',
  imgClass: '',
  iconSize: '1.5rem',
})

const emit = defineEmits<{
  load: []
  error: []
}>()

const hasError = ref(false)

function onError() {
  hasError.value = true
  emit('error')
}

function onLoad() {
  emit('load')
}
</script>
