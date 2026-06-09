export interface DocumentPictureInPictureOptions {
  width?: number
  height?: number
  disallowReturnToOpener?: boolean
  preferInitialWindowPlacement?: boolean
}

export interface DocumentPictureInPictureEvent extends Event {
  window: Window
}

export interface DocumentPictureInPicture {
  window: Window | null
  requestWindow: (options?: DocumentPictureInPictureOptions) => Promise<Window>
  addEventListener: (
    type: 'enter',
    listener: (event: DocumentPictureInPictureEvent) => void,
    options?: boolean | AddEventListenerOptions,
  ) => void
  removeEventListener: (
    type: 'enter',
    listener: (event: DocumentPictureInPictureEvent) => void,
    options?: boolean | EventListenerOptions,
  ) => void
}

declare global {
  interface Window {
    documentPictureInPicture?: DocumentPictureInPicture
  }
}
