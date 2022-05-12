import React, { ReactElement } from 'react';
import Raven from 'raven-js';
import { useQuery } from '@apollo/client';

import { SummaryCountsResponse, SUMMARY_COUNTS } from 'queries/summaryCounts';
import SummaryTileCount from 'Components/SummaryTileCount';

const SummaryCounts = (): ReactElement => {
    const { loading, error, data } = useQuery<SummaryCountsResponse>(SUMMARY_COUNTS, {
        pollInterval: 30000,
    });
    if (error) {
        Raven.captureException(error);
    }
    const { clusterCount, nodeCount, violationCount, deploymentCount, imageCount, secretCount } =
        data || {};
    return (
        <ul className="flex uppercase text-sm p-0">
            <SummaryTileCount label="Cluster" value={clusterCount} loading={loading} />
            <SummaryTileCount label="Node" value={nodeCount} loading={loading} />
            <SummaryTileCount label="Violation" value={violationCount} loading={loading} />
            <SummaryTileCount label="Deployment" value={deploymentCount} loading={loading} />
            <SummaryTileCount label="Image" value={imageCount} loading={loading} />
            <SummaryTileCount label="Secret" value={secretCount} loading={loading} />
        </ul>
    );
};

export default SummaryCounts;
