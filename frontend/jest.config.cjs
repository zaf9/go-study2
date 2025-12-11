const nextJest = require("next/jest");

const createJestConfig = nextJest({ dir: "./" });

const customJestConfig = {
  testEnvironment: "jest-environment-jsdom",
  setupFilesAfterEnv: ["<rootDir>/jest.setup.ts"],
  moduleNameMapper: {
    "^@/(.*)$": "<rootDir>/$1",
    "\\.(css|less|scss|sass)$": "identity-obj-proxy",
  },
  testMatch: ["**/?(*.)+(spec|test).[tj]s?(x)"],
  collectCoverage: true,
  coverageProvider: "v8",
  collectCoverageFrom: [
    "app/**/*.{ts,tsx}",
    "components/**/*.{ts,tsx}",
    "hooks/**/*.ts",
    "lib/**/*.ts",
    "contexts/**/*.tsx",
    "types/**/*.ts",
    "tests/**/*.{ts,tsx}",
  ],
  coveragePathIgnorePatterns: [
    "/node_modules/",
    "<rootDir>/app/(auth)/",
    "<rootDir>/app/(protected)/quiz/[topic]/QuizPageClient.tsx", // 主要为数据桥接薄层
  ],
};

module.exports = createJestConfig(customJestConfig);


