/// <reference types="cypress" />

import {MainPage} from "../pages/mainPage";

describe("Main page title", () => {
    const mainPage = new MainPage()

    beforeEach(() => {
        // GIVEN that I am in the main page
        mainPage.visit()
    });

    it("title reads MyStocks", () => {
        // THEN the title should read "MyStocks"
        mainPage.elements.title()
            .should("be.visible")
            .and("have.text", "MyStocks");
    });
});