import {
  type ApiResponseProps,
  type PaginatedProps,
  type RefreshToken,
  type RequestHooks,
  ResponseStatuses,
  type ResponseTypes,
  type UseRequestConfig,
} from './types'
import { type AxiosError, type AxiosInstance, type AxiosRequestConfig } from 'axios'
import type { ZodError, ZodSafeParseResult } from 'zod'
import z from 'zod'

/**
 * Build a request wrapper that uses the supplied Axios instance
 * and the supplied hooks for all side‑effects.
 */
export function createRequestWrapper(client: AxiosInstance, hooks: RequestHooks = {}) {
  // -----------------------------------------------------------------
  //  Default (no‑op) implementations – guarantees the factory works
  // -----------------------------------------------------------------
  const {
    toast = {
      success: () => {},
      error: () => {},
      info: () => {},
      warning: () => {},
    },
    getMaintenanceSecrets = () => null,
    getAuthToken = () => null,
    setAuthToken = () => {},
    clearAuthToken = () => {},
    redirectToLogin = async () => {},
    resetAuthStore = () => {},
    refreshToken = async (): Promise<RefreshToken> => ({}),
    refreshTokenUrlRejecter = () => false,
    extraHeaders = (): Record<string, string> => ({}),
  } = hooks

  // --------------------------------------------------------------
  //  Interceptors (using the injected hooks)
  // --------------------------------------------------------------
  client.interceptors.request.use((config) => {
    config.headers = config.headers ?? {}
    const isFormData =
      typeof FormData !== 'undefined' && config.data != null && config.data instanceof FormData

    if (isFormData && config.headers) {
      // Let browser set multipart boundary automatically.
      delete (config.headers as Record<string, unknown>)['Content-Type']
    }

    // Attach maintenance secret if present
    const secretCode = getMaintenanceSecrets()

    if (secretCode) {
      config.params = { ...(config.params || {}), maintenance_secret: secretCode }
      if (config.method && config.method.toLowerCase() !== 'get') {
        if (isFormData) {
          ;(config.data as FormData).append('maintenance_secret', secretCode)
        } else {
          config.data = { ...(config.data || {}), maintenance_secret: secretCode }
        }
      }
    }

    // Extra headers
    const extra = extraHeaders()
    for (const key in extra) {
      if (Object.prototype.hasOwnProperty.call(extra, key)) {
        config.headers[key] = extra[key]
      }
    }

    // Attach auth token
    const token = getAuthToken()

    if (token) config.headers.Authorization = `Bearer ${token}`

    return config
  })

  /* -------------------------------------------------------------------------- */
  /*  Response Interceptor – token refresh & forbidden handling                 */
  /* -------------------------------------------------------------------------- */
  let isRefreshing = false
  type QueueItem = {
    resolve: (val?: unknown) => void
    reject: (err: unknown) => void
    config: AxiosRequestConfig
  }
  let failedQueue: QueueItem[] = []

  function processQueue(error: unknown, token: string | null = null) {
    failedQueue.forEach((p) => {
      if (error) {
        p.reject(error)
      } else {
        if (token && p.config.headers) {
          p.config.headers.Authorization = `Bearer ${token}`
        }
        p.resolve(client(p.config))
      }
    })
    failedQueue = []
  }

  client.interceptors.response.use(
    (response) => response,
    async (error: AxiosError) => {
      const originalConfig = error.config as AxiosRequestConfig & { _retry?: boolean }

      if (!error.response) return Promise.reject(error)

      const status = error.response.status

      // 403 – clear auth and redirect to login
      if (status === ResponseStatuses.HTTP_FORBIDDEN) {
        resetAuthStore()
        clearAuthToken()
        await redirectToLogin()

        return Promise.reject(error)
      }

      // 401 – try token refresh once
      if (
        status === ResponseStatuses.HTTP_UNAUTHORIZED &&
        !originalConfig._retry &&
        !refreshTokenUrlRejecter(originalConfig)
      ) {
        if (isRefreshing) {
          return new Promise((resolve, reject) => {
            const retryConfig = { ...originalConfig }
            failedQueue.push({ resolve, reject, config: retryConfig })
          })
        }

        originalConfig._retry = true
        isRefreshing = true

        try {
          const refreshRes = (await refreshToken()) as RefreshToken

          if (refreshRes?.access_token) {
            setAuthToken(refreshRes.access_token)
            processQueue(null, refreshRes.access_token)

            if (originalConfig.headers) {
              originalConfig.headers.Authorization = `Bearer ${refreshRes.access_token}`
            }

            return client(originalConfig)
          }

          throw new Error('No refresh token response')
        } catch (refreshError) {
          processQueue(refreshError, null)

          resetAuthStore()
          clearAuthToken()
          await redirectToLogin()

          return Promise.reject(refreshError)
        } finally {
          isRefreshing = false
        }
      }

      return Promise.reject(error)
    },
  )

  // --------------------------------------------------------------
  //  The reusable `useRequest` implementation
  // --------------------------------------------------------------
  async function useRequest<T = unknown, IsArray extends boolean = false>(
    url: string,
    config?: AxiosRequestConfig | null,
    resultConfig?: UseRequestConfig<IsArray extends true ? T[] : T>,
  ): Promise<IsArray extends true ? T[] : T> {
    type DataType = IsArray extends true ? T[] : T

    function isPaginatedResponse(value: unknown): value is PaginatedProps<T> {
      return (value &&
        typeof value === 'object' &&
        Array.isArray((value as PaginatedProps<T>).items) &&
        !!(value as PaginatedProps<T>).meta) as boolean
    }

    function invalidZodSchema(
      onAnyError: ((data: unknown | null, msg: string) => void) | undefined,
      onError: ((data: unknown, msg: string) => boolean | void) | undefined,
      silent: boolean,
      validation: ZodSafeParseResult<unknown>,
    ): string {
      // Turn a validation failure into a “non‑critical” error
      const zErr = validation.error as ZodError
      const msg = 'خطا در تبدیل داده دریافتی به نوع داده مورد نظر (خطای schema): ' + zErr.message
      if (typeof onAnyError === 'function') onAnyError(null, msg)
      if (typeof onError === 'function') onError(null, msg)
      if (!silent) {
        if (toast?.['error']) toast?.['error'](msg)
        else if (window) alert(msg)
        else console.error(msg)
      }

      if (import.meta.env.DEV) console.error(msg)

      return msg
    }

    return new Promise<DataType>((resolve, reject: (reason: unknown) => void) => {
      // ----- callbacks -------------------------------------------------
      const silent = resultConfig?.silent === true
      const allowEmptyArray =
        resultConfig && 'allowEmptyArray' in resultConfig ? resultConfig.allowEmptyArray : true
      const onBeforeRequest = resultConfig?.beforeRequest
      const onSuccess = resultConfig?.success
      const onCriticalError = resultConfig?.criticalError
      const onError = resultConfig?.error
      const onAnyError = resultConfig?.anyError
      const onFinally = resultConfig?.finally

      if (typeof onBeforeRequest === 'function') onBeforeRequest()

      // ----- request preparation ---------------------------------------
      const requestConfig = (config ?? {}) as AxiosRequestConfig
      requestConfig.method = requestConfig.method ?? 'GET'

      // ----- actual request --------------------------------------------
      client(url, requestConfig)
        .then((response) => {
          // The server may wrap the payload in a `data` field or return it directly.
          const raw = response.data ?? {}
          let dataValue = raw?.data ?? raw

          // Merge pagination metadata when the backend sends data as a plain
          // array with pagination info in a top-level `meta` field.
          if (Array.isArray(dataValue) && raw?.meta && typeof raw.meta === 'object') {
            dataValue = { items: dataValue, ...raw.meta }
          }

          const payload: ApiResponseProps<T | T[] | PaginatedProps<T>> = {
            type: raw?.type as ResponseTypes,
            data: dataValue as T | T[] | PaginatedProps<T>,
            message: raw?.message as string,
            run_time: raw?.run_time as string,
          }

          // ----- zod validation -----------------------------------------
          if (resultConfig?.schema) {
            const schema = resultConfig.schema as z.ZodTypeAny
            const data = payload.data

            // detect if response is paginated
            const paginated = isPaginatedResponse(data)

            if (paginated) {
              const effectiveSchema = z.array(schema).min(allowEmptyArray ? 0 : 1)
              const validation = effectiveSchema.safeParse(data.items)

              if (!validation.success) {
                const msg = invalidZodSchema(onAnyError, onError, silent, validation)
                reject({ message: msg })
                return
              }

              payload.data = {
                items: validation.data,
                meta: data.meta,
              } as DataType
            } else {
              const isArray = Array.isArray(data)
              const effectiveSchema = isArray
                ? z.array(schema).min(allowEmptyArray ? 0 : 1)
                : schema
              const validation = effectiveSchema.safeParse(payload.data)

              if (!validation.success) {
                const msg = invalidZodSchema(onAnyError, onError, silent, validation)
                reject({ message: msg })
                return
              }

              // If it succeeded, replace `payload.data` with the parsed value
              payload.data = validation.data as DataType
            }
          }

          // ----- total count logic --------------------------------------
          let total = 0
          if (raw?.meta?.total) {
            total = raw.meta.total
          } else if (Array.isArray(payload.data)) {
            total = payload.data.length
          } else if (payload.data && typeof payload.data === 'object') {
            total = 1
          }

          // ----- success callback ---------------------------------------
          let continueDefault = true
          if (typeof onSuccess === 'function') {
            const result = onSuccess(payload.data as DataType, total, response)
            if (result === false) continueDefault = false
          }

          // ----- toast handling (unless silent or prevented) ------------
          if (
            !silent &&
            continueDefault &&
            payload.message &&
            response.status !== ResponseStatuses.HTTP_NO_CONTENT
          ) {
            if (payload.type) {
              if (toast?.[payload.type]) toast?.[payload.type]!(payload.message)
              else if (window) alert(payload.message)
              else console.info(payload.message)
            }
          }

          resolve(payload.data as DataType)
        })
        .catch((error) => {
          const errResp = error?.response
          const errRaw = errResp?.data ?? {}
          const errPayload: ApiResponseProps = {
            type: errRaw?.type,
            data: errRaw?.data ?? errRaw,
            message: errRaw?.message ?? errResp?.statusText ?? error?.message ?? 'خطای ناشناخته',
            run_time: errRaw?.run_time,
          }

          // ----- critical‑error detection --------------------------------
          const isCritical =
            (typeof errPayload.data === 'object' &&
              Object.keys(errPayload.data || {}).length === 0) ||
            (Array.isArray(errPayload.data) && errPayload.data.length === 0) ||
            errResp?.status >= ResponseStatuses.HTTP_INTERNAL_SERVER_ERROR ||
            error?.request?.status >= ResponseStatuses.HTTP_INTERNAL_SERVER_ERROR ||
            Number(error?.request?.status) === ResponseStatuses.HTTP_METHOD_NOT_ALLOWED

          if (isCritical) {
            if (errPayload.message?.toLowerCase() !== 'canceled') {
              if (import.meta.env.DEV) console.error(error)

              if (typeof onAnyError === 'function') onAnyError(null, errPayload.message as string)
              if (typeof onCriticalError === 'function')
                onCriticalError(errPayload.message as string)

              const defaultMsg = 'خطا در ارتباط با سرور و دریافت اطلاعات!'
              if (!silent) {
                if (errPayload.type && errPayload.message) {
                  if (toast?.[errPayload.type])
                    toast?.[errPayload.type]!(errPayload.message || defaultMsg)
                  else if (window) alert(errPayload.message || defaultMsg)
                  else console.error(errPayload.message || defaultMsg)
                } else {
                  if (toast?.['error']) toast?.['error'](defaultMsg)
                  else if (window) alert(defaultMsg)
                  else console.error(defaultMsg)
                }
              }
            }
            reject(errPayload)
            return
          }

          // ----- non‑critical error handling -----------------------------
          if (typeof onAnyError === 'function')
            onAnyError(errPayload.data, errPayload.message as string)

          let continueDefault = true
          if (typeof onError === 'function') {
            const result = onError(errPayload.data, errPayload.message as string)
            if (result === false) continueDefault = false
          }

          if (!silent && continueDefault && errPayload.message) {
            if (errPayload.type) {
              if (toast[errPayload.type]) toast[errPayload.type]!(errPayload.message)
              else if (window) alert(errPayload.message)
              else console.error(errPayload.message)
            } else {
              if (toast?.['error']) toast?.['error'](errPayload.message)
              else if (window) alert(errPayload.message)
              else console.error(errPayload.message)
            }
          }

          reject(errPayload)
        })
        .finally(() => {
          if (typeof onFinally === 'function') onFinally()
        })
    })
  }

  // Export what the rest of the app needs
  return { useRequest, client }
}
