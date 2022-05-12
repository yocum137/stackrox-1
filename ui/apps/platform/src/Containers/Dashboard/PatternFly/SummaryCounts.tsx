import React from 'react';
import Raven from 'raven-js';
import pluralize from 'pluralize';
import { useQuery } from '@apollo/client';
import { Spinner, Split, Stack } from '@patternfly/react-core';

import { SummaryCountsResponse, SUMMARY_COUNTS } from 'queries/summaryCounts';

type SummaryTileProps = {
    label: string;
    value: number;
    isLoading: boolean;
    isError: boolean;
};

/*
  TODO
  We have a `TileContent` component that was used in the old version of this
  widget at ./apps/platform/src/Components/TileContent/TileContent.tsx

  Is this a pattern we want to retain in new PatternFly pages?

  PF has a similar looking `Tile` component, but it has different semantics from what
  we are doing here: https://www.patternfly.org/v4/components/tile
*/
function SummaryTile({ label, value, isLoading, isError }: SummaryTileProps) {
    let content;
    if (isLoading) {
        content = <Spinner size="lg" aria-label={`Total ${label} count`} />;
    } else if (isError) {
        content = <p title={`Failed to load ${label} counts`}>- -</p>;
    } else {
        content = (
            <span className="pf-u-font-size-lg-on-sm pf-u-font-size-sm pf-u-font-weight-bold">
                {value}
            </span>
        );
    }
    return (
        <Stack className="pf-u-px-md pf-u-px-lg-on-xl pf-u-py-xs pf-u-py-sm-on-xl pf-u-align-items-center">
            {content}
            <span className="pf-u-font-size-md-on-sm pf-u-font-size-xs">
                {pluralize(label, value)}
            </span>
        </Stack>
    );
}

function SummaryCounts() {
    const { loading, error, data } = useQuery<SummaryCountsResponse>(SUMMARY_COUNTS, {
        pollInterval: 30000,
    });

    if (error) {
        Raven.captureException(error);
    }

    const tileData = {
        Cluster: data?.clusterCount ?? 0,
        Node: data?.nodeCount ?? 0,
        Violation: data?.violationCount ?? 0,
        Deployment: data?.deploymentCount ?? 0,
        Image: data?.imageCount ?? 0,
        Secret: data?.secretCount ?? 0,
    };

    return (
        <Split className="pf-u-flex-wrap">
            {Object.entries(tileData).map(([label, value]) => (
                <SummaryTile
                    key={label}
                    label={label}
                    value={value}
                    isLoading={loading}
                    isError={!!error}
                />
            ))}
        </Split>
    );
}

export default SummaryCounts;
