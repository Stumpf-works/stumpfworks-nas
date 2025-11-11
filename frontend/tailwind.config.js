/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // macOS-inspired color palette
        'macos': {
          'blue': '#007AFF',
          'green': '#34C759',
          'red': '#FF3B30',
          'orange': '#FF9500',
          'purple': '#AF52DE',
          'yellow': '#FFCC00',
          'pink': '#FF2D55',
          'gray': {
            50: '#F9FAFB',
            100: '#F3F4F6',
            200: '#E5E7EB',
            300: '#D1D5DB',
            400: '#9CA3AF',
            500: '#6B7280',
            600: '#4B5563',
            700: '#374151',
            800: '#1F2937',
            900: '#111827',
          },
          'dark': {
            50: '#1E1E1E',
            100: '#2C2C2C',
            200: '#3A3A3A',
            300: '#4A4A4A',
            400: '#6E6E6E',
            500: '#8E8E93',
          },
        },
      },
      backdropBlur: {
        'macos': '40px',
      },
      boxShadow: {
        'macos': '0 10px 40px rgba(0, 0, 0, 0.15)',
        'macos-sm': '0 1px 2px rgba(0, 0, 0, 0.05)',
        'macos-md': '0 4px 6px rgba(0, 0, 0, 0.1)',
        'macos-lg': '0 10px 25px rgba(0, 0, 0, 0.15)',
        'macos-xl': '0 20px 40px rgba(0, 0, 0, 0.2)',
        'window': '0 20px 60px rgba(0, 0, 0, 0.25)',
      },
      borderRadius: {
        'macos': '0.75rem',
        'macos-lg': '1rem',
        'macos-xl': '1.5rem',
      },
      fontFamily: {
        sans: ['system-ui', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'sans-serif'],
        mono: ['SF Mono', 'Monaco', 'Cascadia Code', 'Consolas', 'monospace'],
      },
      animation: {
        'dock-bounce': 'bounce 0.5s ease-in-out',
        'fade-in': 'fadeIn 0.3s ease-in-out',
        'slide-in': 'slideIn 0.3s ease-in-out',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideIn: {
          '0%': { transform: 'translateY(-10px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
      },
    },
  },
  plugins: [],
}
