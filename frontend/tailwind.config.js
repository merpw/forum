/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    container: {
      screens: {
        xl: "800px",
        // TODO: remove when we will have more content to display (e.g. sidebar)
      },
    },
    extend: {},
  },
}
