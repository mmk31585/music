import { z } from 'zod'

const optionalNumber = z.preprocess((v) => {
  if (v === null || v === undefined || v === '') return undefined
  const n = Number(v)
  return Number.isFinite(n) ? n : undefined
}, z.number().optional())

export const ExtractedTagsSchema = z.object({
  title: z.string().optional(),
  artist: z.string().optional(),
  album: z.string().optional(),
  albumArtist: z.string().optional(),
  trackNumber: optionalNumber,
  trackTotal: optionalNumber,
  discNumber: optionalNumber,
  discTotal: optionalNumber,
  year: optionalNumber,
  genre: z.string().optional(),
  comment: z.string().optional(),
  composer: z.string().optional(),
  lyrics: z.string().optional(),
  duration: optionalNumber,
  bitrate: optionalNumber,
  format: z.string().optional(),
  hasCoverArt: z.boolean().optional(),
})

export type ExtractedTags = z.infer<typeof ExtractedTagsSchema>

export const AssetResponseSchema = z.object({
  id: z.string(),
  assetType: z.string(),
  url: z.string(),
  source: z.string(),
  createdAt: z.string(),
})

export type AssetResponse = z.infer<typeof AssetResponseSchema>

export const UploadResponseSchema = z.object({
  draftId: z.string(),
  originalFilename: z.string(),
  fileSize: z.number(),
  format: z.string(),
  durationSeconds: optionalNumber,
  bitrate: optionalNumber,
  status: z.string(),
  coverArtUrl: z.string().optional(),
  extractedMetadata: ExtractedTagsSchema.optional(),
  assets: z.array(AssetResponseSchema).optional(),
  createdAt: z.string(),
})

export type UploadResponse = z.infer<typeof UploadResponseSchema>

export const DraftListItemSchema = z.object({
  id: z.string(),
  originalFilename: z.string(),
  fileSize: z.number(),
  format: z.string(),
  durationSeconds: optionalNumber,
  status: z.string(),
  title: z.string().optional(),
  artist: z.string().optional(),
  album: z.string().optional(),
  coverArtUrl: z.string().optional(),
  hasCoverArt: z.boolean().optional(),
  createdAt: z.string(),
})

export type DraftListItem = z.infer<typeof DraftListItemSchema>

export const EnrichedSuggestionSchema = z.object({
  field: z.string(),
  value: z.any(),
  source: z.string(),
  confidence: z.string(),
})

export type EnrichedSuggestion = z.infer<typeof EnrichedSuggestionSchema>

export const MusicBrainzResultSchema = z.object({
  mbid: z.string().optional(),
  title: z.string().optional(),
  artistName: z.string().optional(),
  artistMbid: z.string().optional(),
  albumName: z.string().optional(),
  albumMbid: z.string().optional(),
  releaseYear: z.number().optional(),
  duration: z.number().optional(),
  genres: z.array(z.string()).catch([]),
})

export const LastFMResultSchema = z.object({
  playCount: z.number().optional(),
  listenerCount: z.number().optional(),
  tags: z.array(z.string()).optional(),
  artistBio: z.string().optional(),
  artistImageUrl: z.string().optional(),
  albumCoverUrl: z.string().optional(),
  similarArtists: z.array(z.string()).optional(),
})

export const SpotifyResultSchema = z.object({
  spotifyId: z.string().optional(),
  previewUrl: z.string().optional(),
  albumCoverUrl: z.string().optional(),
  artistImageUrl: z.string().optional(),
  popularity: z.number().optional(),
})

export const LRCLibResultSchema = z.object({
  id: z.number().optional(),
  trackName: z.string().optional(),
  artistName: z.string().optional(),
  albumName: z.string().optional(),
  duration: z.number().optional(),
  synced: z.boolean().optional(),
  plainLyrics: z.string().optional(),
  syncedLyrics: z.string().optional(),
})

export type LRCLibResult = z.infer<typeof LRCLibResultSchema>

export const EnrichmentResultSchema = z.object({
  musicbrainz: MusicBrainzResultSchema.optional(),
  lastfm: LastFMResultSchema.optional(),
  spotify: SpotifyResultSchema.optional(),
  lrclib: LRCLibResultSchema.optional(),
  suggestions: z.array(EnrichedSuggestionSchema).catch([]),
  enrichment_attempted: z.boolean(),
})

export type EnrichmentResult = z.infer<typeof EnrichmentResultSchema>

export const DraftDetailResponseSchema = z.object({
  id: z.string(),
  uploadedBy: z.string(),
  originalFilename: z.string(),
  fileSize: z.number(),
  format: z.string(),
  durationSeconds: optionalNumber,
  bitrate: optionalNumber,
  status: z.string(),
  extractedMetadata: ExtractedTagsSchema.optional(),
  enrichedMetadata: EnrichmentResultSchema.optional(),
  finalMetadata: z.any().optional(),
  assets: z.array(AssetResponseSchema).optional(),
  createdAt: z.string(),
  updatedAt: z.string(),
})

export type DraftDetailResponse = z.infer<typeof DraftDetailResponseSchema>

export const FinalArtistMetadataSchema = z.object({
  action: z.string(),
  existingId: z.string().optional(),
  name: z.string(),
  bio: z.string().optional(),
  imageUrl: z.string().optional(),
  country: z.string().optional(),
  musicbrainzMbid: z.string().optional(),
})

export type FinalArtistMetadata = z.infer<typeof FinalArtistMetadataSchema>

export const FinalAlbumMetadataSchema = z.object({
  action: z.string(),
  existingId: z.string().optional(),
  title: z.string(),
  releaseYear: z.number().optional(),
  genre: z.string().optional(),
  coverUrl: z.string().optional(),
  musicbrainzReleaseId: z.string().optional(),
})

export type FinalAlbumMetadata = z.infer<typeof FinalAlbumMetadataSchema>

export const FinalTrackMetadataSchema = z.object({
  title: z.string(),
  trackNumber: z.number().optional(),
  durationSeconds: z.number(),
  genre: z.string().optional(),
  lyrics: z.string().optional(),
  explicit: z.boolean(),
  spotifyPreviewUrl: z.string().optional(),
  coverUrl: z.string().optional(),
})

export type FinalTrackMetadata = z.infer<typeof FinalTrackMetadataSchema>

export const SaveFinalMetadataRequestSchema = z.object({
  artist: FinalArtistMetadataSchema,
  album: FinalAlbumMetadataSchema,
  track: FinalTrackMetadataSchema,
})

export type SaveFinalMetadataRequest = z.infer<typeof SaveFinalMetadataRequestSchema>

export const ArtistSearchResultSchema = z.object({
  id: z.string(),
  name: z.string(),
  slug: z.string(),
  bio: z.string().optional(),
  imageUrl: z.string().optional(),
  country: z.string().optional(),
})

export type ArtistSearchResult = z.infer<typeof ArtistSearchResultSchema>

export const AlbumSearchResultSchema = z.object({
  id: z.string(),
  title: z.string(),
  slug: z.string(),
  artistName: z.string().optional(),
  coverUrl: z.string().optional(),
  releaseYear: z.number().optional(),
})

export type AlbumSearchResult = z.infer<typeof AlbumSearchResultSchema>

export const FinalizeResultSchema = z.object({
  artistId: z.string(),
  albumId: z.string().optional(),
  trackId: z.string(),
  audioUrl: z.string(),
  coverUrl: z.string().optional(),
})

export type FinalizeResult = z.infer<typeof FinalizeResultSchema>

export const ListDraftsResponseSchema = z.object({
  items: z.array(DraftListItemSchema),
  total: z.number(),
  page: z.number(),
  limit: z.number(),
})

export type ListDraftsResponse = z.infer<typeof ListDraftsResponseSchema>

export const IngestionStatsSchema = z.object({
  publishedThisMonth: z.number(),
  pendingReview: z.number(),
  totalDrafts: z.number(),
  byStatus: z.record(z.string(), z.number()),
})

export type IngestionStats = z.infer<typeof IngestionStatsSchema>

export const IngestionConfigSchema = z.object({
  maxUploadSize: z.number(),
  enrichmentEnabled: z.boolean(),
})

export type IngestionConfig = z.infer<typeof IngestionConfigSchema>
