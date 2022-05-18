import React, { useEffect, useState } from 'react';
import { Link, useHistory } from 'react-router-dom';
import {
    Chart,
    ChartAxis,
    ChartStack,
    ChartBar,
    ChartTooltip,
    getTheme,
    ChartThemeColor,
    ChartLabel,
    ChartLabelProps,
} from '@patternfly/react-charts';

import { AlertGroup, fetchSummaryAlertCounts } from 'services/AlertsService';
import { severityLabels } from 'messages/common';
import { Card, Title } from '@patternfly/react-core';

type CountsBySeverity = {
    Low: Record<string, number>;
    Medium: Record<string, number>;
    High: Record<string, number>;
    Critical: Record<string, number>;
};

// TODO The PF provided colors do not use the mocks, which ones should we use here?
// Note that if we do not use the PF colors, we can ditch the @patternfly/patternfly dep and import
const colorScale = [
    // TODO Note: The colors on the PF docs site for charts may not be accurate
    'var(--pf-chart-color-black-300)',
    'var(--pf-chart-color-gold-300)',
    'var(--pf-chart-color-orange-300)',
    'var(--pf-chart-color-red-100)',
];

// Clone default PatternFly chart themes
const defaultTheme = getTheme(ChartThemeColor.multi);
const severityTheme = {
    ...defaultTheme,
    stack: {
        ...defaultTheme.stack,
        colorScale,
    },
    legend: {
        ...defaultTheme.legend,
        // We need to flip the theme since the bars run LOW->CRITICAL
        // and the legend runs CRITICAL->LOW
        colorScale: [...colorScale].reverse(),
    },
};

function getCountsBySeverity(groups: AlertGroup[]): CountsBySeverity {
    const result = {
        Low: {},
        Medium: {},
        High: {},
        Critical: {},
    };

    groups.forEach(({ group, counts }) => {
        // TODO These should default to zero
        result.Low[group] = 1;
        result.Medium[group] = 1;
        result.High[group] = 1;
        result.Critical[group] = 1;

        counts.forEach(({ severity, count }) => {
            result[severityLabels[severity]][group] = parseInt(count, 10);
        });
    });

    return result;
}

function LinkLabel(props: ChartLabelProps) {
    // TODO Construct this link in a better way
    const link = `/main/violations?s[Category]=${props.text as string}`;
    return (
        <Link to={link}>
            <ChartLabel {...props} style={{ fill: 'var(--pf-global--link--Color)' }} />
        </Link>
    );
}

function ViolationsByPolicyCategory() {
    const history = useHistory();
    const [violationCountsByPolicyCategories, setViolationCountsByPolicyCategories] = useState<
        AlertGroup[]
    >([]);

    useEffect(() => {
        fetchSummaryAlertCounts({ 'request.query': '', group_by: 'CATEGORY' })
            .then(setViolationCountsByPolicyCategories)
            .catch(() => {
                // TODO
            });
    }, []);

    const countsBySeverity = getCountsBySeverity(violationCountsByPolicyCategories);
    const bars = Object.entries(countsBySeverity).map(([severity, counts]) => {
        const data = Object.entries(counts).map(([group, count]) => ({
            name: severity,
            x: group,
            y: count,
            label: `${severity}: ${count}`,
        }));

        return (
            <ChartBar
                barWidth={16}
                key={severity}
                data={data}
                labelComponent={<ChartTooltip constrainToVisibleArea />}
                // TODO Gross - clean up events
                events={[
                    {
                        target: 'data',
                        eventHandlers: {
                            onClick: () => {
                                return [
                                    {
                                        mutation: (props) => {
                                            const link = `/main/violations?s[Category]=${
                                                props.datum.xName as string
                                            }&sortOption[field]=Severity&sortOption[direction]=asc`;
                                            history.push(link);
                                            return null;
                                        },
                                    },
                                ];
                            },
                        },
                    },
                ]}
            />
        );
    });

    // TODO How to abstract out common PF components for other widgets

    // TODO *Big Picture* - how to abstract out components and data providers to allow users to customize widgets?

    // TODO Handle sizing better
    return (
        <Card>
            <Title headingLevel="h2" className="pf-u-p-md">
                Policy Violations by Category
            </Title>
            <div style={{ height: '245px' }}>
                <Chart
                    ariaDesc="Number of violation by policy category, grouped by severity"
                    ariaTitle="Policy Violations by Category"
                    domainPadding={{ x: [30, 25] }}
                    legendData={[
                        { name: 'Critical' },
                        { name: 'High' },
                        { name: 'Medium' },
                        { name: 'Low' },
                    ]}
                    legendPosition="bottom"
                    // height, width, and padding need to be somewhat dynamic, or at least responsive
                    height={245}
                    width={450}
                    padding={{
                        left: 150, // left padding is dependent on the length of the text on the left axis
                        bottom: 75, // Adjusted to accommodate legend
                    }}
                    theme={severityTheme}
                >
                    <ChartAxis tickLabelComponent={<LinkLabel />} />
                    <ChartAxis dependentAxis showGrid />
                    <ChartStack horizontal>{bars}</ChartStack>
                </Chart>
            </div>
        </Card>
    );
}

export default ViolationsByPolicyCategory;
