import type { Config } from 'tailwindcss'

export default <Partial<Config>>{
  content: [
    './components/**/*.{vue,js,ts}',
    './layouts/**/*.{vue,js,ts}',
    './pages/**/*.{vue,js,ts}',
    './composables/**/*.{js,ts}',
    './plugins/**/*.{js,ts}',
    './app.vue'
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'Nunito', 'system-ui', 'sans-serif'],
        display: ['Outfit', 'Quicksand', 'system-ui', 'sans-serif'],
      },
      colors: {
        // MilkyVerse "Morning Latte" Palette
        primary: '#FFD93D',
        bg: {
          cream: {
            100: '#FFFDF5', // Base Cream
            200: '#FFF7D1', // Butter Yellow
          },
          pink: {
            100: '#FFF0F5', // Soft Pink
          }
        },
        charcoal: {
          100: '#F2F3F5',
          200: '#E0E1E4',
          300: '#C4C5CA', // Borders
          400: '#9A9CA2',
          500: '#6E7077', // Meta
          600: '#4D4E54',
          700: '#3A3B40', // Subheadings
          800: '#25262B',
          900: '#1A1B1E', // Hero / Primary
        },
        accent: {
          yellow: {
            DEFAULT: '#FFD93D',
            soft: '#FFE58A',
          },
          pink: {
            soft: '#FFE5F0',
          }
        },
        status: {
          success: '#5FB36D',
          warning: '#F7C948',
          error: '#E9515C',
        }
      },
      animation: {
        'marquee': 'marquee 24s linear infinite',
        'stretch-in': 'stretchIn 0.8s cubic-bezier(.34,1.56,.64,1) forwards',
        'fade-up': 'fadeUp 0.8s cubic-bezier(0.16, 1, 0.3, 1) forwards',
      },
      keyframes: {
        marquee: {
          '0%': { transform: 'translateX(0)' },
          '100%': { transform: 'translateX(-50%)' },
        },
        stretchIn: {
          '0%': { transform: 'scaleX(0)' },
          '100%': { transform: 'scaleX(1)' },
        },
        fadeUp: {
          '0%': { opacity: '0', transform: 'translateY(20px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
      }
    }
  }
}
