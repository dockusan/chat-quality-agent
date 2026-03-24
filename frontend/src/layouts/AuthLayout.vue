<template>
  <v-container fluid class="fill-height bg-background">
    <v-row align="center" justify="center" class="fill-height">
      <v-col cols="12" sm="8" md="5" lg="4" xl="3">
        <div class="text-center mb-6">
          <v-icon size="48" color="primary" class="mb-2">mdi-chat-processing</v-icon>
          <h1 class="text-h5 font-weight-bold text-primary">{{ $t('app_name') }}</h1>
        </div>
        <slot />
        <div class="text-center mt-4">
          <LanguageSwitcher />
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import LanguageSwitcher from '../components/LanguageSwitcher.vue'
</script>
