import axios from 'axios'

// admin-api 数据接口
const ADMIN_API_BASE_URL = import.meta.env.VITE_ADMIN_API_URL || 'http://localhost:8088/api/v1/admin'

// buyer-api 认证接口（登录/刷新/登出）
const BUYER_API_BASE_URL = import.meta.env.VITE_BUYER_API_URL || 'http://localhost:18080/api/v1/buyers'

const adminClient = axios.create({
  baseURL: ADMIN_API_BASE_URL,
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
})

const authClient = axios.create({
  baseURL: BUYER_API_BASE_URL,
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器：注入 Access Token
function addAuthInterceptor(client: ReturnType<typeof axios.create>) {
  client.interceptors.request.use(
    (config) => {
      const token = localStorage.getItem('access_token')
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
      return config
    },
    (error) => Promise.reject(error),
  )
}

// 响应拦截器：统一错误处理与 Token 刷新
function addRefreshInterceptor(client: ReturnType<typeof axios.create>) {
  client.interceptors.response.use(
    (response) => response,
    async (error) => {
      const originalRequest = error.config

      if (error.response?.status === 401 && !originalRequest._retry) {
        originalRequest._retry = true

        const refreshToken = localStorage.getItem('refresh_token')
        if (refreshToken) {
          try {
            const { data } = await authClient.post('/token/refresh', {
              refresh_token: refreshToken,
            })

            const { access_token, refresh_token: newRefreshToken } = data.data
            localStorage.setItem('access_token', access_token)
            localStorage.setItem('refresh_token', newRefreshToken)

            originalRequest.headers.Authorization = `Bearer ${access_token}`
            return client(originalRequest)
          } catch {
            localStorage.removeItem('access_token')
            localStorage.removeItem('refresh_token')
            window.location.href = '/login'
          }
        } else {
          window.location.href = '/login'
        }
      }

      return Promise.reject(error)
    },
  )
}

addAuthInterceptor(adminClient)
addAuthInterceptor(authClient)
addRefreshInterceptor(adminClient)
addRefreshInterceptor(authClient)

export { adminClient, authClient }
export default adminClient
