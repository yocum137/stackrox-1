/* eslint-disable no-nested-ternary */
/* eslint-disable react/no-array-index-key */
import React, { ReactElement, useState } from 'react';

import { SearchFilter } from 'types/search';
import queryService from 'utils/queryService';
import useURLSort, { SortOption } from 'hooks/patternfly/useURLSort';
import useURLPagination from 'hooks/patternfly/useURLPagination';
import ObservedCVEsTable from './ObservedCVEsTable';
import useImageVulnerabilities from '../useImageVulnerabilities';

type ObservedCVEsProps = {
    imageId: string;
};

const sortFields = ['Severity', 'CVSS', 'Discovered'];
const defaultSortOption: SortOption = {
    field: 'Severity',
    direction: 'desc',
};

function ObservedCVEs({ imageId }: ObservedCVEsProps): ReactElement {
    const [searchFilter, setSearchFilter] = useState<SearchFilter>({});
    const { page, perPage, onSetPage, onPerPageSelect } = useURLPagination();
    const { sortOption, getSortParams } = useURLSort({
        sortFields,
        defaultSortOption,
    });

    const vulnsQuery = queryService.objectToWhereClause({
        ...searchFilter,
        'Vulnerability State': 'OBSERVED',
    });

    const { isLoading, data, refetchQuery } = useImageVulnerabilities({
        imageId,
        vulnsQuery,
        pagination: {
            limit: perPage,
            offset: (page - 1) * perPage,
            sortOption,
        },
    });

    const itemCount = data?.image?.vulnCount || 0;
    const rows = data?.image?.vulns || [];
    const registry = data?.image?.name?.registry || '';
    const remote = data?.image?.name?.remote || '';
    const tag = data?.image?.name?.tag || '';

    return (
        <ObservedCVEsTable
            rows={rows}
            registry={registry}
            remote={remote}
            tag={tag}
            isLoading={isLoading}
            itemCount={itemCount}
            page={page}
            perPage={perPage}
            onSetPage={onSetPage}
            onPerPageSelect={onPerPageSelect}
            getSortParams={getSortParams}
            updateTable={refetchQuery}
            searchFilter={searchFilter}
            setSearchFilter={setSearchFilter}
        />
    );
}

export default ObservedCVEs;
