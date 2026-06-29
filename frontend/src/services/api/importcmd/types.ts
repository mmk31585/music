export interface SearchResult {
  title: string
  artist: string
  album?: string
  url: string
  duration: number
  thumbnail: string
  source: string
  score?: number
  isrc?: string
  external_ids?: Record<string, string>
}

export interface ImportRequest {
  url?: string
  source?: string
  title?: string
  artist?: string
  duration?: number
  isrc?: string
  external_ids?: Record<string, string>
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

// ── Artist Search types ───────────────────────────────────────────

export interface ArtistSearchResponse {
  artist_info: ArtistInfo
  albums: AlbumGroup[]
}

export interface ArtistInfo {
  name: string
  image: string
}

export interface AlbumGroup {
  title: string
  cover: string
  source: string
  tracks: TrackResult[]
}

export interface TrackResult {
  title: string
  duration: number
  source: string
  album?: string
  external_ids?: Record<string, string>
}

// ── Batch Import types ────────────────────────────────────────────

export interface BatchImportItem {
  title: string
  artist: string
  url?: string
  album?: string
  duration?: number
  source?: string
  external_ids?: Record<string, string>
}

export interface BatchImportResponse {
  batchId: string
  jobs: BatchJobResult[]
  message: string
}

export interface BatchJobResult {
  title: string
  artist: string
  jobId?: string
  error?: string
  status?: string
}

export interface BatchProgressResponse {
  batchId: string
  total: number
  completed: number
  failed: number
  inProgress: number
  progressPct: number
}
