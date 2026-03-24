import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api'

interface Tenant {
  id: string
  name: string
  slug: string
  channels_count?: number
  jobs_count?: number
}

export const useTenantStore = defineStore('tenants', () => {
  const tenants = ref<Tenant[]>([])
  const currentTenant = ref<Tenant | null>(null)

  async function fetchTenants() {
    const { data } = await api.get('/tenants')
    tenants.value = data
  }

  async function createTenant(name: string, slug: string) {
    const { data } = await api.post('/tenants', { name, slug })
    tenants.value.push(data)
    return data
  }

  function setCurrentTenant(tenant: Tenant) {
    currentTenant.value = tenant
    localStorage.setItem('cqa_current_tenant', tenant.id)
  }

  function loadSavedTenant() {
    const savedId = localStorage.getItem('cqa_current_tenant')
    if (savedId && tenants.value.length) {
      const found = tenants.value.find((t) => t.id === savedId)
      if (found) currentTenant.value = found
    }
  }

  async function deleteTenant(tenantId: string) {
    await api.delete(`/tenants/${tenantId}`)
    tenants.value = tenants.value.filter(t => t.id !== tenantId)
    if (currentTenant.value?.id === tenantId) {
      currentTenant.value = tenants.value[0] || null
      if (currentTenant.value) {
        localStorage.setItem('cqa_current_tenant', currentTenant.value.id)
      } else {
        localStorage.removeItem('cqa_current_tenant')
      }
    }
  }

  return { tenants, currentTenant, fetchTenants, createTenant, deleteTenant, setCurrentTenant, loadSavedTenant }
})
