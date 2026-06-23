import type { PlaybackTrack } from '@/services/api/player'

interface HistoryEntry {
  trackId: string
  queuePosition: number
}

export class QueueManager {
  private queue: PlaybackTrack[] = []
  private index = -1

  /** Play history stack: tracks the order tracks were actually played,
   *  so that previous() navigates back in the correct order even when
   *  shuffleMode jumps around the queue. */
  private history: HistoryEntry[] = []

  setQueue(tracks: PlaybackTrack[], startIndex = 0) {
    this.queue = [...tracks]
    this.index = Math.max(0, Math.min(startIndex, this.queue.length - 1))
    this.history = []
  }

  /**
   * Mark the given track as "currently playing" and push it onto the history stack.
   * If it's already the current track, no history entry is added.
   */
  setCurrent(track: PlaybackTrack) {
    const existingIndex = this.queue.findIndex((item) => item.id === track.id)

    if (existingIndex >= 0) {
      // Only push to history if we're actually moving to a different track
      if (this.index !== existingIndex) {
        this.pushHistory(track.id, existingIndex)
        this.index = existingIndex
      }
      return
    }

    // Track not in queue — prepend it
    this.queue = [track, ...this.queue]
    this.index = 0
    this.pushHistory(track.id, 0)
  }

  private pushHistory(trackId: string, queuePosition: number) {
    this.history.push({ trackId, queuePosition })
  }

  getCurrent() {
    if (this.index < 0) return null
    return this.queue[this.index] || null
  }

  getNext() {
    if (this.index < 0) return null
    return this.queue[this.index + 1] || null
  }

  getPrevious() {
    if (this.index <= 0) return null
    return this.queue[this.index - 1] || null
  }

  next() {
    // Push current position to history before moving forward
    const current = this.getCurrent()
    if (current) {
      this.pushHistory(current.id, this.index)
    }

    const nextTrack = this.getNext()
    if (!nextTrack) return null

    this.index += 1
    return nextTrack
  }

  /**
   * Navigate back through play history.
   * Pops the last history entry and navigates to that track.
   * Falls back to queue sequential navigation if history is empty.
   */
  previous() {
    // Pop the last history entry to go back
    const lastEntry = this.history.pop()
    if (lastEntry) {
      // Find the track in the current queue
      const trackIndex = this.queue.findIndex((t) => t.id === lastEntry.trackId)
      if (trackIndex >= 0) {
        this.index = trackIndex
        return this.queue[trackIndex]!
      }
    }

    // Fallback: linear previous
    const previousTrack = this.getPrevious()
    if (!previousTrack) return null

    this.index -= 1
    return previousTrack
  }

  /**
   * Returns the full history stack (for debugging)
   */
  getHistory(): readonly HistoryEntry[] {
    return this.history
  }

  /**
   * Clears history (e.g. when shuffle mode changes)
   */
  clearHistory() {
    this.history = []
  }

  all() {
    return [...this.queue]
  }

  currentIndex() {
    return this.index
  }

  replaceAll(tracks: PlaybackTrack[], preserveIndex = true) {
    const currentId = preserveIndex ? this.queue[this.index]?.id : undefined
    this.queue = [...tracks]
    if (currentId !== undefined) {
      const newIndex = this.queue.findIndex((t) => t.id === currentId)
      this.index = newIndex >= 0 ? newIndex : 0
    } else {
      this.index = this.queue.length > 0 ? 0 : -1
    }
    this.history = []
  }

  reorderQueue(oldIndex: number, newIndex: number) {
    const [moved] = this.queue.splice(oldIndex, 1)
    if (!moved) return
    this.queue.splice(newIndex, 0, moved)

    if (oldIndex === this.index) {
      this.index = newIndex
    } else if (oldIndex < this.index && newIndex >= this.index) {
      this.index -= 1
    } else if (oldIndex > this.index && newIndex <= this.index) {
      this.index += 1
    }
  }

  clear() {
    this.queue = []
    this.index = -1
    this.history = []
  }
}

export const queueManager = new QueueManager()
