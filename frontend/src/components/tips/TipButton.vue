<template>
  <div>
    <button
      type="button"
      @click="openDialog"
      class="inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/[0.04] px-4 py-2 text-sm font-medium text-white/60 backdrop-blur transition hover:border-[#e91e63]/30 hover:bg-[#e91e63]/10 hover:text-[#e91e63]"
    >
      <i aria-hidden="true" class="pi pi-heart text-xs" />
      Tip
    </button>

    <Teleport to="body">
      <div
        v-if="visible"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
        @click.self="visible = false"
      >
        <div class="glass-strong mx-4 w-full max-w-sm rounded-2xl p-6">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-bold text-white">Send a Tip</h3>
            <button aria-label="Close dialog" class="text-white/30 hover:text-white/50" @click="visible = false">
              <i aria-hidden="true" class="pi pi-times" />
            </button>
          </div>

          <div v-if="!sent" class="mt-4 space-y-4">
            <div>
              <label class="mb-1.5 block text-xs text-white/40">Amount</label>
              <div class="grid grid-cols-4 gap-2">
                <button
                  v-for="a in amounts"
                  :key="a"
                  class="rounded-xl py-2 text-sm font-medium transition"
                  :class="
                    selectedAmount === a
                      ? 'bg-[#1db954] text-black'
                      : 'bg-white/5 text-white/60 hover:bg-white/10'
                  "
                  @click="selectedAmount = a"
                >
                  {{ formatAmount(a) }}
                </button>
              </div>
              <div class="mt-2">
                <input
                  v-model.number="customAmount"
                  placeholder="Custom"
                  aria-label="Custom tip amount"
                  class="w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2.5 text-sm text-white outline-none focus:border-white/20"
                  type="number"
                  min="1000"
                  step="1000"
                  @input="selectedAmount = 0"
                />
              </div>
            </div>

            <div>
              <label class="mb-1.5 block text-xs text-white/40">Message (optional)</label>
              <textarea
                v-model="message"
                rows="2"
                placeholder="Say something nice..."
                aria-label="Tip message"
                class="w-full resize-none rounded-xl border border-white/10 bg-white/5 px-4 py-2.5 text-sm text-white outline-none focus:border-white/20"
              />
            </div>

            <button
              @click="sendTip"
              class="flex w-full items-center justify-center gap-2 rounded-xl bg-gradient-to-r from-[#e91e63] to-[#ff6b9d] py-3 text-sm font-bold text-white transition hover:scale-[1.02] disabled:opacity-50"
              :disabled="sending"
            >
              <i aria-hidden="true" v-if="sending" class="pi pi-spin pi-spinner" />
              <i aria-hidden="true" v-else class="pi pi-heart" />
              {{ sending ? 'Processing...' : `Send ${formatAmount(finalAmount)}` }}
            </button>

            <p class="text-center text-[10px] text-white/20">
              You'll be redirected to the payment gateway
            </p>
          </div>

          <div v-else class="py-8 text-center">
            <i aria-hidden="true" class="pi pi-check-circle text-4xl text-[#1db954]" />
            <p class="mt-3 text-lg font-bold text-white">Tip Sent!</p>
            <p class="mt-1 text-sm text-white/40">Thank you for supporting the artist.</p>
            <button
              class="mt-6 rounded-xl bg-white/5 px-6 py-2.5 text-sm font-medium text-white/60 hover:bg-white/10"
              @click="visible = false"
            >
              Close
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useTipsApi } from '@/services/api/tips'
import { useToast } from 'primevue/usetoast'

const props = defineProps<{
  artistId: string
  artistName?: string
  trackId?: string
}>()

const api = useTipsApi()
const toast = useToast()

const visible = ref(false)
const sending = ref(false)
const sent = ref(false)
const selectedAmount = ref(50000)
const customAmount = ref(0)
const message = ref('')
const amounts = [10000, 25000, 50000, 100000]

const finalAmount = computed(() =>
  customAmount.value > 0 ? customAmount.value : selectedAmount.value,
)

function formatAmount(cents: number): string {
  if (cents >= 1000) return `${(cents / 1000).toFixed(0)}K`
  return String(cents)
}

function openDialog() {
  selectedAmount.value = 50000
  customAmount.value = 0
  message.value = ''
  sent.value = false
  visible.value = true
}

async function sendTip() {
  if (finalAmount.value < 1000) return
  sending.value = true
  try {
    const res = await api.createTip({
      artist_id: props.artistId,
      track_id: props.trackId,
      amount_cents: finalAmount.value,
      message: message.value,
      callback_url: window.location.origin + window.location.pathname,
    })
    if (res?.redirect_url) {
      sent.value = true
      window.open(res.redirect_url, '_blank')
    }
  } catch (err) {
    console.error('Failed to send tip:', err)
    toast.add({ severity: 'error', summary: 'Tip failed', life: 3000 })
  } finally {
    sending.value = false
  }
}
</script>
