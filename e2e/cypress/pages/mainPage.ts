/**
 * Represents MyStocks main page
 */
export class MainPage {

    elements = {
        title: () => cy.get("h1"),
        table: () => cy.get("table"),
        tableHeaders: () => cy.get("thead tr th"),
        tableRows: () => cy.get("table tbody tr"),
        closeDetailButton: () => cy.get("button")
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
            .then(($headers) => [...$headers].map(header => header.innerText.trim()));
    }
}