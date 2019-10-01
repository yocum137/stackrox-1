import { selectors, url } from '../constants/AccessPage';
import withAuth from '../helpers/basicAuth';

describe('Access Control Page', () => {
    describe('Auth Provider Rules', () => {
        withAuth();

        beforeEach(() => {
            cy.visit(url);
            cy.get(selectors.tabs.authProviders).click();
        });

        it('should open the new auth provider panel', () => {
            cy.get(selectors.authProviders.addProvider).select(selectors.authProviders.newAuth0);
            cy.get(selectors.authProviders.newAuthProviderPanel).contains(
                'Create New auth0 Auth Provider'
            );
        });
    });

    describe('Roles and Permissions', () => {
        withAuth();

        beforeEach(() => {
            cy.visit(url);
            cy.get(selectors.tabs.roles).click();
        });

        const selectRole = roleName => {
            cy.get(selectors.roles)
                .contains(roleName)
                .click();
        };

        const createRole = roleName => {
            cy.get(selectors.addNewRoleButton).click();
            cy.get(selectors.permissionsPanelHeader).contains('Create New Role');
            cy.get(selectors.input.roleName).type(roleName);
            cy.get(selectors.saveButton).click();
            cy.get(selectors.roles)
                .contains(roleName)
                .then($role => {
                    cy.get(selectors.permissionsPanelHeader).contains(
                        `"${$role.text()}" Permissions`
                    );
                });
        };

        it('should have the default roles', () => {
            cy.get(selectors.roles).contains('Admin');
            cy.get(selectors.roles).contains('Analyst');
            cy.get(selectors.roles).contains('Continuous Integration');
            cy.get(selectors.roles).contains('None');
            cy.get(selectors.roles).contains('Sensor Creator');
        });

        it('should automatically select the first role', () => {
            cy.get(selectors.roles)
                .eq(0)
                .then($role => {
                    cy.get(selectors.permissionsPanelHeader).contains(
                        `"${$role.text()}" Permissions`
                    );
                });
        });

        it('should not be able to edit default roles', () => {
            selectRole('Admin');
            cy.get(selectors.editButton).should('not.exist');
            selectRole('Analyst');
            cy.get(selectors.editButton).should('not.exist');
            selectRole('Continuous Integration');
            cy.get(selectors.editButton).should('not.exist');
            selectRole('None');
            cy.get(selectors.editButton).should('not.exist');
            selectRole('Sensor Creator');
            cy.get(selectors.editButton).should('not.exist');
        });

        it('should create a new role', () => {
            const newRoleName = `Role-${new Date().getTime()}`;
            createRole(newRoleName);
        });

        it('should not be able to edit an existing role name', () => {
            const newRoleName = `Role-${new Date().getTime()}`;
            createRole(newRoleName);
            cy.get(selectors.editButton).click();
            cy.get(selectors.input.roleName).then($input => {
                cy.wrap($input).should('have.attr', 'disabled');
            });
        });
    });
});
