import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 10000,
})

// 请求拦截器：自动附加 Token
request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// 响应拦截器：统一错误处理
request.interceptors.response.use(
  (response) => response.data,
  (error) => {
    const { code, message } = error.response?.data || {}

    if (code === 1100 || error.response?.status === 401) {
      // Token 过期，清除登录态，跳转登录
      localStorage.removeItem('token')
      router.push({ name: 'Login' })
      ElMessage.warning('登录已过期，请重新登录')
    } else {
      ElMessage.error(message || '网络错误，请稍后重试')
    }

    return Promise.reject(new Error(message || '网络错误'))
  }
)

export default request
