export interface Plan {
  id: string
  code: string
  name: string
  description?: string | null
  priceCents: number
  currency: string
  interval: string
  isPremium: boolean
  isActive: boolean
  available: boolean
  features: string[]
}

export interface Subscription {
  id: string
  userId: string
  planId: string
  plan?: Plan
  status: string
  currentPeriodStart: string
  currentPeriodEnd?: string | null
  cancelAtPeriodEnd: boolean
  canceledAt?: string | null
  createdAt: string
}

export interface Payment {
  id: string
  userId: string
  subscriptionId?: string | null
  planId?: string | null
  amountCents: number
  currency: string
  status: string
  provider?: string | null
  failureReason?: string | null
  paidAt?: string | null
  createdAt: string
}

export interface CheckoutRequest {
  planId: string
  callbackUrl?: string
  provider?: string
}

export interface CheckoutResponse {
  message: string
  available: boolean
  redirectUrl?: string
  authority?: string
  subscription?: Subscription
  payment?: Payment
}
