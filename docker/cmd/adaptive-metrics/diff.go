package main

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/google/go-cmp/cmp"

	"github.com/grafana/adaptive-metrics-autoapply/docker/internal"
)

func writeDiff(output io.Writer, segment internal.Segment, oldRec, newRec []internal.Recommendation) int {
	type stateChange struct {
		old, new internal.Recommendation
	}

	changesByName := map[string]stateChange{}
	for _, rule := range oldRec {
		changesByName[rule.Metric] = stateChange{old: rule}
	}

	for _, rule := range newRec {
		change, ok := changesByName[rule.Metric]
		if !ok {
			changesByName[rule.Metric] = stateChange{new: rule}
			continue
		}

		change.new = rule
		changesByName[rule.Metric] = change
	}

	var changes int
	var segmentOutput = new(strings.Builder)
	for _, change := range changesByName {
		if generateDiff(segmentOutput, change.old, change.new) {
			changes++
		}
	}

	if changes > 0 {
		diffOutput := segmentOutput.String()
		diffOutput = strings.Trim(diffOutput, "\n")
		fmt.Fprintf(output, "#### Segment %q:\n```diff\n%s\n```\n", segment.Name, diffOutput)
	}

	return changes
}

func generateDiff(output *strings.Builder, a, b internal.Recommendation) bool {

	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)
	rType := aVal.Type()

	metricName := a.Metric
	diffType := "~"
	if aVal.IsZero() {
		diffType = "+"
		metricName = b.Metric
	}
	if bVal.IsZero() {
		diffType = "-"
		metricName = a.Metric
	}

	metricOutput := new(strings.Builder)
	for i := 0; i < aVal.NumField(); i++ {
		fieldType := rType.Field(i)
		if fieldType.Name == "Metric" {
			continue
		}
		aField := aVal.Field(i)
		bField := bVal.Field(i)
		name := strings.Split(fieldType.Tag.Get("json"), ",")[0]

		if aField.IsZero() && bField.IsZero() {
			continue
		}

		if aField.IsZero() {
			fmt.Fprintf(metricOutput, "+\t%s=%q\n", name, bField.Interface())
			continue
		}

		if bField.IsZero() {
			fmt.Fprintf(metricOutput, "-\t%s=%q\n", name, aField.Interface())
			continue
		}

		d := cmp.Diff(aField.Interface(), bField.Interface())
		if d != "" {
			fmt.Fprintf(metricOutput, "~\t%s\n", name)
			fmt.Fprintln(metricOutput, d)
		}
	}

	if metricOutput.Len() > 0 {
		fmt.Fprintf(output, "%s%s\n%s", diffType, metricName, metricOutput.String())
		output.WriteString("\n")
		return true
	}

	return false
}
