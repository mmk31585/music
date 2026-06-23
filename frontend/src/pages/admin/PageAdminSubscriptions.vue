<template>
  <div class="mx-auto w-full max-w-7xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
    title="Subscriptions"
    eyebrow="Billing"
    description="View and manage user subscriptions"
  >
    <template #actions>
      <Button
        label="Refresh"
        icon="pi pi-refresh"
        severity="info"
        text
        size="small"
        :loading="loading"
        @click="fetchData"
      />
    </template>
  </AdminSectionHeader>

  <div class="mt-6 grid gap-6 md:grid-cols-3 reveal-stagger">
    <div class="rounded-xl border border-white/6 bg-white/3 p-5">
      <p class="text-xs font-medium tracking-wider text-slate-500 uppercase">Active Plans</p>
      <p class="mt-2 text-3xl font-bold text-white tabular-nums">{{ plans.length }}</p>
    </div>
    <div class="rounded-xl border border-white/6 bg-white/3 p-5">
      <p class="text-xs font-medium tracking-wider text-slate-500 uppercase">Your Subscription</p>
      <p class="mt-2 text-3xl font-bold text-white tabular-nums">
        {{ currentSub ? 'Active' : 'None' }}
      </p>
    </div>
    <div class="rounded-xl border border-white/6 bg-white/3 p-5">
      <p class="text-xs font-medium tracking-wider text-slate-500 uppercase">Payment History</p>
      <p class="mt-2 text-3xl font-bold text-white tabular-nums">{{ payments.length }}</p>
    </div>
  </div>

  <div class="mt-8 reveal-fade">
    <h3 class="mb-4 text-sm font-bold text-white">Available Plans</h3>
    <div
      v-if="plans.length === 0"
      class="rounded-xl border border-white/6 bg-white/3 p-6 text-center text-sm text-slate-400"
    >
      No subscription plans available
    </div>
    <div v-else class="grid gap-4 md:grid-cols-3">
      <div
        v-for="plan in plans"
        :key="plan.id"
        class="rounded-xl border border-white/6 bg-white/3 p-5"
      >
        <h4 class="text-lg font-bold text-white">{{ plan.name }}</h4>
        <p class="mt-1 text-2xl font-black text-white">
          {{ formatCents(plan.priceCents) }}<span class="text-sm font-normal text-slate-400">/mo</span>
        </p>
        <p v-if="plan.description" class="mt-2 text-xs text-slate-400">{{ plan.description }}</p>
      </div>
    </div>
  </div>

  <div class="mt-8 reveal-fade">
    <h3 class="mb-4 text-sm font-bold text-white">My Subscription</h3>
    <div
      v-if="!currentSub"
      class="rounded-xl border border-white/6 bg-white/3 p-6 text-center text-sm text-slate-400"
    >
      No active subscription
    </div>
    <div v-else class="rounded-xl border border-white/6 bg-white/3 p-5">
      <p class="text-sm text-slate-300">
        Plan: <span class="font-medium text-white">{{ currentSub.plan?.name || 'Unknown' }}</span>
      </p>
      <p class="mt-1 text-sm text-slate-300">
        Status: <span class="font-medium text-green-400">{{ currentSub.status }}</span>
      </p>
      <p class="mt-1 text-sm text-slate-300">
        Expires:
        {{ (currentSub as any).expires_at ? new Date((currentSub as any).expires_at).toLocaleDateString() : 'N/A' }}
      </p>
    </div>
  </div>

  <div class="mt-8 reveal-fade">
    <h3 class="mb-4 text-sm font-bold text-white">Payment History</h3>
    <div
      v-if="payments.length === 0"
      class="rounded-xl border border-white/6 bg-white/3 p-6 text-center text-sm text-slate-400"
    >
      No payment history
    </div>
    <div v-else class="space-y-2">
      <div
        v-for="p in payments"
        :key="p.id"
        class="flex items-center justify-between rounded-xl bg-white/3 px-4 py-3"
      >
        <div>
          <p class="text-sm font-medium text-white">{{ formatCents(p.amountCents) }} {{ p.currency }}</p>
          <p class="text-xs text-slate-500">
            {{ p.status }} · {{ new Date((p as any).created_at).toLocaleDateString() }}
          </p>
        </div>
        <span class="rounded-full bg-green-500/10 px-2 py-0.5 text-xs font-medium text-green-400">{{
          p.status
        }}</span>
      </div>
    </div>
  </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useSubscriptionApi } from '@/services/api/subscription'
import type { Plan, Subscription, Payment } from '@/services/api/subscription'
import { AdminSectionHeader } from '@/components/admin'

const subApi = useSubscriptionApi()

const plans = ref<Plan[]>([])
const currentSub = ref<Subscription | null>(null)
const payments = ref<Payment[]>([])
const loading = ref(true)

function formatCents(cents: number | undefined | null): string {
  if (cents == null) return '$0.00'
  return `$${(cents / 100).toFixed(2)}`
}

async function fetchData() {
  loading.value = true
  try {
    const [plansResult, subResult, paymentsResult] = await Promise.all([
      subApi.listPlans(),
      subApi.currentSubscription(),
      subApi.listPayments(),
    ])
    plans.value = plansResult?.plans ?? []
    currentSub.value = subResult ?? null
    payments.value = paymentsResult?.payments ?? []
  // TODO LOW: Silent catch — should show a toast on failure instead of swallowing errors.
  } catch {
    /* silent */
  } finally {
    loading.value = false
  }
}

onMounted(fetchData)
</script>
