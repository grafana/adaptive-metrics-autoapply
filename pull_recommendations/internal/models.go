package internal

import (
	"github.com/prometheus/common/model"
)

const (
	DefaultSegmentName = "default"
)

var DefaultSegment = Segment{Name: DefaultSegmentName}

type Segment struct {
	Identifier        string `json:"id,omitempty"`
	Name              string `json:"name"`
	Selector          string `json:"selector,omitempty"`
	FallbackToDefault bool   `json:"fallback_to_default,omitempty"`
}

type Recommendation struct {
	Metric    string `json:"metric" yaml:"metric"`
	MatchType string `json:"match_type,omitempty" yaml:"match_type,omitempty"`

	Drop       bool     `json:"drop,omitempty" yaml:"drop,omitempty"`
	KeepLabels []string `json:"keep_labels,omitempty" yaml:"keep_labels,omitempty"`
	DropLabels []string `json:"drop_labels,omitempty" yaml:"drop_labels,omitempty"`

	Aggregations        []string       `json:"aggregations,omitempty" yaml:"aggregations,omitempty"`
	AggregationInterval model.Duration `json:"aggregation_interval,omitempty" yaml:"aggregation_interval,omitempty"`
	AggregationDelay    model.Duration `json:"aggregation_delay,omitempty" yaml:"aggregation_delay,omitempty"`

	Ingest bool `json:"ingest,omitempty" yaml:"ingest,omitempty"`

	ManagedBy string `json:"managed_by,omitempty"`

	RecommendedAction  string `json:"recommended_action,omitempty"`
	UsagesInRules      int    `json:"usages_in_rules,omitempty"`
	UsagesInQueries    int    `json:"usages_in_queries,omitempty"`
	UsagesInDashboards int    `json:"usages_in_dashboards,omitempty"`

	// Used when recommended action is "add"
	KeptLabels                   []string `json:"kept_labels,omitempty"`
	TotalSeriesAfterAggregation  int      `json:"total_series_after_aggregation,omitempty"`
	TotalSeriesBeforeAggregation int      `json:"total_series_before_aggregation,omitempty"`
}
