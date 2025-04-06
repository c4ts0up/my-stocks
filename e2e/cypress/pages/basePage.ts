
export abstract class BasePage {
    protected abstract URL_RESOURCE: string;

    /**
     * Visits the resource
     */
    visit() {
        cy.log(`Visiting ${this.URL_RESOURCE}`)
        cy.visit(this.URL_RESOURCE)
    }
}