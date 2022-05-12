import { gql } from '@apollo/client';

export type SummaryCountsResponse = {
    clusterCount: number;
    nodeCount: number;
    violationCount: number;
    deploymentCount: number;
    imageCount: number;
    secretCount: number;
};

export const SUMMARY_COUNTS = gql`
    query summary_counts {
        clusterCount
        nodeCount
        violationCount
        deploymentCount
        imageCount
        secretCount
    }
`;
