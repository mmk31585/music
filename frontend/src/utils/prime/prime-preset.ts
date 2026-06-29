import { definePreset } from '@primeuix/themes'
import Aura from '@primeuix/themes/aura'

export const AppPreset = definePreset(Aura, {
  semantic: {
    primary: {
      50: '{emerald.50}',
      100: '{emerald.100}',
      200: '{emerald.200}',
      300: '{emerald.300}',
      400: '{emerald.400}',
      500: '{emerald.500}',
      600: '{emerald.600}',
      700: '{emerald.700}',
      800: '{emerald.800}',
      900: '{emerald.900}',
      950: '{emerald.950}',
    },
    colorScheme: {
      light: {
        surface: {
          0: '#ffffff',
          50: '{zinc.50}',
          100: '{zinc.100}',
          200: '{zinc.200}',
          300: '{zinc.300}',
          400: '{zinc.400}',
          500: '{zinc.500}',
          600: '{zinc.600}',
          700: '{zinc.700}',
          800: '{zinc.800}',
          900: '{zinc.900}',
          950: '{zinc.950}',
        },
      },
      dark: {
        surface: {
          0: '#ffffff',
          50: '{slate.50}',
          100: '{slate.100}',
          200: '{slate.200}',
          300: '{slate.300}',
          400: '{slate.400}',
          500: '{slate.500}',
          600: '{slate.600}',
          700: '{slate.700}',
          800: '{slate.800}',
          900: '{slate.900}',
          950: '{slate.950}',
        },
      },
    },
  },
  components: {
    dialog: {
      root: {
        background: 'rgba(0, 0, 0, 0.85)',
        borderColor: 'rgba(255, 255, 255, 0.08)',
        shadow: '0 16px 64px rgba(0, 0, 0, 0.6)',
        backdropFilter: 'blur(32px)',
      },
      mask: {
        background: 'rgba(0, 0, 0, 0.5)',
        backdropFilter: 'blur(4px)',
      },
    } as any,
    button: {
      root: {
        borderRadius: '9999px',
        fontWeight: '700',
        padding: '0.625rem 1.5rem',
      },
      primary: {
        background: '#1db954',
        hoverBackground: '#1ed760',
        activeBackground: '#169c46',
        color: '#000000',
        hoverColor: '#000000',
        activeColor: '#000000',
        borderColor: 'transparent',
        hoverBorderColor: 'transparent',
        activeBorderColor: 'transparent',
      },
    } as any,
    inputtext: {
      root: {
        background: 'rgba(255, 255, 255, 0.06)',
        borderColor: 'rgba(255, 255, 255, 0.08)',
        color: '#ffffff',
        placeholderColor: 'rgba(255, 255, 255, 0.3)',
        focusBorderColor: '#1db954',
        borderRadius: '0.75rem',
        paddingX: '1rem',
        paddingY: '0.625rem',
        transitionDuration: '0.15s',
      },
    },
    select: {
      root: {
        background: 'rgba(255, 255, 255, 0.06)',
        borderColor: 'rgba(255, 255, 255, 0.08)',
        color: '#ffffff',
        focusBorderColor: '#1db954',
        borderRadius: '0.75rem',
      },
    },
    checkbox: {
      root: {
        borderRadius: '0.375rem',
        width: '1.25rem',
        height: '1.25rem',
        background: 'rgba(255, 255, 255, 0.1)',
        checkedBackground: '#1db954',
        checkedBorderColor: '#1db954',
        borderColor: 'rgba(255, 255, 255, 0.2)',
      },
    },
    inputswitch: {
      root: {
        width: '2.5rem',
        height: '1.5rem',
        background: 'rgba(255, 255, 255, 0.15)',
        checkedBackground: '#1db954',
      },
    },
    toast: {
      root: {
        background: 'rgba(0, 0, 0, 0.85)',
        borderColor: 'rgba(255, 255, 255, 0.08)',
        backdropFilter: 'blur(24px)',
        borderRadius: '1rem',
        shadow: '0 8px 32px rgba(0, 0, 0, 0.5)',
      } as any,
    },
  },
})
