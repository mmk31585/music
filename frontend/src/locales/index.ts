import { createI18n } from 'vue-i18n'
import fa from './fa.json'
import en from './en.json'

export type Locale = 'fa' | 'en'

export const SUPPORTED_LOCALES: Locale[] = ['fa', 'en']
export const DEFAULT_LOCALE: Locale = 'fa'

const i18n = createI18n({
  legacy: false,
  locale: DEFAULT_LOCALE,
  fallbackLocale: 'en',
  messages: {
    fa,
    en,
  },
})

export default i18n
