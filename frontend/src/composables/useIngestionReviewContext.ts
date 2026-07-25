import { provide, inject, type Ref, type ComputedRef } from 'vue'
import type {
  SaveFinalMetadataRequest,
  DraftDetailResponse,
  EnrichmentResult,
  EnrichedSuggestion,
  ArtistSearchResult,
  AlbumSearchResult,
} from '@/services/api/ingestion/types'

export interface IngestionReviewContext {
  // Core data
  draftId: string
  finalMetadata: SaveFinalMetadataRequest
  draftDetail: Ref<DraftDetailResponse | null>
  enrichment: Ref<EnrichmentResult | null>

  // Step navigation
  step: Ref<number>
  canProceed: ComputedRef<boolean>
  isValid: ComputedRef<boolean>
  hasUnsavedChanges: Ref<boolean>

  // Loading / publishing state
  loading: Ref<boolean>
  error: Ref<string | null>
  publishing: Ref<boolean>
  publishingError: Ref<string | null>
  rejecting: Ref<boolean>
  published: Ref<boolean>
  publishResult: Ref<{ artistId?: string; albumId?: string; trackId?: string; audioUrl: string } | null>
  enrichmentComplete: Ref<boolean>

  // Enrichment per-section state
  enriching: Ref<boolean>
  enrichingArtist: Ref<boolean>
  enrichingAlbum: Ref<boolean>
  enrichingTrack: Ref<boolean>

  // Computed helpers
  embeddedCover: ComputedRef<string | null>
  artistImageAsset: ComputedRef<string | null>
  albumCoverAsset: ComputedRef<string | null>
  suggestedAlbumCover: ComputedRef<string | null>
  suggestedLastfmAlbumCover: ComputedRef<string | null>
  extractedTitle: ComputedRef<string>
  extractedArtist: ComputedRef<string>
  extractedAlbum: ComputedRef<string>
  trackAudioUrl: ComputedRef<string | null>
  lyricsType: ComputedRef<string>

  // Audio preview state
  audioCurrentTime: Ref<number>
  audioDuration: Ref<number>
  audioPlaying: Ref<boolean>
  showSyncedLyrics: Ref<boolean>
  lyricsExpanded: Ref<boolean>

  // Search state
  artistSearchQuery: Ref<string>
  albumSearchQuery: Ref<string>
  artistSearchResults: Ref<ArtistSearchResult[]>
  albumSearchResults: Ref<AlbumSearchResult[]>
  artistSearching: Ref<boolean>
  albumSearching: Ref<boolean>
  selectedArtistId: Ref<string | undefined>
  selectedAlbumId: Ref<string | undefined>

  // Actions
  findSuggestion(field: string): EnrichedSuggestion | undefined
  sourceLabel(source: string): string
  sourceSeverity(source: string): string | undefined
  confidenceSeverity(confidence: string): string | undefined
  truncateBio(bio: string): string
  stepperClass(i: number): string
  stepIconClass(i: number): string
  seekAudio(seconds: number): void

  addFeaturedArtist(): void
  removeFeaturedArtist(index: number): void
  moveFeaturedArtist(fromIndex: number, direction: -1 | 1): void
  selectExistingArtist(a: ArtistSearchResult, index?: number): void
  selectExistingAlbum(a: AlbumSearchResult): void
  uploadArtistImage(e: Event, index?: number): void
  uploadAlbumCover(e: Event): void
  fetchSection(section: 'artist' | 'album' | 'track'): void
  debouncedArtistSearch(): void
  debouncedAlbumSearch(): void
  triggerEnrich(): void
  publish(): void
  confirmReject(): void
  goBack(): void
}

const CONTEXT_KEY = 'ingestion-review'

export function provideIngestionReviewContext(ctx: IngestionReviewContext) {
  provide(CONTEXT_KEY, ctx)
}

export function useIngestionReviewContext(): IngestionReviewContext {
  const ctx = inject<IngestionReviewContext>(CONTEXT_KEY)
  if (!ctx) throw new Error('useIngestionReviewContext must be used within a provider')
  return ctx
}
