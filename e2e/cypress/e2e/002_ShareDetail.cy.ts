/// <reference types="cypress" />

import {MainPage} from "../pages/mainPage";

describe("Stock Detail", () => {
    const mainPage = new MainPage();

    beforeEach(() => {
        // GIVEN that I am in the main page
        mainPage.visit()
    })

    // 002-001
    it("row click opens stock detail", () => {
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

    // 002-003
    it("stock details are complete", () => {
        // GIVEN that I am in the stock detail
        mainPage.clickRow(0);
        // THEN the stock details are shown
        mainPage.elements.stockDetail().should("be.visible");
        mainPage.elements.stockDetail().find('[data-testid="ticker"]').should("exist");
        mainPage.elements.stockDetail().find('[data-testid="companyName"]').should("exist");
        mainPage.elements.stockDetail().find('[data-testid="currentPrice"]').should("exist");
        mainPage.elements.stockDetail().find('[data-testid="recommendationTag"]').should("exist");
        mainPage.elements.stockDetail().find('[data-testid="ratings"]').should("exist");
    });
})