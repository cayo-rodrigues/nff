/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./static/**.js",
    "./static/**.css",
    "./components/**.templ",
    "./components/forms/**.templ",
    "./components/pages/**.templ",
    "./components/layouts/**.templ",
    "./components/scripts/**.templ",
    "./components/shared/**.templ",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
};

