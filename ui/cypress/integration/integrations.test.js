import { selectors } from '../constants/IntegrationsPage';
import * as api from '../constants/apiEndpoints';
import withAuth from '../helpers/basicAuth';

describe('Integrations page', () => {
    withAuth();

    beforeEach(() => {
        cy.visit('/');
        cy.get(selectors.configure).click();
        cy.get(selectors.navLink).click({ force: true });
    });

    it('Plugin tiles should all be the same height', () => {
        let value = null;
        cy.get(selectors.plugins).each($el => {
            if (value) expect($el[0].clientHeight).to.equal(value);
            else value = $el[0].clientHeight;
        });
    });

    it('should have selected item in nav bar', () => {
        cy.get(selectors.configure).should('have.class', 'bg-primary-700');
    });

    it('should allow integration with Slack', () => {
        cy.get('div.ReactModalPortal').should('not.exist');

        cy.get(selectors.slackTile).click();
        cy.get('div.ReactModalPortal');
    });

    it('should add an integration with DockerHub', () => {
        cy.get(selectors.dockerRegistryTile).click();
        cy.get(selectors.buttons.delete).should('not.exist');
        cy.get(selectors.buttons.new).click();

        const name = `Docker Registry ${Math.random()
            .toString(36)
            .substring(7)}`;
        cy.get(selectors.dockerRegistryForm.nameInput).type(name);

        cy.get(
            `${selectors.dockerRegistryForm.typesSelect} .react-select__dropdown-indicator`
        ).click();
        cy.get('.react-select__menu-list > div:contains("Registry")').click();

        // test that validation error happens when form is incomplete
        cy.get(selectors.buttons.test).click();
        cy.get('div').contains('error');

        cy.get(selectors.dockerRegistryForm.endpointInput).type('registry-1.docker.io');

        cy.get(selectors.buttons.create).click();

        // delete the integration after to clean up
        cy.get(`.rt-tr:contains("${name}") .rt-td input[type="checkbox"]`).check();
        cy.get(selectors.buttons.delete).click({ force: true });
        cy.get(selectors.buttons.confirm).click();
        cy.get(`.rt-tr:contains("${name}")`).should('not.exist');
    });
});

describe('API Token Creation Flow', () => {
    withAuth();

    const randomTokenName = `Token${Math.random()
        .toString(36)
        .substring(7)}`;

    beforeEach(() => {
        cy.visit('/');
        cy.get(selectors.configure).click();
        cy.get(selectors.navLink).click({ force: true });
    });

    it('should pop up API Token Modal', () => {
        cy.get('div.ReactModalPortal').should('not.exist');

        cy.get(selectors.apiTokenTile).click();
        cy.get('div.ReactModalPortal');
    });

    it('should be able to generate an API token', () => {
        cy.get(selectors.apiTokenTile).click();
        cy.get(selectors.buttons.generate).click();
        cy.get(selectors.apiTokenForm.nameInput).type(randomTokenName);
        cy.get(`${selectors.apiTokenForm.roleSelect} .react-select__dropdown-indicator`).click();
        cy.get('.react-select__menu-list > div:contains("Admin")').click();
        cy.get(selectors.buttons.generate).click();
        cy.get(selectors.apiTokenBox);
        cy.get(selectors.apiTokenDetailsDiv).contains(`Name:${randomTokenName}`);
        cy.get(selectors.apiTokenDetailsDiv).contains('Role:Admin');
    });

    it('should show the generated API token in the table, and be clickable', () => {
        cy.get(selectors.apiTokenTile).click();
        cy.get(`.rt-tr:contains("${randomTokenName}")`).click();
        cy.get(selectors.apiTokenDetailsDiv).contains(`Name:${randomTokenName}`);
        cy.get(selectors.apiTokenDetailsDiv).contains('Role:Admin');
    });

    it('should be able to revoke the API token', () => {
        cy.get(selectors.apiTokenTile).click();
        cy.get(`.rt-tr:contains("${randomTokenName}") input`).check();
        cy.get(selectors.buttons.revoke).click({ force: true });
        cy.get(selectors.buttons.confirm).click();
        cy.get(`.rt-td:contains("${randomTokenName}")`).should('not.exist');
    });
});

// TODO: change this test suite to use the new Cluster pages, instead of an Integrations modal,
//       after ROX_SENSOR_AUTOUPGRADE=true becomes the default;
//       then, un-skip the test suite
xdescribe('Cluster Creation Flow', () => {
    withAuth();

    beforeEach(() => {
        cy.server();
        cy.fixture('clusters/single.json').as('singleCluster');
        cy.route('GET', api.clusters.list, '@singleCluster').as('clusters');
        cy.route('POST', api.clusters.zip, {}).as('download');
        cy.route('POST', api.clusters.list).as('addCluster');
        cy.visit('/');
        cy.get(selectors.configure).click();
        cy.get(selectors.navLink).click({ force: true });
        cy.wait('@clusters');
    });

    it('Should show a confirmation dialog when trying to delete clusters', () => {
        cy.get(selectors.kubernetesTile).click();
        cy.get(selectors.dialog).should('not.exist');
        cy.get(selectors.checkboxes).check();
        cy.get(selectors.buttons.delete).click({ force: true });
        cy.get(selectors.dialog);
    });

    it('Should show the remote cluster when clicking the Kubernetes tile', () => {
        cy.get(selectors.kubernetesTile).click();

        cy.get(selectors.clusters.k8sCluster0);
    });

    it.skip('Should be able to fill out the Kubernetes form, download config files and see cluster checked-in', () => {
        cy.get(selectors.kubernetesTile).click();

        cy.get(selectors.buttons.new).click();

        const clusterName = 'Kubernetes Cluster TestInstance';
        cy.get(selectors.clusterForm.nameInput).type(clusterName);
        // The image name should be pre-populated, so we don't type it in to test that the prepopulation works.
        // (The backend WILL error out if the image is empty.)
        cy.get(selectors.clusterForm.endpointInput)
            .clear()
            .type('central.stackrox:443');

        cy.get(selectors.buttons.next).click();
        cy.wait('@addCluster')
            .its('responseBody')
            .then(response => {
                const clusterId = response.cluster.id;

                cy.get(selectors.buttons.downloadYAML).click();
                cy.wait('@download');

                cy.get('div:contains("Waiting for the cluster to check in successfully...")');

                // make cluster to "check-in" by adding "lastContact"
                cy.route('GET', `${api.clusters.list}/${clusterId}`, {
                    cluster: {
                        id: clusterId,
                        status: {
                            lastContact: '2018-06-25T19:12:44.955289Z'
                        }
                    }
                }).as('getCluster');
                cy.wait('@getCluster');
                cy.get(
                    'div:contains("Success! The cluster has been recognized properly by StackRox. You may now save the configuration.")'
                );

                cy.get(selectors.buttons.closePanel).click();

                // clean up after the test by deleting the cluster
                cy.get(`.rt-tr:contains("${clusterName}") .rt-td input[type="checkbox"]`).check();
                cy.get(selectors.buttons.delete).click();
                cy.get(selectors.buttons.confirm).click();
                cy.get(`.rt-tr:contains("${clusterName}")`).should('not.exist');
            });
    });
});
