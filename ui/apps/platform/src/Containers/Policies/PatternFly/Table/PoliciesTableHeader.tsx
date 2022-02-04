import React, { useState } from 'react';
import { useHistory } from 'react-router-dom';
import {
    Button,
    Divider,
    Dropdown,
    DropdownItem,
    DropdownSeparator,
    DropdownToggle,
    PageSection,
    Pagination,
    Title,
    Toolbar,
    ToolbarContent,
    ToolbarGroup,
    ToolbarItem,
    Tooltip,
} from '@patternfly/react-core';
import { CaretDownIcon } from '@patternfly/react-icons';

import { policiesBasePathPatternFly as policiesBasePath } from 'routePaths';
import { reassessPolicies } from 'services/PoliciesService';
import SearchFilterInput from 'Components/SearchFilterInput';
import useToasts from 'hooks/patternfly/useToasts';
import { ListPolicy } from 'types/policy.proto';
import { SearchFilter } from 'types/search';
import ImportPolicyJSONModal from '../Modal/ImportPolicyJSONModal';

import './PoliciesTable.css';

type PoliciesTableProps = {
    hasWriteAccessForPolicy: boolean;
    deletePoliciesHandler: (ids) => void;
    enablePoliciesHandler: (ids) => void;
    disablePoliciesHandler: (ids) => void;
    fetchPoliciesWithQuery: () => void;
    handleChangeSearchFilter: (searchFilter: SearchFilter) => void;
    searchFilter?: SearchFilter;
    searchOptions: string[];
    policiesCount: number;
    perPage: number;
    currentPage: number;
    setCurrentPage: (page) => void;
    setPerPage: (perPage) => void;
    selectedIds: string[];
    selectedPolicies: ListPolicy[];
};

function PoliciesTableHeader({
    hasWriteAccessForPolicy,
    deletePoliciesHandler,
    enablePoliciesHandler,
    disablePoliciesHandler,
    fetchPoliciesWithQuery,
    handleChangeSearchFilter,
    searchFilter,
    searchOptions,
    policiesCount,
    currentPage,
    setCurrentPage,
    perPage,
    setPerPage,
    selectedIds,
    selectedPolicies,
}: PoliciesTableProps): React.ReactElement {
    const history = useHistory();
    const { addToast } = useToasts();

    // Handle Bulk Actions dropdown state.
    const [isActionsOpen, setIsActionsOpen] = useState(false);

    const [isImportModalOpen, setIsImportModalOpen] = useState(false);

    function onToggleActions(toggleOpen) {
        setIsActionsOpen(toggleOpen);
    }

    function onSelectActions() {
        setIsActionsOpen(false);
    }

    // Handle page changes.
    function changePage(e, newPage) {
        if (newPage !== currentPage) {
            setCurrentPage(newPage);
        }
    }

    function changePerPage(e, newPerPage) {
        setPerPage(newPerPage);
    }

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

    let numEnabled = 0;
    let numDisabled = 0;
    let numDeletable = 0;
    selectedPolicies.forEach(({ disabled, isDefault }) => {
        if (disabled) {
            numDisabled += 1;
        } else {
            numEnabled += 1;
        }
        if (!isDefault) {
            numDeletable += 1;
        }
    });

    // TODO: https://stack-rox.atlassian.net/browse/ROX-8613
    // isDisabled={!hasSelections}
    // dropdownItems={hasWriteAccessForPolicy ? [Enable, Disable, Export, Delete] : [Export]} see PolicyDetail.tsx
    return (
        <>
            <PageSection
                variant="light"
                id="policies-table-header"
                padding={{ default: 'noPadding' }}
            >
                <Toolbar>
                    <ToolbarContent>
                        <ToolbarItem>
                            <Title headingLevel="h1">Policies</Title>
                        </ToolbarItem>
                        <ToolbarGroup
                            alignment={{ default: 'alignRight' }}
                            spaceItems={{ default: 'spaceItemsSm' }}
                            variant="button-group"
                        >
                            <ToolbarItem>
                                <Button variant="primary" onClick={onClickCreatePolicy}>
                                    Create policy
                                </Button>
                            </ToolbarItem>
                            <ToolbarItem>
                                <Button variant="secondary" onClick={onClickImportPolicy}>
                                    Import policy
                                </Button>
                            </ToolbarItem>
                        </ToolbarGroup>
                    </ToolbarContent>
                </Toolbar>
                <Divider component="div" />
                <Toolbar>
                    <ToolbarContent>
                        <ToolbarItem
                            variant="search-filter"
                            className="pf-u-flex-grow-1 pf-u-flex-shrink-1"
                        >
                            <SearchFilterInput
                                className="w-full theme-light"
                                handleChangeSearchFilter={handleChangeSearchFilter}
                                placeholder="Filter policies"
                                searchCategory="POLICIES"
                                searchFilter={searchFilter ?? {}}
                                searchOptions={searchOptions}
                            />
                        </ToolbarItem>
                        <ToolbarGroup
                            spaceItems={{ default: 'spaceItemsSm' }}
                            variant="button-group"
                        >
                            <ToolbarItem>
                                <Dropdown
                                    data-testid="policies-bulk-actions-dropdown"
                                    onSelect={onSelectActions}
                                    toggle={
                                        <DropdownToggle
                                            isDisabled={
                                                !hasWriteAccessForPolicy || selectedIds.length === 0
                                            }
                                            isPrimary
                                            onToggle={onToggleActions}
                                            toggleIndicator={CaretDownIcon}
                                        >
                                            Bulk actions
                                        </DropdownToggle>
                                    }
                                    isOpen={isActionsOpen}
                                    dropdownItems={[
                                        <DropdownItem
                                            key="Enable policies"
                                            component="button"
                                            isDisabled={numDisabled === 0}
                                            onClick={() => enablePoliciesHandler(selectedIds)}
                                        >
                                            {`Enable policies (${numDisabled})`}
                                        </DropdownItem>,
                                        <DropdownItem
                                            key="Disable policies"
                                            component="button"
                                            isDisabled={numEnabled === 0}
                                            onClick={() => disablePoliciesHandler(selectedIds)}
                                        >
                                            {`Disable policies (${numEnabled})`}
                                        </DropdownItem>,
                                        // TODO: https://stack-rox.atlassian.net/browse/ROX-8613
                                        // Export policies to JSON
                                        // onClick={() => exportPoliciesHandler(selectedIds, onClearAll)}
                                        // {`Export policies to JSON (${numSelected})`}
                                        <DropdownSeparator key="Separator" />,
                                        <DropdownItem
                                            key="Delete policy"
                                            component="button"
                                            isDisabled={numDeletable === 0}
                                            onClick={() =>
                                                deletePoliciesHandler(
                                                    selectedPolicies
                                                        .filter(({ isDefault }) => !isDefault)
                                                        .map(({ id }) => id)
                                                )
                                            }
                                        >
                                            {`Delete policies (${numDeletable})`}
                                        </DropdownItem>,
                                    ]}
                                />
                            </ToolbarItem>
                            <ToolbarItem>
                                <Tooltip content="Manually enrich external data">
                                    <Button variant="secondary" onClick={onClickReassessPolicies}>
                                        Reassess all
                                    </Button>
                                </Tooltip>
                            </ToolbarItem>
                        </ToolbarGroup>
                        <ToolbarItem variant="pagination" alignment={{ default: 'alignRight' }}>
                            <Pagination
                                isCompact
                                page={currentPage}
                                perPage={perPage}
                                onPerPageSelect={changePerPage}
                                onSetPage={changePage}
                                itemCount={policiesCount}
                            />
                        </ToolbarItem>
                    </ToolbarContent>
                </Toolbar>
            </PageSection>
            <ImportPolicyJSONModal
                isOpen={isImportModalOpen}
                cancelModal={() => {
                    setIsImportModalOpen(false);
                }}
                fetchPoliciesWithQuery={fetchPoliciesWithQuery}
            />
        </>
    );
}

export default PoliciesTableHeader;
