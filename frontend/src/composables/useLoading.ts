import { computed, ref, type Ref, unref } from 'vue'
import isObject from 'lodash.isobject'
import type { ZodType } from 'zod'
import { useRequest } from '@/composables'
import type { ApiResponseProps, MetaProps, PaginatedProps } from '@/plugins/client/types.ts'

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type LoadingOptionsType<T, TBody = unknown, TQuery = Record<string, any>> = {
  immediate?: boolean
  key?: string
  schema?: ZodType<T>
  allowEmptyArray?: boolean

  /** HTTP method: GET, POST, PUT, PATCH, DELETE, ... */
  method?: string

  /** Request body; can be plain value or a Ref */
  body?: TBody | Ref<TBody>

  /** Query params; can be plain value or a Ref */
  query?: TQuery | Ref<TQuery>

  parameters?: any[]
}

type LoadingResult<T> =
  T extends PaginatedProps<infer U>
    ? { data: Ref<T | null>; items: Ref<U[]>; meta: Ref<MetaProps> }
    : T extends Array<infer U>
      ? { data: Ref<T | null>; items: Ref<U[]>; meta: Ref<MetaProps> }
      : { data: Ref<T | null>; items: Ref<T[]>; meta: Ref<MetaProps> }

export function useLoading<
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  T extends object | any[] | PaginatedProps<any>,
  TBody = unknown,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  TQuery = Record<string, any>,
>(
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  fetcher: ((...parameters: any) => Promise<T>) | string,
  options?: LoadingOptionsType<T, TBody, TQuery>,
): LoadingResult<T> & {
  pending: Ref<boolean>
  error: Ref<Error | ApiResponseProps | null>
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  load: (...parameters: any) => Promise<void>
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  reload: (...parameters: any) => Promise<void>
} {
  const defaultMeta: MetaProps = {
    current_page: 1,
    last_page: 1,
    path: '',
    per_page: 0,
    total: 0,
  }

  const data: Ref<T | null> = ref(null)
  const pending: Ref<boolean> = ref(false)
  const error: Ref<Error | ApiResponseProps | null> = ref(null)

  type ItemsType =
    T extends PaginatedProps<infer U>
      ? U[]
      : T extends Array<infer U>
        ? U[]
        : T extends object
          ? T[]
          : any[]

  const items = computed(() => {
    const val = data.value
    if (val && isObject(val) && 'items' in val) return (val as PaginatedProps<ItemsType[0]>).items
    if (Array.isArray(val)) return val as any as ItemsType
    return (val ? [val] : []) as ItemsType
  })

  const meta = computed<MetaProps>(() => {
    const val = data.value
    if (val && isObject(val) && 'meta' in val) return (val as PaginatedProps<ItemsType[0]>).meta
    return defaultMeta
  })

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const load = async (...parameters: any) => {
    if (pending.value) return
    pending.value = true
    error.value = null

    try {
      const res: T | ApiResponseProps =
        typeof fetcher === 'string'
          ? await useRequest<T>(
              fetcher,
              {
                method: options?.method,
                data: options?.body ? unref(options.body as TBody) : undefined,
                params: options?.query ? unref(options.query as TQuery) : undefined,
              },
              {
                schema: options?.schema,
                allowEmptyArray: (options?.allowEmptyArray ?? true) as true,
              },
            )
          : await fetcher(...parameters)

      if (typeof res === 'object' && 'type' in res && res.type === 'error') {
        error.value = res as ApiResponseProps
      } else {
        data.value = res as T
      }
    } catch (err) {
      error.value =
        err instanceof Error
          ? err
          : isObject(err)
            ? (err as ApiResponseProps)
            : new Error(String(err))
    } finally {
      pending.value = false
    }
  }

  if (options?.immediate) void load(...(options?.parameters || []))

  return {
    data,
    items,
    meta,
    pending,
    error,
    load,
    reload: load,
  } as any as LoadingResult<T> & {
    pending: Ref<boolean>
    error: Ref<Error | ApiResponseProps | null>
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    load: (...parameters: any) => Promise<void>
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    reload: (...parameters: any) => Promise<void>
  }
}
