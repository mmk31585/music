export class PreloadManager {
  private preloadAudio: HTMLAudioElement | null = null
  private currentUrl: string | null = null

  preload(url: string) {
    if (!url || this.currentUrl === url) return

    this.cancel()

    this.currentUrl = url
    this.preloadAudio = new Audio()
    this.preloadAudio.preload = 'metadata'
    this.preloadAudio.src = url
    this.preloadAudio.load()
  }

  cancel() {
    if (!this.preloadAudio) return

    this.preloadAudio.pause()
    this.preloadAudio.removeAttribute('src')
    this.preloadAudio.load()

    this.preloadAudio = null
    this.currentUrl = null
  }
}

export const preloadManager = new PreloadManager()
