import { createI18n } from 'vue-i18n'
import vi from './vi'
import en from './en'

const savedLocale = localStorage.getItem('cqa_locale') || 'vi'

export default createI18n({
  legacy: false,
  locale: savedLocale,
  fallbackLocale: 'vi',
  messages: { vi, en },
})
