/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        telegram: {
          blue: '#2481cc',
          bg: '#f0f2f5',
          chat: '#ffffff',
          bubble: '#effdde',
        }
      }
    },
  },
  plugins: [],
}
