export interface SearchResult {
  title: string
  artist: string
  url: string
  duration: number
  thumbnail: string
  source: string
}

export interface ImportRequest {
  url: string
}

export interface ImportResponse {
  draftId: string
  title: string
  artist: string
  duration: number
  message: string
}
