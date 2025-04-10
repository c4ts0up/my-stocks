/// <reference types="cypress" />

import {MainPage} from "../pages/mainPage";

describe("Main page", () => {
    const mainPage = new MainPage();

    beforeEach(() => {
        // GIVEN that I am in the main page
        mainPage.visit()
    });

    // 001-001
    it("title reads MyStocks", () => {
        // THEN the title should read "MyStocks"
        mainPage.elements.title()
            .should("be.visible")
            .and("have.text", "MyStocks");
    });

    // 001-002
    it("table headers check", () => {
        const expectedHeaders = ["Ticker", "Company Name", "Current Price", "Recommendation"];

        // THEN the table has the headers Ticker, Company Name, Current Price and Recommendation
        mainPage.getTableHeaders().then((actualHeaders) => {
            expect(actualHeaders).to.deep.equal(expectedHeaders);
        });
    })

    // 001-003
    it("all rows have non-empty cells", () => {
        cy.get("table tbody tr").each(($row) => {
            cy.wrap($row)
                .find('[data-testid="ticker-cell"]')
                .should('not.be.empty');

            cy.wrap($row)
                .find('[data-testid="companyName-cell"]')
                .should('not.be.empty');

            cy.wrap($row)
                .find('[data-testid="currentPrice-cell"]')
                .should('not.be.empty');
        });
    })
});