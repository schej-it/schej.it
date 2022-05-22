const colors = require('tailwindcss/colors')

module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  important: true,
  theme: {
    extend: {},
    colors: {
      transparent: 'transparent',
      current: 'currentColor',
      'light-green': '#29BC68',
      'green': '#219653',
      'dark-green': '#1C7D45',
      'light-blue': '#6FACCF',
      'blue': '#2F80ED',
      'white': '#F2F2F2',
      'black': '#4F4F4F',
      'gray': '#BDBDBD',
      'light-gray': '#f3f4f6',
      'avail-green': colors.emerald, // The green used for marking availability
    },
  },
  plugins: [
  ],
  prefix: 'tw-',
}
