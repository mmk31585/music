export enum AIApiRoutes {
  GENERATE_EMBEDDING = '/ai/embeddings',
  ANALYZE_MOOD = '/ai/moods',
  GET_MOOD = '/ai/moods/:trackId',
  GENERATE_PLAYLIST = '/ai/playlists/generate',
  SIMILAR_BY_MOOD = '/ai/similar/mood/:trackId',
  SIMILAR_BY_EMBEDDING = '/ai/similar/embedding/:trackId',
}
