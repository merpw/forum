/** @type {import("tailwindcss").Config} */
module.exports = {
  content: ["./src/**/*.{js,ts,jsx,tsx}"],
  future: {
    hoverOnlyWhenSupported: true,
  },
  theme: {
    container: {
      screens: {
        xl: "1200px",
        // TODO: remove when we will have more content to display (e.g. sidebar)
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
        // TODO: choose colors
        light_theme: {
          "base-100": "#e0e0ea", // text
          primary: "#6176cb", // username
          secondary: "#ffac8d", // user's reaction
          accent: "#69f19f", // online status
          neutral: "#bebbc7", // break lines, background etc.
          info: "#757575", // reactions, comments button etc.
        },
        dark_theme: {
          "base-100": "#2c2c2c",
          primary: "#505bdc",
          secondary: "#ffa874",
          accent: "#00e16a",
          neutral: "#444444",
          info: "#797979",
        },
      },
    ],
  },
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
}
