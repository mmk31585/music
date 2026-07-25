<template>
  <div class="space-y-5">
    <div>
      <p class="mb-2 text-[11px] font-semibold tracking-widest uppercase text-tertiary">Theme</p>
      <div class="flex gap-1 rounded-xl bg-surface-overlay p-1">
        <button
          v-for="t in themeOptions"
          :key="t.value"
          type="button"
          :aria-label="t.label"
          :class="[
            'flex flex-1 items-center justify-center gap-1.5 rounded-lg px-3 py-2 text-xs font-medium transition-all duration-200',
            theme.state.mode === t.value
              ? 'bg-surface-active text-primary shadow-sm'
              : 'text-tertiary hover:text-secondary',
          ]"
          @click="theme.setMode(t.value)"
        >
          <span>{{ t.icon }}</span>
          <span>{{ t.label }}</span>
        </button>
      </div>
    </div>

    <div>
      <p class="mb-2 text-[11px] font-semibold tracking-widest uppercase text-tertiary">Primary Color</p>
      <div class="flex flex-wrap gap-2">
        <button
          v-for="color in theme.availablePrimaries"
          :key="color"
          type="button"
          :aria-label="color"
          :title="color"
          :class="[
            'h-7 w-7 rounded-full ring-1 ring-border-default transition-all duration-200 hover:scale-110 focus-visible:outline-2 focus-visible:outline-offset-2',
            theme.state.primary === color ? 'scale-110 ring-2 ring-primary' : '',
          ]"
          :style="{ backgroundColor: theme.primaryPalettes[color][500] }"
          @click="theme.setPrimary(color)"
        />
      </div>
    </div>

    <div>
      <p class="mb-2 text-[11px] font-semibold tracking-widest uppercase text-tertiary">Surface</p>
      <div class="grid grid-cols-2 gap-2">
        <button
          v-for="s in theme.availableSurfaces"
          :key="s"
          type="button"
          :class="[
            'flex items-center gap-2 rounded-xl px-3 py-2.5 text-xs font-medium transition-all duration-200',
            theme.state.surface === s
              ? 'bg-surface-active text-primary ring-1 ring-border-strong'
              : 'bg-surface-hover text-tertiary hover:bg-surface-active hover:text-secondary',
          ]"
          @click="theme.setSurface(s)"
        >
          <span class="h-4 w-4 rounded-full ring-1 ring-border-default" :style="{ backgroundColor: theme.surfacePalettes[s][700] }" />
          <span class="capitalize">{{ s }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTheme } from '@/composables/useTheme'

const theme = useTheme()

const themeOptions = [
  { value: 'dark' as const, label: 'Dark', icon: '🌙' },
  { value: 'light' as const, label: 'Light', icon: '☀' },
  { value: 'system' as const, label: 'System', icon: '🖥' },
]
</script>
