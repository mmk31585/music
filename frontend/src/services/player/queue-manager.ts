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

  clear() {
    this.queue = []
    this.index = -1
  }
}

export const queueManager = new QueueManager()
