import React, { ReactElement, useEffect } from 'react';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { Select, SelectOption } from '@patternfly/react-core';

import { selectors } from 'reducers';
import { actions as graphActions, networkGraphClusters } from 'reducers/network/graph';
import { actions as pageActions } from 'reducers/network/page';

import useSelectToggle from 'hooks/patternfly/useSelectToggle';
import useURLCluster from 'hooks/useURLCluster';
import { Cluster } from 'types/cluster.proto';

type ClusterSelectProps = {
    id?: string;
    selectClusterId: (clusterId: string) => void;
    closeSidePanel: () => void;
    clusters: Cluster[];
    selectedClusterId?: string;
    isDisabled?: boolean;
};

const ClusterSelect = ({
    id,
    selectClusterId,
    closeSidePanel,
    clusters,
    selectedClusterId = '',
    isDisabled = false,
}: ClusterSelectProps): ReactElement => {
    const { closeSelect, isOpen, onToggle } = useSelectToggle();
    const { cluster, setCluster } = useURLCluster(selectedClusterId);

    useEffect(() => {
        selectClusterId(cluster || '');
    }, [cluster, selectClusterId]);

    // Set the clusterId in the URL to the default selected cluster, if it doesn't
    // yet exist in the URL
    useEffect(() => {
        if (!cluster) {
            setCluster(selectedClusterId);
        }
    }, [cluster, setCluster, selectedClusterId]);

    function changeCluster(_e, clusterId) {
        selectClusterId(clusterId);
        closeSelect();
        closeSidePanel();
        setCluster(clusterId);
    }

    return (
        <Select
            id={id}
            isOpen={isOpen}
            onToggle={onToggle}
            isDisabled={isDisabled || !clusters.length}
            selections={selectedClusterId}
            placeholderText="Select a cluster"
            onSelect={changeCluster}
        >
            {clusters
                .filter((c) => networkGraphClusters[c.type])
                .map(({ id: clusterId, name }) => (
                    <SelectOption
                        isSelected={clusterId === cluster}
                        key={clusterId}
                        value={clusterId}
                    >
                        {name}
                    </SelectOption>
                ))}
        </Select>
    );
};

const mapStateToProps = createStructuredSelector({
    clusters: selectors.getClusters,
    selectedClusterId: selectors.getSelectedNetworkClusterId,
});

const mapDispatchToProps = {
    selectClusterId: graphActions.selectNetworkClusterId,
    closeSidePanel: pageActions.closeSidePanel,
};

export default connect(mapStateToProps, mapDispatchToProps)(ClusterSelect);
