package effectiveaccessscope

import (
	"github.com/stackrox/rox/generated/storage"
)

// treeNodeAttributes stores additional information for a tree node.
type treeNodeAttributes struct {
	ID     string
	Name   string
	Labels map[string]string
}

func (t *treeNodeAttributes) copy() *treeNodeAttributes {
	labels := make(map[string]string, len(t.Labels))
	for k, v := range t.Labels {
		labels[k] = v
	}
	return &treeNodeAttributes{
		ID:     t.ID,
		Name:   t.Name,
		Labels: labels,
	}
}

func nodeAttributesForCluster(cluster *storage.Cluster, detail storage.ComputeEffectiveAccessScopeRequest_Detail) treeNodeAttributes {
	attributes := treeNodeAttributes{
		ID: cluster.GetId(),
	}
	if detail != storage.ComputeEffectiveAccessScopeRequest_MINIMAL {
		attributes.Name = cluster.GetName()
	}
	if detail == storage.ComputeEffectiveAccessScopeRequest_HIGH {
		attributes.Labels = cluster.GetLabels()
	}
	return attributes
}

func nodeAttributesForNamespace(namespace *storage.NamespaceMetadata, detail storage.ComputeEffectiveAccessScopeRequest_Detail) treeNodeAttributes {
	attributes := treeNodeAttributes{
		ID: namespace.GetId(),
	}
	if detail != storage.ComputeEffectiveAccessScopeRequest_MINIMAL {
		attributes.Name = namespace.GetName()
	}
	if detail == storage.ComputeEffectiveAccessScopeRequest_HIGH {
		attributes.Labels = namespace.GetLabels()
	}
	return attributes
}
