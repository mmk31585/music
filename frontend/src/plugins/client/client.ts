import axios, { type AxiosInstance } from 'axios'

const BASE_URL = import.meta.env.VITE_API_BASE_URL
const REQUEST_TIMEOUT = Number(import.meta.env.VITE_API_TIMEOUT_MS ?? 1500)

if (import.meta.env.PROD && !BASE_URL) {
  console.error(
    '[Muse] VITE_API_BASE_URL is not set. All API calls will go to the Vercel origin ' +
    'and fail with 405. Set it in Vercel Dashboard → Environment Variables.',
  )
}

const axiosClient: AxiosInstance = axios.create({
  baseURL: BASE_URL,
  timeout: REQUEST_TIMEOUT,
  withCredentials: true,
})

axiosClient.defaults.headers.common['Content-Type'] = 'application/json'

export default axiosClient
