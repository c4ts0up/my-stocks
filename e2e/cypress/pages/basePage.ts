
export abstract class BasePage {
    protected abstract URL_RESOURCE: string;

    /**
     * Visits the resource
     */
    visit() {
        cy.visit(this.URL_RESOURCE)
    }
}