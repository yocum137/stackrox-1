import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Button, Flex, PageSection, Truncate } from '@patternfly/react-core';
import { TableComposable, Thead, Tbody, Tr, Th, Td } from '@patternfly/react-table';
import { CheckCircleIcon } from '@patternfly/react-icons';

import { ListPolicy } from 'types/policy.proto';
import { ActionItem } from 'Containers/Violations/ViolationsTablePanel';
import { TableColumn, SortDirection } from 'hooks/useTableSort';
import { policiesBasePathPatternFly as policiesBasePath } from 'routePaths';
import { NotifierIntegration } from 'types/notifier.proto';

import {
    LabelAndNotifierIdsForType,
    formatLifecycleStages,
    formatNotifierCountsWithLabelStrings,
    getLabelAndNotifierIdsForTypes,
} from '../policies.utils';
import PolicySeverityLabel from '../PolicySeverityLabel';

import './PoliciesTable.css';

type PoliciesTableProps = {
    notifiers: NotifierIntegration[];
    policies?: ListPolicy[];
    hasWriteAccessForPolicy: boolean;
    deletePoliciesHandler: (ids) => void;
    exportPoliciesHandler: (ids, onClearAll?) => void;
    enablePoliciesHandler: (ids) => void;
    disablePoliciesHandler: (ids) => void;
    columns: TableColumn[];
    activeSortIndex: number;
    setActiveSortIndex: (idx) => void;
    activeSortDirection: SortDirection;
    setActiveSortDirection: (dir) => void;
    selected: boolean[];
    allRowsSelected: boolean;
    onSelect: (
        event: React.FormEvent<HTMLInputElement>,
        isSelected: boolean,
        rowId: number
    ) => void;
    onSelectAll: (event: React.FormEvent<HTMLInputElement>, isSelected: boolean) => void;
};

function PoliciesTable({
    notifiers,
    policies = [],
    hasWriteAccessForPolicy,
    deletePoliciesHandler,
    exportPoliciesHandler,
    enablePoliciesHandler,
    disablePoliciesHandler,
    columns,
    activeSortIndex,
    setActiveSortIndex,
    activeSortDirection,
    setActiveSortDirection,
    selected,
    allRowsSelected,
    onSelect,
    onSelectAll,
}: PoliciesTableProps): React.ReactElement {
    const [labelAndNotifierIdsForTypes, setLabelAndNotifierIdsForTypes] = useState<
        LabelAndNotifierIdsForType[]
    >([]);

    useEffect(() => {
        setLabelAndNotifierIdsForTypes(getLabelAndNotifierIdsForTypes(notifiers));
    }, [notifiers]);

    function onSort(e, index, direction) {
        setActiveSortIndex(index);
        setActiveSortDirection(direction);
    }

    // TODO: https://stack-rox.atlassian.net/browse/ROX-8613
    // isDisabled={!hasSelections}
    // dropdownItems={hasWriteAccessForPolicy ? [Enable, Disable, Export, Delete] : [Export]} see PolicyDetail.tsx
    return (
        <PageSection
            isFilled
            padding={{ default: 'noPadding' }}
            hasOverflowScroll
            id="policies-table"
        >
            <TableComposable isStickyHeader>
                <Thead>
                    <Tr>
                        <Th
                            select={{
                                onSelect: onSelectAll,
                                isSelected: allRowsSelected,
                            }}
                        />
                        {columns.map(({ Header, width, sortField }, columnIndex) => {
                            const sortParams = sortField
                                ? {
                                      sort: {
                                          sortBy: {
                                              index: activeSortIndex,
                                              direction: activeSortDirection,
                                          },
                                          onSort,
                                          columnIndex,
                                      },
                                  }
                                : {};
                            return (
                                <Th key={Header} modifier="wrap" width={width} {...sortParams}>
                                    {Header}
                                </Th>
                            );
                        })}
                        <Th />
                    </Tr>
                </Thead>
                <Tbody>
                    {policies.map((policy, rowIndex) => {
                        const {
                            description,
                            disabled,
                            id,
                            isDefault,
                            lifecycleStages,
                            name,
                            notifiers: notifierIds,
                            severity,
                        } = policy;
                        const notifierCountsWithLabelStrings = formatNotifierCountsWithLabelStrings(
                            labelAndNotifierIdsForTypes,
                            notifierIds
                        );
                        const exportPolicyAction: ActionItem = {
                            title: 'Export policy to JSON',
                            onClick: () => exportPoliciesHandler([id]),
                        };
                        const actionItems = hasWriteAccessForPolicy
                            ? [
                                  disabled
                                      ? {
                                            title: 'Enable policy',
                                            onClick: () => enablePoliciesHandler([id]),
                                        }
                                      : {
                                            title: 'Disable policy',
                                            onClick: () => disablePoliciesHandler([id]),
                                        },
                                  exportPolicyAction,
                                  {
                                      isSeparator: true,
                                  },
                                  {
                                      title: 'Delete policy',
                                      onClick: () => deletePoliciesHandler([id]),
                                      disabled: isDefault,
                                  },
                              ]
                            : [exportPolicyAction];
                        return (
                            <Tr key={id}>
                                <Td
                                    select={{
                                        rowIndex,
                                        onSelect,
                                        isSelected: selected[rowIndex],
                                    }}
                                />
                                <Td dataLabel="Policy">
                                    <Button
                                        variant="link"
                                        isInline
                                        component={(props) => (
                                            <Link {...props} to={`${policiesBasePath}/${id}`} />
                                        )}
                                    >
                                        {name}
                                    </Button>
                                </Td>
                                <Td dataLabel="Description">
                                    <Truncate content={description || '-'} tooltipPosition="top" />
                                    {/* {description || '-'} */}
                                </Td>
                                <Td dataLabel="Status">
                                    {disabled ? (
                                        'Disabled'
                                    ) : (
                                        <Flex className="pf-u-info-color-200">
                                            <CheckCircleIcon className="pf-u-mr-sm pf-m-align-self-center" />
                                            Enabled
                                        </Flex>
                                    )}
                                </Td>
                                <Td dataLabel="Notifiers">
                                    {notifierCountsWithLabelStrings.length === 0 ? (
                                        '-'
                                    ) : (
                                        <>
                                            {notifierCountsWithLabelStrings.map(
                                                (notifierCountWithLabelString) => (
                                                    <div
                                                        key={notifierCountWithLabelString}
                                                        className="pf-u-text-nowrap"
                                                    >
                                                        {notifierCountWithLabelString}
                                                    </div>
                                                )
                                            )}
                                        </>
                                    )}
                                </Td>
                                <Td dataLabel="Severity">
                                    <PolicySeverityLabel severity={severity} />
                                </Td>
                                <Td dataLabel="Lifecycle">
                                    {formatLifecycleStages(lifecycleStages)}
                                </Td>
                                <Td
                                    actions={{
                                        items: actionItems,
                                    }}
                                />
                            </Tr>
                        );
                    })}
                </Tbody>
            </TableComposable>
        </PageSection>
    );
}

export default PoliciesTable;
