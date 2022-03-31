import { selectors } from '../../../constants/PoliciesPagePatternFly';
import withAuth from '../../../helpers/basicAuth';
import {
    visitPolicies,
    doPolicyRowAction,
    editFirstPolicyFromTable,
    cloneFirstPolicyFromTable,
    goToStep2,
    goToStep4,
} from '../../../helpers/policiesPatternFly';

function goToPoliciesAndCloneToStep4() {
    visitPolicies();
    cloneFirstPolicyFromTable();
    goToStep4();
}

function goToPoliciesAndCloneToStep2() {
    visitPolicies();
    cloneFirstPolicyFromTable();
    goToStep2();
}

describe('Policy wizard, Step 4 Policy Scope', () => {
    withAuth();

    describe('Inclusion card', () => {
        it('should allow the user to add and delete a inclusion scope card', () => {
            goToPoliciesAndCloneToStep4();

            cy.get(selectors.step4.inclusionScope.addBtn).click();
            cy.get(selectors.step4.inclusionScope.addBtn).click();
            cy.get(selectors.step4.inclusionScope.addBtn).click();
            cy.get(selectors.step4.inclusionScope.cards).should('have.length', 3);
            cy.get(selectors.step4.inclusionScope.deleteBtn).first().click();
            cy.get(selectors.step4.inclusionScope.cards).should('have.length', 2);
        });

        it('should have cluster select and respect changed value', () => {
            goToPoliciesAndCloneToStep4();

            cy.get(selectors.step4.inclusionScope.addBtn).click();
            cy.get(selectors.step4.inclusionScope.clusterSelect).should('have.value', '');
            cy.get(selectors.step4.inclusionScope.clusterSelect).click();
            cy.get(selectors.step4.inclusionScope.clusterSelectOption)
                .first()
                .then((option) => {
                    cy.wrap(option).click();
                    cy.get(selectors.step4.inclusionScope.clusterSelect).contains(option.text());
                });
        });

        it('should populate cluster select with existing value', () => {
            goToPoliciesAndCloneToStep4();
        });

        it('should have namespace input and respect changed value', () => {
            goToPoliciesAndCloneToStep4();

            cy.get(selectors.step4.inclusionScope.addBtn).click();
            cy.get(selectors.step4.inclusionScope.namespaceInput).should('have.value', '');
            cy.get(selectors.step4.inclusionScope.namespaceInput).type('test');
            cy.get(selectors.step4.inclusionScope.namespaceInput).should('have.value', 'test');
        });

        it('should populate namespace input with existing value', () => {
            goToPoliciesAndCloneToStep4();
        });

        it('should not have deployment select', () => {
            goToPoliciesAndCloneToStep4();

            cy.get(selectors.step4.inclusionScope.addBtn).click();
            cy.get(selectors.step4.inclusionScope.deploymentSelect).should('not.exist');
        });

        it('should have label key/value input and respect changed value', () => {
            goToPoliciesAndCloneToStep4();

            cy.get(selectors.step4.inclusionScope.addBtn).click();
            cy.get(selectors.step4.inclusionScope.labelKeyInput).should('have.value', '');
            cy.get(selectors.step4.inclusionScope.labelKeyInput).type('hello');
            cy.get(selectors.step4.inclusionScope.labelKeyInput).should('have.value', 'hello');
            cy.get(selectors.step4.inclusionScope.labelValueInput).should('have.value', '');
            cy.get(selectors.step4.inclusionScope.labelValueInput).type('world');
            cy.get(selectors.step4.inclusionScope.labelValueInput).should('have.value', 'world');
        });

        it('should populate label key/value input with existing values', () => {
            goToPoliciesAndCloneToStep4();
        });

        it('should disable label key/value inputs if Event Source is Audit Log', () => {
            goToPoliciesAndCloneToStep2();
        });
    });

    describe('Exclusion card', () => {
        it('should allow the user to add and delete a exclusion scope card', () => {
            goToPoliciesAndCloneToStep4();

            cy.get(selectors.step4.exclusionScope.addBtn).click();
            cy.get(selectors.step4.exclusionScope.addBtn).click();
            cy.get(selectors.step4.exclusionScope.addBtn).click();
            cy.get(selectors.step4.exclusionScope.cards).should('have.length', 3);
            cy.get(selectors.step4.exclusionScope.deleteBtn).first().click();
            cy.get(selectors.step4.exclusionScope.cards).should('have.length', 2);
        });

        it('should have cluster select and respect changed value', () => {
            goToPoliciesAndCloneToStep4();

            cy.get(selectors.step4.exclusionScope.addBtn).click();
            cy.get(selectors.step4.exclusionScope.clusterSelect).should('have.value', '');
            cy.get(selectors.step4.exclusionScope.clusterSelect).click();
            cy.get(selectors.step4.exclusionScope.clusterSelectOption)
                .first()
                .then((option) => {
                    cy.wrap(option).click();
                    cy.get(selectors.step4.exclusionScope.clusterSelect).contains(option.text());
                });
        });

        it('should populate cluster select with existing value', () => {
            goToPoliciesAndCloneToStep4();
        });

        it('should have namespace input and respect changed value', () => {
            goToPoliciesAndCloneToStep4();

            cy.get(selectors.step4.exclusionScope.addBtn).click();
            cy.get(selectors.step4.exclusionScope.namespaceInput).should('have.value', '');
            cy.get(selectors.step4.exclusionScope.namespaceInput).type('test');
            cy.get(selectors.step4.exclusionScope.namespaceInput).should('have.value', 'test');
        });

        it('should populate namespace input with existing value', () => {
            goToPoliciesAndCloneToStep4();
        });

        it('should have deployment select and respect changed value', () => {
            goToPoliciesAndCloneToStep4();

            cy.get(selectors.step4.exclusionScope.addBtn).click();
            cy.get(selectors.step4.exclusionScope.deploymentSelect).should('have.value', '');
            cy.get(selectors.step4.exclusionScope.deploymentSelect).click();
            cy.get(selectors.step4.exclusionScope.deploymentSelectOption)
                .first()
                .then((option) => {
                    cy.wrap(option).click();
                    cy.get(selectors.step4.exclusionScope.deploymentSelect).contains(option.text());
                });
        });

        it('should populate deployment select with existing value', () => {
            goToPoliciesAndCloneToStep4();
        });

        it('should have label key/value input and respect changed value', () => {
            goToPoliciesAndCloneToStep4();

            cy.get(selectors.step4.exclusionScope.addBtn).click();
            cy.get(selectors.step4.exclusionScope.labelKeyInput).should('have.value', '');
            cy.get(selectors.step4.exclusionScope.labelKeyInput).type('hello');
            cy.get(selectors.step4.exclusionScope.labelKeyInput).should('have.value', 'hello');
            cy.get(selectors.step4.exclusionScope.labelValueInput).should('have.value', '');
            cy.get(selectors.step4.exclusionScope.labelValueInput).type('world');
            cy.get(selectors.step4.exclusionScope.labelValueInput).should('have.value', 'world');
        });

        it('should populate label key/value input with existing values', () => {
            goToPoliciesAndCloneToStep4();
        });

        it('should disable label key/value inputs and deployment dropdown if Event Source is Audit Log', () => {
            goToPoliciesAndCloneToStep2();
        });
    });

    describe('Exclude images', () => {
        it('should populate dropdown with existing values', () => {
            goToPoliciesAndCloneToStep4();
        });

        it.only('exclude images dropdown should be enabled and allow the user to add/delete/create image exclusions if lifecycle stage includes BUILD', () => {
            goToPoliciesAndCloneToStep2();

            cy.get(selectors.step2.lifecycleStage.buildCheckbox).then((buildCheckbox) => {
                if (buildCheckbox.prop('checked') === false) {
                    cy.wrap(buildCheckbox).click();
                }
            });
            goToStep4();
            cy.get(selectors.step4.excludeImages.multiselect).should(
                'not.have.class',
                'pf-m-disabled'
            );
            cy.get(selectors.step4.excludeImages.multiselect).click();
            cy.get(selectors.step4.excludeImages.multiselectOption)
                .first()
                .then((option) => {
                    cy.wrap(option).click();
                    cy.get(selectors.step4.excludeImages.multiselect).contains(option.text());
                    cy.get(selectors.step4.excludeImages.multiselectOptionDeleteBtn)
                        .first()
                        .click();
                    cy.get(selectors.step4.excludeImages.multiselect).should(
                        'not.contain',
                        option.text()
                    );
                });
            cy.get(selectors.step4.excludeImages.multiselectInput).type('hello{enter}');
            cy.get(selectors.step4.excludeImages.multiselect).should('contain', 'hello');
        });

        it('exclude images dropdown should be cleared and disabled if lifecycle stage does not include BUILD', () => {
            goToPoliciesAndCloneToStep2();

            cy.get(selectors.step2.lifecycleStage.buildCheckbox).then((buildCheckbox) => {
                if (buildCheckbox.prop('checked') === true) {
                    cy.wrap(buildCheckbox).click();
                }
            });
            goToStep4();
            cy.get(selectors.step4.excludeImages.multiselect).should('have.class', 'pf-m-disabled');
        });
    });
});
