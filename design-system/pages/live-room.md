# Live Room (/rooms/:id)

> Redesign of existing `LiveRoomStage.vue` and related social components
> Fixes: emoji icons, hover-only interactions, keyboard focus, glass border contrast

---

## Layout (Desktop)

```
┌─────────────────────────────────────────────────────┐
│  Header: Room name + participant count + share btn   │
│  ┌─────────────────────────────────────────────────┐│
│  │  ┌──────────────────────┐  ┌──────────────────┐ ││
│  │  │     STAGE AREA        │  │   CHAT SIDEBAR   │ ││
│  │  │                       │  │                  │ ││
│  │  │   ( 👑 ) Host         │  │  User: پیام تست  │ ││
│  │  │   Host Name           │  │  User2: سلام     │ ││
│  │  │   [HOST BADGE]        │  │                  │ ││
│  │  │                       │  │  ──────────────  │ ││
│  │  │   (S1) (S2) (S3)     │  │  [Input field]    │ ││
│  │  │   Speakers grid       │  │  [Send ▶]        │ ││
│  │  │                       │  │                  │ ││
│  │  │   Now Playing Hero    │  │                  │ ││
│  │  └──────────────────────┘  └──────────────────┘ ││
│  └─────────────────────────────────────────────────┘│
│                                                       │
│  Queue (bottom panel, collapsible)                    │
│  ┌─────────────────────────────────────────────────┐│
│  │ [⏮] [▶/⏸] [⏭] [Queue] [Add to queue]  volume ││
│  └─────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────┘
```

## Component Specs

### Stage Area — Fixed Anti-Patterns

**1. Emoji Icon → SVG Icon (CRITICAL)**
```vue
<!-- BEFORE (emoji) -->
<div v-else class="flex h-full w-full items-center justify-center bg-white/10 text-2xl">
  👑
</div>

<!-- AFTER (Lucide Crown icon) -->
<div v-else class="flex h-full w-full items-center justify-center bg-white/10">
  <svg class="h-6 w-6 text-yellow-400" viewBox="0 0 24 24" fill="none"
       stroke="currentColor" stroke-width="2" aria-hidden="true">
    <path d="M2 4l3 12h14l3-12-6 7-4-7-4 7-6-7z" />
    <path d="M3 20h18" />
  </svg>
</div>
```

**2. Hover-only actions → Click + Keyboard**
```vue
<!-- BEFORE: hover-only mute button -->
<div @mouseenter="hoveredSpeaker = s.user_id"
     @mouseleave="hoveredSpeaker = null">

<!-- AFTER: always visible actions with focus-visible -->
<div :key="speaker.user_id"
     class="group relative flex flex-col items-center gap-2">
  <!-- ...avatar -->
  <button aria-label="Mute speaker"
          class="absolute -top-1 -right-1 flex h-6 w-6 items-center justify-center
                 rounded-full bg-red-500/90 text-white opacity-0 group-hover:opacity-100
                 group-focus-within:opacity-100 transition-opacity
                 focus-visible:opacity-100 focus-visible:outline-2 focus-visible:outline-white"
          @click="toggleMute(speaker.user_id)">
    <Icon name="microphone-off" class="h-3 w-3" aria-hidden="true" />
  </button>
</div>
```

**3. Glass border contrast**
```vue
<!-- BEFORE -->
<div class="rounded-2xl bg-white/[0.04] p-6 ring-1 ring-white/[0.07]">

<!-- AFTER: stronger ring for min 3:1 non-text contrast -->
<div class="rounded-2xl bg-white/[0.06] p-6 ring-1 ring-white/[0.10]">
```

### Host Badge
```vue
<span class="inline-flex items-center gap-1 rounded-full bg-yellow-500/10
             px-2.5 py-0.5 text-[10px] font-semibold text-yellow-400">
  <Icon name="crown" class="h-3 w-3" aria-hidden="true" />
  Host
</span>
```

### Speaking Pulse Animation
```vue
<div v-if="speakingUserIds.has(speaker.user_id)"
     class="absolute inset-0 rounded-full
            motion-safe:animate-pulse ring-2 ring-[#1db954]
            ring-offset-2 ring-offset-transparent"
     role="status" aria-label="Currently speaking" />
```
- `motion-safe:` — respects `prefers-reduced-motion` ✅
- `role="status"` — announces to screen readers ✅

### Chat Sidebar
```vue
<div class="flex h-full flex-col rounded-2xl bg-white/[0.06] ring-1 ring-white/[0.10]"
     role="log" aria-label="Live chat" aria-live="polite">

  <div class="flex-1 overflow-y-auto p-4 space-y-3">
    <ChatMessage v-for="msg in messages" :key="msg.id" :message="msg" />
  </div>

  <form class="flex items-center gap-2 border-t border-white/[0.06] p-3"
        @submit.prevent="sendMessage">
    <input type="text"
           :placeholder="$t('chat.type_message')"
           class="flex-1 rounded-lg bg-white/[0.06] px-3 py-2 text-sm text-white
                  placeholder:text-white/30 outline-none
                  focus-visible:ring-2 focus-visible:ring-[#1db954]"
           v-model="messageText"
           aria-label="Type a message">
    <button type="submit"
            aria-label="Send message"
            class="flex h-10 w-10 items-center justify-center rounded-full
                   bg-[#1db954] text-black transition hover:scale-105
                   active:scale-95 focus-visible:outline-2 focus-visible:outline-white
                   disabled:opacity-50 disabled:cursor-not-allowed"
            :disabled="!messageText.trim()">
      <Icon name="send" class="h-4 w-4" aria-hidden="true" />
    </button>
  </form>
</div>
```

## States

| State | Stage Area | Chat |
|-------|-----------|------|
| **Loading** | Skeleton circles (3 large, 6 small) + shimmer | "Loading messages..." |
| **Empty Room** | "Be the first to join!" + Join button | "No messages yet" |
| **Speaking** | Green pulse ring on speaker's avatar | — |
| **Muted** | Red microphone-off badge | — |
| **Full Room** | All 12 speaker slots filled | — |
| **Error** | "Connection lost" banner + Reconnect | Messages queued locally |
| **RTL** | Layout flips right → left | Chat input on left side |

## Data Flow
```
LiveRoomStage.vue (PagePartyDetail.vue)
├── socialStore/room/:id → room state via WebSocket
├── socialStore.stageState → host, speakers, listeners
├── socialStore.speakingUserIds → from WebSocket voice events
├── chatStore.messages → from WebSocket chat events
├── playerStore.currentTrack → now playing display
├── WebSocket lifecycle:
│   ├── connect(roomId) → join room
│   ├── on('stage:update') → mutate stageState
│   ├── on('chat:message') → push to messages
│   ├── on('speaking:start') → add to speakingUserIds
│   └── on('disconnect') → set error state
└── cleanup onBeforeUnmount → disconnect socket
```
