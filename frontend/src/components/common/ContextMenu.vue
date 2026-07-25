<template>
  <Teleport to="body">
    <!-- ── Desktop: floating context menu ────────────────────── -->
    <Transition name="ctx-scale">
      <div
        v-if="visible && !isMobile"
        ref="backdropRef"
        class="ctx-backdrop"
        @click="close"
        @contextmenu.prevent="close"
      >
        <div
          ref="menuRef"
          class="ctx-menu"
          :style="positionStyle"
          role="menu"
          :aria-label="ariaLabel || 'Context menu'"
          tabindex="-1"
          @keydown="onKeydown"
        >
          <!-- Header -->
          <div v-if="header" class="ctx-header">
            <img
              v-if="header.coverUrl"
              :src="header.coverUrl"
              :alt="header.title"
              class="ctx-header-cover"
              loading="eager"
            />
            <div class="ctx-header-info">
              <div class="ctx-header-title">
                <span class="ctx-truncate">{{ header.title }}</span>
                <span v-if="header.explicit" class="ctx-badge">E</span>
                <span v-if="header.isHighQuality" class="ctx-badge ctx-badge-accent">HQ</span>
              </div>
              <div class="ctx-header-subtitle">
                <span class="ctx-truncate">{{ header.artistName }}</span>
                <template v-if="header.albumName">
                  <span class="ctx-dot" />
                  <span class="ctx-truncate">{{ header.albumName }}</span>
                </template>
                <template v-if="header.duration">
                  <span class="ctx-dot" />
                  <span class="ctx-tabular">{{ formatDuration(header.duration) }}</span>
                </template>
              </div>
              <div class="ctx-header-indicators">
                <span v-if="header.isPlaying" class="ctx-playing-indicator" aria-label="Playing">
                  <span class="ctx-eq-bar" /><span class="ctx-eq-bar ctx-eq-delay-1" /><span class="ctx-eq-bar ctx-eq-delay-2" />
                </span>
                <Heart v-if="header.isLiked" class="ctx-indicator-icon" :size="14" fill="currentColor" aria-label="Liked" />
                <Download v-if="header.isDownloaded" class="ctx-indicator-icon" :size="14" aria-label="Downloaded" />
              </div>
            </div>
          </div>

          <!-- Sections -->
          <template v-for="(section, si) in visibleSections" :key="section.id">
            <div v-if="section.label" class="ctx-section-label" role="presentation">
              {{ section.label }}
            </div>
            <div class="ctx-section" role="group" :aria-label="section.label">
              <template v-for="(action, ai) in section.items" :key="action.id">
                <button
                  v-if="!action.hidden"
                  :data-menu-item="flatIndex(si, ai)"
                  role="menuitem"
                  class="ctx-item"
                  :class="{
                    'ctx-item-danger': action.danger,
                    'ctx-item-checked': action.checked,
                    'ctx-item-disabled': action.disabled,
                  }"
                  :style="itemStyle"
                  :tabindex="focusedIndex === flatIndex(si, ai) ? 0 : -1"
                  :disabled="action.disabled"
                  :aria-checked="action.checked"
                  @mouseenter="focusedIndex = flatIndex(si, ai)"
                  @focus="focusedIndex = flatIndex(si, ai)"
                  @click="activateItem(action)"
                >
                  <!-- Left accent bar on hover -->
                  <span class="ctx-item-accent" :style="accentBorderStyle" />
                  <!-- Icon -->
                  <component
                    :is="resolveIcon(action.icon)"
                    class="ctx-item-icon"
                    :class="{ 'ctx-item-icon-active': action.checked }"
                    :size="18"
                    :fill="action.checked ? (accentColor || 'currentColor') : 'none'"
                    :stroke-width="action.checked ? 0 : 1.5"
                    aria-hidden="true"
                  />
                  <!-- Label -->
                  <span class="ctx-item-label">
                    {{ action.label }}
                    <span v-if="action.subtitle" class="ctx-item-subtitle">{{ action.subtitle }}</span>
                  </span>
                  <!-- Right side: badge / shortcut -->
                  <span v-if="action.badge" class="ctx-item-badge">{{ action.badge }}</span>
                  <kbd v-if="action.shortcut" class="ctx-kbd">{{ action.shortcut }}</kbd>
                  <Check v-if="action.checked" class="ctx-check-icon" :size="16" :stroke-width="2.5" aria-hidden="true" />
                </button>
              </template>
            </div>
            <!-- Section separator -->
            <div v-if="si < visibleSections.length - 1" class="ctx-separator" role="separator" />
          </template>
        </div>
      </div>
    </Transition>

    <!-- ── Mobile: bottom action sheet ───────────────────────── -->
    <Transition name="ctx-sheet">
      <div
        v-if="visible && isMobile"
        ref="sheetBackdropRef"
        class="ctx-sheet-backdrop"
        @click="close"
      >
        <div
          ref="sheetRef"
          class="ctx-sheet"
          role="dialog"
          :aria-label="ariaLabel || 'Actions'"
          @click.stop
          @touchstart.passive="onSheetTouchStart"
          @touchmove.passive="onSheetTouchMove"
          @touchend.passive="onSheetTouchEnd"
        >
          <!-- Drag handle -->
          <div class="ctx-sheet-handle" aria-hidden="true">
            <div class="ctx-sheet-handle-bar" />
          </div>

          <!-- Header -->
          <div v-if="header" class="ctx-header">
            <img
              v-if="header.coverUrl"
              :src="header.coverUrl"
              :alt="header.title"
              class="ctx-header-cover"
              loading="eager"
            />
            <div class="ctx-header-info">
              <div class="ctx-header-title">
                <span class="ctx-truncate">{{ header.title }}</span>
                <span v-if="header.explicit" class="ctx-badge">E</span>
                <span v-if="header.isHighQuality" class="ctx-badge ctx-badge-accent">HQ</span>
              </div>
              <div class="ctx-header-subtitle">
                <span class="ctx-truncate">{{ header.artistName }}</span>
                <template v-if="header.albumName">
                  <span class="ctx-dot" />
                  <span class="ctx-truncate">{{ header.albumName }}</span>
                </template>
                <template v-if="header.duration">
                  <span class="ctx-dot" />
                  <span class="ctx-tabular">{{ formatDuration(header.duration) }}</span>
                </template>
              </div>
              <div class="ctx-header-indicators">
                <span v-if="header.isPlaying" class="ctx-playing-indicator" aria-label="Playing">
                  <span class="ctx-eq-bar" /><span class="ctx-eq-bar ctx-eq-delay-1" /><span class="ctx-eq-bar ctx-eq-delay-2" />
                </span>
                <Heart v-if="header.isLiked" class="ctx-indicator-icon" :size="14" fill="currentColor" aria-label="Liked" />
                <Download v-if="header.isDownloaded" class="ctx-indicator-icon" :size="14" aria-label="Downloaded" />
              </div>
            </div>
          </div>

          <!-- Scrollable sections -->
          <div class="ctx-sheet-scroll">
            <template v-for="(section, si) in visibleSections" :key="section.id">
              <div v-if="section.label" class="ctx-section-label" role="presentation">
                {{ section.label }}
              </div>
              <div class="ctx-section" role="group" :aria-label="section.label">
                <template v-for="action in section.items" :key="action.id">
                  <button
                    v-if="!action.hidden"
                    role="menuitem"
                    class="ctx-item ctx-item-sheet"
                    :class="{
                      'ctx-item-danger': action.danger,
                      'ctx-item-checked': action.checked,
                      'ctx-item-disabled': action.disabled,
                    }"
                    :disabled="action.disabled"
                    :aria-checked="action.checked"
                    @click="activateItem(action)"
                  >
                    <component
                      :is="resolveIcon(action.icon)"
                      class="ctx-item-icon"
                      :class="{ 'ctx-item-icon-active': action.checked }"
                      :size="20"
                      :fill="action.checked ? (accentColor || 'currentColor') : 'none'"
                      :stroke-width="action.checked ? 0 : 1.5"
                      aria-hidden="true"
                    />
                    <span class="ctx-item-label">
                      {{ action.label }}
                      <span v-if="action.subtitle" class="ctx-item-subtitle">{{ action.subtitle }}</span>
                    </span>
                    <span v-if="action.badge" class="ctx-item-badge">{{ action.badge }}</span>
                    <Check v-if="action.checked" class="ctx-check-icon" :size="18" :stroke-width="2.5" aria-hidden="true" />
                  </button>
                </template>
              </div>
              <div v-if="si < visibleSections.length - 1" class="ctx-separator" role="separator" />
            </template>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, ref, watch, onMounted, onUnmounted, nextTick } from 'vue'
import {
  Heart,
  Download,
  Play,
  SkipForward,
  ListMusic,
  Radio,
  PlusCircle,
  UserRound,
  Disc3,
  Music2,
  Share2,
  Link2,
  Info,
  Flag,
  Monitor,
  Clock,
  Settings2,
  MoreHorizontal,
  Check,
  type Component,
} from 'lucide-vue-next'
import type {
  ContextMenuSection,
  ContextMenuHeader,
  ContextMenuAction,
} from '@/types/context-menu'
import { useContextMenuState } from '@/composables/useContextMenuState'
import { useContextMenuPosition } from '@/composables/useContextMenuPosition'
import { useContextMenuKeyboard } from '@/composables/useContextMenuKeyboard'

// ── Props ───────────────────────────────────────────────────────

const props = withDefaults(defineProps<{
  visible: boolean
  sections: ContextMenuSection[]
  header?: ContextMenuHeader
  position?: { x: number; y: number }
  accentColor?: string
  anchorEl?: HTMLElement | null
  ariaLabel?: string
}>(), {
  ariaLabel: 'Context menu',
})

const emit = defineEmits<{
  'update:visible': [value: boolean]
  close: []
}>()

// ── State ───────────────────────────────────────────────────────

const { open, close: releaseState, forceClose } = useContextMenuState()

const menuRef = ref<HTMLElement | null>(null)
const backdropRef = ref<HTMLElement | null>(null)
const sheetRef = ref<HTMLElement | null>(null)
const sheetBackdropRef = ref<HTMLElement | null>(null)

const isMobile = ref(false)
const mobileQuery = typeof window !== 'undefined'
  ? window.matchMedia('(max-width: 767px)')
  : null

function updateMobile() {
  isMobile.value = mobileQuery?.matches ?? false
}

// ── Positioning ─────────────────────────────────────────────────

const positionRef = computed(() => props.position)
const anchorRef = computed(() => props.anchorEl ?? null)

const { style: positionStyle, recompute: recomputePosition } = useContextMenuPosition({
  position: positionRef,
  anchorEl: anchorRef,
  menuEl: menuRef,
})

// ── Flattened items for keyboard nav ────────────────────────────

const flatItems = computed<ContextMenuAction[]>(() => {
  const items: ContextMenuAction[] = []
  for (const section of visibleSections.value) {
    for (const action of section.items) {
      if (!action.hidden) items.push(action)
    }
  }
  return items
})

function flatIndex(sectionIdx: number, itemIdx: number): number {
  let count = 0
  for (let si = 0; si < sectionIdx; si++) {
    count += visibleSections.value[si]!.items.filter((a) => !a.hidden).length
  }
  return count + itemIdx
}

// ── Keyboard navigation ─────────────────────────────────────────

const activeMenuEl = computed(() => isMobile.value ? sheetRef.value : menuRef.value)

const { focusedIndex, onKeydown } = useContextMenuKeyboard({
  items: flatItems,
  menuEl: activeMenuEl,
  onClose: close,
  onActivate: activateItem,
})

// ── Computed ────────────────────────────────────────────────────

const visibleSections = computed(() =>
  props.sections
    .map((section) => ({
      ...section,
      items: section.items.filter((a) => !a.hidden),
    }))
    .filter((section) => section.items.length > 0),
)

const itemStyle = computed(() => {
  if (!props.accentColor) return {}
  return {
    '--ctx-accent': props.accentColor,
  } as Record<string, string>
})

const accentBorderStyle = computed(() => {
  if (!props.accentColor) return {}
  return {
    background: props.accentColor,
  } as Record<string, string>
})

// ── Icon registry ───────────────────────────────────────────────

const iconMap: Record<string, typeof Component> = {
  Play,
  SkipForward,
  ListMusic,
  Radio,
  Heart,
  PlusCircle,
  UserRound,
  Disc3,
  Music2,
  Share2,
  Link2,
  Download,
  Info,
  Flag,
  Monitor,
  Clock,
  Settings2,
  MoreHorizontal,
}

function resolveIcon(name: string): typeof Component {
  return iconMap[name] ?? MoreHorizontal
}

// ── Actions ─────────────────────────────────────────────────────

function activateItem(action: ContextMenuAction) {
  if (action.disabled) return
  action.action()
  close()
}

function close() {
  emit('update:visible', false)
  emit('close')
  releaseState()
}

// ── Mobile sheet touch handling (swipe down to dismiss) ──────────

let sheetStartY = 0
let sheetDelta = 0
const sheetTransform = ref(0)

function onSheetTouchStart(e: TouchEvent) {
  sheetStartY = e.touches[0]!.clientY
  sheetDelta = 0
}

function onSheetTouchMove(e: TouchEvent) {
  const dy = e.touches[0]!.clientY - sheetStartY
  if (dy > 0) {
    sheetDelta = dy
    sheetTransform.value = dy
  }
}

function onSheetTouchEnd() {
  if (sheetDelta > 100) {
    close()
  }
  sheetTransform.value = 0
  sheetDelta = 0
}

// ── Format duration ─────────────────────────────────────────────

function formatDuration(seconds: number): string {
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}

// ── Watch visibility ────────────────────────────────────────────

watch(() => props.visible, (v) => {
  if (v) {
    open()
    if (isMobile.value) {
      // Mobile sheet doesn't need position calculation
      return
    }
    nextTick(() => {
      recomputePosition()
      // Focus first item
      const first = menuRef.value?.querySelector('[role="menuitem"]:not(:disabled)') as HTMLElement
      first?.focus()
    })
  }
})

// ── Global close listeners ──────────────────────────────────────

function onGlobalClose() {
  if (props.visible) close()
}

function onKeyHandler(e: KeyboardEvent) {
  if (e.key === 'Escape' && props.visible) {
    e.preventDefault()
    close()
  }
}

onMounted(() => {
  updateMobile()
  mobileQuery?.addEventListener('change', updateMobile)
  window.addEventListener('keydown', onKeyHandler)
  window.addEventListener('ctx-menu:close-all', onGlobalClose as EventListener)
})

onUnmounted(() => {
  mobileQuery?.removeEventListener('change', updateMobile)
  window.removeEventListener('keydown', onKeyHandler)
  window.removeEventListener('ctx-menu:close-all', onGlobalClose as EventListener)
})
</script>

<style scoped>
/* ── Backdrop ────────────────────────────────────────────── */

.ctx-backdrop {
  position: fixed;
  inset: 0;
  z-index: 300;
}

/* ── Desktop Menu ────────────────────────────────────────── */

.ctx-menu {
  position: fixed;
  min-width: 300px;
  max-width: 340px;
  max-height: calc(100vh - 24px);
  overflow-y: auto;
  overflow-x: hidden;
  overscroll-behavior: contain;
  border-radius: 20px;
  background: var(--surface-popover);
  border: 1px solid var(--surface-border);
  box-shadow: var(--shadow-overlay);
  backdrop-filter: blur(32px);
  -webkit-backdrop-filter: blur(32px);
  padding: 8px;
  z-index: 301;
  outline: none;
}

/* ── Header ──────────────────────────────────────────────── */

.ctx-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px 12px;
  min-height: 56px;
}

.ctx-header-cover {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  object-fit: cover;
  flex-shrink: 0;
}

.ctx-header-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.ctx-header-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  line-height: 1.3;
}

.ctx-header-subtitle {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--text-secondary);
  line-height: 1.3;
}

.ctx-header-indicators {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 2px;
}

.ctx-indicator-icon {
  color: var(--accent);
  width: 14px;
  height: 14px;
}

/* ── Truncation ──────────────────────────────────────────── */

.ctx-truncate {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* ── Badges ──────────────────────────────────────────────── */

.ctx-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 9px;
  font-weight: 700;
  line-height: 1;
  padding: 2px 5px;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.12);
  color: var(--text-secondary);
  flex-shrink: 0;
  letter-spacing: 0.5px;
}

.ctx-badge-accent {
  background: var(--accent-subtle);
  color: var(--accent);
}

.ctx-dot {
  width: 3px;
  height: 3px;
  border-radius: 50%;
  background: var(--text-muted);
  flex-shrink: 0;
}

.ctx-tabular {
  font-variant-numeric: tabular-nums;
}

/* ── Section labels ──────────────────────────────────────── */

.ctx-section-label {
  padding: 6px 12px 4px;
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--text-muted);
  user-select: none;
}

/* ── Separator ───────────────────────────────────────────── */

.ctx-separator {
  height: 1px;
  margin: 6px 8px;
  background: var(--border-subtle);
}

/* ── Section ─────────────────────────────────────────────── */

.ctx-section {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

/* ── Items (desktop) ─────────────────────────────────────── */

.ctx-item {
  position: relative;
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 8px 12px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
  background: transparent;
  border: none;
  cursor: pointer;
  text-align: start;
  transition: color 120ms ease, background 120ms ease;
  outline: none;
  -webkit-tap-highlight-color: transparent;
}

.ctx-item:hover,
.ctx-item:focus-visible {
  background: var(--surface-hover);
  color: var(--text-primary);
}

.ctx-item:active {
  transform: scale(0.98);
}

.ctx-item:focus-visible {
  box-shadow: 0 0 0 2px var(--ring-focus);
}

.ctx-item-danger:hover,
.ctx-item-danger:focus-visible {
  color: var(--danger);
  background: var(--danger-subtle);
}

.ctx-item-disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.ctx-item-disabled:hover {
  background: transparent;
  color: var(--text-secondary);
}

/* ── Left accent bar ─────────────────────────────────────── */

.ctx-item-accent {
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 0;
  border-radius: 0 2px 2px 0;
  background: var(--ctx-accent, var(--accent));
  opacity: 0;
  transition: height 120ms ease, opacity 120ms ease;
}

[dir="rtl"] .ctx-item-accent {
  left: auto;
  right: 0;
  border-radius: 2px 0 0 2px;
}

.ctx-item:hover .ctx-item-accent,
.ctx-item:focus-visible .ctx-item-accent {
  height: 20px;
  opacity: 1;
}

/* ── Item sub-elements ───────────────────────────────────── */

.ctx-item-icon {
  flex-shrink: 0;
  color: var(--text-tertiary);
  transition: color 120ms ease;
}

.ctx-item:hover .ctx-item-icon,
.ctx-item:focus-visible .ctx-item-icon {
  color: var(--text-secondary);
}

.ctx-item-icon-active {
  color: var(--accent) !important;
}

.ctx-item-label {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 1px;
  text-align: start;
}

.ctx-item-subtitle {
  font-size: 12px;
  font-weight: 400;
  color: var(--text-tertiary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ctx-item-badge {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-tertiary);
  font-variant-numeric: tabular-nums;
  flex-shrink: 0;
}

.ctx-kbd {
  font-size: 10px;
  font-family: ui-monospace, 'SF Mono', 'Menlo', monospace;
  color: var(--text-muted);
  padding: 1px 5px;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border-subtle);
  line-height: 1.4;
  flex-shrink: 0;
}

.ctx-check-icon {
  flex-shrink: 0;
  color: var(--accent);
  margin-inline-start: auto;
}

/* ── Playing indicator (eq bars) ─────────────────────────── */

.ctx-playing-indicator {
  display: inline-flex;
  align-items: flex-end;
  gap: 1.5px;
  height: 12px;
}

.ctx-eq-bar {
  width: 2.5px;
  border-radius: 999px;
  background: var(--accent);
  animation: ctx-eq 850ms ease-in-out infinite alternate;
  will-change: transform, opacity;
}

.ctx-eq-bar:nth-child(1) { height: 60%; animation-delay: 0ms; }
.ctx-eq-bar:nth-child(2) { height: 100%; animation-delay: 150ms; }
.ctx-eq-bar:nth-child(3) { height: 40%; animation-delay: 300ms; }

@keyframes ctx-eq {
  from { transform: scaleY(0.45); opacity: 0.6; }
  to { transform: scaleY(1); opacity: 1; }
}

@media (prefers-reduced-motion: reduce) {
  .ctx-eq-bar { animation: none; }
}

/* ── Desktop scale transition ────────────────────────────── */

.ctx-scale-enter-active {
  transition: opacity 150ms cubic-bezier(0.16, 1, 0.3, 1),
              transform 150ms cubic-bezier(0.16, 1, 0.3, 1);
}

.ctx-scale-leave-active {
  transition: opacity 100ms ease,
              transform 100ms ease;
}

.ctx-scale-enter-from {
  opacity: 0;
  transform: scale(0.95) translateY(4px);
}

.ctx-scale-leave-to {
  opacity: 0;
  transform: scale(0.96);
}

@media (prefers-reduced-motion: reduce) {
  .ctx-scale-enter-active,
  .ctx-scale-leave-active {
    transition: opacity 100ms ease;
  }
  .ctx-scale-enter-from,
  .ctx-scale-leave-to {
    transform: none;
  }
}

/* ── Mobile sheet ────────────────────────────────────────── */

.ctx-sheet-backdrop {
  position: fixed;
  inset: 0;
  z-index: 300;
  background: var(--bg-overlay);
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.ctx-sheet {
  width: 100%;
  max-width: 480px;
  max-height: 85vh;
  background: var(--surface-popover);
  border-radius: 20px 20px 0 0;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  box-shadow: 0 -4px 32px rgba(0, 0, 0, 0.4);
  padding-bottom: env(safe-area-inset-bottom, 0px);
}

.ctx-sheet-handle {
  display: flex;
  justify-content: center;
  padding: 10px 0 4px;
}

.ctx-sheet-handle-bar {
  width: 36px;
  height: 4px;
  border-radius: 2px;
  background: rgba(255, 255, 255, 0.15);
}

.ctx-sheet-scroll {
  overflow-y: auto;
  overscroll-behavior: contain;
  padding: 4px 8px 8px;
  -webkit-overflow-scrolling: touch;
}

.ctx-item-sheet {
  padding: 14px 16px;
  font-size: 16px;
  gap: 14px;
  border-radius: 14px;
}

/* ── Sheet slide transition ──────────────────────────────── */

.ctx-sheet-enter-active {
  transition: opacity 200ms ease, transform 300ms cubic-bezier(0.16, 1, 0.3, 1);
}

.ctx-sheet-leave-active {
  transition: opacity 150ms ease, transform 200ms ease;
}

.ctx-sheet-enter-from {
  opacity: 0;
  transform: translateY(100%);
}

.ctx-sheet-leave-to {
  opacity: 0;
  transform: translateY(60%);
}

@media (prefers-reduced-motion: reduce) {
  .ctx-sheet-enter-active,
  .ctx-sheet-leave-active {
    transition: opacity 150ms ease;
  }
  .ctx-sheet-enter-from,
  .ctx-sheet-leave-to {
    transform: none;
  }
}
</style>
