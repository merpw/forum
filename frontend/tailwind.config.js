/** @type {import("tailwindcss").Config} */
module.exports = {
  content: ["./src/**/*.{js,ts,jsx,tsx}"],
  future: {
    hoverOnlyWhenSupported: true,
  },
  theme: {
    container: {
      screens: {
        xl: "800px",
        // TODO: remove when we will have more content to display (e.g. sidebar)
      },
    },
    extend: {
      fontFamily: {
        Alatsi: ["Alatsi", "sans-serif"],
        Yesteryear: ["Yesteryear", "cursive"],
      },
      typography: {
        DEFAULT: {
          css: {
            maxWidth: "100ch", // add required value here
          },
        },
      },
    },
  },
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
}
