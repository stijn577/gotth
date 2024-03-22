/** @type {import('tailwindcss').Config} */
const defaultTheme = require("@catppuccin/tailwindcss");
module.exports = {
  content: [
    "./internal/**/*.html",
    "./internal/**/*.templ",
    "./internal/**/*.go",
    "./internal/**/*.js",
    "./internal/**/*.css",
  ],
  theme: {
    screens: {
      tablet: "620px",
      laptop: "1280px",
      desktop: "1700px",
    },
    fontFamily: {
      custom: ["JetBrainsMono Nerd Font"],
      mono: ["FiraCode"],
    },
    extend: {},
  },
  plugins: [
    require("@catppuccin/tailwindcss"),
    require("tailwind-scrollbar")({
      nocompatible: true,
    }),
  ],
};
