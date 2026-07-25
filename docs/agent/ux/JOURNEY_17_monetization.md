# Journey 17: Monetization — Free vs Premium, Subscriptions, Tipping

> Full trace: tier gating → subscription management → creator tipping → payment flows.

---

## Tier Architecture

### Two-Tier Model

| Feature | Free | Premium |
|---------|------|---------|
| Track playback | ✅ Unlimited | ✅ Unlimited |
| Audio quality | 128kbps | 320kbps (or FLAC) |
| Ad-free | ❌ Ads every 3-5 tracks | ✅ |
| Offline downloads | ❌ | ✅ (up to 10,000 tracks) |
| Background play (mobile) | ❌ | ✅ |
| Skip limit | 6 skips/hour | Unlimited |
| Repeat one track | ❌ | ✅ |
| Queue management | Basic (20 tracks) | Full (500 tracks) |
| Radio mode | Standard | Unlimited skips |
| Sound quality settings | ❌ | ✅ (auto, high, lossless) |
| Crossfade | ❌ | ✅ |
| Lyrics view | Limited | Full sync |
| Themes | Default only | All themes |

### Tier Determination

```typescript
// stores/user.ts (or subscriptionStore.ts)
const tier = computed(() => {
  if (user.value?.subscription?.status === 'active') return 'premium'
  if (user.value?.subscription?.status === 'trial') return 'premium_trial'
  return 'free'
})

// Guard: feature gate composable
function useFeature(feature: PremiumFeature): boolean {
  return tier.value === 'premium' || tier.value === 'premium_trial'
}
```

---

## Gating UX Pattern

### Feature Gate Component

```typescript
// FeatureGate.vue — wraps premium-only content
// Props: feature: PremiumFeature, fallback?: 'upgrade' | 'hide' | 'disabled'

// Upgrade mode: Show full content with overlay + "Upgrade to Premium" CTA
// Hide mode: Don't render content at all
// Disabled mode: Show content but disabled/locked state
```

### Upgrade Prompts (Touchpoints)

| Location | Trigger | Gate Type | Content |
|----------|---------|-----------|---------|
| NowPlayingBar quality badge | Free user clicks "128kbps" | Upgrade overlay | "Go Premium for 320kbps" |
| Settings > Audio Quality | Free user taps quality selector | Disabled with lock icon | "Available in Premium" |
| Download button (album/track) | Free user taps download | Upgrade modal | "Download with Premium" |
| Background play (mobile) | Free user locks screen | Interrupt + upgrade toast | "Premium required to play in background" |
| Skip limit hit | 6th skip within hour | Toast + upgrade CTA | "Out of skips — go Premium for unlimited" |
| Visualizer > advanced modes | Free user taps advanced theme | Upgrade overlay | "Unlock all themes with Premium" |
| Queue > 20 tracks | Free user adds 21st track | Toast | "Queue limit reached — Premium holds 500" |
| Player settings > Crossfade | Free user toggles on | Disabled + upgrade link | "Crossfade is a Premium feature" |

### Upgrade Modal (`UpgradeModal.vue`)

```
Triggered from any premium gate point.

Modal content:
  ├── Header: "Go Premium" with logo
  ├── Hero: Feature comparison (with emphasis on the feature user just hit)
  ├── Pricing section:
  │     ├── Monthly: $X.99/month
  │     ├── Yearly: $Y.99/year (save Z%)
  │     └── Lifetime (if available)
  ├── CTA: "Start Free Trial" (if available) or "Subscribe Now"
  ├── Testimonials or stats: "Join 1M+ Premium listeners"
  └── Footer: "No commitment — cancel anytime"
```

---

## Subscription Management

### Subscription Store

```typescript
// stores/subscriptionStore.ts
const subscription = ref<Subscription | null>(null)
const plans = ref<Plan[]>([])
const isLoading = ref(false)

async function fetchSubscription() {
  const { data } = await api.get('/api/v1/subscription')
  subscription.value = data
}

async function fetchPlans() {
  const { data } = await api.get('/api/v1/subscription/plans')
  plans.value = data
}

async function subscribe(planId: string, paymentMethodId: string) {
  const { data } = await api.post('/api/v1/subscription', { plan_id: planId, payment_method_id: paymentMethodId })
  // Redirect to checkout URL if needed, or update subscription in-place
}

async function cancelSubscription() {
  await api.post('/api/v1/subscription/cancel')
  subscription.value.status = 'canceling' // end of billing period
}

async function reactivateSubscription() {
  await api.post('/api/v1/subscription/reactivate')
  subscription.value.status = 'active'
}
```

### Subscription Page (`PageSubscription.vue` → `/settings/subscription`)

| Section | Content |
|---------|---------|
| **Current Plan** | Tier name, status (active/trial/canceling/expired), next billing date |
| **Plan Details** | Quality, features, limits for current tier |
| **Available Plans** | Monthly, Yearly, Lifetime cards with comparison |
| **Payment Methods** | Saved cards, add new card form |
| **Billing History** | Past invoices (last 12 months) |
| **Cancel/Reactivate** | Cancel button (if active), Reactivate button (if canceling) |

### Payment Flow

```
1. User clicks "Subscribe" / "Upgrade"
2. Payment method selection:
   a. Saved card (if exists)
   b. Add new card (Stripe Elements / iframe)
3. Review order: plan, price, tax
4. Confirm:
   a. Successful → subscription becomes active, redirect to confirmation page
   b. Failed (card declined) → error inline with specific message
   c. Failed (3D Secure) → redirect for authentication, then back
5. Confirmation: "Welcome to Premium!" with celebration animation
```

### Cancellation Flow

```
1. User clicks "Cancel Subscription"
2. Confirmation dialog:
   "Are you sure? You'll lose access to Premium features at the end of your billing period."
   Options: "Keep Premium" (primary), "Cancel Anyway" (danger)
3. If "Cancel Anyway" → show retention offer (if eligible):
   "Wait! Here's 1 month free to stay."
4. If declined → subscription marked 'canceling'
5. Grace period: remains active until end of billing period
6. Expiry: downgrade to free, library/favorites preserved, playlists preserved
```

### Subscription States Visual

| State | Badge | Period Remaining | Actions |
|-------|-------|-----------------|---------|
| Active | Green "Premium" | N/A | Cancel, Change plan |
| Trial | Blue "Trial" | "14 days remaining" | Subscribe, Cancel |
| Canceling | Yellow "Ending" | "Until June 30" | Reactivate |
| Expired | Gray "Expired" | "Expired on June 30" | Resubscribe |
| Grace | Orange "Past due" | "Payment failed — update method" | Update payment |
| Free | — | — | Subscribe |

---

## Tipping (Creator Economy)

### Tip Flow

```typescript
// POST /api/v1/creators/:id/tip { amount, currency, message?, is_anonymous }
// GET /api/v1/creators/:id/tips/top → top tippers for the month

interface TipPayload {
  creator_id: string
  track_id?: string  // optional — tip for a specific track
  amount: number     // in smallest currency unit (cents)
  currency: string   // ISO 4217
  message?: string   // optional message to creator
  is_anonymous: boolean
}
```

### Tipping Entry Points

| Location | Trigger | UI |
|----------|---------|-----|
| Artist page | "Support" button | Tip amount selector modal |
| Creator profile | "Tip" button | Tip modal |
| Stage room | "🎁 Tip" button | Quick tip (preset amounts) |
| FullscreenPlayer | Artist name dropdown → "Tip" | Small tip modal |
| Post-queue (track end) | "Support the artist" toast | Tip suggestion |

### Tip Modal (`TipModal.vue`)

```
├── Artist avatar + name
├── Preset amounts: $1, $3, $5, $10, custom
├── Message field (optional, 200 char max)
├── Anonymous toggle
├── Total with platform fee shown
└── CTA: "Send Tip" → payment confirmation
```

### Top Tippers

```
Creator profile sidebar:
  → "Top Tippers this month"
  → Shows top 3 (if not anonymous)
  → Crown icon on #1
  → "See all" link to full list
```

---

## Promotional Codes

```typescript
// POST /api/v1/subscription/apply-promo { code: string }
// Returns: { discount_percent, discount_amount, new_price, valid }
```

| CTA | Entry Point |
|-----|-------------|
| Subscription page | "Have a promo code?" → inline input |
| Upgrade modal | Same |
| Checkout | Promo code field |

---

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-1701 | ⚠️ MAJOR | `FeatureGate.vue` | Premium gates use **component-level checks** scattered across 15+ files — no centralized feature registry | Hard to maintain, hard to audit, easy to miss gates | Create centralized `FeatureRegistry` with all premium features listed |
| F-1702 | ⚠️ MAJOR | `PageSubscription.vue` | **No "Restore Purchases"** for Apple/Google in-app purchases | Mobile web users who subscribed via app can't restore on web | Add restore purchases flow |
| F-1703 | ⚠️ MAJOR | `PageSubscription.vue` | **Cancel flow has no retention offers** — user cancels and that's it | No chance to retain user | Add "1 month free" or "discount" retention offer |
| F-1704 | 💡 IMPROVE | `TipModal.vue` | Tips are **currency-agnostic** — always USD, no conversion | Persian users want to tip in Tomans or Rials | Add local currency support |
| F-1705 | 💡 IMPROVE | `PageSubscription.vue` | No **family plan** tier | Competing platforms offer family sharing | Add family plan (up to 6 accounts) |
| F-1706 | 💡 IMPROVE | `FeatureGate.vue` | **Gate type inconsistency** — some use overlay, some use disabled, some use toast | Inconsistent UX across premium features | Standardize gate types per feature category |
| F-1707 | 💡 IMPROVE | Subscription | **No trial-to-paid reminder** | Users who sign up for trial and forget lose access with no notice | Send email + push notification 3 days before trial ends |
| F-1708 | 💡 IMPROVE | `TipModal.vue` | **No tip receipt/invoice** | Users who tip for business can't expense it | Email receipt on successful tip |
| F-1709 | 💡 IMPROVE | `TipModal.vue` | **No suggested tip amount** based on listening time | User unsure how much to tip | Show "You've listened to X hours of this artist's music" |
| F-1710 | 💡 IMPROVE | `PageSubscription.vue` | **No subscription gifting** — "Gift Premium to a friend" | Users can't gift subscriptions | Add gift subscription flow |

## RTL / A11y / Mobile Notes

- ✅ Upgrade modal uses RTL layout properly
- ❌ Tip modal has **no `aria-describedby`** linking amount to description → screen reader just says "3" not "$3.00"
- ✅ Pricing tables use proper `aria-label` for currency symbols
- ❌ Cancel confirmation dialog has **no focus trap** — Tab can escape the dialog
- ✅ Subscription status badges have accessible color indicators (text + icon, not just color)
