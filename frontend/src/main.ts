import App from './App.vue'

import { createApp } from 'vue'
import PrimeVue from 'primevue/config'
import Aura from '@primeuix/themes/aura'
import 'primeicons/primeicons.css'
import './style.scss'
import { definePreset } from '@primeuix/themes'

const NoirX = definePreset(Aura, {
  semantic: {
    primary: {
      50: '#EFF6FF',
      100: '#DBEAFE',
      200: '#BFDBFE',
      300: '#93C5FD',
      400: '#60A5FA',
      500: '#3B82F6',
      600: '#2563EB',
      700: '#1D4ED8',
      800: '#1E40AF',
      900: '#1E3A8A',
      950: '#172554',
    },
    accent: {
      red: {
        500: '#EF4444',
        600: '#DC2626',
      },
      green: {
        500: '#10B981',
        600: '#059669',
      },
      yellow: {
        500: '#FACC15',
        600: '#EAB308',
      },
      purple: {
        500: '#A855F7',
        600: '#9333EA',
      },
    },
    colorScheme: {
      light: {
        primary: {
          color: '{primary.600}',
          inverseColor: '#ffffff',
          hoverColor: '{primary.500}',
          activeColor: '{primary.700}',
        },
        highlight: {
          background: '{primary.600}',
          focusBackground: '{primary.400}',
          color: '#ffffff',
          focusColor: '#ffffff',
        },
        accent: {
          color: '{accent.purple.500}',
          hoverColor: '{accent.purple.400}',
          activeColor: '{accent.purple.600}',
        },
      },
      dark: {
        primary: {
          color: '{primary.100}',
          inverseColor: '{primary.950}',
          hoverColor: '{primary.200}',
          activeColor: '{primary.300}',
        },
        highlight: {
          background: 'rgba(250, 250, 250, .16)',
          focusBackground: 'rgba(250, 250, 250, .24)',
          color: 'rgba(255,255,255,.87)',
          focusColor: 'rgba(255,255,255,.87)',
        },
        accent: {
          color: '{accent.red.500}',
          hoverColor: '{accent.red.400}',
          activeColor: '{accent.red.600}',
        },
      },
    },
  },
})

const app = createApp(App)
app.use(PrimeVue, {
  theme: {
    preset: NoirX,
    options: {
      darkModeSelector: '.dark-mode',
    },
  },
})

app.mount('#app')
