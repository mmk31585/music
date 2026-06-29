import z from 'zod'
import type { AxiosRequestConfig } from 'axios'

export interface RefreshToken {
  access_token?: string
}

export interface RequestHooks {
  /** UI toast helpers – can be real toast or a no‑op object */
  toast?: {
    success?: (msg: string) => void
    error?: (msg: string) => void
    info?: (msg: string) => void
    warning?: (msg: string) => void
  }
  /** Maintenance helpers */
  getMaintenanceSecrets?: () => string | null
  /** Auth helpers */
  getAuthToken?: () => string | null
  setAuthToken?: (t: string) => void
  clearAuthToken?: () => void
  /** Router / redirect helper */
  redirectToLogin?: (currentPath?: string) => Promise<void>
  /** Store reset helper */
  resetAuthStore?: () => void
  /** Refresh token */
  refreshToken?: () => Promise<RefreshToken> | null
  refreshTokenUrlRejecter?: (config: AxiosRequestConfig) => boolean
  /** Extras */
  extraHeaders?: () => Record<string, string>
}

export interface MetaProps {
  current_page: number
  from?: number | null
  last_page: number
  path: string
  per_page: number
  to?: number | null
  total: number
}

export interface PaginatedProps<T = unknown> {
  items: T[]
  meta?: MetaProps
  total?: number
  page?: number
  limit?: number
  page_size?: number
  total_count?: number
  count?: number
  offset?: number
  has_more?: boolean
}

/**
 * Generic API response shape.
 * `T` is the type of the `data` field that the caller expects.
 */
export interface ApiResponseProps<T = unknown> {
  type?: ResponseTypes
  data?: T | PaginatedProps<T>
  message?: string
  run_time?: string
}

/**
 * Configuration for the reusable request wrapper.
 * `T` is the type of the successful payload.
 */
export interface UseRequestConfig<T = unknown> {
  silent?: boolean
  allowEmptyArray?: true
  schema?: z.ZodTypeAny
  beforeRequest?: () => void
  /** Called on success – receives the parsed payload (`T`). */
  success?: (data: T, total: number, response: any) => boolean | void
  criticalError?: (msg: string) => void
  /** Called on a non‑critical error – receives the raw error payload (may be `unknown`). */
  error?: (data: any, msg: string) => boolean | void
  anyError?: (data: any | null, msg: string) => void
  finally?: () => void
}

export type ResponseTypes = 'success' | 'info' | 'warning' | 'error'

export enum ResponseStatuses {
  HTTP_CONTINUE = 100,
  HTTP_SWITCHING_PROTOCOLS = 101,
  HTTP_PROCESSING = 102,
  HTTP_EARLY_HINTS = 103,
  HTTP_OK = 200,
  HTTP_CREATED = 201,
  HTTP_ACCEPTED = 202,
  HTTP_NON_AUTHORITATIVE_INFORMATION = 203,
  HTTP_NO_CONTENT = 204,
  HTTP_RESET_CONTENT = 205,
  HTTP_PARTIAL_CONTENT = 206,
  HTTP_MULTI_STATUS = 207,
  HTTP_ALREADY_REPORTED = 208,
  HTTP_IM_USED = 226,
  HTTP_MULTIPLE_CHOICES = 300,
  HTTP_MOVED_PERMANENTLY = 301,
  HTTP_FOUND = 302,
  HTTP_SEE_OTHER = 303,
  HTTP_NOT_MODIFIED = 304,
  HTTP_USE_PROXY = 305,
  HTTP_RESERVED = 306,
  HTTP_TEMPORARY_REDIRECT = 307,
  HTTP_PERMANENTLY_REDIRECT = 308,
  HTTP_BAD_REQUEST = 400,
  HTTP_UNAUTHORIZED = 401,
  HTTP_PAYMENT_REQUIRED = 402,
  HTTP_FORBIDDEN = 403,
  HTTP_NOT_FOUND = 404,
  HTTP_METHOD_NOT_ALLOWED = 405,
  HTTP_NOT_ACCEPTABLE = 406,
  HTTP_PROXY_AUTHENTICATION_REQUIRED = 407,
  HTTP_REQUEST_TIMEOUT = 408,
  HTTP_CONFLICT = 409,
  HTTP_GONE = 410,
  HTTP_LENGTH_REQUIRED = 411,
  HTTP_PRECONDITION_FAILED = 412,
  HTTP_REQUEST_ENTITY_TOO_LARGE = 413,
  HTTP_REQUEST_URI_TOO_LONG = 414,
  HTTP_UNSUPPORTED_MEDIA_TYPE = 415,
  HTTP_REQUESTED_RANGE_NOT_SATISFIABLE = 416,
  HTTP_EXPECTATION_FAILED = 417,
  HTTP_I_AM_A_TEAPOT = 418,
  HTTP_MISDIRECTED_REQUEST = 421,
  HTTP_UNPROCESSABLE_ENTITY = 422,
  HTTP_LOCKED = 423,
  HTTP_FAILED_DEPENDENCY = 424,
  HTTP_TOO_EARLY = 425,
  HTTP_UPGRADE_REQUIRED = 426,
  HTTP_PRECONDITION_REQUIRED = 428,
  HTTP_TOO_MANY_REQUESTS = 429,
  HTTP_REQUEST_HEADER_FIELDS_TOO_LARGE = 431,
  HTTP_UNAVAILABLE_FOR_LEGAL_REASONS = 451,
  HTTP_INTERNAL_SERVER_ERROR = 500,
  HTTP_NOT_IMPLEMENTED = 501,
  HTTP_BAD_GATEWAY = 502,
  HTTP_SERVICE_UNAVAILABLE = 503,
  HTTP_GATEWAY_TIMEOUT = 504,
  HTTP_VERSION_NOT_SUPPORTED = 505,
  HTTP_VARIANT_ALSO_NEGOTIATES_EXPERIMENTAL = 506,
  HTTP_INSUFFICIENT_STORAGE = 507,
  HTTP_LOOP_DETECTED = 508,
  HTTP_NOT_EXTENDED = 510,
  HTTP_NETWORK_AUTHENTICATION_REQUIRED = 511,
}
