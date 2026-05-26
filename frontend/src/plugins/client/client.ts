import axios, { type AxiosInstance } from 'axios'

/* -------------------------------------------------------------------------- */
/*  Axios Instance                                                            */
/* -------------------------------------------------------------------------- */
const BASE_URL = import.meta.env.VITE_API_BASE_URL
const REQUEST_TIMEOUT = Number(import.meta.env.VITE_API_TIMEOUT_MS ?? 1500)

const axiosClient: AxiosInstance = axios.create({
  baseURL: BASE_URL,
  timeout: REQUEST_TIMEOUT,
  withCredentials: true,
})

axiosClient.defaults.headers.common['X-Requested-With'] = 'XMLHttpRequest'
axiosClient.defaults.headers.common['Content-Type'] = 'application/json'

export default axiosClient
