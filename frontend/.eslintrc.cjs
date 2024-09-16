/* eslint-env node */
require('@rushstack/eslint-patch/modern-module-resolution')

module.exports = {
  root: true,
  'extends': [
    'plugin:vue/vue3-essential',
    'eslint:recommended',
    '@vue/eslint-config-typescript',
    '@vue/eslint-config-prettier/skip-formatting',
    'plugin:import/recommended',
    'plugin:import/typescript',
  ],
  parserOptions: {
    ecmaVersion: 'latest'
  },
  env: {
    browser: true,
    node: true
  },
  rules: {
    'vue/multi-word-component-names': 0 ,
    'import/order': 'error'
  },
  settings: {
    'import/resolver': {
      typescript: true
    },
  },
}
