import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api'

export interface TenantUser {
  user_id: string
  email: string
  name: string
  role: string
  permissions: string
}

export const useUserStore = defineStore('users', () => {
  const users = ref<TenantUser[]>([])

  async function fetchUsers(tenantId: string) {
    const { data } = await api.get(`/tenants/${tenantId}/users`)
    users.value = data || []
  }

  async function inviteUser(tenantId: string, payload: { email: string; role: string; password?: string }) {
    const { data } = await api.post(`/tenants/${tenantId}/users/invite`, payload)
    users.value.push(data)
    return data
  }

  async function updateRole(tenantId: string, userId: string, role: string) {
    await api.put(`/tenants/${tenantId}/users/${userId}/role`, { role })
    const u = users.value.find(u => u.user_id === userId)
    if (u) u.role = role
  }

  async function removeUser(tenantId: string, userId: string) {
    await api.delete(`/tenants/${tenantId}/users/${userId}`)
    users.value = users.value.filter(u => u.user_id !== userId)
  }

  return { users, fetchUsers, inviteUser, updateRole, removeUser }
})
