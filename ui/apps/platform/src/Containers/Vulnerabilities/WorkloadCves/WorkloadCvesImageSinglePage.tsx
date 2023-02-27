import React, { ReactNode, useState } from 'react';
import {
    Breadcrumb,
    BreadcrumbItem,
    Bullseye,
    Divider,
    EmptyState,
    EmptyStateBody,
    EmptyStateIcon,
    EmptyStateVariant,
    Flex,
    Label,
    LabelGroup,
    PageSection,
    Skeleton,
    Tab,
    Tabs,
    TabsComponent,
    TabsProps,
    TabTitleText,
    Title,
    Tooltip,
} from '@patternfly/react-core';
import CopyIcon from '@patternfly/react-icons/dist/js/icons/copy-icon';
import { useQuery } from '@apollo/client';
import { useParams } from 'react-router-dom';

import BreadcrumbItemLink from 'Components/BreadcrumbItemLink';
import { vulnerabilitiesWorkloadCvesPath } from 'routePaths';
import { getDateTime, getDistanceStrictAsPhrase } from 'utils/dateUtils';
import EmptyStateTemplate from 'Components/PatternFly/EmptyStateTemplate';
import { getAxiosErrorMessage } from 'utils/responseErrorUtils';
import { ExclamationCircleIcon } from '@patternfly/react-icons';
import {
    imageDetailsQuery,
    ImageDetailsResponse,
    ImageDetailsVariables,
} from './queries/imageDetails';
import ImageSingleVulnerabilities from './ImageSingleVulnerabilities';
import ImageSingleResources from './ImageSingleResources';
import useDetailsTabParameter from './hooks/useDetailsTabParameter';
import { isDetailsTab } from './types';

const workloadCveOverviewImagePath = `${vulnerabilitiesWorkloadCvesPath}?cveStatusTab=Observed&entityTab=Image`;

function ImageDetailBadges({ imageData }: { imageData: ImageDetailsResponse['image'] }) {
    const [hasSuccessfulCopy, setHasSuccessfulCopy] = useState(false);

    const { deploymentCount, operatingSystem, metadata, scan } = imageData;
    const created = metadata?.v1?.created;
    const sha = metadata?.v1?.digest;

    function copyToClipboard(imageSha: string) {
        navigator.clipboard
            .writeText(imageSha)
            .then(() => setHasSuccessfulCopy(true))
            .catch(() => {
                // Permission is not required to write to the clipboard in secure contexts when initiated
                // via a user event so this Promise should not reject
            })
            .finally(() => {
                setTimeout(() => setHasSuccessfulCopy(false), 2000);
            });
    }

    return (
        <LabelGroup numLabels={Infinity}>
            {deploymentCount > 0 ? (
                <Label isCompact color="green">
                    Active
                </Label>
            ) : (
                <Label isCompact>Inactive</Label>
            )}
            {operatingSystem && <Label isCompact>OS: {operatingSystem}</Label>}
            {created && (
                <Label isCompact>Age: {getDistanceStrictAsPhrase(created, new Date())}</Label>
            )}
            {scan && (
                <Label isCompact>
                    Scan time: {getDateTime(scan.scanTime)} by {scan.dataSource.name}
                </Label>
            )}
            {sha && (
                <Tooltip content="Copy image SHA to clipboard">
                    <Label
                        style={{ cursor: 'pointer' }}
                        icon={<CopyIcon />}
                        isCompact
                        color={hasSuccessfulCopy ? 'green' : 'grey'}
                        onClick={() => copyToClipboard(sha)}
                    >
                        {hasSuccessfulCopy ? 'Copied!' : 'SHA'}
                    </Label>
                </Tooltip>
            )}
        </LabelGroup>
    );
}

function WorkloadCvesImageSinglePage() {
    const { imageId } = useParams();
    const { data, error } = useQuery<ImageDetailsResponse, ImageDetailsVariables>(
        imageDetailsQuery,
        {
            variables: { id: imageId },
        }
    );

    const [activeTabKey, setActiveTabKey] = useDetailsTabParameter();

    const imageData = data && data.image;
    const imageName = imageData?.name?.fullName ?? 'NAME UNKNOWN';

    const handleTabClick: TabsProps['onSelect'] = (e, tabKey) => {
        if (isDetailsTab(tabKey)) {
            setActiveTabKey(tabKey);
        }
    };

    let mainContent: ReactNode | null = null;

    if (error) {
        mainContent = (
            <PageSection variant="light">
                <Bullseye>
                    <EmptyState variant={EmptyStateVariant.large}>
                        <EmptyStateIcon
                            className="pf-u-danger-color-100"
                            icon={ExclamationCircleIcon}
                        />
                        <Title headingLevel="h2">{getAxiosErrorMessage(error)}</Title>
                    </EmptyState>
                </Bullseye>
            </PageSection>
        );
    } else {
        mainContent = (
            <>
                <PageSection variant="light">
                    {imageData ? (
                        <Flex direction={{ default: 'column' }}>
                            <Title headingLevel="h1" className="pf-u-mb-sm">
                                {imageName}
                            </Title>
                            <ImageDetailBadges imageData={imageData} />
                        </Flex>
                    ) : (
                        <Flex
                            direction={{ default: 'column' }}
                            spaceItems={{ default: 'spaceItemsXs' }}
                            className="pf-u-w-50"
                        >
                            <Skeleton fontSize="2xl" />
                            <Skeleton fontSize="sm" />
                        </Flex>
                    )}
                </PageSection>
                <PageSection variant="light" padding={{ default: 'noPadding' }}>
                    <Tabs
                        activeKey={activeTabKey}
                        onSelect={handleTabClick}
                        component={TabsComponent.nav}
                        className="pf-u-pl-md"
                        mountOnEnter
                        unmountOnExit
                    >
                        <Tab
                            eventKey="Vulnerabilities"
                            title={<TabTitleText>Vulnerabilities</TabTitleText>}
                        >
                            <ImageSingleVulnerabilities />
                        </Tab>
                        <Tab
                            eventKey="Resources"
                            title={<TabTitleText>Resources</TabTitleText>}
                            isDisabled
                        >
                            <ImageSingleResources />
                        </Tab>
                    </Tabs>
                </PageSection>
            </>
        );
    }

    return (
        <>
            <PageSection variant="light" className="pf-u-py-md">
                <Breadcrumb>
                    <BreadcrumbItemLink to={workloadCveOverviewImagePath}>
                        Images
                    </BreadcrumbItemLink>
                    {!error && (
                        <BreadcrumbItem isActive>
                            {imageData ? (
                                imageName
                            ) : (
                                <Skeleton screenreaderText="Loading image name" width="200px" />
                            )}
                        </BreadcrumbItem>
                    )}
                </Breadcrumb>
            </PageSection>
            <Divider component="div" />
            {mainContent}
        </>
    );
}

export default WorkloadCvesImageSinglePage;
