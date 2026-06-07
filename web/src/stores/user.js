import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
  // 状态
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(null)

  // 计算属性
  const isLoggedIn = computed(() => !!token.value)

  // 初始化：如果已有 token，自动获取用户信息
  async function init() {
    if (token.value) {
      try {
        const res = await authApi.getMe()
        if (res.code === 0) {
          userInfo.value = res.data
        } else {
          clearAuth()
        }
      } catch {
        clearAuth()
      }
    }
  }

  // 登录
  async function login(credentials) {
    const res = await authApi.login(credentials)
    if (res.code === 0) {
      token.value = res.data.token
      localStorage.setItem('token', res.data.token)
      userInfo.value = {
        id: res.data.user_id,
        nickname: res.data.nickname,
        avatar_url: res.data.avatar_url,
        school_id: res.data.school_id,
      }
      return true
    }
    throw new Error(res.message || '登录失败')
  }

  // 注册
  async function register(data) {
    const res = await authApi.register(data)
    if (res.code === 0) {
      token.value = res.data.token
      localStorage.setItem('token', res.data.token)
      userInfo.value = {
        id: res.data.user_id,
      }
      // 注册成功后获取完整信息
      try {
        const meRes = await authApi.getMe()
        if (meRes.code === 0) {
          userInfo.value = meRes.data
        }
      } catch { /* ignore */ }
      return true
    }
    throw new Error(res.message || '注册失败')
  }

  // 退出登录
  function logout() {
    clearAuth()
  }

  // 获取用户信息
  async function fetchMe() {
    const res = await authApi.getMe()
    if (res.code === 0) {
      userInfo.value = res.data
      return res.data
    }
    throw new Error(res.message || '获取用户信息失败')
  }

  // 清除认证信息
  function clearAuth() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    init,
    login,
    register,
    logout,
    fetchMe,
    clearAuth,
  }
})
