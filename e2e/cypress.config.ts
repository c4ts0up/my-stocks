import { defineConfig } from "cypress";

export default defineConfig({
  e2e: {
    baseUrl: "http://localhost:5173", // Set the correct base URL
    specPattern: "cypress/e2e/**/*.cy.{js,jsx,ts,tsx}", // Ensure .ts and .tsx are included
  },
});
