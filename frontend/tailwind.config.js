/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.vue'],
  theme: {
    extend: {},
  },
  plugins: [require('daisyui'), require('@tailwindcss/typography'), require('prettier-plugin-tailwindcss')],
  daisyui: {
    themes: ['corporate']
  }
}
