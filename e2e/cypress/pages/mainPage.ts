import {BasePage} from "./basePage";
import {BACKEND_GET_STOCKS_ENDPOINT} from "../support/e2e";

/**
 * Represents MyStocks main page
 */
export class MainPage extends BasePage {
    protected readonly URL_RESOURCE = "/";

    elements = {
        title: () => cy.get("h1"),
        table: () => cy.get("table"),
        tableHeaders: () => cy.get("thead tr th", { timeout: 10000 }),
        tableRows: () => cy.get("table tbody tr"),
        closeDetailButton: () => cy.get("button")
    }

    visit() {
        cy.intercept(
            BACKEND_GET_STOCKS_ENDPOINT.method,
            BACKEND_GET_STOCKS_ENDPOINT.resource
        ).as("getStocks");

        super.visit()

        cy.wait("@getStocks").its("response.body").should("not.be.empty");
    }

    /**
     * Clicks the n-th row in the table.
     * @param row row number. 0-indexed
     */
    clickRow(row: number) {
        this.elements.tableRows()
            .eq(row)
            .should("exist")
            .click()
    }

    /**
     * Closes the current detail to display the main menu again.
     */
    closeDetail() {
        this.elements.closeDetailButton()
            .should("exist")
            .click()
    }

    /**
     * Returns the extracted table headers into an interable
     */
    getTableHeaders(): Cypress.Chainable<string[]> {
        return this.elements.tableHeaders()
            .should("have.length.at.least", 1)
            .then(($headers) => {
                return Cypress._.map($headers, (header) => header.innerText.trim());
            });
    }
}