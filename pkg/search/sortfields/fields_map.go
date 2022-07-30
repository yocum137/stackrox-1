package sortfields

import (
	"github.com/stackrox/rox/generated/aux"
	"github.com/stackrox/rox/pkg/search"
)

// SortFieldMapper represents helper function that returns an array of query sort options to fulfill sorting by incoming sort option.
type SortFieldMapper func(option *auxpb.QuerySortOption) []*auxpb.QuerySortOption

var (
	// SortFieldsMap represents the mapping from searchable fields to sort field helper function
	SortFieldsMap = map[search.FieldLabel]SortFieldMapper{
		search.PolicyName: func(option *auxpb.QuerySortOption) []*auxpb.QuerySortOption {
			return []*auxpb.QuerySortOption{
				{
					Field:    search.SORTPolicyName.String(),
					Reversed: option.GetReversed(),
				},
			}
		},
		search.ImageName: func(option *auxpb.QuerySortOption) []*auxpb.QuerySortOption {
			return []*auxpb.QuerySortOption{
				{
					Field:    search.ImageRegistry.String(),
					Reversed: option.GetReversed(),
				},
				{
					Field:    search.ImageRemote.String(),
					Reversed: option.GetReversed(),
				},
				{
					Field:    search.ImageTag.String(),
					Reversed: option.GetReversed(),
				},
			}
		},
		search.Component: func(option *auxpb.QuerySortOption) []*auxpb.QuerySortOption {
			return []*auxpb.QuerySortOption{
				{
					Field:    search.Component.String(),
					Reversed: option.GetReversed(),
				},
				{
					Field:    search.ComponentVersion.String(),
					Reversed: option.GetReversed(),
				},
			}
		},
		search.LifecycleStage: func(option *auxpb.QuerySortOption) []*auxpb.QuerySortOption {
			return []*auxpb.QuerySortOption{
				{
					Field:    search.SORTLifecycleStage.String(),
					Reversed: option.GetReversed(),
				},
			}
		},
		search.NodePriority: func(option *auxpb.QuerySortOption) []*auxpb.QuerySortOption {
			return []*auxpb.QuerySortOption{
				{
					Field:    search.NodeRiskScore.String(),
					Reversed: !option.GetReversed(),
				},
			}
		},
		search.DeploymentPriority: func(option *auxpb.QuerySortOption) []*auxpb.QuerySortOption {
			return []*auxpb.QuerySortOption{
				{
					Field:    search.DeploymentRiskScore.String(),
					Reversed: !option.GetReversed(),
				},
			}
		},
		search.ImagePriority: func(option *auxpb.QuerySortOption) []*auxpb.QuerySortOption {
			return []*auxpb.QuerySortOption{
				{
					Field:    search.ImageRiskScore.String(),
					Reversed: !option.GetReversed(),
				},
			}
		},
		search.ComponentPriority: func(option *auxpb.QuerySortOption) []*auxpb.QuerySortOption {
			return []*auxpb.QuerySortOption{
				{
					Field:    search.ComponentRiskScore.String(),
					Reversed: !option.GetReversed(),
				},
			}
		},
	}
)
