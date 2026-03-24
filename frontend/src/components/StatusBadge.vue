<template>
  <v-chip :color="color" :size="size" variant="tonal">
    <v-icon v-if="icon" start size="small">{{ icon }}</v-icon>
    {{ label }}
  </v-chip>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  status: string
  size?: string
}>()

const { t } = useI18n()

const color = computed(() => {
  switch (props.status) {
    case 'success': case 'active': case 'sent': return 'success'
    case 'error': case 'failed': return 'error'
    case 'warning': case 'CAN_CAI_THIEN': return 'warning'
    case 'critical': case 'NGHIEM_TRONG': return 'error'
    case 'running': case 'info': return 'info'
    default: return 'grey'
  }
})

const icon = computed(() => {
  switch (props.status) {
    case 'success': case 'sent': return 'mdi-check-circle'
    case 'error': case 'failed': return 'mdi-close-circle'
    case 'warning': return 'mdi-alert'
    case 'running': return 'mdi-loading mdi-spin'
    default: return ''
  }
})

const label = computed(() => {
  switch (props.status) {
    case 'success': return t('success')
    case 'error': return t('error')
    case 'active': return t('active')
    case 'inactive': return t('inactive')
    case 'NGHIEM_TRONG': return t('severity_critical')
    case 'CAN_CAI_THIEN': return t('severity_warning')
    default: return props.status
  }
})
</script>
