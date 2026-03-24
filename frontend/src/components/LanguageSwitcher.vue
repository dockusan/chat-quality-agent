<template>
  <v-menu>
    <template #activator="{ props }">
      <v-btn variant="text" size="small" v-bind="props" class="text-none">
        <v-icon start size="small">mdi-translate</v-icon>
        {{ currentLabel }}
        <v-icon end size="small">mdi-menu-down</v-icon>
      </v-btn>
    </template>
    <v-list density="compact">
      <v-list-item
        v-for="lang in languages"
        :key="lang.code"
        :active="locale === lang.code"
        @click="switchLocale(lang.code)"
      >
        <v-list-item-title>{{ lang.label }}</v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { locale } = useI18n()

const languages = [
  { code: 'vi', label: 'Tiếng Việt' },
  { code: 'en', label: 'English' },
]

const currentLabel = computed(() => languages.find((l) => l.code === locale.value)?.label || 'Tiếng Việt')

function switchLocale(code: string) {
  locale.value = code
  localStorage.setItem('cqa_locale', code)
}
</script>
