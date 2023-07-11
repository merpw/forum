/** @type {import("tailwindcss").Config} */
module.exports = {
  content: ["./src/**/*.{js,ts,jsx,tsx}"],
  future: {
    hoverOnlyWhenSupported: true,
  },
  theme: {
    container: {
      screens: {
        xl: "1280px",
      },
    },
    extend: {
      fontFamily: {
        Alatsi: "var(--font-alatsi)",
        Yesteryear: "var(--font-yesteryear)",
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
  daisyui: {
    styled: true,
    base: false,
    utils: true,
    logs: true,
    rtl: false,
    prefix: "",
    darkTheme: "dark_theme",
    themes: [
      {
        light_theme: {
          "base-100": "#ededf1", // element footer background, react-select background
          primary: "#6176cb", // username
          secondary: "#d37ffc", // user's reaction
          "secondary-content": "#e6d9f1", // message background
          accent: "#69f19f", // online status
          neutral: "#bebbc7", // break lines, background etc.
          info: "#656565", // reactions, comments button etc.
          error: "#ffc3cd", // error message
          "error-content": "#b73945", // error message content
        },
        dark_theme: {
          "base-100": "#1d1d26",
          primary: "#538bff",
          secondary: "#b869d7",
          "secondary-content": "#4d3257",
          accent: "#00e16a",
          neutral: "#444444",
          info: "#969696",
          error: "#7a3d43",
          "error-content": "#fdc9c9",
        },
      },
    ],
  },
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
}
