const production = !process.env.ROLLUP_WATCH

module.exports = {
  purge: [],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {},
  },
  variants: {
    extend: {},
  },
  purge: {
    content: [
      './src/**/*.svelte',
      './src/**/*.html',
    ],
    enabled: production, // disable purge in dev
  },

  plugins: [],
}
