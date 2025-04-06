/// <reference types="cypress" />

import {MainPage} from "../pages/mainPage";

describe("Stock Detail", () => {
    const mainPage = new MainPage();

    beforeEach(() => {
        // GIVEN that I am in the main page
        mainPage.visit()
    })

    // 002-001
    it("row click open stock detail", () => {
        // WHEN I click a row
        mainPage.clickRow(0);
        // THEN the stock detail is shown
        mainPage.elements.stockDetail().should("exist");
    });

    // 002-002
    it("close button closes detail", () => {
        // GIVEN that I am in the stock detail
        mainPage.clickRow(0);
        // WHEN I click the Close button
        mainPage.closeDetail();
        // THEN the stock detail is removed
        mainPage.elements.stockDetail().should("not.exist");
    });
})