/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./static/**.js",
    "./static/**.css",
    "./ui/**.templ",
    "./ui/forms/**.templ",
    "./ui/pages/**.templ",
    "./ui/layouts/**.templ",
    "./ui/scripts/**.templ",
    "./ui/shared/**.templ",
    "./ui/components/**.templ",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
};

