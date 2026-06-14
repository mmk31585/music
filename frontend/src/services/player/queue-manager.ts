import type { PlaybackTrack } from '@/services/api/player'

export class QueueManager {
  private queue: PlaybackTrack[] = []
  private index = -1

  setQueue(tracks: PlaybackTrack[], startIndex = 0) {
    this.queue = [...tracks]
    this.index = Math.max(0, Math.min(startIndex, this.queue.length - 1))
  }

  setCurrent(track: PlaybackTrack) {
    const existingIndex = this.queue.findIndex((item) => item.id === track.id)

    if (existingIndex >= 0) {
      this.index = existingIndex
      return
    }

    this.queue = [track, ...this.queue]
    this.index = 0
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
    const nextTrack = this.getNext()
    if (!nextTrack) return null

    this.index += 1
    return nextTrack
  }

  previous() {
    const previousTrack = this.getPrevious()
    if (!previousTrack) return null

    this.index -= 1
    return previousTrack
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
  }
}

export const queueManager = new QueueManager()
