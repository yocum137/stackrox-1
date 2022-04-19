import React, { useState, useEffect } from 'react';
import { useHistory } from 'react-router-dom';
import {
    PageSection,
    Bullseye,
    Alert,
    Spinner,
    AlertGroup,
    AlertActionCloseButton,
    AlertVariant,
} from '@patternfly/react-core';
import pluralize from 'pluralize';
import { useMutation, useQuery, useQueryClient } from 'react-query';

import { policiesBasePath } from 'routePaths';
import {
    getPolicies,
    reassessPolicies,
    deletePolicies,
    exportPolicies,
    updatePoliciesDisabledState,
} from 'services/PoliciesService';
import { fetchNotifierIntegrations } from 'services/NotifierIntegrationsService';
import useToasts, { Toast } from 'hooks/patternfly/useToasts';
import { getSearchOptionsForCategory } from 'services/SearchService';
import { NotifierIntegration } from 'types/notifier.proto';
import { SearchFilter } from 'types/search';
import { getAxiosErrorMessage } from 'utils/responseErrorUtils';
import { getRequestQueryStringForSearchFilter } from 'utils/searchUtils';

import ImportPolicyJSONModal from '../Modal/ImportPolicyJSONModal';
import PoliciesTable from './PoliciesTable';

type PoliciesTablePageProps = {
    hasWriteAccessForPolicy: boolean;
    handleChangeSearchFilter: (searchFilter: SearchFilter) => void;
    searchFilter?: SearchFilter;
};

type DisableStateVars = { ids: string[]; isDisabled: boolean };

function PoliciesTablePage({
    hasWriteAccessForPolicy,
    handleChangeSearchFilter,
    searchFilter,
}: PoliciesTablePageProps): React.ReactElement {
    const query = searchFilter ? getRequestQueryStringForSearchFilter(searchFilter) : '';
    const history = useHistory();

    const [notifiers, setNotifiers] = useState<NotifierIntegration[]>([]);
    const { toasts, addToast, removeToast } = useToasts();

    const [searchOptions, setSearchOptions] = useState<string[]>([]);

    const [isImportModalOpen, setIsImportModalOpen] = useState(false);

    const queryClient = useQueryClient();

    // Note that this will automatically refetch when
    // - the window is refocused
    // - any search parameter changes (is this right??? Why doesn't it cache?)
    const {
        data: policies,
        error: policiesError,
        isLoading,
        refetch,
    } = useQuery(['policies', query], () => getPolicies(query));

    const policyDisabledMutation = useMutation(
        ({ ids, isDisabled }: DisableStateVars) => updatePoliciesDisabledState(ids, isDisabled),
        {
            onSuccess: (_data, { ids, isDisabled }) => {
                const policyText = pluralize('policy', ids.length);
                const stateText = isDisabled ? 'disabled' : 'enabled';
                addToast(`Successfully ${stateText} ${policyText}`, 'success');
            },
            onError: (stateError: Error, { ids, isDisabled }) => {
                const policyText = pluralize('policy', ids.length);
                const stateText = isDisabled ? 'disable' : 'enable';
                addToast(`Could not ${stateText} the ${policyText}`, 'danger', stateError.message);
            },
            onSettled: () => queryClient.invalidateQueries('policies'),
        }
    );

    const errorMessage = policiesError ? getAxiosErrorMessage(policiesError) : '';

    function onClickCreatePolicy() {
        history.push(`${policiesBasePath}/?action=create`);
    }

    function onClickImportPolicy() {
        setIsImportModalOpen(true);
    }

    function onClickReassessPolicies() {
        return reassessPolicies()
            .then(() => {
                addToast('Successfully reassessed policies', 'success');
            })
            .catch(({ response }) => {
                addToast('Could not reassess policies', 'danger', response.data.message);
            });
    }

    function fetchPolicies(query: string) {
        console.log('TODO');
        refetch();
    }

    function deletePoliciesHandler(ids: string[]): Promise<void> {
        const policyText = pluralize('policy', ids.length);
        return deletePolicies(ids)
            .then(() => {
                fetchPolicies(query);
                addToast(`Successfully deleted ${policyText}`, 'success');
            })
            .catch(({ response }) => {
                addToast(`Could not delete ${policyText}`, 'danger', response.data.message);
            });
    }

    function exportPoliciesHandler(ids: string[], onClearAll?: () => void) {
        const policyText = pluralize('policy', ids.length);
        exportPolicies(ids)
            .then(() => {
                addToast(`Successfully exported ${policyText}`, 'success');
                if (onClearAll) {
                    onClearAll();
                }
            })
            .catch((error) => {
                const message = getAxiosErrorMessage(error);
                addToast(`Could not export the ${policyText}`, 'danger', message);
            });
    }

    function enablePoliciesHandler(ids: string[]) {
        policyDisabledMutation.mutate({ ids, isDisabled: false });
    }

    function disablePoliciesHandler(ids: string[]) {
        policyDisabledMutation.mutate({ ids, isDisabled: true });
    }

    useEffect(() => {
        fetchNotifierIntegrations()
            .then((data) => {
                setNotifiers(data as NotifierIntegration[]);
            })
            .catch(() => {
                setNotifiers([]);
            });
    }, []);

    useEffect(() => {
        getSearchOptionsForCategory('POLICIES')
            .then((options) => {
                setSearchOptions(options);
            })
            .catch(() => {
                // TODO
            });
    }, []);

    if (isLoading) {
        return (
            <PageSection variant="light" isFilled id="policies-table-loading">
                <Bullseye>
                    <Spinner isSVG />
                </Bullseye>
            </PageSection>
        );
    }

    return (
        <>
            {errorMessage ? (
                <PageSection variant="light" isFilled id="policies-table-error">
                    <Bullseye>
                        <Alert variant="danger" title={errorMessage} />
                    </Bullseye>
                </PageSection>
            ) : (
                <PoliciesTable
                    notifiers={notifiers}
                    policies={policies}
                    fetchPoliciesHandler={() => fetchPolicies(query)}
                    addToast={addToast}
                    hasWriteAccessForPolicy={hasWriteAccessForPolicy}
                    deletePoliciesHandler={deletePoliciesHandler}
                    exportPoliciesHandler={exportPoliciesHandler}
                    enablePoliciesHandler={enablePoliciesHandler}
                    disablePoliciesHandler={disablePoliciesHandler}
                    handleChangeSearchFilter={handleChangeSearchFilter}
                    onClickCreatePolicy={onClickCreatePolicy}
                    onClickImportPolicy={onClickImportPolicy}
                    onClickReassessPolicies={onClickReassessPolicies}
                    searchFilter={searchFilter}
                    searchOptions={searchOptions}
                />
            )}
            <ImportPolicyJSONModal
                isOpen={isImportModalOpen}
                cancelModal={() => {
                    setIsImportModalOpen(false);
                }}
                fetchPoliciesWithQuery={() => fetchPolicies(query)}
            />
            <AlertGroup isToast isLiveRegion>
                {toasts.map(({ key, variant, title, children }: Toast) => (
                    <Alert
                        variant={AlertVariant[variant]}
                        title={title}
                        timeout={4000}
                        onTimeout={() => removeToast(key)}
                        actionClose={
                            <AlertActionCloseButton
                                title={title}
                                variantLabel={`${variant} alert`}
                                onClose={() => removeToast(key)}
                            />
                        }
                        key={key}
                    >
                        {children}
                    </Alert>
                ))}
            </AlertGroup>
        </>
    );
}

export default PoliciesTablePage;
