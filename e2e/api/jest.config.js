/** @type {import('ts-jest').JestConfigWithTsJest} **/
module.exports = {
  testEnvironment: 'node',
  transform: {
    '^.+.tsx?$': ['<rootDir>/e2e/api/node_modules/ts-jest/preprocessor.js', {}]
  },
  testMatch: ['**/_api.spec.ts'],
  collectCoverage: true,
  collectCoverageFrom: [
    'frontend/src/api/*.ts',
  ],
  rootDir: '../../',
  coverageDirectory: 'e2e/api/coverage'
}