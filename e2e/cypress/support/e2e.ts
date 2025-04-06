// ***********************************************************
// This example support/e2e.js is processed and
// loaded automatically before your test files.
//
// This is a great place to put global configuration and
// behavior that modifies Cypress.
//
// You can change the location of this file or turn off
// automatically serving support files with the
// 'supportFile' configuration option.
//
// You can read more here:
// https://on.cypress.io/configuration
// ***********************************************************

// Import commands.js using ES2015 syntax:
import './commands'

/// <reference types="cypress" />
export const BACKEND_URL = "http://localhost:8080";
export const BACKEND_GET_STOCKS_ENDPOINT = {
    method: "GET",
    resource: BACKEND_URL + "/stocks"
};