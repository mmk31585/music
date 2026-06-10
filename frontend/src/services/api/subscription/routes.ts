import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { SubscriptionApiRoutes } from './enums'
import type { Plan, Subscription, Payment, CheckoutRequest, CheckoutResponse } from './types'

export const useSubscriptionApi = () => {
  const listPlans = async (config?: UseRequestConfig<{ plans: Plan[] }>) => {
    return useRequest<{ plans: Plan[] }>(
      SubscriptionApiRoutes.PLANS,
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const currentSubscription = async (config?: UseRequestConfig<Subscription>) => {
    return useRequest<Subscription>(
      SubscriptionApiRoutes.ME,
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const listSubscriptions = async (
    config?: UseRequestConfig<{ subscriptions: Subscription[] }>,
  ) => {
    return useRequest<{ subscriptions: Subscription[] }>(
      SubscriptionApiRoutes.HISTORY,
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const checkout = async (
    payload: CheckoutRequest,
    config?: UseRequestConfig<CheckoutResponse>,
  ) => {
    return useRequest<CheckoutResponse>(
      SubscriptionApiRoutes.CHECKOUT,
      { method: 'POST', data: payload },
      { silent: false, ...config },
    )
  }

  const cancel = async (config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      SubscriptionApiRoutes.CANCEL,
      { method: 'POST' },
      { silent: false, ...config },
    )
  }

  const listPayments = async (config?: UseRequestConfig<{ payments: Payment[] }>) => {
    return useRequest<{ payments: Payment[] }>(
      SubscriptionApiRoutes.PAYMENTS,
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  return { listPlans, currentSubscription, listSubscriptions, checkout, cancel, listPayments }
}
