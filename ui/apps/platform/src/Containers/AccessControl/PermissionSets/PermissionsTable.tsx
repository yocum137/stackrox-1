import React, { ReactElement } from 'react';
import { Badge, SelectOption } from '@patternfly/react-core';
import { TableComposable, Tbody, Td, Thead, Th, Tr } from '@patternfly/react-table';

import SelectSingle from 'Components/SelectSingle';
import { accessControl as accessTypeLabels } from 'messages/common';
import { PermissionsMap } from 'services/RolesService';

import { ReadAccessIcon, WriteAccessIcon } from './AccessIcons';
import { getReadAccessCount, getWriteAccessCount } from './permissionSets.utils';
import ResourceDescription from './ResourceDescription';

export type PermissionsTableProps = {
    resourceToAccess: PermissionsMap;
    setResourceValue: (resource: string, value: string) => void;
    isDisabled: boolean;
};

export function PermissionsTable({
    resourceToAccess,
    setResourceValue,
    isDisabled,
}: PermissionsTableProps): ReactElement {
    const resourceToAccessEntries = Object.entries(resourceToAccess);

    return (
        <TableComposable variant="compact" isStickyHeader>
            <Thead>
                <Tr>
                    <Th width={20}>
                        Resource
                        <Badge isRead className="pf-u-ml-sm">
                            {resourceToAccessEntries.length}
                        </Badge>
                    </Th>
                    <Th width={40}>Description</Th>
                    <Th width={10}>
                        Read
                        <Badge isRead className="pf-u-ml-sm">
                            {getReadAccessCount(resourceToAccess)}
                        </Badge>
                    </Th>
                    <Th width={10}>
                        Write
                        <Badge isRead className="pf-u-ml-sm">
                            {getWriteAccessCount(resourceToAccess)}
                        </Badge>
                    </Th>
                    <Th width={20}>Access level</Th>
                </Tr>
            </Thead>
            <Tbody>
                {resourceToAccessEntries.map(([resource, accessLevel]) => (
                    <Tr key={resource}>
                        {isNew(resource) ? (
                            <Td dataLabel="Resource"><b>NEW!</b> {resource}</Td>
                        ) : (
                            <Td dataLabel="Resource">{resource}</Td>
                        )}
                        <Td dataLabel="Description">
                            <ResourceDescription resource={resource} />
                        </Td>
                        <Td dataLabel="Read" data-testid="read">
                            <ReadAccessIcon accessLevel={accessLevel} />
                        </Td>
                        <Td dataLabel="Write" data-testid="write">
                            <WriteAccessIcon accessLevel={accessLevel} />
                        </Td>
                        <Td dataLabel="Access level">
                            <SelectSingle
                                id={resource}
                                value={accessLevel}
                                handleSelect={setResourceValue}
                                isDisabled={isDisabled}
                            >
                                {Object.entries(accessTypeLabels).map(([id, name]) => (
                                    <SelectOption key={id} value={id}>
                                        {name}
                                    </SelectOption>
                                ))}
                            </SelectSingle>
                        </Td>
                    </Tr>
                ))}
            </Tbody>
        </TableComposable>
    );
}

export function SplitResourcesByDeprecation(resourceToAccess: PermissionsMap): [PermissionsMap, PermissionsMap] {
    let deprecated: PermissionsMap = {}
    let current: PermissionsMap = {}

    for (let r in resourceToAccess) {
        if (isDeprecated(r)) {
            deprecated[r] = resourceToAccess[r]
        } else {
            current[r] = resourceToAccess[r]
        }
    }

    return [current, deprecated]
}

function isDeprecated(resource: string): boolean {
    let deprecated = new Set([
        "AuthPlugin",
        "AuthProvider",             
        "Group",                    
        "Licenses",                 
        "Role",                     
        "User",                     
        "APIToken",                 
        "BackupPlugins",            
        "ImageIntegration",         
        "Notifier",                 
        "ComplianceRunSchedule",    
        "ComplianceRuns",           
        "AllComments",              
        "Config",                   
        "DebugLogs",                
        "NetworkGraphConfig",       
        "ProbeUpload",              
        "ScannerBundle",            
        "ScannerDefinitions",       
        "SensorUpgradeConfig",      
        "ServiceIdentity",          
        "Detection",                
        "NetworkBaseline",          
        "ProcessWhitelist",         
        "Risk",                     
        "WatchedImage",
    ])

    return deprecated.has(resource)
}

function isNew(resource: string): boolean {
    let newres = new Set([
        "Access",
        "Administration",
        "DeploymentExtension",
        "Integration",
    ])

    return newres.has(resource)
}