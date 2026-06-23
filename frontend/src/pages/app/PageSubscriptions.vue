<template>
  <div class="mx-auto w-full max-w-4xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <!-- Loading -->
    <div v-if="loading" class="space-y-4">
      <SkeletonLoader variant="lines" :lines="1" class="w-48" />
      <div class="grid grid-cols-1 gap-4 md:grid-cols-3">
        <SkeletonLoader v-for="i in 3" :key="i" variant="card" />
      </div>
    </div>

    <template v-else>
      <!-- Header -->
      <div class="mb-8">
        <p class="text-xs font-bold tracking-[0.25em] text-spotify uppercase">Subscription</p>
        <h1 class="mt-2 text-3xl font-black text-white">Plans &amp; Pricing</h1>
        <p class="mt-1 text-sm text-white/40">Choose the plan that fits your listening habits</p>
      </div>

      <!-- Current subscription banner -->
      <div
        v-if="currentSub && currentSub.plan"
        class="mb-8 overflow-hidden rounded-2xl border border-white/6 bg-linear-to-r from-spotify/5 to-transparent p-5"
      >
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div>
            <p class="text-xs text-white/40">Current Plan</p>
            <p class="text-lg font-bold text-white">{{ currentSub.plan.name }}</p>
            <p class="text-xs text-white/30 capitalize">
              Status: <span class="font-medium text-white/60">{{ currentSub.status }}</span>
            </p>
          </div>
          <div class="flex gap-2">
            <button
              v-if="currentSub.status === 'pending'"
              class="rounded-xl bg-amber-500/10 px-5 py-2 text-sm font-medium text-amber-400"
            >
              Awaiting payment
            </button>
            <button
              v-else-if="currentSub.status === 'active' && currentSub.plan.code !== 'free'"
              class="rounded-xl bg-red-500/10 px-5 py-2 text-sm font-medium text-red-400 transition hover:bg-red-500/20"
              @click="cancelSubscription"
            >
              Cancel
            </button>
          </div>
        </div>
        <div v-if="currentSub.currentPeriodEnd" class="mt-2 text-xs text-white/30">
          Renews {{ formatDate(currentSub.currentPeriodEnd) }}
        </div>
      </div>

      <!-- Plans grid -->
      <div class="grid grid-cols-1 gap-4 md:grid-cols-3">
        <div
          v-for="plan in plans"
          :key="plan.id"
          class="group relative overflow-hidden rounded-2xl border border-white/6 bg-white/3 p-6 transition hover:bg-white/5"
          :class="{ 'border-spotify/30 bg-spotify/3': isCurrentPlan(plan) }"
        >
          <!-- Premium badge -->
          <div
            v-if="plan.isPremium"
            class="absolute top-3 right-3 rounded-full bg-linear-to-r from-amber-400 to-spotify px-2.5 py-0.5 text-[10px] font-bold text-black"
          >
            Premium
          </div>

          <p class="text-sm font-bold text-white/40">{{ plan.name }}</p>
          <div class="mt-2 flex items-baseline gap-1">
            <span class="text-3xl font-black text-white">{{ formatPrice(plan) }}</span>
            <span class="text-xs text-white/30">/ {{ plan.interval }}</span>
          </div>
          <p v-if="plan.description" class="mt-2 text-xs text-white/40">{{ plan.description }}</p>

          <ul class="mt-4 space-y-2">
            <li
              v-for="feat in plan.features"
              :key="feat"
              class="flex items-center gap-2 text-xs text-white/60"
            >
              <i aria-hidden="true" class="pi pi-check text-[10px] text-spotify" /> {{ formatFeature(feat) }}
            </li>
          </ul>

          <button
            class="mt-6 w-full rounded-xl py-3 text-sm font-bold transition"
            :class="
              isCurrentPlan(plan)
                ? 'cursor-default bg-white/5 text-white/40'
                : plan.priceCents === 0
                  ? 'bg-white/10 text-white hover:bg-white/20'
                  : 'bg-spotify text-black hover:bg-spotify-hover'
            "
            :disabled="isCurrentPlan(plan)"
            @click="selectPlan(plan)"
          >
            {{
              isCurrentPlan(plan)
                ? 'Current Plan'
                : plan.priceCents === 0
                  ? 'Get Started'
                  : 'Subscribe'
            }}
          </button>
        </div>
      </div>

      <!-- Payment History -->
      <section v-if="payments.length > 0" class="mt-12">
        <h3 class="mb-4 text-lg font-bold text-white">Payment History</h3>
        <div class="space-y-2">
          <div
            v-for="p in payments"
            :key="p.id"
            class="flex items-center justify-between rounded-xl bg-white/3 px-4 py-3 transition hover:bg-white/5"
          >
            <div>
              <p class="text-sm font-medium text-white">${{ (p.amountCents / 100).toFixed(2) }}</p>
              <p class="text-xs text-slate-500">{{ formatDate(p.createdAt) }}</p>
            </div>
            <span
              class="rounded-full px-2.5 py-0.5 text-[10px] font-bold tracking-wider uppercase"
              :class="
                p.status === 'paid'
                  ? 'bg-green-500/10 text-green-400'
                  : p.status === 'pending'
                    ? 'bg-amber-500/10 text-amber-400'
                    : 'bg-red-500/10 text-red-400'
              "
            >
              {{ p.status }}
            </span>
          </div>
        </div>
      </section>
    </template>

    <!-- Checkout redirect dialog -->
    <Teleport to="body">
      <div
        v-if="showCheckout"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs"
        @click.self="showCheckout = false"
      >
        <div class="glass-strong mx-4 w-full max-w-md rounded-2xl p-8 text-center">
          <i aria-hidden="true" class="pi pi-external-link text-4xl text-spotify" />
          <h3 class="mt-4 text-xl font-bold text-white">Redirecting to Payment</h3>
          <p class="mt-2 text-sm text-white/40">
            You'll be redirected to the payment gateway to complete your subscription.
          </p>
          <div class="mt-6 flex gap-3">
            <button
              class="flex-1 rounded-xl bg-white/5 py-3 text-sm font-medium text-white/50 hover:bg-white/10"
              @click="showCheckout = false"
            >
              Cancel
            </button>
            <a
              :href="checkoutUrl"
              target="_blank"
              class="block flex-1 rounded-xl bg-spotify py-3 text-center text-sm font-bold text-black hover:bg-spotify-hover"
            >
              Proceed to Pay
            </a>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { SkeletonLoader } from '@/components/common'
import { useSubscriptionApi } from '@/services/api/subscription'
import { useToast } from 'primevue/usetoast'
import type { Plan, Subscription, Payment } from '@/services/api/subscription/types'

const api = useSubscriptionApi()
const toast = useToast()

const loading = ref(true)
const plans = ref<Plan[]>([])
const currentSub = ref<Subscription | null>(null)
const payments = ref<Payment[]>([])
const showCheckout = ref(false)
const checkoutUrl = ref('')

async function loadData() {
  loading.value = true
  try {
    const [plansRes, subRes, paymentsRes] = await Promise.all([
      api.listPlans(),
      api.currentSubscription(),
      api.listPayments(),
    ])
    plans.value = plansRes?.plans ?? []
    currentSub.value = subRes
    payments.value = paymentsRes?.payments ?? []
  } catch (err) {
    console.error('Failed to load subscription data:', err)
    toast.add({ severity: 'error', summary: 'Failed to load subscription data', life: 3000 })
  } finally {
    loading.value = false
  }
}

function isCurrentPlan(plan: Plan): boolean {
  return currentSub.value?.planId === plan.id
}

async function selectPlan(plan: Plan) {
  if (isCurrentPlan(plan)) return
  if (plan.priceCents === 0) {
    try {
      await api.checkout({ planId: plan.id, callbackUrl: window.location.origin + '/subscription' })
      toast.add({ severity: 'success', summary: `${plan.name} activated!`, life: 3000 })
      await loadData()
  } catch (err) {
    console.error('Failed to activate plan:', err)
    toast.add({ severity: 'error', summary: 'Failed to activate plan', life: 3000 })
    }
    return
  }

  try {
    const res = await api.checkout({
      planId: plan.id,
      callbackUrl: window.location.origin + '/subscription',
    })
    if (res?.redirectUrl) {
      checkoutUrl.value = res.redirectUrl
      showCheckout.value = true
    }
  } catch (err) {
    console.error('Checkout failed:', err)
    toast.add({ severity: 'error', summary: 'Checkout failed', life: 3000 })
  }
}

async function cancelSubscription() {
  try {
    await api.cancel()
    toast.add({ severity: 'info', summary: 'Subscription canceled', life: 3000 })
    await loadData()
  } catch (err) {
    console.error('Cancel failed:', err)
    toast.add({ severity: 'error', summary: 'Cancel failed', life: 3000 })
  }
}

function formatPrice(plan: Plan): string {
  if (plan.priceCents === 0) return 'Free'
  return '$' + (plan.priceCents / 100).toFixed(plan.currency === 'IRR' ? 0 : 2)
}

function formatDate(d: string): string {
  return new Date(d).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

function formatFeature(f: string): string {
  return f.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase())
}

onMounted(loadData)
</script>
