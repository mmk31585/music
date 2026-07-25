export interface Tip {
  id: string
  sender_id: string
  artist_id: string
  track_id?: string | null
  amount_cents: number
  currency: string
  message?: string | null
  status: string
  provider?: string | null
  provider_pay_id?: string | null
  paid_at?: string | null
  created_at: string
}

export interface CreateTipRequest {
  artist_id: string
  track_id?: string
  amount_cents: number
  currency?: string
  message?: string
  callback_url: string
}
