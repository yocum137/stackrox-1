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
    Divider,
} from '@patternfly/react-core';
import pluralize from 'pluralize';

import {
    getPolicies,
    getPoliciesCount,
    deletePolicies,
    exportPolicies,
    updatePoliciesDisabledState,
} from 'services/PoliciesService';
import { fetchNotifierIntegrations } from 'services/NotifierIntegrationsService';
import useToasts, { Toast } from 'hooks/patternfly/useToasts';
import { getSearchOptionsForCategory } from 'services/SearchService';
import { ListPolicy } from 'types/policy.proto';
import { NotifierIntegration } from 'types/notifier.proto';
import { SearchFilter } from 'types/search';
import { getAxiosErrorMessage } from 'utils/responseErrorUtils';
import useTableSort from 'hooks/useTableSort';
import useTableSelection from 'hooks/useTableSelection';
import { policiesBasePathPatternFly } from 'routePaths';

import { getSearchStringForFilter, getRequestQueryStringForSearchFilter } from '../policies.utils';
import PoliciesTable from './PoliciesTable';
import PoliciesTableHeader from './PoliciesTableHeader';

const columns = [
    {
        Header: 'Policy',
        accessor: 'name',
        sortField: 'SORT_Policy',
        width: 20 as const,
    },
    {
        Header: 'Description',
        accessor: 'description',
        width: 40 as const,
    },
    {
        Header: 'Status',
        accessor: 'disabled',
        sortField: 'Disabled',
        width: 15 as const,
    },
    {
        Header: 'Notifiers',
        accessor: 'notifiers',
    },
    {
        Header: 'Severity',
        accessor: 'severity',
        sortField: 'Severity',
    },
    {
        Header: 'Lifecycle',
        accessor: 'lifecycleStages',
        sortField: 'Lifecycle Stage',
    },
];

type PoliciesTablePageProps = {
    hasWriteAccessForPolicy: boolean;
    searchFilter?: SearchFilter;
};

function PoliciesTablePage({
    hasWriteAccessForPolicy,
    searchFilter,
}: PoliciesTablePageProps): React.ReactElement {
    const history = useHistory();

    const [notifiers, setNotifiers] = useState<NotifierIntegration[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [policies, setPolicies] = useState<ListPolicy[]>([]);
    const [errorMessage, setErrorMessage] = useState('');
    const { toasts, addToast, removeToast } = useToasts();

    const [searchOptions, setSearchOptions] = useState<string[]>([]);

    // Handle changes to applied search options.
    const [isViewFiltered, setIsViewFiltered] = useState(false);

    // Handle changes in the current table page.
    const [currentPage, setCurrentPage] = useState(1);
    const [perPage, setPerPage] = useState(50);
    const [policiesCount, setPoliciesCount] = useState(0);
    const defaultSort = {
        field: 'Policy',
        reversed: true,
    };
    const {
        activeSortIndex,
        setActiveSortIndex,
        activeSortDirection,
        setActiveSortDirection,
        sortOption,
    } = useTableSort(columns, defaultSort);

    // const hasExecutableFilter = searchOptions.length && !searchOptions[searchOptions.length - 1];
    // const hasNoFilter = !searchOptions.length;

    // if (hasExecutableFilter && !isViewFiltered) {
    //     setIsViewFiltered(true);
    //     setCurrentPage(1);
    // } else if (hasNoFilter && isViewFiltered) {
    //     setIsViewFiltered(false);
    //     setCurrentPage(1);
    // }

    // Handle selected rows in table
    const {
        selected,
        allRowsSelected,
        onSelect,
        onSelectAll,
        // onClearAll,
        getSelectedIds,
    } = useTableSelection(policies);

    function fetchPolicies(query: string) {
        setIsLoading(true);
        getPolicies(query, sortOption, currentPage - 1, perPage)
            .then((data) => {
                setPolicies(data);
                setErrorMessage('');
            })
            .catch((error) => {
                setPolicies([]);
                setErrorMessage(getAxiosErrorMessage(error));
            })
            .finally(() => setIsLoading(false));
        getPoliciesCount(query)
            .then((data) => {
                setPoliciesCount(data.length);
                setErrorMessage('');
            })
            .catch((error) => {
                setPoliciesCount(0);
                setErrorMessage(getAxiosErrorMessage(error));
            })
            .finally(() => setIsLoading(false));
    }

    const query = searchFilter ? getRequestQueryStringForSearchFilter(searchFilter) : '';

    function fetchPoliciesWithQuery() {
        fetchPolicies(query);
    }

    function deletePoliciesHandler(ids: string[]) {
        const policyText = pluralize('policy', ids.length);
        deletePolicies(ids)
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
        const policyText = pluralize('policy', ids.length);
        updatePoliciesDisabledState(ids, false)
            .then(() => {
                fetchPolicies(query);
                addToast(`Successfully enabled ${policyText}`, 'success');
            })
            .catch(({ response }) => {
                addToast(`Could not enable the ${policyText}`, 'danger', response.data.message);
            });
    }

    function disablePoliciesHandler(ids: string[]) {
        const policyText = pluralize('policy', ids.length);
        updatePoliciesDisabledState(ids, true)
            .then(() => {
                fetchPolicies(query);
                addToast(`Successfully disabled ${policyText}`, 'success');
            })
            .catch(({ response }) => {
                addToast(`Could not disable the ${policyText}`, 'danger', response.data.message);
            });
    }

    function handleChangeSearchFilter(changedSearchFilter: SearchFilter) {
        // Browser history has only the most recent search filter.
        history.replace({
            pathname: policiesBasePathPatternFly,
            search: getSearchStringForFilter(changedSearchFilter),
        });
    }

    const selectedIds = getSelectedIds();
    const selectedPolicies = policies.filter(({ id }) => selectedIds.includes(id));

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

    useEffect(() => {
        fetchPoliciesWithQuery();
    }, [query, sortOption, currentPage, perPage]);

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
                <>
                    <PoliciesTableHeader
                        hasWriteAccessForPolicy={hasWriteAccessForPolicy}
                        deletePoliciesHandler={deletePoliciesHandler}
                        enablePoliciesHandler={enablePoliciesHandler}
                        disablePoliciesHandler={disablePoliciesHandler}
                        fetchPoliciesWithQuery={fetchPoliciesWithQuery}
                        handleChangeSearchFilter={handleChangeSearchFilter}
                        searchFilter={searchFilter}
                        searchOptions={searchOptions}
                        policiesCount={policiesCount}
                        perPage={perPage}
                        setPerPage={setPerPage}
                        currentPage={currentPage}
                        setCurrentPage={setCurrentPage}
                        selectedIds={selectedIds}
                        selectedPolicies={selectedPolicies}
                    />
                    <Divider component="div" />
                    <PoliciesTable
                        notifiers={notifiers}
                        policies={policies}
                        hasWriteAccessForPolicy={hasWriteAccessForPolicy}
                        deletePoliciesHandler={deletePoliciesHandler}
                        exportPoliciesHandler={exportPoliciesHandler}
                        enablePoliciesHandler={enablePoliciesHandler}
                        disablePoliciesHandler={disablePoliciesHandler}
                        activeSortIndex={activeSortIndex}
                        setActiveSortIndex={setActiveSortIndex}
                        activeSortDirection={activeSortDirection}
                        setActiveSortDirection={setActiveSortDirection}
                        columns={columns}
                        selected={selected}
                        allRowsSelected={allRowsSelected}
                        onSelect={onSelect}
                        onSelectAll={onSelectAll}
                    />
                </>
            )}
            <AlertGroup isToast isLiveRegion>
                {toasts.map(({ key, variant, title, children }: Toast) => (
                    <Alert
                        variant={AlertVariant[variant]}
                        title={title}
                        timeout={4000}
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
