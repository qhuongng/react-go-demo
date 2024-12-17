/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  daisyui: {
    themes: [{
      winter: {
        ...require("daisyui/src/theming/themes")["winter"],
        primary: "#394e6a",
        primaryContent: "#ffffff",
        neutral: "#394e6a",
        neutralFocus: "#2a2e37",
        neutralContent: "#ffffff",
      },
    }],
  },
  plugins: [require("daisyui")],
}

