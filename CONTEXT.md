# Muse — Domain Glossary

## Core Domain Concepts

### Player Engine
The module that composes audio playback, queue management, track preloading, and Media Session API integration behind a single seam. Exposes actions (`play`, `pause`, `seek`, `next`, `previous`) and emits events (`trackchange`, `playstate`, `timeupdate`, `buffering`, `error`). Accepts optional provider interfaces for network-dependent operations (track fetching, play history). Testable by substituting providers and the audio backend.

### Track
A single audio entity with id, title, artist, album, duration, streamUrl, and coverUrl.

### PlaybackTrack
The runtime representation of a Track used by the Player Engine. Adds normalized media URLs and runtime-only fields.

### Queue
An ordered list of PlaybackTracks awaiting playback. Managed by QueueManager within the Player Engine.

### Shuffle
Four modes: off (sequential), queue (shuffled order of current queue), catalog (random tracks from full catalog), similar (random tracks similar to current).

### Repeat
Three modes: off, one (repeat current track), all (loop queue).

## Architecture Terms

### Seam
A boundary where one module can be replaced without changing its consumers. In the frontend, the PlayerEngine interface is a seam between playback orchestration and UI state.

### Deep Module
A module whose interface is significantly simpler than its implementation. The PlayerEngine is deep: 12 action methods + 6 events replace 10+ direct singleton dependencies spread across the store.

### Locality
When related logic lives in one module rather than scattered. Playback orchestration (audio + queue + preload + media-session + shuffle + repeat) has locality inside the Player Engine.
