import { useToast } from 'primevue/usetoast'

interface Shareable {
  id: string | number
  title: string
  type: 'track' | 'album' | 'artist' | 'playlist' | 'video'
  artistName?: string | null
  coverUrl?: string | null
}

export function useSocialShare() {
  const toast = useToast()

  function getShareUrl(item: Shareable): string {
    const base = window.location.origin
    switch (item.type) {
      case 'track':
        return `${base}/track/${item.id}`
      case 'album':
        return `${base}/album/${item.id}`
      case 'artist':
        return `${base}/artist/${item.id}`
      case 'playlist':
        return `${base}/playlists/${item.id}`
      case 'video':
        return `${base}/video/${item.id}`
    }
    return ''
  }

  function getShareText(item: Shareable): string {
    switch (item.type) {
      case 'track':
        return `🎵 ${item.title}${item.artistName ? ` — ${item.artistName}` : ''}`
      case 'album':
        return `💿 ${item.title}${item.artistName ? ` — ${item.artistName}` : ''}`
      case 'artist':
        return `🎤 ${item.title}`
      case 'playlist':
        return `📋 ${item.title}`
      case 'video':
        return `🎬 ${item.title}`
    }
    return ''
  }

  async function copyLink(item: Shareable): Promise<boolean> {
    const url = getShareUrl(item)
    try {
      await navigator.clipboard.writeText(url)
      toast.add({
        severity: 'success',
        summary: 'Link copied',
        detail: `${item.title} link copied to clipboard`,
        life: 2500,
      })
      return true
    } catch {
      toast.add({
        severity: 'error',
        summary: 'Failed to copy',
        detail: 'Could not copy link to clipboard',
        life: 3000,
      })
      return false
    }
  }

  function shareNative(item: Shareable) {
    const url = getShareUrl(item)
    const text = getShareText(item)
    if (navigator.share) {
      navigator.share({ title: item.title, text, url }).catch(() => {})
    } else {
      copyLink(item)
    }
  }

  return { copyLink, shareNative, getShareUrl, getShareText }
}
