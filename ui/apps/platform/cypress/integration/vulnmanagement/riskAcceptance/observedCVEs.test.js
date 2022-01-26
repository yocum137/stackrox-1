import * as api from '../../../constants/apiEndpoints';
import withAuth from '../../../helpers/basicAuth';
import { hasFeatureFlag } from '../../../helpers/features';

const imageEntityPage =
    '/main/vulnerability-management/image/sha256:5469b2315904f5f720034495c3938a4d6f058ec468ce4eca0b1a9291c616c494';

function aliasImageVulnerabilitiesQuery(req, vulnsQuery, alias) {
    const { body } = req;
    const matchesQuery = body?.variables?.vulnsQuery === vulnsQuery;
    if (matchesQuery) {
        req.alias = alias;
    }
}

function submitDeferralForm() {
    cy.get('input[value="2 weeks"]').check();
    cy.get('input[value="All tags within image"]').check();
    cy.get('textarea[id="comment"]').type('Defer for 2 weeks');
    cy.get('button:contains("Request approval")').click();
    cy.wait('@deferVulnerability');
}

function submitFalsePositiveForm() {
    cy.get('input[value="All tags within image"]').check();
    cy.get('textarea[id="comment"]').type('Marked as false positive');
    cy.get('button:contains("Request approval")').click();
    cy.wait('@markVulnerabilityFalsePositive');
}

function selectBulkAction(actionText) {
    cy.get('button:contains("Bulk actions")').click();
    cy.get(`li[role="menuitem"] button:contains("${actionText}")`);
}

function getRowActionsByRowIndex(rowIndex) {
    return cy.get(
        `table[aria-label="Observed CVEs Table"] tbody tr:nth(${rowIndex}) button[aria-label="Actions"]`
    );
}

function getCheckboxByRowIndex(rowIndex) {
    return cy.get(
        `table[aria-label="Observed CVEs Table"] tbody tr:nth(${rowIndex}) input[type="checkbox"]`
    );
}

function getPendingApprovalIconByRowIndex(rowIndex) {
    return cy.get(
        `table[aria-label="Observed CVEs Table"] tbody tr:nth(${rowIndex}) svg[aria-label="Pending approval icon"]`
    );
}

function getRowActionItem(actionText) {
    return cy.get(`li[role="menuitem"] button:contains("${actionText}")`);
}

describe('Vulnmanagement Risk Acceptance', () => {
    before(function beforeHook() {
        if (!hasFeatureFlag('ROX_VULN_RISK_MANAGEMENT')) {
            this.skip();
        }
    });

    withAuth();

    describe('Observed CVEs', () => {
        beforeEach(() => {
            cy.intercept('POST', api.riskAcceptance.getImageVulnerabilities, (req) => {
                aliasImageVulnerabilitiesQuery(
                    req,
                    'Vulnerability State:OBSERVED',
                    'getObservedCVEs'
                );
            });
            cy.intercept('POST', api.riskAcceptance.getImageVulnerabilities, (req) => {
                aliasImageVulnerabilitiesQuery(
                    req,
                    'Vulnerability State:DEFERRED',
                    'getDeferredCVEs'
                );
            });
            cy.intercept('POST', api.riskAcceptance.getImageVulnerabilities, (req) => {
                aliasImageVulnerabilitiesQuery(
                    req,
                    'Vulnerability State:FALSE_POSITIVE',
                    'getFalsePositiveCVEs'
                );
            });
            cy.intercept('POST', api.riskAcceptance.deferVulnerability).as('deferVulnerability');
            cy.intercept('POST', api.riskAcceptance.markVulnerabilityFalsePositive).as(
                'markVulnerabilityFalsePositive'
            );
        });

        it('should be able to defer a CVE', () => {
            cy.visit(imageEntityPage);
            cy.wait('@getObservedCVEs');

            getRowActionsByRowIndex(0).click();
            getRowActionItem('Defer CVE').click();
            submitDeferralForm();

            getPendingApprovalIconByRowIndex(0);
        });

        it('should be able to mark a CVE as false positive', () => {
            cy.visit(imageEntityPage);
            cy.wait('@getObservedCVEs');

            getRowActionsByRowIndex(1).click();
            getRowActionItem('Mark as False Positive').click();
            submitFalsePositiveForm();

            getPendingApprovalIconByRowIndex(1);
        });

        it('should be able to defer CVEs in bulk', () => {
            cy.visit(imageEntityPage);
            cy.wait('@getObservedCVEs');

            getCheckboxByRowIndex(2).check({ force: true });
            getCheckboxByRowIndex(3).check({ force: true });
            selectBulkAction('Defer CVE (2)');
            submitDeferralForm();

            getPendingApprovalIconByRowIndex(2);
            getPendingApprovalIconByRowIndex(3);
        });

        it('should be able to mark CVEs as false positive in bulk', () => {
            cy.visit(imageEntityPage);
            cy.wait('@getObservedCVEs');

            getCheckboxByRowIndex(4).check({ force: true });
            getCheckboxByRowIndex(5).check({ force: true });
            selectBulkAction('Mark false positive (2)');
            submitFalsePositiveForm();

            getPendingApprovalIconByRowIndex(4);
            getPendingApprovalIconByRowIndex(5);
        });
    });
});
