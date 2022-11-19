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
      'white': '#FFFFFF',
      'off-white': '#F2F2F2',
      'black': '#000000',
      'gray': '#BDBDBD',
      'dark-gray': '#7E7E7E',
      'very-dark-gray': '#4F4F4F',
      'light-gray': '#f3f4f6',
      'avail-green': colors.emerald, // The green used for marking availability
      'red': '#B13C3C',
      'bright-red': '#DB1616',
    }
  },
  plugins: [
  ],
  prefix: 'tw-',
}
