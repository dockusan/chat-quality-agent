<template>
  <div>
    <v-radio-group v-model="frequency" density="compact">
      <v-radio value="daily" :label="$t('cron_daily')" />
      <v-radio value="weekly" :label="$t('cron_weekly')" />
      <v-radio value="monthly" :label="$t('cron_monthly')" />
    </v-radio-group>

    <!-- Time picker -->
    <v-row dense class="mt-2 align-center">
      <v-col cols="auto">
        <span class="text-body-2">{{ $t('cron_at') }}</span>
      </v-col>
      <v-col cols="3">
        <v-select
          v-model="hour"
          :items="hours"
          density="compact"
          variant="outlined"
          hide-details
        />
      </v-col>
      <v-col cols="auto"><span class="text-body-2">:</span></v-col>
      <v-col cols="3">
        <v-select
          v-model="minute"
          :items="minutes"
          density="compact"
          variant="outlined"
          hide-details
        />
      </v-col>
    </v-row>

    <!-- Weekly: day of week -->
    <div v-if="frequency === 'weekly'" class="mt-3">
      <div class="text-body-2 mb-2">{{ $t('cron_select_days') }}</div>
      <v-btn-toggle v-model="selectedDays" multiple density="compact" color="primary" variant="outlined">
        <v-btn v-for="(day, idx) in dayLabels" :key="idx" :value="idx" size="small">
          {{ day }}
        </v-btn>
      </v-btn-toggle>
    </div>

    <!-- Monthly: day of month -->
    <div v-if="frequency === 'monthly'" class="mt-3">
      <v-select
        v-model="dayOfMonth"
        :items="daysOfMonth"
        :label="$t('cron_day_of_month')"
        density="compact"
        variant="outlined"
        hide-details
      />
    </div>

    <!-- Preview -->
    <v-alert v-if="preview" type="info" variant="tonal" density="compact" class="mt-3 text-body-2">
      {{ preview }}
    </v-alert>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const modelValue = defineModel<string>({ default: '0 7 * * *' })

const frequency = ref('daily')
const hour = ref(7)
const minute = ref(0)
const selectedDays = ref<number[]>([1, 2, 3, 4, 5]) // Mon-Fri
const dayOfMonth = ref(1)

const hours = Array.from({ length: 24 }, (_, i) => ({ title: String(i).padStart(2, '0'), value: i }))
const minutes = [0, 5, 10, 15, 20, 30, 45].map(m => ({ title: String(m).padStart(2, '0'), value: m }))
const daysOfMonth = Array.from({ length: 31 }, (_, i) => ({ title: `${t('cron_day')} ${i + 1}`, value: i + 1 }))

const dayLabels = ['CN', 'T2', 'T3', 'T4', 'T5', 'T6', 'T7']

const cronExpression = computed(() => {
  const m = minute.value
  const h = hour.value
  switch (frequency.value) {
    case 'daily':
      return `${m} ${h} * * *`
    case 'weekly': {
      const days = selectedDays.value.length > 0 ? selectedDays.value.sort().join(',') : '*'
      return `${m} ${h} * * ${days}`
    }
    case 'monthly':
      return `${m} ${h} ${dayOfMonth.value} * *`
    default:
      return `${m} ${h} * * *`
  }
})

const preview = computed(() => {
  const timeStr = `${String(hour.value).padStart(2, '0')}:${String(minute.value).padStart(2, '0')}`
  switch (frequency.value) {
    case 'daily':
      return `${t('cron_preview_daily')} ${timeStr}`
    case 'weekly': {
      const dayNames = selectedDays.value.sort().map(d => dayLabels[d]).join(', ')
      return `${t('cron_preview_weekly')} ${dayNames} ${t('cron_at')} ${timeStr}`
    }
    case 'monthly':
      return `${t('cron_preview_monthly')} ${dayOfMonth.value} ${t('cron_at')} ${timeStr}`
    default:
      return ''
  }
})

// Parse initial value
function parseCron(cron: string) {
  const parts = cron.trim().split(/\s+/)
  if (parts.length !== 5) return
  const [min, hr, dom, , dow] = parts
  minute.value = parseInt(min) || 0
  hour.value = parseInt(hr) || 7
  if (dom !== '*') {
    frequency.value = 'monthly'
    dayOfMonth.value = parseInt(dom) || 1
  } else if (dow !== '*') {
    frequency.value = 'weekly'
    selectedDays.value = dow.split(',').map(Number).filter(n => !isNaN(n))
  } else {
    frequency.value = 'daily'
  }
}

// Init from model
if (modelValue.value) {
  parseCron(modelValue.value)
}

// Emit changes
watch(cronExpression, (val) => {
  modelValue.value = val
})
</script>
