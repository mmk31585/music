export interface SearchResult {
  title: string
  artist: string
  url: string
  duration: number
  thumbnail: string
  source: string
  score?: number
  isrc?: string
  external_ids?: Record<string, string>
}

export interface ImportRequest {
  url: string
  source?: string
}

export interface ImportResponse {
  jobId?: string
  draftId?: string
  title?: string
  artist?: string
  duration?: number
  message: string
}

export interface ProgressResponse {
  jobId: string
  status: string
  progress: number
  stage: string
  error?: string
  draftId?: string
}
