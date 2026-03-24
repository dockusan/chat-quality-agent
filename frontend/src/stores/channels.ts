import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api'

export interface Channel {
  id: string
  tenant_id: string
  channel_type: string
  name: string
  external_id: string
  is_active: boolean
  metadata: string
  last_sync_at: string | null
  last_sync_status: string
  conversation_count: number
  created_at: string
}

export const useChannelStore = defineStore('channels', () => {
  const channels = ref<Channel[]>([])

  async function fetchChannels(tenantId: string) {
    const { data } = await api.get(`/tenants/${tenantId}/channels`)
    channels.value = data
  }

  async function createChannel(tenantId: string, payload: Record<string, unknown>) {
    const { data } = await api.post(`/tenants/${tenantId}/channels`, payload)
    channels.value.unshift(data)
    return data
  }

  async function updateChannel(tenantId: string, channelId: string, payload: Record<string, unknown>) {
    await api.put(`/tenants/${tenantId}/channels/${channelId}`, payload)
    const idx = channels.value.findIndex((c) => c.id === channelId)
    if (idx >= 0) Object.assign(channels.value[idx], payload)
  }

  async function deleteChannel(tenantId: string, channelId: string) {
    await api.delete(`/tenants/${tenantId}/channels/${channelId}`)
    channels.value = channels.value.filter((c) => c.id !== channelId)
  }

  async function testConnection(tenantId: string, channelId: string) {
    const { data } = await api.post(`/tenants/${tenantId}/channels/${channelId}/test`)
    return data
  }

  async function syncChannel(tenantId: string, channelId: string) {
    const { data } = await api.post(`/tenants/${tenantId}/channels/${channelId}/sync`)
    return data
  }

  const currentChannel = ref<Channel | null>(null)
  const syncHistory = ref<any[]>([])
  const syncHistoryTotal = ref(0)

  async function fetchChannel(tenantId: string, channelId: string) {
    const { data } = await api.get(`/tenants/${tenantId}/channels/${channelId}`)
    currentChannel.value = data
    return data
  }

  async function fetchSyncHistory(tenantId: string, channelId: string, page = 1) {
    const { data } = await api.get(`/tenants/${tenantId}/channels/${channelId}/sync-history`, { params: { page, per_page: 10 } })
    syncHistory.value = data.data || []
    syncHistoryTotal.value = data.total || 0
    return data
  }

  async function purgeConversations(tenantId: string, channelId: string) {
    const { data } = await api.delete(`/tenants/${tenantId}/channels/${channelId}/conversations`)
    return data
  }

  return { channels, currentChannel, syncHistory, syncHistoryTotal, fetchChannels, fetchChannel, createChannel, updateChannel, deleteChannel, testConnection, syncChannel, fetchSyncHistory, purgeConversations }
})
