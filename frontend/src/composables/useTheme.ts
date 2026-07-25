import { reactive, watch } from 'vue'

export type ThemeMode = 'light' | 'dark' | 'system'
export type PrimaryColor = 'purple' | 'blue' | 'emerald' | 'rose' | 'orange' | 'amber' | 'pink' | 'indigo' | 'neutral'
export type SurfaceColor = 'slate' | 'graphite' | 'midnight' | 'neutral' | 'warm'

export interface ThemeState {
  mode: ThemeMode
  primary: PrimaryColor
  surface: SurfaceColor
}

const STORAGE_KEY = 'muse-theme'

const primaryPalettes: Record<PrimaryColor, Record<string, string>> = {
  purple: {
    50: '#faf5ff', 100: '#f3e8ff', 200: '#e9d5ff', 300: '#d8b4fe',
    400: '#c084fc', 500: '#a855f7', 600: '#9333ea', 700: '#7e22ce',
    800: '#6b21a8', 900: '#581c87', 950: '#3b0764',
  },
  blue: {
    50: '#eff6ff', 100: '#dbeafe', 200: '#bfdbfe', 300: '#93c5fd',
    400: '#60a5fa', 500: '#3b82f6', 600: '#2563eb', 700: '#1d4ed8',
    800: '#1e40af', 900: '#1e3a8a', 950: '#172554',
  },
  emerald: {
    50: '#ecfdf5', 100: '#d1fae5', 200: '#a7f3d0', 300: '#6ee7b7',
    400: '#34d399', 500: '#1db954', 600: '#059669', 700: '#047857',
    800: '#065f46', 900: '#064e3b', 950: '#022c22',
  },
  rose: {
    50: '#fff1f2', 100: '#ffe4e6', 200: '#fecdd3', 300: '#fda4af',
    400: '#fb7185', 500: '#f43f5e', 600: '#e11d48', 700: '#be123c',
    800: '#9f1239', 900: '#881337', 950: '#4c0519',
  },
  orange: {
    50: '#fff7ed', 100: '#ffedd5', 200: '#fed7aa', 300: '#fdba74',
    400: '#fb923c', 500: '#f97316', 600: '#ea580c', 700: '#c2410c',
    800: '#9a3412', 900: '#7c2d12', 950: '#431407',
  },
  amber: {
    50: '#fffbeb', 100: '#fef3c7', 200: '#fde68a', 300: '#fcd34d',
    400: '#fbbf24', 500: '#f59e0b', 600: '#d97706', 700: '#b45309',
    800: '#92400e', 900: '#78350f', 950: '#451a03',
  },
  pink: {
    50: '#fdf2f8', 100: '#fce7f3', 200: '#fbcfe8', 300: '#f9a8d4',
    400: '#f472b6', 500: '#ec4899', 600: '#db2777', 700: '#be185d',
    800: '#9d174d', 900: '#831843', 950: '#500724',
  },
  indigo: {
    50: '#eef2ff', 100: '#e0e7ff', 200: '#c7d2fe', 300: '#a5b4fc',
    400: '#818cf8', 500: '#6366f1', 600: '#4f46e5', 700: '#4338ca',
    800: '#3730a3', 900: '#312e81', 950: '#1e1b4b',
  },
  neutral: {
    50: '#fafafa', 100: '#f5f5f5', 200: '#e5e5e5', 300: '#d4d4d4',
    400: '#a3a3a3', 500: '#737373', 600: '#525252', 700: '#404040',
    800: '#262626', 900: '#171717', 950: '#0a0a0a',
  },
}

const surfacePalettes: Record<SurfaceColor, Record<string, string>> = {
  slate: {
    0: '#ffffff', 50: '#f8fafc', 100: '#f1f5f9', 200: '#e2e8f0',
    300: '#cbd5e1', 400: '#94a3b8', 500: '#64748b', 600: '#475569',
    700: '#334155', 800: '#1e293b', 900: '#0f172a', 950: '#020617',
  },
  graphite: {
    0: '#ffffff', 50: '#f5f5f5', 100: '#e8e8e8', 200: '#d1d1d1',
    300: '#b3b3b3', 400: '#8c8c8c', 500: '#6b6b6b', 600: '#525252',
    700: '#3d3d3d', 800: '#2a2a2a', 900: '#1a1a1a', 950: '#0d0d0d',
  },
  midnight: {
    0: '#ffffff', 50: '#f0f1f5', 100: '#d5d8e3', 200: '#b0b5cc',
    300: '#868db0', 400: '#636b94', 500: '#4a5078', 600: '#3b4060',
    700: '#2e334d', 800: '#21253b', 900: '#161928', 950: '#0c0e1a',
  },
  neutral: {
    0: '#ffffff', 50: '#f7f7f7', 100: '#ebebeb', 200: '#d6d6d6',
    300: '#b8b8b8', 400: '#969696', 500: '#7a7a7a', 600: '#626262',
    700: '#4d4d4d', 800: '#383838', 900: '#242424', 950: '#121212',
  },
  warm: {
    0: '#ffffff', 50: '#faf8f5', 100: '#f0ebe3', 200: '#e0d6c8',
    300: '#ccbea8', 400: '#b3a086', 500: '#9a876c', 600: '#7f6e58',
    700: '#665847', 800: '#4f4437', 900: '#3a3128', 950: '#221d17',
  },
}

function loadState(): ThemeState {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved) {
      const parsed = JSON.parse(saved) as Partial<ThemeState>
      return {
        mode: parsed.mode || 'dark',
        primary: parsed.primary || 'emerald',
        surface: parsed.surface || 'slate',
      }
    }
  } catch { }
  return { mode: 'dark', primary: 'emerald', surface: 'slate' }
}

function saveState(state: ThemeState) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(state))
}

function getSystemDark(): boolean {
  return window.matchMedia('(prefers-color-scheme: dark)').matches
}

function applyTheme(state: ThemeState) {
  const root = document.documentElement

  const isDark = state.mode === 'dark' || (state.mode === 'system' && getSystemDark())
  root.classList.toggle('app-dark', isDark)

  const primary = primaryPalettes[state.primary]
  for (const [shade, value] of Object.entries(primary)) {
    root.style.setProperty(`--p-${shade}`, value)
  }

  const surface = surfacePalettes[state.surface]
  for (const [shade, value] of Object.entries(surface)) {
    root.style.setProperty(`--s-${shade}`, value)
  }
}

const loadedState = loadState()
const state = reactive<ThemeState>({ ...loadedState })

applyTheme(state)

watch(state, (s) => {
  applyTheme(s)
  saveState(s)
}, { deep: true })

if (state.mode === 'system') {
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
    applyTheme(state)
  })
}

export function useTheme() {
  function setMode(mode: ThemeMode) {
    state.mode = mode
  }

  function setPrimary(primary: PrimaryColor) {
    state.primary = primary
  }

  function setSurface(surface: SurfaceColor) {
    state.surface = surface
  }

  function toggleMode() {
    const cycle: ThemeMode[] = ['dark', 'light', 'system']
    const idx = cycle.indexOf(state.mode)
    state.mode = cycle[(idx + 1) % cycle.length]!
  }

  function registerTheme(name: string, shades: Record<string, string>) {
    primaryPalettes[name as PrimaryColor] = shades
  }

  return {
    state,
    setMode,
    setPrimary,
    setSurface,
    toggleMode,
    registerTheme,
    primaryPalettes,
    surfacePalettes,
    availablePrimaries: Object.keys(primaryPalettes) as PrimaryColor[],
    availableSurfaces: Object.keys(surfacePalettes) as SurfaceColor[],
  }
}
