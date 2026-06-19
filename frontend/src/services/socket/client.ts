import { ref } from 'vue'
import { useUserAuthStore } from '@/stores/user-auth'

const WS_URL = import.meta.env.VITE_WS_URL || undefined

type MessageHandler = (data: any) => void

let socket: WebSocket | null = null
let reconnectTimer: ReturnType<typeof setTimeout> | null = null
let pingTimer: ReturnType<typeof setInterval> | null = null
let messageQueue: { type: string; payload?: any }[] = []
let intentionalClose = false
let reconnectAttempts = 0
const maxReconnectAttempts = 20

const handlers = new Map<string, Set<MessageHandler>>()
const globalHandlers = new Set<MessageHandler>()
const subscriptions = new Set<string>()

const isConnected = ref(false)
let unsubscribeTokenWatch: (() => void) | null = null

function getToken(): string | null {
  const auth = useUserAuthStore()
  return auth.token
}

function connectInner() {
  if (
    socket &&
    (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING)
  ) {
    return
  }

  intentionalClose = false
  const token = getToken()
  if (!token) return

  try {
    const base = WS_URL || `${location.protocol === 'https:' ? 'wss' : 'ws'}://${location.host}/api/v1/ws`
    socket = new WebSocket(`${base}?token=${encodeURIComponent(token)}`)
  } catch {
    scheduleReconnect()
    return
  }

  socket.onopen = () => {
    isConnected.value = true
    reconnectAttempts = 0
    resubscribeAll()
    flushQueue()
    startPing()
  }

  socket.onmessage = (event) => {
    try {
      const msg = JSON.parse(event.data)
      const typeHandlers = handlers.get(msg.type)
      if (typeHandlers) {
        typeHandlers.forEach((h) => h(msg))
      }
      globalHandlers.forEach((h) => h(msg))
    } catch {
      /* ignore parse errors */
    }
  }

  socket.onclose = () => {
    isConnected.value = false
    stopPing()
    if (!intentionalClose) {
      scheduleReconnect()
    }
  }

  socket.onerror = () => {
    socket?.close()
  }
}

function resubscribeAll() {
  if (subscriptions.size === 0) return
  const batch = Array.from(subscriptions)
  subscriptions.clear()
  batch.forEach((ch) => {
    subscriptions.add(ch)
    rawSend('subscribe', ch)
  })
}

function flushQueue() {
  if (messageQueue.length === 0) return
  const batch = messageQueue.splice(0)
  batch.forEach((m) => rawSend(m.type, m.payload))
}

function rawSend(type: string, payload?: any) {
  if (socket?.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify({ type, payload }))
  }
}

function scheduleReconnect() {
  if (reconnectTimer || reconnectAttempts >= maxReconnectAttempts) return
  reconnectAttempts++
  const delay = Math.min(1000 * Math.pow(2, reconnectAttempts), 30000)
  reconnectTimer = setTimeout(() => {
    reconnectTimer = null
    connectInner()
  }, delay)
}

function startPing() {
  stopPing()
  pingTimer = setInterval(() => {
    if (socket?.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify({ type: 'ping' }))
    }
  }, 25000)
}

function stopPing() {
  if (pingTimer) {
    clearInterval(pingTimer)
    pingTimer = null
  }
}

export const wsClient = {
  connect() {
    connectInner()
    this.setupTokenWatch()
  },

  disconnect() {
    intentionalClose = true
    stopPing()
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    if (unsubscribeTokenWatch) {
      unsubscribeTokenWatch()
      unsubscribeTokenWatch = null
    }
    messageQueue = []
    subscriptions.clear()
    handlers.clear()
    globalHandlers.clear()
    socket?.close()
    socket = null
    isConnected.value = false
  },

  send(type: string, payload?: any) {
    if (socket?.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify({ type, payload }))
    } else {
      messageQueue.push({ type, payload })
    }
  },

  subscribe(channel: string) {
    subscriptions.add(channel)
    this.send('subscribe', channel)
  },

  unsubscribe(channel: string) {
    subscriptions.delete(channel)
    this.send('unsubscribe', channel)
  },

  on(type: string, handler: MessageHandler) {
    if (!handlers.has(type)) {
      handlers.set(type, new Set())
    }
    handlers.get(type)!.add(handler)
    return () => {
      handlers.get(type)?.delete(handler)
    }
  },

  onAny(handler: MessageHandler) {
    globalHandlers.add(handler)
    return () => {
      globalHandlers.delete(handler)
    }
  },

  get connected() {
    return isConnected.value
  },

  connectedRef: isConnected,

  reconnect() {
    this.disconnect()
    this.connect()
  },

  setupTokenWatch() {
    if (unsubscribeTokenWatch) return
    const auth = useUserAuthStore()
    let prev = auth.token
    unsubscribeTokenWatch = auth.$subscribe(() => {
      if (auth.token && auth.token !== prev) {
        prev = auth.token
        if (intentionalClose) return
        this.reconnect()
      }
    })
  },
}
