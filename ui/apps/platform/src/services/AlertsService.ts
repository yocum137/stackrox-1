import queryString from 'qs';

import { Alert, ListAlert } from 'Containers/Violations/types/violationTypes';

import { SearchFilter } from 'types/search';
import { getRequestQueryStringForSearchFilter } from 'utils/searchUtils';
import axios from './instance';
import { RestSortOption } from './sortOption';

const baseUrl = '/v1/alerts';
const baseCountUrl = '/v1/alertscount';

// TODO import RestSearchOption and RestSortOption from searchUtils when it is TypeScript.
export type RestSearchOption = {
    label?: string;
    type?: string; // for example, 'categoryOption'
    value: string | string[];
};

// TODO import Severity from PoliciesService when it is TypeScript.
export type Severity = 'LOW_SEVERITY' | 'MEDIUM_SEVERITY' | 'HIGH_SEVERITY' | 'CRITICAL_SEVERITY';

export type AlertEventType = 'CREATED' | 'REMOVED';

export type AlertEvent = {
    time: string; // int64
    type: AlertEventType;
    id: string;
};

export type AlertEventsBySeverity = {
    severity: Severity;
    events: AlertEvent[];
};

export type ClusterAlert = {
    cluster: string;
    severities: AlertEventsBySeverity[];
};

type AlertsByTimeseriesFilters = {
    query: string;
};

/*
 * Fetch alerts by time for timeseries.
 */
export function fetchAlertsByTimeseries(
    filters: AlertsByTimeseriesFilters
): Promise<{ response: { clusters: ClusterAlert[] } }> {
    const params = queryString.stringify(filters);

    // set higher timeout for this call to handle known backend scale issues with dashboard
    return axios
        .get<{ clusters: ClusterAlert[] }>(`${baseUrl}/summary/timeseries?${params}`, {
            timeout: 59999,
        })
        .then((response) => ({
            response: response.data,
        }));
}

export type AlertCountBySeverity = {
    severity: Severity;
    count: number;
};

export type AlertGroup = {
    group: string;
    counts: AlertCountBySeverity[];
};

type SummaryAlertCountsFilters = {
    'request.query': string;
    group_by: string;
};

/*
 * Fetch severity counts.
 */
export function fetchSummaryAlertCounts(
    filters: SummaryAlertCountsFilters
): Promise<{ response: { groups: AlertGroup[] } }> {
    const params = queryString.stringify(filters);

    // set higher timeout for this call to handle known backend scale issues with dashboard
    return axios
        .get<{ groups: AlertGroup[] }>(`${baseUrl}/summary/counts?${params}`, { timeout: 59999 })
        .then((response) => ({
            response: response.data,
        }));
}

type CancelCallback = () => void;
type CancelablePromise<T> = Promise<T> & { cancel: CancelCallback };
export class CanceledPromiseError extends Error {
    constructor(message = 'The network request has been canceled.') {
        super(message);
        this.name = 'CanceledPromise';
    }
}

function makeCancelableAxiosRequest<T>(
    promiseProvider: (signal: AbortSignal) => Promise<T>
): CancelablePromise<T> {
    const abortController = new AbortController();
    const cancel = () => abortController.abort();
    const promise = promiseProvider(abortController.signal).catch((err) => {
        return err.message === 'canceled'
            ? Promise.reject(new CanceledPromiseError())
            : Promise.reject(err);
    });
    return Object.assign(promise, { cancel });
}

/*
 * Fetch a page of list alert objects.
 */
export function fetchAlerts(
    searchFilter: SearchFilter,
    sortOption: RestSortOption,
    page: number,
    pageSize: number
): CancelablePromise<ListAlert[]> {
    const offset = page > 0 ? page * pageSize : 0;
    const query = getRequestQueryStringForSearchFilter(searchFilter);
    const params = queryString.stringify(
        {
            query,
            pagination: {
                offset,
                limit: pageSize,
                sortOption,
            },
        },
        { arrayFormat: 'repeat', allowDots: true }
    );

    return makeCancelableAxiosRequest<ListAlert[]>((signal) =>
        axios
            .get<{ alerts: ListAlert[] }>(`${baseUrl}?${params}`, { signal })
            .then((response) => response?.data?.alerts ?? [])
    );
}

/*
 * Fetch count of alerts.
 */
export function fetchAlertCount(searchFilter: SearchFilter): CancelablePromise<number> {
    const params = queryString.stringify(
        { query: getRequestQueryStringForSearchFilter(searchFilter) },
        { arrayFormat: 'repeat' }
    );
    return makeCancelableAxiosRequest<number>((signal) =>
        axios
            .get<{ count: number }>(`${baseCountUrl}?${params}`, { signal })
            .then((response) => response?.data?.count ?? 0)
    );
}

/*
 * Fetch a specified alert.
 */
export function fetchAlert(id: string): Promise<Alert> {
    if (!id) {
        throw new Error('Alert ID must be specified');
    }
    return axios.get<Alert>(`${baseUrl}/${id}`).then((response) => response.data);
}

/*
 * Resolve an alert given an alert ID.
 */
export function resolveAlert(
    alertId: string,
    addToBaseline = false
): Promise<Record<string, never>> {
    return axios
        .patch<Record<string, never>>(`${baseUrl}/${alertId}/resolve`, { addToBaseline })
        .then((response) => response.data);
}

/*
 * Resolve a list of alerts by alert ID.
 */
export function resolveAlerts(
    alertIds: string[] = [],
    addToBaseline = false
): Promise<Record<string, never>[]> {
    return Promise.all(alertIds.map((id) => resolveAlert(id, addToBaseline)));
}
