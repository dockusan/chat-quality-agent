import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../api'

interface User {
  id: string
  email: string
  name: string
  is_admin: boolean
  language: string
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const accessToken = ref(localStorage.getItem('cqa_access_token') || '')
  const isAuthenticated = computed(() => !!accessToken.value)

  async function login(email: string, password: string) {
    const { data } = await api.post('/auth/login', { email, password })
    accessToken.value = data.access_token
    localStorage.setItem('cqa_access_token', data.access_token)
    // Refresh token is now set as HttpOnly cookie by backend
    await fetchProfile()
  }

  async function register(name: string, email: string, password: string) {
    const { data } = await api.post('/auth/register', { name, email, password })
    accessToken.value = data.access_token
    localStorage.setItem('cqa_access_token', data.access_token)
    await fetchProfile()
  }

  async function fetchProfile() {
    const { data } = await api.get('/profile')
    user.value = data
  }

  async function updateProfile(name: string) {
    const { data } = await api.put('/profile', { name })
    if (user.value) {
      user.value.name = data.name || name
    }
  }

  async function logout() {
    try { await api.post('/auth/logout') } catch { /* ignore */ }
    user.value = null
    accessToken.value = ''
    localStorage.removeItem('cqa_access_token')
    localStorage.removeItem('cqa_refresh_token') // cleanup legacy
  }

  return { user, accessToken, isAuthenticated, login, register, fetchProfile, updateProfile, logout }
})
