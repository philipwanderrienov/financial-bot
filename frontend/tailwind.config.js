/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        'finance-bg': '#020617',
        'finance-panel': 'rgba(15, 23, 42, 0.78)',
        'finance-border': 'rgba(148, 163, 184, 0.16)'
      },
      boxShadow: {
        glow: '0 20px 60px rgba(15, 23, 42, 0.35)'
      }
    }
  },
  plugins: []
}