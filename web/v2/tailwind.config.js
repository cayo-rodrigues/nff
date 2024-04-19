/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./static/js/**.js",
    "./static/css/**.css",
    "./ui/**.templ",
    "./ui/forms/**.templ",
    "./ui/pages/**.templ",
    "./ui/layouts/**.templ",
    "./ui/scripts/**.templ",
    "./ui/shared/**.templ",
    "./ui/components/**.templ",
  ],
  theme: {
    extend: {
      animation: {
        'spin-slow': 'spin 2s linear infinite',
      }
    },
  },
  plugins: [],
};

