import { computed, ref, type Ref, unref } from 'vue'
import isObject from 'lodash.isobject'
import type { ZodType } from 'zod'
import { useRequest } from '@/composables'
import type { AxiosRequestConfig } from 'axios'
import type { ApiResponseProps, MetaProps, PaginatedProps } from '@/plugins/client/types'

export type LoadingOptionsType<T, TBody = unknown, TQuery = Record<string, unknown>> = {
  immediate?: boolean
  key?: string
  schema?: ZodType<T>
  allowEmptyArray?: boolean
  method?: string
  body?: TBody | Ref<TBody>
  query?: TQuery | Ref<TQuery>
  parameters?: unknown[]
}

type PaginatedItems<T> = T extends PaginatedProps<infer U> ? U[] : never
type ArrayItems<T> = T extends Array<infer U> ? U[] : never

type LoadingResult<T> = {
  data: Ref<T | null>
  items: Ref<T extends PaginatedProps<infer U> ? U[] : T extends Array<infer U> ? U[] : T[]>
  meta: Ref<MetaProps>
  pending: Ref<boolean>
  error: Ref<Error | ApiResponseProps | null>
  load: (...args: unknown[]) => Promise<void>
  reload: (...args: unknown[]) => Promise<void>
}

/**
 * A composable that wraps data fetching with loading/error/data state management.
 *
 * Accepts either a URL string (uses `useRequest`) or a fetcher function.
 */
export function useLoading<
  TData,
  TBody = unknown,
  TQuery = Record<string, unknown>,
>(
  fetcher: ((...args: unknown[]) => Promise<TData>) | string,
  options?: LoadingOptionsType<TData, TBody, TQuery>,
): LoadingResult<TData> {
  const defaultMeta: MetaProps = {
    current_page: 1,
    last_page: 1,
    path: '',
    per_page: 0,
    total: 0,
  }

  const data = ref<TData | null>(null) as Ref<TData | null>
  const pending = ref(false)
  const error = ref<Error | ApiResponseProps | null>(null)

  const items = computed(() => {
    const val = data.value
    if (val && typeof val === 'object' && 'items' in val) {
      return (val as PaginatedProps<unknown>).items as TData extends PaginatedProps<infer U> ? U[] : TData[]
    }
    if (Array.isArray(val)) return val as unknown as TData[]
    return val ? [val] : []
  })

  const meta = computed<MetaProps>(() => {
    const val = data.value
    if (val && typeof val === 'object' && 'meta' in val) {
      return (val as Record<string, unknown>).meta as MetaProps
    }
    return defaultMeta
  })

  const load = async (...args: unknown[]) => {
    if (pending.value) return
    pending.value = true
    error.value = null

    try {
      let res: TData | ApiResponseProps

      if (typeof fetcher === 'string') {
        res = await useRequest<TData>(
          fetcher,
          {
            method: options?.method,
            data: options?.body ? unref(options.body as TBody) : undefined,
            params: options?.query ? unref(options.query as TQuery) : undefined,
          } as AxiosRequestConfig,
          {
            schema: options?.schema,
            allowEmptyArray: options?.allowEmptyArray ?? true,
          },
        )
      } else {
        res = await fetcher(...args)
      }

      if (res && typeof res === 'object' && 'type' in res && (res as ApiResponseProps).type === 'error') {
        error.value = res as ApiResponseProps
      } else {
        data.value = res as TData
      }
    } catch (err: unknown) {
      if (err instanceof Error) {
        error.value = err
      } else if (isObject(err)) {
        error.value = err as ApiResponseProps
      } else {
        error.value = new Error(String(err))
      }
    } finally {
      pending.value = false
    }
  }

  if (options?.immediate) void load(...(options?.parameters ?? []))

  return {
    data,
    items,
    meta,
    pending,
    error,
    load,
    reload: load,
  }
}
