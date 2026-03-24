<template>
  <div>
    <h3 class="text-h6 mb-2">{{ $t('job_summary') }}</h3>
    <div class="text-body-2 text-grey-darken-1 mb-4">Kiểm tra lại thông tin trước khi tạo công việc.</div>

    <v-table density="compact">
      <tbody>
        <tr>
          <td class="font-weight-bold" width="200">{{ $t('job_name') }}</td>
          <td>{{ form.name }}</td>
        </tr>
        <tr>
          <td class="font-weight-bold">{{ $t('select_channels') }}</td>
          <td>{{ form.input_channel_ids?.length || 0 }} kênh</td>
        </tr>
        <tr>
          <td class="font-weight-bold">{{ $t('job_wizard_step_output') }}</td>
          <td>{{ outputCount }} đầu ra</td>
        </tr>
        <tr>
          <td class="font-weight-bold">{{ $t('job_wizard_step_output_schedule') }}</td>
          <td>
            <span v-if="form.output_schedule === 'instant'">{{ $t('output_instant') }}</span>
            <span v-else-if="form.output_schedule === 'cron'">{{ $t('output_scheduled') }}: {{ form.output_cron }}</span>
            <span v-else>{{ $t('output_once') }}: {{ form.output_at }}</span>
          </td>
        </tr>
        <tr>
          <td class="font-weight-bold">{{ $t('job_wizard_step_analysis_schedule') }}</td>
          <td>
            <span v-if="form.schedule_type === 'cron'">{{ $t('schedule_cron') }}: {{ form.schedule_cron }}</span>
            <span v-else-if="form.schedule_type === 'after_sync'">{{ $t('schedule_after_sync') }}</span>
            <span v-else>{{ $t('schedule_manual') }}</span>
          </td>
        </tr>
      </tbody>
    </v-table>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const form = defineModel<Record<string, any>>('form', { required: true })

const outputCount = computed(() => {
  try {
    return JSON.parse(form.value.outputs || '[]').length
  } catch {
    return 0
  }
})
</script>
