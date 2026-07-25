import { definePreset } from '@primeuix/themes'
import Aura from '@primeuix/themes/aura'

/**
 * Muse App Preset — extends PrimeVue Aura theme to match Muse design tokens.
 *
 * All component overrides reference CSS variables from main.css so they
 * automatically adapt to dark/light mode and dynamic album colors.
 */
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
        background: 'var(--surface-popover)',
        borderColor: 'var(--border-default)',
        shadow: 'var(--shadow-dialog)',
        borderRadius: 'var(--radius-xl)',
        backdropFilter: 'blur(32px)',
      },
      mask: {
        background: 'var(--bg-overlay)',
        backdropFilter: 'blur(4px)',
      },
    } as any,
    button: {
      root: {
        borderRadius: '9999px',
        fontWeight: '600',
        padding: '0.625rem 1.5rem',
        transitionDuration: '0.15s',
        transitionProperty: 'background, color, border-color, box-shadow, transform',
        transitionTimingFunction: 'cubic-bezier(0, 0, 0.2, 1)',
      },
      primary: {
        background: 'var(--accent)',
        hoverBackground: 'var(--accent-hover)',
        activeBackground: 'var(--accent-active)',
        color: 'var(--accent-text)',
        hoverColor: 'var(--accent-text)',
        activeColor: 'var(--accent-text)',
        borderColor: 'transparent',
        hoverBorderColor: 'transparent',
        activeBorderColor: 'transparent',
      },
      secondary: {
        background: 'var(--surface-2)',
        hoverBackground: 'var(--surface-3)',
        activeBackground: 'var(--surface-4)',
        color: 'var(--text-secondary)',
        hoverColor: 'var(--text-primary)',
        activeColor: 'var(--text-primary)',
        borderColor: 'var(--border-default)',
        hoverBorderColor: 'var(--border-strong)',
        activeBorderColor: 'var(--border-strong)',
      },
      outlined: {
        background: 'transparent',
        hoverBackground: 'var(--surface-hover)',
        activeBackground: 'var(--surface-active)',
        color: 'var(--text-secondary)',
        hoverColor: 'var(--accent)',
        activeColor: 'var(--accent)',
        borderColor: 'var(--border-strong)',
        hoverBorderColor: 'var(--accent)',
        activeBorderColor: 'var(--accent)',
      },
      text: {
        background: 'transparent',
        hoverBackground: 'var(--accent-subtle)',
        activeBackground: 'var(--accent-muted)',
        color: 'var(--accent)',
        hoverColor: 'var(--accent)',
        activeColor: 'var(--accent)',
      },
      ghost: {
        background: 'transparent',
        hoverBackground: 'var(--surface-hover)',
        activeBackground: 'var(--surface-active)',
        color: 'var(--text-secondary)',
        hoverColor: 'var(--text-primary)',
        activeColor: 'var(--text-primary)',
      },
      info: {
        background: 'var(--info)',
        hoverBackground: 'var(--info-hover)',
        color: '#ffffff',
      },
      success: {
        background: 'var(--success)',
        hoverBackground: 'var(--success-hover)',
        color: '#ffffff',
      },
      warn: {
        background: 'var(--warning)',
        hoverBackground: 'var(--warning-hover)',
        color: '#000000',
      },
      danger: {
        background: 'var(--danger)',
        hoverBackground: 'var(--danger-hover)',
        color: '#ffffff',
      },
      contrast: {
        background: 'var(--surface-900)',
        hoverBackground: 'var(--surface-800)',
        color: 'var(--text-primary)',
      },
    } as any,
    inputtext: {
      root: {
        background: 'var(--surface-hover)',
        borderColor: 'var(--border-default)',
        color: 'var(--text-primary)',
        placeholderColor: 'var(--text-muted)',
        hoverBorderColor: 'var(--border-strong)',
        focusBorderColor: 'var(--accent)',
        borderRadius: 'var(--radius-md)',
        paddingX: '0.875rem',
        paddingY: '0.625rem',
        transitionDuration: '0.15s',
        transitionProperty: 'border-color, box-shadow',
        transitionTimingFunction: 'cubic-bezier(0, 0, 0.2, 1)',
      },
      colorScheme: {
        dark: {
          root: {
            background: 'var(--surface-hover)',
          },
        },
      },
    },
    textarea: {
      root: {
        background: 'var(--surface-hover)',
        borderColor: 'var(--border-default)',
        color: 'var(--text-primary)',
        placeholderColor: 'var(--text-muted)',
        hoverBorderColor: 'var(--border-strong)',
        focusBorderColor: 'var(--accent)',
        borderRadius: 'var(--radius-md)',
        paddingX: '0.875rem',
        paddingY: '0.625rem',
      },
    },
    select: {
      root: {
        background: 'var(--surface-hover)',
        borderColor: 'var(--border-default)',
        color: 'var(--text-primary)',
        hoverBorderColor: 'var(--border-strong)',
        focusBorderColor: 'var(--accent)',
        borderRadius: 'var(--radius-md)',
      },
    },
    dropdown: {
      root: {
        background: 'var(--surface-hover)',
        borderColor: 'var(--border-default)',
        color: 'var(--text-primary)',
        hoverBorderColor: 'var(--border-strong)',
        focusBorderColor: 'var(--accent)',
        borderRadius: 'var(--radius-md)',
      },
      panel: {
        background: 'var(--surface-popover)',
        borderColor: 'var(--border-default)',
        shadow: 'var(--shadow-floating)',
        borderRadius: 'var(--radius-lg)',
      },
      item: {
        focusBackground: 'var(--surface-hover)',
        focusColor: 'var(--text-primary)',
        color: 'var(--text-secondary)',
        borderRadius: 'var(--radius-md)',
      },
    },
    checkbox: {
      root: {
        borderRadius: 'var(--radius-sm)',
        width: '1.25rem',
        height: '1.25rem',
        background: 'var(--surface-active)',
        checkedBackground: 'var(--accent)',
        checkedBorderColor: 'var(--accent)',
        borderColor: 'var(--border-strong)',
        hoverBorderColor: 'var(--accent)',
      },
    },
    inputswitch: {
      root: {
        width: '2.5rem',
        height: '1.5rem',
        background: 'var(--surface-active)',
        checkedBackground: 'var(--accent)',
      },
    },
    toast: {
      root: {
        background: 'var(--surface-elevated)',
        borderColor: 'var(--border-default)',
        backdropFilter: 'blur(24px)',
        borderRadius: 'var(--radius-lg)',
        shadow: 'var(--shadow-elevated)',
      } as any,
    },
    tabs: {
      tablist: {
        borderColor: 'var(--border-default)',
      },
      tab: {
        color: 'var(--text-tertiary)',
        activeColor: 'var(--accent)',
        activeBorderColor: 'var(--accent)',
        fontWeight: '500',
      },
    },
    menu: {
      root: {
        background: 'var(--surface-popover)',
        borderColor: 'var(--border-default)',
        borderRadius: 'var(--radius-lg)',
        shadow: 'var(--shadow-floating)',
        padding: '4px',
      },
      item: {
        focusBackground: 'var(--surface-hover)',
        focusColor: 'var(--text-primary)',
        color: 'var(--text-secondary)',
        borderRadius: 'var(--radius-md)',
      },
    } as any,
    panel: {
      root: {
        background: 'var(--surface-1)',
        borderColor: 'var(--border-default)',
        borderRadius: 'var(--radius-lg)',
      },
      header: {
        borderColor: 'var(--border-subtle)',
      },
    },
    skeleton: {
      root: {
        background: 'var(--surface-2)',
        borderRadius: 'var(--radius-md)',
      },
    },
    progressbar: {
      root: {
        background: 'var(--surface-3)',
        borderRadius: '999px',
        height: '6px',
      },
      value: {
        background: 'var(--accent)',
        borderRadius: '999px',
      },
    },
    chip: {
      root: {
        background: 'var(--surface-2)',
        borderColor: 'var(--border-default)',
        borderRadius: '9999px',
        color: 'var(--text-secondary)',
        paddingX: '0.75rem',
        paddingY: '0.25rem',
      },
    },
    badge: {
      root: {
        background: 'var(--accent)',
        color: 'var(--accent-text)',
        borderRadius: '9999px',
        fontSize: '0.7rem',
        fontWeight: '700',
        minWidth: '1.25rem',
        height: '1.25rem',
      },
    },
    slider: {
      root: {
        track: {
          background: 'var(--surface-3)',
          borderRadius: '999px',
          height: '4px',
        },
        range: {
          background: 'var(--accent)',
          borderRadius: '999px',
        },
        handle: {
          background: 'var(--text-primary)',
          borderColor: 'var(--accent)',
          width: '14px',
          height: '14px',
          borderRadius: '999px',
        },
      },
    } as any,
    paginator: {
      root: {
        background: 'transparent',
        borderColor: 'var(--border-subtle)',
      },
      pageButton: {
        color: 'var(--text-secondary)',
        hoverColor: 'var(--text-primary)',
        borderRadius: 'var(--radius-md)',
      },
      activePage: {
        background: 'var(--accent)',
        color: 'var(--accent-text)',
      },
    } as any,
    speeddial: {
      button: {
        background: 'var(--surface-elevated)',
        borderColor: 'var(--border-default)',
        color: 'var(--text-primary)',
        hoverBackground: 'var(--accent)',
        hoverColor: 'var(--accent-text)',
        borderRadius: '9999px',
        shadow: 'var(--shadow-floating)',
      },
    } as any,
    calendar: {
      root: {
        background: 'var(--surface-hover)',
        borderColor: 'var(--border-default)',
        borderRadius: 'var(--radius-md)',
      },
      dropdown: {
        background: 'var(--surface-hover)',
        borderColor: 'var(--border-default)',
      },
    },
    tree: {
      root: {
        background: 'transparent',
        borderColor: 'var(--border-default)',
      },
      node: {
        focusBackground: 'var(--surface-hover)',
        color: 'var(--text-secondary)',
        borderRadius: 'var(--radius-md)',
      },
    } as any,
    confirmPopup: {
      root: {
        background: 'var(--surface-popover)',
        borderColor: 'var(--border-default)',
        borderRadius: 'var(--radius-lg)',
        shadow: 'var(--shadow-dialog)',
        backdropFilter: 'blur(16px)',
      },
    },
    tooltip: {
      root: {
        background: 'var(--surface-4)',
        color: 'var(--text-primary)',
        borderRadius: 'var(--radius-sm)',
        shadow: 'var(--shadow-floating)',
      },
    },
    menubar: {
      root: {
        background: 'var(--surface-1)',
        borderColor: 'var(--border-default)',
        borderRadius: 'var(--radius-lg)',
        padding: '4px 8px',
      },
      item: {
        focusBackground: 'var(--surface-hover)',
        focusColor: 'var(--text-primary)',
        color: 'var(--text-secondary)',
        borderRadius: 'var(--radius-md)',
      },
    },
    accordion: {
      root: {
        borderColor: 'var(--border-default)',
      },
      content: {
        background: 'transparent',
        borderColor: 'var(--border-subtle)',
      },
    },
    datatable: {
      root: {
        background: 'transparent',
      },
      headerRow: {
        background: 'var(--surface-1)',
      },
      headerCell: {
        borderColor: 'var(--border-default)',
        color: 'var(--text-tertiary)',
        fontWeight: '600',
        fontSize: '0.6875rem',
        letterSpacing: '0.04em',
        textTransform: 'uppercase',
        padding: '0.625rem 1rem',
      },
      bodyRow: {
        borderColor: 'var(--border-subtle)',
        transitionDuration: '0.15s',
        transitionProperty: 'background',
      },
      bodyCell: {
        color: 'var(--text-secondary)',
        fontSize: 'var(--text-sm)',
        padding: '0.75rem 1rem',
      },
    },
  },
})
