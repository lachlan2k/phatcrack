import { globalIgnores } from 'eslint/config'
import { defineConfigWithVueTs, vueTsConfigs } from '@vue/eslint-config-typescript'
import pluginVue from 'eslint-plugin-vue'
import importPlugin from 'eslint-plugin-import'
import skipFormatting from '@vue/eslint-config-prettier/skip-formatting'

// To allow more languages other than `ts` in `.vue` files, uncomment the following lines:
// import { configureVueProject } from '@vue/eslint-config-typescript'
// configureVueProject({ scriptLangs: ['ts', 'tsx'] })
// More info at https://github.com/vuejs/eslint-config-typescript/#advanced-setup

export default defineConfigWithVueTs(
  {
    name: 'app/files-to-lint',
    files: ['**/*.{ts,mts,tsx,vue}'],
  },

  globalIgnores(['**/dist/**', '**/dist-ssr/**', '**/coverage/**', './tailwind.config.js']),

  pluginVue.configs['flat/essential'],
  vueTsConfigs.recommended,
  importPlugin.flatConfigs.recommended,
  importPlugin.flatConfigs.typescript,
  skipFormatting,

  {
    languageOptions: {
      ecmaVersion: 'latest',
      sourceType: 'module'
    },

    settings: {
      'import/resolver': {
        typescript: true,
      },
    },

    rules: {
      'vue/multi-word-component-names': 'off',

      '@typescript-eslint/no-explicit-any': 'off',
      '@typescript-eslint/no-empty-object-type': 'off',

        'import/order': ['error', {
            'newlines-between': 'always',

            pathGroups: [{
                pattern: '@/components/**',
                group: 'internal',
                position: 'before',
            }, {
                pattern: '@/api/**',
                group: 'internal',
                position: 'before',
            }, {
                pattern: '@/composables/**',
                group: 'internal',
                position: 'before',
            }, {
                pattern: '@/stores/**',
                group: 'internal',
                position: 'before',
            }, {
                pattern: '@/util/**',
                group: 'internal',
                position: 'before',
            }],
        }],

        'import/newline-after-import': 'error',
        'import/first': 'error',
    }
  }
)
