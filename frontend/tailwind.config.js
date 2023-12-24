const colors = require("tailwindcss/colors")

module.exports = {
  content: ["./index.html", "./src/**/*.{vue,js,ts,jsx,tsx}"],
  important: true,
  theme: {
    extend: {},
    colors: {
      transparent: "transparent",
      current: "currentColor",
      "pale-green": "#CDEBDC",
      "light-green": "#29BC68",
      "ligher-green": "#EBF7EF",
      green: "#00994C",
      "dark-green": "#1C7D45",
      "darkest-green": "#007F36",
      "light-blue": "#6FACCF",
      blue: "#2F80ED",
      white: "#FFFFFF",
      "off-white": "#F2F2F2",
      black: "#000000",
      gray: "#BDBDBD",
      "dark-gray": "#7E7E7E",
      "very-dark-gray": "#4F4F4F",
      "light-gray": "#f3f4f6",
      "light-gray-stroke": "#dfdfdf",
      "avail-green": colors.emerald, // The green used for marking availability
      red: "#B13C3C",
      "bright-red": "#DB1616",
    },
    screens: {
      sm: "640px",
      md: "768px",
      mdlg: "896px",
      lg: "1024px",
      xl: "1280px",
      "2xl": "1536px",
    },
  },
  plugins: [],
  prefix: "tw-",
}
