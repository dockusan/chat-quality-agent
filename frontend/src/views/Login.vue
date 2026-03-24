<template>
  <v-card class="pa-6" elevation="2">
    <v-card-title class="text-h6 text-center pb-4">{{ $t('login_title') }}</v-card-title>
    <v-alert v-if="errorMsg" type="error" variant="tonal" density="compact" class="mb-4">{{ errorMsg }}</v-alert>
    <v-form @submit.prevent="handleLogin">
      <v-text-field
        v-model="email"
        :label="$t('email')"
        type="email"
        prepend-inner-icon="mdi-email"
        required
        class="mb-2"
      />
      <v-text-field
        v-model="password"
        :label="$t('password')"
        :type="showPass ? 'text' : 'password'"
        prepend-inner-icon="mdi-lock"
        :append-inner-icon="showPass ? 'mdi-eye-off' : 'mdi-eye'"
        @click:append-inner="showPass = !showPass"
        required
        class="mb-4"
      />
      <v-btn type="submit" color="primary" block size="large" :loading="loading">
        {{ $t('login') }}
      </v-btn>
    </v-form>
  </v-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const { t } = useI18n()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const showPass = ref(false)
const loading = ref(false)
const errorMsg = ref('')

async function handleLogin() {
  loading.value = true
  errorMsg.value = ''
  try {
    await authStore.login(email.value, password.value)
    router.push('/')
  } catch {
    errorMsg.value = t('invalid_credentials')
  } finally {
    loading.value = false
  }
}
</script>
