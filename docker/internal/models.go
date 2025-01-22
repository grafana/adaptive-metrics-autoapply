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

type RuleData struct {
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
}

type Recommendation struct {
	RuleData

	RecommendedAction  string `json:"recommended_action,omitempty"`
	UsagesInRules      int    `json:"usages_in_rules,omitempty"`
	UsagesInQueries    int    `json:"usages_in_queries,omitempty"`
	UsagesInDashboards int    `json:"usages_in_dashboards,omitempty"`

	// Used when recommended action is "add"
	KeptLabels []string `json:"kept_labels,omitempty"`

	RawSeriesCount         int `json:"raw_series_count,omitempty"`         // Number of series the client is sending.
	CurrentSeriesCount     int `json:"current_series_count,omitempty"`     // Number of series stored in the database.
	RecommendedSeriesCount int `json:"recommended_series_count,omitempty"` // Number of series after applying the recommendation.
}

func ConvertVerboseToRules(recs []Recommendation) []RuleData {
	rules := make([]RuleData, 0, len(recs))
	for _, rec := range recs {
		if rec.RecommendedAction == "remove" {
			continue
		}
		rules = append(rules, rec.RuleData)
	}
	return rules
}
