/**
 * Unified context menu type contract.
 *
 * Both the desktop floating menu and the mobile bottom action sheet
 * consume this same structure. One source of truth, two presentations.
 */

// ── Actions ────────────────────────────────────────────────────────

export interface ContextMenuAction {
  /** Unique identifier (used for keyboard navigation, typeahead, and testing). */
  id: string
  /** Display label (supports Persian and English). */
  label: string
  /** Lucide icon name string (rendered by the component). */
  icon: string
  /** Optional subtitle shown below the label (e.g. "Artist page", "Lossless"). */
  subtitle?: string
  /** Optional keyboard shortcut hint shown on the right (e.g. "S", "Ctrl+C"). */
  shortcut?: string
  /** Optional badge / counter shown on the right (e.g. "12 songs", "3 playlists"). */
  badge?: string
  /** Disabled state — action exists but is currently unavailable. */
  disabled?: boolean
  /** Danger style — destructive actions (delete, remove, report, block). */
  danger?: boolean
  /** Checked state — toggled-on indicator (filled heart, checkmark, etc.). */
  checked?: boolean
  /** If true, a visual separator is rendered BEFORE this action. */
  separator?: boolean
  /** If true, the action is hidden (context-aware visibility). */
  hidden?: boolean
  /** The action callback. */
  action: () => void
}

// ── Sections ───────────────────────────────────────────────────────

export interface ContextMenuSection {
  /** Unique section id (used for keying). */
  id: string
  /** Section header label (e.g. "PLAYBACK", "LIBRARY"). Rendered as 11px uppercase. */
  label?: string
  /** Ordered list of actions in this section. */
  items: ContextMenuAction[]
}

// ── Header (premium track info) ────────────────────────────────────

export interface ContextMenuHeader {
  coverUrl?: string
  title: string
  artistName: string
  artistId?: string | number
  albumName?: string
  albumId?: string | number
  duration?: number
  explicit?: boolean
  isPlaying?: boolean
  isLiked?: boolean
  isDownloaded?: boolean
  isHighQuality?: boolean
}

// ── Positioning ────────────────────────────────────────────────────

export interface ContextMenuPosition {
  x: number
  y: number
}

// ── Component props ────────────────────────────────────────────────

export interface ContextMenuProps {
  /** Whether the menu is visible. Use v-model:visible for two-way binding. */
  visible: boolean
  /** Ordered sections of actions to render. */
  sections: ContextMenuSection[]
  /** Optional premium header shown at the top (track info). */
  header?: ContextMenuHeader
  /** Cursor / anchor position (desktop only). */
  position?: ContextMenuPosition
  /** Dynamic album accent color for hover tints and focus rings. */
  accentColor?: string
  /** Optional HTML element to anchor the menu to (e.g. overflow button). */
  anchorEl?: HTMLElement | null
  /** Accessible label for the menu. */
  ariaLabel?: string
}
