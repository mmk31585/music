<template>
  <Teleport to="body">
    <Transition name="context-fade">
      <div
        v-if="visible"
        ref="backdrop"
        class="fixed inset-0 z-[300]"
        @click="close"
        @contextmenu.prevent="close"
      >
        <div
          ref="menuRef"
          class="absolute min-w-45 overflow-hidden rounded-xl border border-white/8 bg-[#1a1a2e] py-1 shadow-2xl backdrop-blur-2xl"
          :style="{ left: `${pos.x}px`, top: `${pos.y}px` }"
          role="menu"
          :aria-label="label"
          @click.stop
          @keydown="onKeydown"
        >
          <button
            v-for="(item, i) in items"
            :key="i"
            ref="itemRefs"
            role="menuitem"
            class="flex w-full items-center gap-3 px-4 py-2.5 text-left text-sm text-slate-200 transition hover:bg-white/8 hover:text-white focus-visible:bg-white/8 focus-visible:outline-hidden"
            :class="{ 'border-t border-white/6': item.separator }"
            :tabindex="focusedIndex === i ? 0 : -1"
            @click="handleAction(item)"
            @mouseenter="focusedIndex = i"
          >
            <i aria-hidden="true" v-if="item.icon" :class="item.icon" class="w-4 text-xs text-slate-500" />
            <span>{{ item.label }}</span>
          </button>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, onUnmounted } from 'vue'

interface ContextMenuItem {
  label: string
  icon?: string
  action: () => void
  separator?: boolean
}

const props = withDefaults(defineProps<{
  visible: boolean
  items: ContextMenuItem[]
  label?: string
  position?: { x: number; y: number }
}>(), {
  label: 'Context menu',
  position: () => ({ x: 0, y: 0 }),
})

const emit = defineEmits<{
  close: []
  action: [item: ContextMenuItem]
}>()

const backdrop = ref<HTMLElement | null>(null)
const menuRef = ref<HTMLElement | null>(null)
const itemRefs = ref<(HTMLElement | null)[]>([])
const focusedIndex = ref(-1)
const pos = ref({ x: 0, y: 0 })

watch(() => props.visible, (v) => {
  if (v) {
    pos.value = { ...props.position }
    focusedIndex.value = -1
    nextTick(() => {
      constrainPosition()
      menuRef.value?.focus()
    })
  }
})

function constrainPosition() {
  const el = menuRef.value
  if (el!) return
  const rect = el.getBoundingClientRect()
  const maxX = window.innerWidth - rect.width - 8
  const maxY = window.innerHeight - rect.height - 8
  if (pos.value.x > maxX) pos.value.x = maxX
  if (pos.value.y > maxY) pos.value.y = maxY
  if (pos.value.x < 8) pos.value.x = 8
  if (pos.value.y < 8) pos.value.y = 8
}

function handleAction(item: ContextMenuItem) {
  item.action()
  emit('action', item)
  close()
}

function close() {
  emit('close')
}

function onKeydown(e: KeyboardEvent) {
  const len = props.items.length
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    focusedIndex.value = (focusedIndex.value + 1) % len
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    focusedIndex.value = (focusedIndex.value - 1 + len) % len
  } else if (e.key === 'Enter' || e.key === ' ') {
    e.preventDefault()
    if (focusedIndex.value >= 0 && focusedIndex.value < len) {
      const item = props.items[focusedIndex.value]
      if (item) handleAction(item)
    }
  } else if (e.key === 'Escape') {
    close()
  }
}

onUnmounted(() => {
  close()
})
</script>

<style scoped>
.context-fade-enter-active,
.context-fade-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}
.context-fade-enter-from,
.context-fade-leave-to {
  opacity: 0;
  transform: scale(0.96);
}
</style>
